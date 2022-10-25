package internal

import (
	"bufio"
	"net"
	"strings"
	"fmt"
)

type client struct {
	conn net.Conn
	name string
	commands chan<- command
}

func (c *client)readInput() {
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return
	}

	msg = strings.Trim(msg, "\r\n")
	args := strings.Split(msg, " ")

	cmd := strings.TrimSpace(args[0])

	switch cmd {
	case "/name": c.commands <- command{
		id: CMD_NAME,
		client: c,
		args: args,
	}
	case "/quit": c.commands <- command{
		id: CMD_QUIT,
		client: c,
		args: args,
	}
	case "/msg": c.commands <- command{
		id: CMD_MSG,
		client: c,
		args: args,
	}
	default:
		c.err(fmt.Errorf("unknown command: %s", cmd))
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}