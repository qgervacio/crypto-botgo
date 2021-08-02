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
	biapiAk  = flag.String("biapiAk", "", "Binance API key")
	biapiSk  = flag.String("biapiSk", "", "Binance secret key")
	taapiSk  = flag.String("taapiSk", "", "TAAPI secret key")
	emUser   = flag.String("emUser", "", "Email username")
	emPass   = flag.String("emPass", "", "Email password")
	notEm    = flag.String("notEm", "", "Emails to receive notifications")
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

	ss, err := cbg.NewSpec(yb, cbg.ArgsSpec{
		BiapiAk: *biapiAk,
		BiapiSk: *biapiSk,
		TaapiSk: *taapiSk,
		EmUser:  *emUser,
		EmPass:  *emPass,
		NotEm:   *notEm,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	bs := cbg.NewBotSvc(ss)
	bs.Run()

	forever := make(chan bool)
	<-forever
}
