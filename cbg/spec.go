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
	DryRun    bool       `yaml:"dry_run"`
	EmailSpec *EmailSpec `validate:"required"  yaml:"email"`
	BiapiSpec *BiapiSpec `validate:"required"  yaml:"biapi"`
	TaapiSpec *TaapiSpec `validate:"required"  yaml:"taapi"`
	StratSpec *StratSpec `validate:"required"  yaml:"strategy"`
	SpotSpec  *SpotSpec  `validate:"required"  yaml:"spot"`
	EndSpec   *EndSpec   `validate:"required"  yaml:"end"`
	CredSpec  *CredSpec  `validate:"required"  yaml:"cred"`
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
	Coin            string `validate:"required" yaml:"coin"`
	Market          string `validate:"required" yaml:"market"`
	PeriodMin       int    `validate:"required" yaml:"period_min"`
	CoinPrecision   string `validate:"required" yaml:"coin_precision"`
	MarketPrecision string `validate:"required" yaml:"market_precision"`
}

type EndSpec struct {
	DummySymbol     string    `validate:"required" yaml:"dummy_symbol"`
	DummyOrder      []float64 `validate:"required" yaml:"dummy_order"`
}

type CredSpec struct {
	BiapiAk string `validate:"required" yaml:"biapi_ak"`
	BiapiSk string `validate:"required" yaml:"biapi_sk"`
	TaapiSk string `validate:"required" yaml:"taapi_sk"`
	EmUser  string `validate:"required" yaml:"em_user"`
	EmPass  string `validate:"required" yaml:"em_pass"`
	NotEm   string `validate:"required" yaml:"noti_em"`
}

type SpecSvc struct {
	Spec Spec
}

func NewSpec(s []byte) (*SpecSvc, error) {
	spec := Spec{}

	if err := yaml.Unmarshal(s, &spec); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(&spec); err != nil {
		return nil, err
	}

	return &SpecSvc{
		Spec: spec,
	}, nil
}
