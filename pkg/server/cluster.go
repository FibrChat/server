package server

import (
	"log"
	"net/url"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

// clusterOpts returns the ClusterOpts for this server
func (o *Options) clusterOpts() natsserver.ClusterOpts {
	if o.ClusterPort == 0 {
		return natsserver.ClusterOpts{}
	}

	clusterOpts := natsserver.ClusterOpts{
		Name:     "cluster-" + o.Domain,
		Password: o.ClusterPassword,
		Port:     o.ClusterPort,
		Username: "cluster",
	}

	return clusterOpts
}

// clusterRoutes returns the list of cluster route URLs for this server
func (o *Options) clusterRoutes() []*url.URL {
	if len(o.ClusterPeers) == 0 || o.ClusterPort == 0 {
		return nil
	}

	routes := make([]*url.URL, 0, len(o.ClusterPeers))
	for _, peer := range o.ClusterPeers {
		u, err := url.Parse("nats-route://cluster:" + o.ClusterPassword + "@" + peer)
		if err != nil {
			log.Fatalf("[server] invalid cluster peer URL %q: %v", peer, err)
		}
		routes = append(routes, u)
	}

	return routes
}
