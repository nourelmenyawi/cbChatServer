package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type client struct {
	conn net.Conn
	name string
	commands chan<- command
}

func NewClient() {
	wg := sync.WaitGroup{}
	conn , err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Printf("unable to accept new connection: %s", err.Error())
	}
	wg.Add(2)
	go sendMsg(conn)
	go readMsg(conn, &wg)
	wg.Wait()
}

func sendMsg(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}

		_, e := conn.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			fmt.Println(e)
		}
	}
}

func readMsg(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for{
		messageBuffer := make([]byte, 1024)
		messageLength, messageError := conn.Read(messageBuffer)

		if messageError != nil {
			os.Exit(1)
		}

		fmt.Printf("%s\n", string(messageBuffer[:messageLength]))
	}
}


func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
