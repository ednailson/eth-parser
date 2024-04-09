package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ednailson/eth-parser/internal/domain"
	"github.com/ednailson/eth-parser/pkg/blockchain"
	"github.com/ednailson/eth-parser/pkg/http"
	"github.com/ednailson/eth-parser/pkg/store"
)

var (
	port = 8080
	host = "https://cloudflare-eth.com"
)

func main() {
	flag.StringVar(&host, "host", host, "blockchain json rpc server host")
	flag.IntVar(&port, "port", port, "http server port")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})))

	bc := blockchain.NewService(host)
	s := store.NewMemory()

	logic := domain.NewLogic(bc, s)

	server := http.NewServer(port, logic)

	slog.Info("HTTP server started!")
	err := server.Serve()
	if err != nil {
		slog.Error("Failed to serve http", "error", err)
		os.Exit(1)
	}

	slog.Info("Server has been shut down!")
}
