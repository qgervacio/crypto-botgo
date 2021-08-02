// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package main

import (
	"flag"
	"github.com/qgervacio/crypto-botgo/cbg"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

var (
	specFile = flag.String("specFile", "", "Spec YAML file")
	logLevel = flag.String("logLevel", "info", "Log level")
)

func main() {
	flag.Parse()

	logFmt := new(log.JSONFormatter)
	logFmt.TimestampFormat = time.RFC3339
	log.SetReportCaller(true)
	log.SetFormatter(logFmt)
	ll, _ := log.ParseLevel(*logLevel)
	log.SetLevel(ll)

	yb, err := ioutil.ReadFile(*specFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	ss, err := cbg.NewSpec(yb)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	bs := cbg.NewBotSvc(ss)
	bs.Run()

	forever := make(chan bool)
	<-forever
}
