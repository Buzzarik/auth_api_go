package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 		string 			`yaml:"env"`
	Redis 		ConfigRedis 	`yaml:"redis"`
	Postgres 	ConfigPostgres 	`yaml:"postgres"`
	Server 		ConfigServer	`yaml:"server_auth"`
};

type ConfigPostgres struct {
	Driver      	string			`yaml:"driver"`
	Host        	string			`yaml:"host"`
	Port        	int64			`yaml:"port"`
	UserName    	string			`yaml:"username"`
	Password    	string			`yaml:"password"`
	Sslmode     	string			`yaml:"sslmode"`
	DbName      	string			`yaml:"db_name"`
	MaxIdleTime 	time.Duration 	`yaml:"max_idle_time"`
	MaxOpenConns 	int64			`yaml:"max_open_conns"`
	MaxIdleConns 	int64			`yaml:"max_idle_conns"`
};

type ConfigRedis struct {
	Host 		string 	`yaml:"host"`
	Port 		int64 	`yaml:"port"`
	Password 	string 	`yaml:"password"`
	Db			int 	`yaml:"db"`
};

type ConfigServer struct {
	Host 		string 			`yaml:"host"`
	Port 		int64 			`yaml:"port"`
	Timeout 	time.Duration 	`yaml:"timeout"`
	IdleTimeout time.Duration 	`yaml:"idle_timeout"`
	Secret		string 			`yaml:"secret"`
	TokenTTL 	time.Duration 	`yaml:"token_ttl"`
	IdAPI 		int64			`yaml:"id_api"`
};

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH");
	if configPath == ""{
		log.Fatal("CONFIG_PATH is not set");
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config;

	if err := cleanenv.ReadConfig(configPath, &config); err != nil{
		log.Fatalf("cannot read config: %s", err)
	}

	return &config;
}
