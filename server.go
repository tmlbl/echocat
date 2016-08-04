package main

import (
	"crypto/tls"
	"encoding/gob"
	"log"
	"net"
)

var chatContext = map[string]Room{}

func emitMessage(room string, m Message) {

}

func startServer(keyPath, certPath string) {
	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":7337", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	log.Println("Listening on port 7337...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		var m Message
		err := gob.NewDecoder(conn).Decode(&m)
		if err != nil {
			break
		}
		log.Println(m)
		emitMessage("lobby", m)
	}
}
