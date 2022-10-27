package internal

type commandID int

const (
	CMD_NAME commandID = iota
	CMD_MSG
	CMD_QUIT
	CMD_SHOUT
	CMD_SPAM
	CMD_WHISPER
	CMD_LIST
	CMD_HELP
)

type command struct {
	id commandID
	client *client 
	args []string
}