// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	stratSvc = NewStratSvc(
		botSvc.SpecSvc.Spec.StratSpec,
		taapiSvc,
	)
)

func Test_Strat_NewStratSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, stratSvc)
}

func Test_Strat_Macd_When_Ok_Then_Pass(t *testing.T) {
	signal := stratSvc.Macd("BTC", "USDT", "15m", []int{12, 26, 9})
	assert.NotNil(t, signal)
}

func Test_Strat_Adx_When_Ok_Then_Pass(t *testing.T) {
	signal := stratSvc.Adx("BTC", "USDT", "15m", 14)
	assert.NotNil(t, signal)
}

func Test_Strat_SuperTrend_When_Ok_Then_Pass(t *testing.T) {
	signal := stratSvc.SuperTrend("BTC", "USDT", "15m", [][]int{{12, 3}, {11, 2}, {10, 1}})
	assert.NotNil(t, signal)
}
