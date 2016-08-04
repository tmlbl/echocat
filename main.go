package main

import (
	"crypto/tls"
	"encoding/gob"
	"flag"
	"log"
	"time"

	"github.com/chzyer/readline"
)

var insecure bool
var tlsConn *tls.Conn

type Message struct {
	Time    time.Time
	User    string
	Message string
}

type Room struct {
	Name  string
	Lines []Message
}

type MessagePayload struct {
	Room    string
	Message Message
}

func sendMessage(room string, m Message) {
	log.Println("Sending message", m)
	err := gob.NewEncoder(tlsConn).Encode(m)
	if err != nil {
		log.Fatal("Error sending message")
	}
}

func startChat(room string) {
	rl, err := readline.NewEx(&readline.Config{
		UniqueEditLine: true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	rl.SetPrompt("username: ")
	username, err := rl.Readline()
	if err != nil {
		return
	}
	rl.ResetHistory()
	log.SetOutput(rl.Stderr())
	rl.SetPrompt(username + "> ")

	for {
		ln := rl.Line()
		sendMessage(room, Message{
			User:    username,
			Message: ln.Line,
			Time:    time.Now(),
		})
		if ln.CanContinue() {
			continue
		} else if ln.CanBreak() {
			break
		}
	}

	rl.Clean()
}

func startConnection(addr string) *tls.Conn {
	log.Printf("Connecting to %s...\n", addr)
	conn, err := tls.Dial("tcp", addr+":7337", &tls.Config{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		log.Fatalf("Could not reach the chat server at %s: %s\n", addr, err)
	}
	log.Println("Connection success")
	return conn
}

func main() {
	shouldServe := flag.Bool("serve", false, "Start a server")
	serverAddr := flag.String("host", "0.0.0.0", "Chat server to connect to")
	keyPath := flag.String("key", "server.key", "Path to key file")
	certPath := flag.String("crt", "server.crt", "Path to certificate file")
	room := flag.String("room", "lobby", "Name of room to join")
	flag.BoolVar(&insecure, "k", false, "OH NO")

	flag.Parse()

	if *shouldServe {
		startServer(*keyPath, *certPath)
	} else {
		tlsConn = startConnection(*serverAddr)
		startChat(*room)
	}
}
