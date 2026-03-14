package subject

const (
	PresenceSubject = "presence"
	PublishSubject  = "publish"
	UsersSubject    = "users"

	InboxPrefix     = "inbox"
	NATSInboxPrefix = "_INBOX"
)

// Inbox returns the subject for a user's inbox
func Inbox(username string) string {
	return InboxPrefix + "." + username
}

// NATSInbox returns the NATS reply inbox subject
func NATSInbox(username string) string {
	return NATSInboxPrefix + "." + username
}
