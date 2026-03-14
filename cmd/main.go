package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/fibrchat/server/pkg/server"
)

func main() {
	domain := mustEnv("DOMAIN")
	wsPort := envInt("PORT", 4222)
	clusterPort := envInt("CLUSTER_PORT", 0)
	workerPassword := mustEnv("WORKER_PASSWORD")
	clusterPassword := envOr("CLUSTER_PASSWORD", "")

	var clusterPeers []string
	if v := os.Getenv("CLUSTER_PEERS"); v != "" {
		for p := range strings.SplitSeq(v, ",") {
			if p = strings.TrimSpace(p); p != "" {
				clusterPeers = append(clusterPeers, p)
			}
		}
	}

	ns, err := server.Start(server.Options{
		Domain:          domain,
		Port:            wsPort,
		ClusterPort:     clusterPort,
		ClusterPeers:    clusterPeers,
		ClusterPassword: clusterPassword,
		WorkerPassword:  workerPassword,
	})
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")
	ns.Stop()
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid value for %s: %v", key, err)
	}
	return n
}
