package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DailyGoal int

	PomodoroDuration  time.Duration
	BreakDuration     time.Duration
	LongBreakDuration time.Duration
}

func Load(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("pomo.daily-goal", 0)

	viper.SetDefault("timer.pomodoro", "25m")
	viper.SetDefault("timer.break", "5m")
	viper.SetDefault("timer.long-break", "15m")

	err := viper.SafeWriteConfig()
	if err != nil {
		var alreadyExistsErr viper.ConfigFileAlreadyExistsError
		if !errors.As(err, &alreadyExistsErr) {
			return Config{}, fmt.Errorf("initializing config: %w", err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}

	return Config{
		DailyGoal: viper.GetInt("pomo.daily-goal"),

		PomodoroDuration:  viper.GetDuration("timer.pomodoro"),
		BreakDuration:     viper.GetDuration("timer.break"),
		LongBreakDuration: viper.GetDuration("timer.long-break"),
	}, nil
}
