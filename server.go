package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	commands chan command
}

func newServer() *server {
	return &server{
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("%s has connected to the sever", conn.RemoteAddr().String())

	c := &client{
		conn: conn,
		name: "anonymous",
		commands: s.commands,
	}

	c.readInput()
} 

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
}


func(s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	c.msg(c.name+": "+msg)
}

func (s *server) quit(c *client, args []string) {
	log.Printf("%s has left the chat", c.conn.RemoteAddr().String())
	c.msg("You have left the server")
	c.conn.Close()
}
