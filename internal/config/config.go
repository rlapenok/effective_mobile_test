package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/client"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/repository"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	Server   serverConfig   `mapstructure:",squash"`
	Postgres postgresConfig `mapstructure:",squash"`
	Logging  loggingConfig  `mapstructure:",squash"`
	Client   clientConfig   `mapstructure:",squash"`
}

type loggingConfig struct {
	Level string `mapstructure:"LOGGING_LEVEL"`
}

type serverConfig struct {
	Port int `mapstructure:"PORT"`
}

type clientConfig struct {
	Url string `mapstructure:"URL"`
}

type postgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     uint   `mapstructure:"POSTGRES_PORT"`
	Db       string `mapstructure:"POSTGRES_DB"`
	Login    string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Migrate  string `mapstructure:"POSTGRES_MIGRATE"`
}

func (cfg postgresConfig) toSqlxDb() *sqlx.DB {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Login,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db)
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		zap.L().Error("Error while connect to Postgres", zap.Error(err))
		os.Exit(1)
	}
	return db
}

func Load(path string) AppConfig {
	var config AppConfig

	if err := godotenv.Load(); err != nil {
		slog.Error("Error while load .env", "err", err)
		os.Exit(1)
	}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error while read config from file", "err", err)
		os.Exit(1)

	}
	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Error while map config to strcut", "err", err)
		os.Exit(1)
	}
	return config

}

func (cfg AppConfig) GetLevel() *string {
	return &cfg.Logging.Level
}

func (cfg AppConfig) createRepository() *repository.Repository {
	db := cfg.Postgres.toSqlxDb()
	return repository.New(db)
}

func (cfg AppConfig) ToState() {
	repo := cfg.createRepository()
	repo.RunMigration(cfg.Postgres.Migrate)
	client := client.New(&cfg.Client.Url)
	state.New(repo, client)
}
