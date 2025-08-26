package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Host  string
	Debug bool
}

func ReadConfig() Config {
	var host string
	flag.StringVar(&host, "host", ":8080", "server host")
	debug := flag.Bool("debug", false, "enable debug logging lvl")
	flag.Parse()
	hostEnv := os.Getenv("SERVER_HOST") //хост взяли из переменной окружения (echo $SERVER_HOST)
	if hostEnv != "" && host == ":8080" {
		host = hostEnv
	}
	log.Println("host: ", host)
	return Config{
		Host:  host,
		Debug: *debug,
	}
}
