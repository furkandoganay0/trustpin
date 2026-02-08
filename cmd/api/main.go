package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"trustpin/internal/audit"
	"trustpin/internal/cache"
	"trustpin/internal/config"
	"trustpin/internal/http"
	"trustpin/internal/migrate"
	"trustpin/internal/push"
	"trustpin/internal/store"
	"trustpin/internal/totp"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pg, err := store.Open(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer pg.Close()
	if err := migrate.Apply(ctx, pg.DB); err != nil {
		log.Fatalf("db migrate failed: %v", err)
	}

	redis := cache.Open(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer redis.Close()

	auditLogger := audit.Logger{Store: pg}
	server := httpapi.New(cfg, pg, redis, auditLogger, push.Mock{}, totp.Disabled{})
	log.Printf("trustpin listening on %s", cfg.HTTPAddr)
	if err := server.Serve(ctx); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
