package internal

import (
	"fmt"
	"log"
	"net"
	"strings"

	//	"sync"
	"bufio"
	// "os"
)

type server struct {
	members  map[net.Addr]*client
	commands chan command
}

func NewServer() {
	s := createServer()
	go s.run()

	//wg:= sync.WaitGroup{}

	//Start TCP server
	listener, err := net.Listen("tcp", ":8000")
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
		members:  make(map[net.Addr]*client),
		commands: make(chan command),
	}
}

func (s *server) createClient(conn net.Conn) {
	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}

	for {
		name := s.addName(c)
		if s.checkName(name) {
			c.name = name
			c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
			break
		} else {
			c.msg("This name already exists, please enter another name")
		}
	}
	for {
		if s.checkPassword(c) {
			log.Printf("%s has connected to the sever", c.name)
			s.members[c.conn.RemoteAddr()] = c
			s.readInput(c)
			break
		}
	}
}

func(s *server) addName(c *client) string{
	c.msg("What should I call you?")
	name, err := bufio.NewReader(c.conn).ReadString('\n')

	name = strings.Trim(name, "\r\n")
	args := strings.Split(name, " ")
	name = strings.TrimSpace(args[0])

	if err != nil {
		fmt.Println(err)
	}
	return name
}

func(s *server) checkName(name string) bool{
	for _, list := range s.members{
		if name == list.name {
			return false
		} 
	}
	return true
}

func (s *server) checkPassword(c *client) bool {
	c.msg("Enter Password: ")
	password, err := bufio.NewReader(c.conn).ReadString('\n')

	password = strings.Trim(password, "\r\n")
	args := strings.Split(password, " ")

	pass := strings.TrimSpace(args[0])

	if err != nil {
		fmt.Println(err)
	}

	if pass == "Password" {
		c.msg("Correct Password, Welcome to the Server!")
		return true
	} else {
		c.msg("Incorrect Password, Try again")
		return false
	}
}

func (s *server) readInput(c *client) {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")

		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/name":
			c.commands <- command{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/shout":
			c.commands <- command{
				id:     CMD_SHOUT,
				client: c,
				args:   args,
			}
		case "/spam":
			c.commands <- command{
				id:     CMD_SPAM,
				client: c,
				args:   args,
			}
		case "/whisper":
			c.commands <- command{
				id:     CMD_WHISPER,
				client: c,
				args:   args,
			}
		case "/list":
			c.commands <- command{
				id:     CMD_LIST,
				client: c,
				args:   args,
			}
		case "/help":
			c.commands <- command{
				id:     CMD_HELP,
				client: c,
				args:   args,
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
			s.quit(cmd.client)
		case CMD_SHOUT:
			s.shout(cmd.client, cmd.args)
		case CMD_SPAM:
			s.spam(cmd.client, cmd.args)
		case CMD_WHISPER:
			s.whisper(cmd.client, cmd.args)
		case CMD_LIST:
			s.list(cmd.client)
		case CMD_HELP:
			s.help(cmd.client)
		}
	}
}

func (s *server) name(c *client, args []string) {
	if s.checkName(args[1]) {
		c.name = args[1]
		c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
		log.Printf("%s has named themselves %s", c.conn.RemoteAddr().String(), c.name)
	} else {
		c.msg("This name already exists, please enter another name")
	}
}

func (s *server) msg(c *client, args []string) {
	msg := c.name + ": " + strings.Join(args[1:], " ")
	s.broadcast(c, msg)
}

func (s *server) quit(c *client) {
	message := c.name + " has left the chat"
	c.msg("You have left the server")
	c.conn.Close()
	s.broadcast(c, message)
	log.Print(message)
}

func (s *server) shout(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	shoutMsg := c.name + ": " + strings.ToUpper(msg)
	s.broadcast(c, shoutMsg)
}

func (s *server) spam(c *client, args []string) {
	msg := c.name + ": " + strings.Join(args[1:], " ")
	for i:= 0 ; i < 5; i++ { 
		s.broadcast(c, msg)
	}
}

func (s *server) broadcast(sender *client, msg string) {
	for _, m := range s.members {
		m.conn.Write([]byte("> " + msg + "\n"))
	}
}

func (s *server) whisper(sender *client, args []string) {
	msg := sender.name + ": " + strings.Join(args[2:], " ")
	for _, m := range s.members {
		if args[1] == m.name {
			m.conn.Write([]byte(">(whisper) " + msg + "\n"))
			sender.conn.Write([]byte(">(whisper) " + msg + "\n"))
		} 
	}
	if s.checkName(args[1]){
		sender.conn.Write([]byte("Could not find " + args[1]))
	}
}

func (s *server) list(sender *client) {
	for _, m := range s.members {
		sender.conn.Write([]byte("> " + m.name + "\n"))
	}
}

func (s *server) help(sender *client) {
	sender.conn.Write([]byte("\n/name\n/msg\n/shout\n/spam\n/whisper\n/list\n/quit\n/help\n"))
}