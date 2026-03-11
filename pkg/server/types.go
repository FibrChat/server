package server

import (
	natsserver "github.com/nats-io/nats-server/v2/server"
)

type Options struct {
	Domain          string
	Port            int
	ClusterPort     int
	ClusterPeers    []string
	ClusterPassword string
	WorkerPassword  string
	RemotePassword  string
}

type Server struct {
	ns   *natsserver.Server
	Opts Options
}

type authHandler struct {
	workerPassword string
	remotePassword string
}
