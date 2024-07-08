# Websockets using TLS client certificates

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

This repo contains sample code for a server and a client that use and
enforce TLS client certificates. The server uses a certificate chain
in pki/pems/localhost.chain.pem. This chain also contains the
top-level CA which signed the client certificates. The client then
uses the pki/pems/user.client.pem which contains both the certificate
and key for that user. The interesting TLS related bits are summarized
in the below sections.

## Server

```
import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

...

var e error
var pem []byte
var pemFile = "pki/pems/localhost.chain.pem"
var pool *x509.CertPool = x509.NewCertPool()
var tcfg *tls.Config

// Read in server + CA chain
if pem, e = ioutil.ReadFile(pemFile); e != nil {
    ...
}

// Add certs to pool
if !pool.AppendCertsFromPEM(pem) {
    ...
}

// Require client certificates for connections
tcfg = &tls.Config{
    ClientAuth: tls.RequireAndVerifyClientCert,
    ClientCAs:  pool,
}

...
```

## Client

```
import "crypto/tls"

...

var e error
var pem tls.Certificate
var pemFile = "pki/pems/user.client.pem"
var tcfg *tls.Config

// Read in x509 PEM certificate (cert + key)
if pem, e = tls.LoadX509KeyPair(pemFile, pemFile); e != nil {
    ...
}

// Use client cert for connections but skip verification b/c
// self-signed CA
tcfg = &tls.Config{
    Certificates:       []tls.Certificate{pem},
    InsecureSkipVerify: true,
}

...
```
