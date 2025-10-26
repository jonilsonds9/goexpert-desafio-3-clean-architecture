package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/configs"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/internal/event/handler"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/internal/infra/graph"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/internal/infra/grpc/pb"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/internal/infra/grpc/service"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/internal/infra/web/webserver"
	"github.com/jonilsonds9/goexpert-desafio-3-clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	rootDir := findProjectRoot()
	configs, err := configs.LoadConfig(rootDir)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(rootDir, db); err != nil {
		panic(err)
	}

	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("POST /order", webOrderHandler.Create)
	webserver.AddHandler("GET /orders", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func findProjectRoot() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "" // chegou na raiz do FS
		}
		dir = parent
	}
}

func runMigrations(rootDir string, db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	migrationsPath := filepath.Join(rootDir, "internal/infra/database/migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations executed successfully")
	return nil
}

func getRabbitMQChannel(cfg *configs.Conf) *amqp.Channel {
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
