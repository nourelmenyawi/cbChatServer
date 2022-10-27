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
	members  map[string]*client
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
	//returns pointer of type server containing the members map and commands channel
	return &server{
		members:  make(map[string]*client),
		commands: make(chan command),
	}
}

func (s *server) createClient(conn net.Conn) {
	//initialises c pointer for every client 
	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}

	//Requests clients to enter a name and saves it into the members map
	for {
		name := s.inputName(c)
		if s.checkNameIsUnique(name) {
			c.name = name
			c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
			break
		} else {
			c.msg("This name already exists, please enter another name")
		}
	}

	//Requests password for server
	for  {
		if s.checkPassword(c) {
			log.Printf("%s has connected to the sever", c.name)
			s.members[c.name] = c
			s.readInput(c)
			break
		}
	}
}

func(s *server) inputName(c *client) string{
	//Asks the user to enter a name and formats it
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

func(s *server) checkNameIsUnique(name string) bool{
	//check if the name exists in the map 
	return s.members[name] == nil
}

func (s *server) checkPassword(c *client) bool {
	//Asks user to enter password and checks it 
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
		c.msg("Incorrect Password, Please try again")
		return false
	}
}

func (s *server) readInput(c *client) {
	//Read input coming from each client's connection
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		//Splits the input from the client into a command and a message
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")

		cmd := strings.TrimSpace(args[0])

		//Based on the commands assign the type of process to the client commands channel
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
	//check the command sent by the client and call the function corresponding to it 
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

//Setting a name for the user
func (s *server) name(c *client, args []string) {
	if s.checkNameIsUnique(args[1]) {
		c.name = args[1]
		c.msg(fmt.Sprintf("All right, I will call you %s", c.name))
		log.Printf("%s has named themselves %s", c.conn.RemoteAddr().String(), c.name)
	} else {
		c.msg("This name already exists, please enter another name")
	}
}

//Sends a message to the entire server
func (s *server) msg(c *client, args []string) {
	msg := c.name + ": " + strings.Join(args[1:], " ")
	s.broadcast(c, msg)
}

//Quits the user from the server
func (s *server) quit(c *client) {
	message := c.name + " has left the chat"
	c.msg("You have left the server")
	c.conn.Close()
	s.broadcast(c, message)
	log.Print(message)
}

//Shouts a message on the server (Upper case all letters)
func (s *server) shout(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	shoutMsg := c.name + ": " + strings.ToUpper(msg)
	s.broadcast(c, shoutMsg)
}

//Sends the same message 5 times
func (s *server) spam(c *client, args []string) {
	msg := c.name + ": " + strings.Join(args[1:], " ")
	for i:= 0 ; i < 5; i++ { 
		s.broadcast(c, msg)
	}
}

//writes a message to every client connection
func (s *server) broadcast(sender *client, msg string) {
	for _, m := range s.members {
		m.conn.Write([]byte("> " + msg + "\n"))
	}
}

//Allows user to send a message to a specific user
func (s *server) whisper(sender *client, args []string) {
	msg := sender.name + ": " + strings.Join(args[2:], " ")

	//checks if the recepient set by the sender exists, if not reply back with could not find
	if s.members[args[1]] != nil {
		s.members[args[1]].conn.Write([]byte(">(whisper) " + msg + "\n"))
		sender.conn.Write([]byte(">(whisper) " + msg + "\n"))
	} else {
		sender.conn.Write([]byte("Could not find " + args[1]))
	}

}

//Lists all members connected to the server
func (s *server) list(sender *client) {
	for _, m := range s.members {
		sender.conn.Write([]byte("> " + m.name + "\n"))
	}
}

//Lists all commands
func (s *server) help(sender *client) {
	sender.conn.Write([]byte("\n/name\n/msg\n/shout\n/spam\n/whisper\n/list\n/quit\n/help\n"))
}