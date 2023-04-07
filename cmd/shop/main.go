package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"xsolla/cmd/shop/internal/adapters/address"
	"xsolla/cmd/shop/internal/adapters/cron"
	"xsolla/cmd/shop/internal/adapters/queue"
	"xsolla/cmd/shop/internal/adapters/repo"
	"xsolla/cmd/shop/internal/api/http"
	"xsolla/cmd/shop/internal/app"
	"xsolla/libs/serve"

	_ "github.com/lib/pq"
)

type (
	config struct {
		Server           server                 `yaml:"server"`
		DB               dbConfig               `yaml:"db"`
		Queue            queueConfig            `yaml:"queue"`
		Scheduler        schedulerConfig        `yaml:"scheduler"`
		AddressValidator addressValidatorConfig `yaml:"address_validator"`
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
	addressValidatorConfig struct {
		BasePath string `yaml:"base_path"`
		APIKey   string `yaml:"api_key"`
	}
)

var (
	cfgFile = flag.String("cfg", "./cmd/shop/config.yml", "path to config file")
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
		log.Fatal("shutdown", err)
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
	db, err := sqlx.Open(cfg.DB.Driver, fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", cfg.DB.User,
		cfg.DB.Password, cfg.DB.Dbname, cfg.DB.Mode, cfg.DB.HostDb, cfg.DB.PortDb))
	if err != nil {
		return fmt.Errorf("sqlx.Open: %w", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Println("close db", err)
		}
	}()

	r := repo.New(db)

	c := cron.New(cron.Config{
		TimeFetch: strconv.FormatInt(int64(cfg.Scheduler.TimeFetch), 10),
		Limit:     cfg.Scheduler.Limit,
	})

	q, err := queue.New(ctx, namespace, queue.Config{
		URLs:     cfg.Queue.URLS,
		Username: cfg.Queue.Username,
		Password: cfg.Queue.Password,
	})
	if err != nil {
		return fmt.Errorf("queue.New: %w", err)
	}
	defer func() {
		err := q.Close()
		if err != nil {
			log.Println("close queue connection", err)
		}
	}()

	a := address.New(address.Config{
		BasePath: cfg.AddressValidator.BasePath,
		APIKey:   cfg.AddressValidator.APIKey,
	})

	module := app.New(r, q, c, a)
	api := http.New(module)

	// it can be grpc, metrics, gateways...
	return serve.Start(ctx,
		serve.HTTP(cfg.Server.Host, cfg.Server.Port.HTTP, api),
		module.Process,
		q.Process)
}

func forceShutdown(ctx context.Context) {
	const shutdownDelay = 15 * time.Second

	<-ctx.Done()
	time.Sleep(shutdownDelay)

	log.Fatal("failed to graceful shutdown")
}
