//  Copyright (c) 2021. Quirino Gervacio
//  MIT License. All Rights Reserved

package cbg

import (
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
)

type Spec struct {
	Name      string     `validate:"required" yaml:"name"`
	Delayed   bool       `yaml:"delayed"`
	EmailSpec *EmailSpec `validate:"required"  yaml:"email"`
	BiapiSpec *BiapiSpec `validate:"required"  yaml:"biapi"`
	TaapiSpec *TaapiSpec `validate:"required"  yaml:"taapi"`
	StratSpec *StratSpec `validate:"required"  yaml:"strategy"`
	SpotSpec  *SpotSpec  `validate:"required"  yaml:"spot"`
	EndSpec   *EndSpec   `validate:"required"  yaml:"end"`
}

type EmailSpec struct {
	Server             string `validate:"required" yaml:"server"`
	Port               int    `validate:"required" yaml:"port"`
	Name               string `validate:"required" yaml:"name"`
	ConnTimeoutSec     int    `validate:"required" yaml:"conn_timeout_sec"`
	SendTimeoutSec     int    `validate:"required" yaml:"send_timeout_sec"`
	InsecureSkipVerify bool   `validate:"required" yaml:"insecure_skip_verify"`
}

type BiapiSpec struct {
	Test       bool   `json:"test"`
	RecvWindow int64  `validate:"required" yaml:"receive_window"`
	Url        string `validate:"required" yaml:"url"`
}

type TaapiSpec struct {
	Url       string `validate:"required" yaml:"url"`
	Exchange  string `validate:"required" yaml:"exchange"`
	Chart     string `validate:"required" yaml:"chart"`
	Backtrack int    `validate:"required" yaml:"backtrack"`
}

type StratSpec struct {
	Adx        int     `validate:"required" yaml:"adx"`
	Macd       []int   `validate:"required" yaml:"macd"`
	SuperTrend [][]int `validate:"required" yaml:"super_trend"`
}

type SpotSpec struct {
	Coin              string  `validate:"required" yaml:"coin"`
	Market            string  `validate:"required" yaml:"market"`
	PeriodMin         int     `validate:"required" yaml:"period_min"`
	CoinPrecision     string  `validate:"required" yaml:"coin_precision"`
	MarketPrecision   string  `validate:"required" yaml:"market_precision"`
	InitialMarketFund float64 `validate:"required" yaml:"initial_market_fund"`
	OrderTimeoutSec   int     `validate:"required" yaml:"order_timeout_sec"`
	Slippage          float64 `validate:"required" yaml:"slippage"`
	BeforeSlippage    int     `validate:"required" yaml:"before_slippage"`
}

type EndSpec struct {
	ConsecutiveLoss int       `validate:"required" yaml:"consecutive_loss"`
	DummySymbol     string    `validate:"required" yaml:"dummy_symbol"`
	DummyOrder      []float64 `validate:"required" yaml:"dummy_order"`
}

type ArgsSpec struct {
	BiapiAk string `validate:"required"`
	BiapiSk string `validate:"required"`
	TaapiSk string `validate:"required"`
	EmUser  string `validate:"required"`
	EmPass  string `validate:"required"`
	NotEm   string `validate:"required"`
}

type SpecSvc struct {
	Spec     Spec
	ArgsSpec ArgsSpec
}

func NewSpec(s []byte, a ArgsSpec) (*SpecSvc, error) {
	spec := Spec{}

	if err := yaml.Unmarshal(s, &spec); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(&spec); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(&a); err != nil {
		return nil, err
	}

	return &SpecSvc{
		Spec:     spec,
		ArgsSpec: a,
	}, nil
}
