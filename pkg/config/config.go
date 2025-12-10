package config

import (
	"github.com/joho/godotenv"

	kenv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

func LoadConfig(config any) error {

	_ = godotenv.Load(".env")

	k := koanf.New(".")

	k.Load(kenv.Provider("", ".", func(s string) string {
		return s
	}), nil)

	if err := k.Unmarshal("", &config); err != nil {
		return err
	}

	return nil
}
