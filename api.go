// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type API struct {
	config    Config
	addresses []string
}

func newAPI(config Config, addresses []string) (*API, error) {
	return &API{
		config:    config,
		addresses: addresses,
	}, nil
}

func (api *API) listenersHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/listeners" {
		type listenersResponse struct {
			Listeners map[string]string
		}

		response := listenersResponse{
			Listeners: map[string]string{},
		}
		for i, address := range api.addresses {
			response.Listeners[fmt.Sprintf("http://%s:%d", api.config.Hostname, api.config.ProxyPortStart+i)] = address
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	http.NotFound(w, r)
}

func (api *API) Run(ctx context.Context) {
	log.Printf("API listening on <%s:%d>", api.config.Address, api.config.APIPort)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", api.config.Address, api.config.APIPort),
		Handler: allowHandler(api.config.Allow, http.HandlerFunc(api.listenersHandler)),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		<-ctx.Done()

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
