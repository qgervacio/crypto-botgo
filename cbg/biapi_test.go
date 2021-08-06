// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	biapiSvc = NewBiapiSvc(
		botSvc.SpecSvc.Spec.BiapiSpec,
		botSvc.SpecSvc.Spec.CredSpec.BiapiAk,
		botSvc.SpecSvc.Spec.CredSpec.BiapiSk,
	)
)

func Test_Biapi_NewBiapiSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, biapiSvc)
}

func Test_Biapi_GetBalance_When_Ok_Then_Pass(t *testing.T) {
	bal, err := biapiSvc.GetBalance("BTC")
	assert.Nil(t, err)
	assert.NotEqual(t, "", bal)
}

func Test_Biapi_BuySellMarket_When_Ok_Then_Pass(t *testing.T) {
	r, err := biapiSvc.BuyMarketMT("BTC", "USDT", "1")
	assert.Nil(t, err)
	assert.NotNil(t, "", r)

	r, err = biapiSvc.SellCoinMT("BTC", "USDT", "1")
	assert.Nil(t, err)
	assert.NotNil(t, "", r)
}
