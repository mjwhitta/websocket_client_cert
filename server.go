package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	ws "github.com/gorilla/websocket"
	"gitlab.com/mjwhitta/log"
)

var addr = "localhost:8443"
var wsup = ws.Upgrader{}

func generateTLSConfig() *tls.Config {
	var e error
	var pem []byte
	var pemFile = "pki/pems/localhost.chain.pem"
	var pool *x509.CertPool = x509.NewCertPool()

	// Read in server + CA chain
	if pem, e = ioutil.ReadFile(pemFile); e != nil {
		log.Err("Failed read CA: " + e.Error())
	}

	// Add certs to pool
	if !pool.AppendCertsFromPEM(pem) {
		log.Err("Failed to add CA")
	}

	// Require client certificates for connections
	return &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}
}

func main() {
	var e error
	var msg []byte
	var mt int
	var s = &http.Server{
		Addr:      addr,
		TLSConfig: generateTLSConfig(),
	}

	// Define / handler
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			var ack string
			var c *ws.Conn
			var e error

			// Upgrade client to websocket
			if c, e = wsup.Upgrade(w, r, nil); e != nil {
				log.Err("Upgrade: " + e.Error())
				return
			}
			defer c.Close()

			for {
				// Read message from client
				if mt, msg, e = c.ReadMessage(); e != nil {
					log.Err("Read: " + e.Error())
					break
				}

				log.Good("Recv: " + string(msg))

				// Send ACK
				ack = "ACK - " + string(msg)
				if e = c.WriteMessage(mt, []byte(ack)); e != nil {
					log.Err("Write: " + e.Error())
					break
				}
			}
		},
	)

	// Start server
	log.Info("Listening on " + addr)
	e = s.ListenAndServeTLS(
		"pki/pems/localhost.chain.pem",
		"pki/pems/localhost.key.pem",
	)
	if e != nil {
		panic(e)
	}
}
