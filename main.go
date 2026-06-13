package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/YagoSchramm/GoNotify/service"
)

func main() {
	router, cleanup, err := service.Build()
	if err != nil {
		slog.Error("failed to build service", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	slog.Info("starting http server", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
