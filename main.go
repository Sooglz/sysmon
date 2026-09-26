// Command sysmon 是轻量级 Linux 系统资源监控代理的入口。
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"sysmon/collector"
	"sysmon/config"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "/etc/sysmon/sysmon.yaml", "配置文件路径")
	listenAddr := flag.String("listen", "", "监听地址，覆盖配置文件")
	interval := flag.Int("interval", 0, "采集间隔秒，覆盖配置文件")
	showVersion := flag.Bool("version", false, "显示版本")
	flag.Parse()

	if *showVersion {
		fmt.Println("sysmon", version)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("load config failed", "err", err)
		os.Exit(1)
	}
	if *listenAddr != "" {
		cfg.ListenAddr = *listenAddr
	}
	if *interval > 0 {
		cfg.Interval = *interval
	}
	if err := cfg.Validate(); err != nil {
		slog.Error("invalid config", "err", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var current atomic.Value
	current.Store(collector.Metrics{})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
		defer ticker.Stop()
		collect := func() {
			m, err := collector.Collect(cfg.DiskPath)
			if err != nil {
				logger.Error("collect failed", "err", err)
				return
			}
			current.Store(m)
			logger.Info("collected", "cpu", m.CPUPercent, "mem", m.MemoryPercent, "disk", m.DiskPercent)
		}
		collect()
		for {
			select {
			case <-ticker.C:
				collect()
			case <-ctx.Done():
				return
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(current.Load().(collector.Metrics))
	})

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		logger.Info("sysmon started", "addr", cfg.ListenAddr, "interval", cfg.Interval)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
