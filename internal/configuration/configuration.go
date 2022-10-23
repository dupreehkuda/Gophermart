package configuration

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type config struct {
	Address        string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabasePath   string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New(logger *zap.Logger) *config {
	var config = config{}
	var err = env.Parse(&config)
	if err != nil {
		logger.Error("Error occurred when parsing config", zap.Error(err))
	}

	flag.StringVar(&config.Address, "a", config.Address, "getting launch address")
	flag.StringVar(&config.DatabasePath, "d", config.DatabasePath, "getting path to database")
	flag.StringVar(&config.AccrualAddress, "r", config.AccrualAddress, "getting accrual system address")
	flag.Parse()

	os.Setenv("secret", "vladimir")

	return &config
}
