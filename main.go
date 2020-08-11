// Copyright 2020 Noah Kantrowitz
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"golang.org/x/net/http2"
)

var config struct {
	listenPort     int
	listenTLS      bool
	certPath       string
	keyPath        string
	redirectScheme string
	redirectHost   string
	redirectPath   string
	redirectStatus int
}

func init() {
	flag.IntVar(&config.listenPort, "port", 0, "Port to listen on")
	flag.BoolVar(&config.listenTLS, "tls", false, "Listen using TLS")
	flag.StringVar(&config.certPath, "cert", "", "Path to tls.crt")
	flag.StringVar(&config.keyPath, "key", "", "Path to tls.key")
	flag.StringVar(&config.redirectScheme, "scheme", "https", "URI scheme to redirect to")
	flag.StringVar(&config.redirectHost, "host", "", "URI host to redirect to")
	flag.StringVar(&config.redirectPath, "path", "", "URI path to redirect to")
	flag.IntVar(&config.redirectStatus, "status", 302, "HTTP status to redirect with")
}

func redirect(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if config.redirectPath != "" {
		path = config.redirectPath
	}
	target := config.redirectScheme + "://" + config.redirectHost + path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, config.redirectStatus)
}

func main() {
	// Handle config options.
	flag.Parse()
	if config.redirectHost == "" {
		log.Fatal("--host argument is required")
	}
	if config.listenTLS && (config.certPath == "" || config.keyPath == "") {
		log.Fatal("--cert and --key are required for TLS")
	}
	if config.listenPort == 0 {
		if config.listenTLS {
			config.listenPort = 8443
		} else {
			config.listenPort = 8080
		}
	}

	// Better TLS, just to be safe.
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// Create the HTTP server.
	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", config.listenPort),
		Handler:   handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(redirect)),
		TLSConfig: cfg,
	}
	http2.ConfigureServer(srv, nil)

	// Listen forever.
	log.Printf("Starting bouncycastle on %s", srv.Addr)
	if config.listenTLS {
		log.Fatal(srv.ListenAndServeTLS(config.certPath, config.keyPath))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}
