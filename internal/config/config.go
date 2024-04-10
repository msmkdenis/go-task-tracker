package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	URLServer string
	TokenName string
	Secret    string
	TokenTTL  int
	Salt      string
	DBFile    string
}

func New() *Config {
	c := Config{}
	c.parseFlags()
	c.parseEnv()
	return &c
}

func (c *Config) parseFlags() {
	var URLServer string
	flag.StringVar(&URLServer, "s", ":7540", "Enter URLServer as ip_address:port Or use SERVER_ADDRESS env")
	var TokenName string
	flag.StringVar(&TokenName, "t", "token", "Enter token name Or use TOKEN_NAME env")
	var SecretKey string
	flag.StringVar(&SecretKey, "k", "secret", "Enter secret Or use SECRET_KEY env")
	var TokenTTL int
	flag.IntVar(&TokenTTL, "l", 64, "Enter token ttl in hours Or use TOKEN_TTL env")
	var Salt string
	flag.StringVar(&Salt, "a", "super-salty-salt", "Enter salt Or use SALT env")
	var DBFile string
	flag.StringVar(&DBFile, "d", "scheduler.db", "Enter salt Or use SALT env")
	flag.Parse()
	c.URLServer = URLServer
	c.TokenName = TokenName
	c.Secret = SecretKey
	c.TokenTTL = TokenTTL
	c.Salt = Salt
	c.DBFile = DBFile
}

func (c *Config) parseEnv() {
	if envURLServer := os.Getenv("SERVER_ADDRESS"); envURLServer != "" {
		c.URLServer = envURLServer
	}
	if envTokenName := os.Getenv("TOKEN_NAME"); envTokenName != "" {
		c.TokenName = envTokenName
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		c.Secret = envSecretKey
	}
	if envTokenTTL := os.Getenv("TOKEN_TTL"); envTokenTTL != "" {
		tokenTTL, err := strconv.Atoi(envTokenTTL)
		if err != nil {
			c.TokenTTL = 8
		}
		c.TokenTTL = tokenTTL
	}
	if envSalt := os.Getenv("SALT"); envSalt != "" {
		c.Salt = envSalt
	}
	if envDBFile := os.Getenv("DB_FILE"); envDBFile != "" {
		c.DBFile = envDBFile
	}
}
