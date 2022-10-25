package internal

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	members map[net.Addr] *client
	commands chan command
}

func NewServer() {
	s := createServer()

	go s.run()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to start the server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Started server on :8080")
}

func NewClient() {
	for {
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			log.Printf("unable to accept new connection: %s", err.Error())
			continue
		}

		go createClient(conn)
	}
	
}

func createServer() *server {
	return &server{
		members: make(map[net.Addr]*client),
		commands: make(chan command),
	}
}

func createClient(conn net.Conn) {
	log.Printf("%s has connected to the sever", conn.RemoteAddr().String())
	s := createServer()
	c := &client{
		conn: conn,
		name: "anonymous",
		commands: s.commands,
	}
		s.addMember(c)
	for {
		c.readInput()
	}
}

func (s* server) addMember(c *client) {
	s.members[c.conn.RemoteAddr()] = c
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

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
}


func(s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	c.msg(c.name+": "+msg)
	s.broadcast(c, c.name+": "+msg)
}

func (s *server) quit(c *client, args []string) {
	log.Printf("%s has left the chat", c.name)
	c.msg("You have left the server")
	c.conn.Close()
}

func (s *server) broadcast(sender *client, msg string) {
	for addr, m := range s.members {
		if addr != sender.conn.RemoteAddr(){
			m.msg(msg)
		}
	}
} 