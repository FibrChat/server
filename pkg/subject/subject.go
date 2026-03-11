package subject

const Send = "chat.send"
const Remote = "chat.remote"

// DM returns the NATS subject for a user's direct messages
func DM(username string) string {
	return "chat.dm." + username
}

// Inbox returns the NATS subject for a user's inbox
func Inbox(username string) string {
	return "_INBOX." + username
}
