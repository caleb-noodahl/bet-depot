package config

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

type APIConf struct {
	BaseUrl         string `yaml:"base_url"`
	Port            int    `yaml:"port"`
	DBHost          string `yaml:"db_host"`
	DBPort          int    `yaml:"db_port"`
	DBName          string `yaml:"db_name"`
	DBUser          string `yaml:"db_user"`
	DBPassword      string `yaml:"db_password"`
	DiscordClientID string `yaml:"discord_client_id"`
	DiscordSecret   string `yaml:"discord_secret"`
	DiscordRedirect string `yaml:"discord_redirect"`
	DiscordBotToken string `yaml:"discord_bot_token"`
	DiscordGuild    string `yaml:"discord_guild"`
	OpenApiKey      string `yaml:"open_api_key"`
	FinnhubApiKey   string `yaml:"finnhub_api_key"`
}

type GamesConfig struct {
	AdminID string   `yaml:"admin_id"`
	Symbols []string `yaml:"symbols"`
}

func ParseGamesConf(gamesConfigBytes []byte) (*GamesConfig, error) {
	c := new(GamesConfig)
	return c, yaml.Unmarshal(gamesConfigBytes, &c)
}

func ParseAPIConf(apiConfBytes []byte) (*APIConf, error) {
	c := new(APIConf)
	return c, yaml.Unmarshal(apiConfBytes, &c)
}
