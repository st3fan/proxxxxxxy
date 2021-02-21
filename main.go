// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

func waitForControlC() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func init() {
	log.SetOutput(os.Stdout) // For systemd?
}

func main() {
	configFileFlag := flag.String("config", "/etc/proxxxxxxy/proxxxxxxy.json", "path to the config file")
	flag.Parse()

	config, err := newConfigFromFile(*configFileFlag)
	if err != nil {
		log.Fatalf("Could not load config <%s>: %v", *configFileFlag, err)
	}

	log.Printf("%+v", config)

	application, err := newApplication(config)
	if err != nil {
		log.Fatalf("Could not create application: %v", err)
	}

	if err := application.run(); err != nil {
		log.Fatalf("Could not start application: %v", err)
	}
}
