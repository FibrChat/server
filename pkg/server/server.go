package server

import (
	"fmt"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// Start creates and starts the NATS server
func Start(o Options) (*Server, error) {
	if o.Domain == "" {
		return nil, fmt.Errorf("Domain is required")
	}

	if o.WorkerPassword == "" {
		return nil, fmt.Errorf("WorkerPassword is required")
	}

	if o.RemotePassword == "" {
		return nil, fmt.Errorf("RemotePassword is required")
	}

	if o.Port == 0 {
		o.Port = 4222
	}

	opts := &natsserver.Options{
		ServerName: fmt.Sprintf("worker-%d", time.Now().UTC().UnixMilli()),
		Routes:     o.clusterRoutes(),
		Cluster:    o.clusterOpts(),
		Port:       o.Port + 1,
		NoLog:      true,
		NoSigs:     true,
		CustomClientAuthentication: &authHandler{
			workerPassword: o.WorkerPassword,
			remotePassword: o.RemotePassword,
		},
		Websocket: natsserver.WebsocketOpts{
			Port:  o.Port,
			NoTLS: true,
		},
	}

	ns, err := natsserver.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create NATS server: %w", err)
	}

	ns.Start()
	if !ns.ReadyForConnections(10 * time.Second) {
		return nil, fmt.Errorf("NATS failed to start within 10s")
	}

	return &Server{ns: ns, Opts: o}, nil
}

// Stop gracefully stops the NATS server.
func (s *Server) Stop() {
	s.ns.Shutdown()
	s.ns.WaitForShutdown()
}
