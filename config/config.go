package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PomodoroDuration  time.Duration
	BreakDuration     time.Duration
	LongBreakDuration time.Duration
}

func Load(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("timer.pomodoro", "25m")
	viper.SetDefault("timer.break", "5m")
	viper.SetDefault("timer.long-break", "15m")

	err := viper.SafeWriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return Config{}, fmt.Errorf("initializing config: %w", err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}

	return Config{
		PomodoroDuration:  viper.GetDuration("timer.pomodoro"),
		BreakDuration:     viper.GetDuration("timer.break"),
		LongBreakDuration: viper.GetDuration("timer.long-break"),
	}, nil
}
