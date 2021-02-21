// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

type Proxy struct {
	proxyAddress  string
	clientAddress string
	verbose       bool
	allow         []string
	proxy         *goproxy.ProxyHttpServer
}

func (p *Proxy) Run(ctx context.Context) {
	log.Printf("Proxy listening on <%s> with request address <%s>", p.proxyAddress, p.clientAddress)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = p.verbose
	proxy.Tr.DisableKeepAlives = true
	proxy.Tr.Dial = func(network, address string) (net.Conn, error) {
		dialer := net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(p.clientAddress),
			},
		}
		return dialer.Dial(network, address)
	}

	server := &http.Server{
		Addr:    p.proxyAddress,
		Handler: allowHandler(p.allow, proxy),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Printf("Shutting down proxy <%s> with request address <%s>", p.proxyAddress, p.clientAddress)
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	<-idleConnsClosed
}

func newProxy(proxyAddress string, clientAddress string, allow []string, verbose bool) (*Proxy, error) {
	return &Proxy{
		proxyAddress:  proxyAddress,
		clientAddress: clientAddress,
		allow:         allow,
		verbose:       verbose,
	}, nil
}
