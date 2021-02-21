// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"net"
	"net/http"
	"strings"
)

func hostIsAllowed(remoteAddr string, allowed []string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return false
	}

	for _, addr := range allowed {
		if addr == host {
			return true
		}
	}

	return false
}

func allowHandler(allow []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hostIsAllowed(r.RemoteAddr, allow) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func discoverInterfaceAddresses(name string, prefix string) ([]string, error) {
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addresses, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, address := range addresses {
		if strings.HasPrefix(address.String(), prefix) && strings.HasSuffix(address.String(), "/128") {
			matches = append(matches, strings.TrimSuffix(address.String(), "/128"))
		}
	}

	return matches, nil
}
