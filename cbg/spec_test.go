//  Copyright (c) 2021. Quirino Gervacio
//  MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

var (
	specSvc = newSpecSvc()
)

func loadTestFile(f string) []byte {
	spec, err := ioutil.ReadFile(fmt.Sprintf("../test/%s", f))
	if err != nil {
		spec, _ = ioutil.ReadFile(fmt.Sprintf("./test/%s", f))
	}
	return spec
}

func newSpecSvc() *SpecSvc {
	s, _ := NewSpec(
		loadTestFile("spec-test.yaml"),
		ArgsSpec{
			BiapiAk: string(loadTestFile("biapiAk-test.test")),
			BiapiSk: string(loadTestFile("biapiSk-test.test")),
			TaapiSk: string(loadTestFile("taapiSk.test")),
			EmUser:  string(loadTestFile("emUser.test")),
			EmPass:  string(loadTestFile("emPass.test")),
			NotEm:   string(loadTestFile("emUser.test")),
		})

	logFmt := new(log.TextFormatter)
	logFmt.TimestampFormat = time.RFC3339
	log.SetFormatter(logFmt)
	ll, _ := log.ParseLevel("debug")
	log.SetLevel(ll)

	return s
}

func Test_Spec_NewSpec_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, specSvc)
	assert.NotNil(t, specSvc.Spec)
}
