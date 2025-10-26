package configs

import "github.com/spf13/viper"

type Conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	RabbitMQHost      string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQUser      string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPassword  string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQPort      string `mapstructure:"RABBITMQ_PORT"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf

	// Set config file path
	viper.SetConfigFile(path + "/.env")
	viper.SetConfigType("env")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Bind environment variables to config keys
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("RABBITMQ_HOST")
	viper.BindEnv("RABBITMQ_USER")
	viper.BindEnv("RABBITMQ_PASSWORD")
	viper.BindEnv("RABBITMQ_PORT")
	viper.BindEnv("WEB_SERVER_PORT")
	viper.BindEnv("GRPC_SERVER_PORT")
	viper.BindEnv("GRAPHQL_SERVER_PORT")

	_ = viper.ReadInConfig()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
