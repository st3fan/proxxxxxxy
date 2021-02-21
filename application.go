// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/st3fan/service"
)

type Application struct {
	config Config
}

func newApplication(config Config) (*Application, error) {
	return &Application{
		config: config,
	}, nil
}

func (app *Application) run() error {
	interfaceAddresses, err := discoverInterfaceAddresses(app.config.Interface, app.config.Prefix)
	if err != nil {
		return errors.Wrapf(err, "Could not discover addresses on <%s>", app.config.Interface)
	}

	var runners []service.ServiceRunner

	for i, interfaceAddress := range interfaceAddresses {
		proxy, err := newProxy(fmt.Sprintf("%s:%d", app.config.Address, app.config.ProxyPortStart+i), interfaceAddress, app.config.Allow, app.config.Verbose)
		if err != nil {
			return errors.Wrapf(err, "Could not create Proxy on port <%d>", app.config.ProxyPortStart+i)
		}
		runners = append(runners, proxy)
	}

	api, err := newAPI(app.config, interfaceAddresses)
	if err != nil {
		return errors.Wrapf(err, "Could not create API on port <%d>", app.config.APIPort)
	}
	runners = append(runners, api)

	service.Run(context.Background(), runners...)

	return nil
}
