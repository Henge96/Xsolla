package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type (
	config struct {
		Server    server          `yaml:"server"`
		DB        dbConfig        `yaml:"db"`
		Queue     queueConfig     `yaml:"queue"`
		Scheduler schedulerConfig `yaml:"scheduler"`
	}
	server struct {
		Host string `yaml:"host"`
		Port ports  `yaml:"port"`
	}
	ports struct {
		HTTP uint16 `yaml:"http"`
	}
	dbConfig struct {
		Driver   string `yaml:"driver"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		HostDb   string `yaml:"hostdb"`
		PortDb   string `yaml:"portdb"`
		Dbname   string `yaml:"dbname"`
		Mode     string `yaml:"mode"`
	}
	queueConfig struct {
		URLS     []string `yaml:"urls"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	}
	schedulerConfig struct {
		TimeFetch uint16 `yaml:"time_fetch"`
		Limit     uint16 `yaml:"limit"`
	}
)

var (
	cfgFile = flag.String("cfg", "config.yml", "path to config file")
)

func main() {
	flag.Parse()

	appName := filepath.Base(os.Args[0])

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	defer cancel()
	go forceShutdown(ctx)

	err := start(ctx, *cfgFile, appName)
	if err != nil {
		// Because if we have fatal we will close servers anyway.
		log.Fatal("shutdown")
	}
}

func start(ctx context.Context, configPath string, appName string) error {
	cFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}

	var cfg config
	err = yaml.NewDecoder(cFile).Decode(&cfg)
	if err != nil {
		return fmt.Errorf("yaml.NewDecoder.Decode: %w", err)
	}

	return run(ctx, cfg, appName)
}

func run(ctx context.Context, cfg config, namespace string) error {
	return nil
}

func forceShutdown(ctx context.Context) {
	const shutdownDelay = 15 * time.Second

	<-ctx.Done()
	time.Sleep(shutdownDelay)

	log.Fatal("failed to graceful shutdown")
}
