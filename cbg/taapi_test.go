// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	taapiSvc = NewTaapiSvc(
		botSvc.SpecSvc.Spec.TaapiSpec,
		botSvc.SpecSvc.ArgsSpec.TaapiSk,
	)
)

func Test_Taapi_NewTaapiSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, taapiSvc)
}

func Test_Taapi_Macd_When_Ok_Then_Pass(t *testing.T) {
	time.Sleep(1 * time.Second) // don't choke taapi subscription
	out, err := taapiSvc.Macd("BTC", "USDT", "15m", 12, 36, 9)
	assert.Nil(t, err, nil)
	assert.NotNil(t, out, nil)
}

func Test_Taapi_Adx_When_Ok_Then_Pass(t *testing.T) {
	time.Sleep(1 * time.Second) // don't choke taapi subscription
	out, err := taapiSvc.Adx("BTC", "USDT", "15m", 14)
	assert.Nil(t, err, nil)
	assert.NotNil(t, out, nil)
}

func Test_Taapi_SuperTrend_When_Ok_Then_Pass(t *testing.T) {
	time.Sleep(1 * time.Second) // don't choke taapi subscription
	out, err := taapiSvc.SuperTrend("BTC", "USDT", "15m", 10, 3)
	assert.Nil(t, err, nil)
	assert.NotNil(t, out, nil)
}
