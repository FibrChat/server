package server

import (
	"time"

	"github.com/fibrchat/server/pkg/subject"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

// Check implements custom authentication logic for NATS clients
func (h *authHandler) Check(client natsserver.ClientAuthentication) bool {
	o := client.GetOpts()

	switch o.Username {
	case "worker":
		return h.handleWorker(client, o.Password)

	case "remote":
		return h.handleRemote(client, o.Password)

	default:
		return h.handleClient(client, o.Username, o.Password)
	}
}

// handleWorker registers a worker client (full access)
func (h *authHandler) handleWorker(client natsserver.ClientAuthentication, password string) bool {
	if password != h.workerPassword {
		return false
	}

	client.RegisterUser(&natsserver.User{
		Username: "worker",
		Password: h.workerPassword,
		Permissions: &natsserver.Permissions{
			Publish:   &natsserver.SubjectPermission{Allow: []string{">"}},
			Subscribe: &natsserver.SubjectPermission{Allow: []string{">"}},
			Response: &natsserver.ResponsePermission{
				MaxMsgs: 1,
				Expires: 5 * time.Minute,
			},
		},
	})

	return true
}

// handleRemote registers a remote server client (send only)
func (h *authHandler) handleRemote(client natsserver.ClientAuthentication, password string) bool {
	if password != h.remotePassword {
		return false
	}

	client.RegisterUser(&natsserver.User{
		Username: "remote",
		Password: h.remotePassword,
		Permissions: &natsserver.Permissions{
			Publish:   &natsserver.SubjectPermission{Allow: []string{subject.Remote}},
			Subscribe: &natsserver.SubjectPermission{Allow: []string{}},
		},
	})

	return true
}

// handleClient registers a regular client
func (h *authHandler) handleClient(client natsserver.ClientAuthentication, username, password string) bool {
	// Temporary auth until implementation
	if username == "" || password != "password" {
		return false
	}

	client.RegisterUser(&natsserver.User{
		Username: username,
		Permissions: &natsserver.Permissions{
			Publish: &natsserver.SubjectPermission{Allow: []string{subject.Send}},
			Subscribe: &natsserver.SubjectPermission{Allow: []string{
				subject.DM(username),
				subject.Inbox(username) + ".>",
			}},
		},
	})

	return true
}
