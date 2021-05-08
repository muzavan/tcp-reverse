package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var typeVar string

func init() {
	flag.StringVar(&typeVar, "type", "client", "client/server")
	flag.Parse()
}

var (
	addr = "127.0.0.1:3333"
)

// Message format
// <message-length><actual-message>

func client() {
	log.Println("client running...")
	buff := bufio.NewScanner(os.Stdin)
	for {
		buff.Scan()
		cmd := buff.Text()
		if cmd == "EXIT" {
			break
		}

		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("net.Dial: %s", err.Error())
		}

		c.Write([]byte{byte(len(cmd))})
		c.Write([]byte(cmd))

		replyLenByte := make([]byte, 1)
		c.Read(replyLenByte)

		expLen := int(replyLenByte[0])
		reply := make([]byte, expLen)
		c.Read(reply)

		c.Close()
		fmt.Println(string(reply))

		if cmd == "SHUTDOWN" {
			break
		}
	}

	log.Println("client finished.")
}

func reverse(msg string) string {
	n := len(msg)
	var sb strings.Builder

	for i := n - 1; i >= 0; i-- {
		sb.WriteByte(msg[i])
	}

	return sb.String()
}

func server() {
	log.Println("server running...")
	c, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.ListenTCP: %s", err.Error())
	}
	defer c.Close()

	for {
		conn, err := c.Accept()
		if err != nil {
			log.Fatalf("c.Accept: %s", err.Error())
		}

		msgLenByte := make([]byte, 1)
		conn.Read(msgLenByte)

		msgLen := int(msgLenByte[0])
		actMsg := make([]byte, msgLen)
		conn.Read(actMsg)

		msg := string(actMsg)
		reply := []byte(fmt.Sprintf("%s received", reverse(msg)))
		replyLen := len(reply)
		conn.Write([]byte{byte(replyLen)})
		conn.Write(reply)
		conn.Close()

		if msg == "SHUTDOWN" {
			break
		}
	}

	log.Println("server finished.")
}

func main() {
	if typeVar == "client" {
		client()
	} else if typeVar == "server" {
		server()
	} else {
		log.Fatal("not a valid type")
	}
}
