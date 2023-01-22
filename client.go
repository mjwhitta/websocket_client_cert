package main

import (
	"bufio"
	"crypto/tls"
	"os"
	"strings"

	ws "github.com/gorilla/websocket"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/log"
)

var addr = "localhost:8443"

func generateTLSConfig() *tls.Config {
	var e error
	var pem tls.Certificate
	var pemFile = "pki/pems/user.client.pem"

	// Read in x509 PEM certificate (cert + key)
	if pem, e = tls.LoadX509KeyPair(pemFile, pemFile); e != nil {
		log.Err("Failed read client cert: " + e.Error())
	}

	// Use client cert for connections but skip verification b/c
	// self-signed CA
	return &tls.Config{
		Certificates:       []tls.Certificate{pem},
		InsecureSkipVerify: true,
	}
}

func main() {
	var ack []byte
	var c *ws.Conn
	var d *ws.Dialer = ws.DefaultDialer
	var e error
	var in string
	var r *bufio.Reader = bufio.NewReader(os.Stdin)

	// Configure Dialer
	d.TLSClientConfig = generateTLSConfig()

	// Connect to server
	log.Info("Connecting to " + addr)
	if c, _, e = d.Dial("wss://"+addr, nil); e != nil {
		log.Err("Dial: " + e.Error())
		os.Exit(1)
	}
	defer c.Close()

	for {
		// Read user input
		hl.Print("Enter text: ")
		if in, e = r.ReadString('\n'); e != nil {
			log.Err("Input: " + e.Error())
			continue
		}

		// Remove trailing newline
		in = strings.TrimSpace(in)

		// Exit if requested
		switch in {
		case "bye", "exit", "q", "quit":
			// Send close to server
			e = c.WriteMessage(
				ws.CloseMessage,
				ws.FormatCloseMessage(ws.CloseNormalClosure, ""),
			)
			if e != nil {
				log.Err("Close: " + e.Error())
			}
			os.Exit(0)
		}

		// Send input to server
		e = c.WriteMessage(ws.TextMessage, []byte(in))
		if e != nil {
			log.Err("Write: " + e.Error())
			continue
		}

		// Receive ACK
		if _, ack, e = c.ReadMessage(); e != nil {
			log.Err("Read " + e.Error())
			continue
		}

		log.Good("Recv: " + string(ack))
	}
}
