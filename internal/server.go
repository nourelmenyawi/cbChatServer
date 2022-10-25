package internal

import (
	"fmt"
	"log"
	"net"
	"strings"
//	"sync"
	"bufio"
//	"os"
)

type server struct {
	members map[net.Addr] *client
	commands chan command
}

func NewServer() {
	s:= createServer()
	go s.run()

	//wg:= sync.WaitGroup{} 

	//Start TCP server
	listener, err := net.Listen("tcp", ":8800")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("Started server on :8800")

	//Accept new clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		}
		go s.createClient(conn)
	}
}

func createServer() *server {
	return &server{
		members: make(map[net.Addr]*client),
		commands: make(chan command),
	}
}

func (s *server)createClient(conn net.Conn) {
	log.Printf("%s has connected to the sever", conn.RemoteAddr().String())
	c := &client{
		conn: conn,
		name: "anonymous",
		commands: s.commands,
	}
//	if s.checkPassword(c) == true {	
		s.members[c.conn.RemoteAddr()] = c	
		s.readInput(c)
//	}
}

// func (s *server) checkPassword(c *client) bool{
// 	c.msg("Enter Password: ")
// 	scanner := bufio.NewScanner(os.Stdin)

// 	if scanner == "Password"{
// 		return true
// 	} else return false	
// }

func (s *server)readInput(c *client) {
	for{
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
	s.broadcast(c, c.name+": "+msg)
}

func (s *server) quit(c *client, args []string) {
	message := c.name + " has left the chat"
	c.msg("You have left the server")
	c.conn.Close()
	s.broadcast(c, message)
}

func (s *server) broadcast(sender *client, msg string) {
	for _, m := range s.members {
		m.conn.Write([]byte("> " + msg + "\n"))
	}
} 