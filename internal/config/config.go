package config

import (
	"github.com/chris102994/global-entry-appointment-bot/internal/cron"
	"github.com/chris102994/global-entry-appointment-bot/internal/logging"
	"github.com/chris102994/global-entry-appointment-bot/internal/lookup"
	"github.com/chris102994/global-entry-appointment-bot/internal/notify"
	"github.com/chris102994/global-entry-appointment-bot/internal/run"
	"github.com/spf13/viper"
)

type Config struct {
	Logging *logging.Logging `mapstructure:",omitempty"`
	Run     *run.Run         `mapstructure:",omitempty"`
	Cron    *cron.Cron       `mapstructure:",omitempty"`
	Notify  *notify.Notify   `mapstructure:",omitempty"`
	Lookup  *lookup.Lookup   `mapstructure:",omitempty"`
}

func (c *Config) LoadConfig() error {
	err := viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	return nil
}
