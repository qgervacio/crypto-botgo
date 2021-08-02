// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type TaapiSvc struct {
	TaapiSpec *TaapiSpec
	SecretKey string
}

func NewTaapiSvc(ts *TaapiSpec, sk string) *TaapiSvc {
	return &TaapiSvc{
		TaapiSpec: ts,
		SecretKey: sk,
	}
}

func (s *TaapiSvc) Macd(
	coin, market, interval string,
	fast int, slow int, signal int) (string, error) {
	p := "optInFastPeriod=%d&optInSlowPeriod=%d&optInSignalPeriod=%d"
	return s.call(
		"macd",
		coin, market, interval,
		fmt.Sprintf(p, fast, slow, signal))
}

func (s *TaapiSvc) Adx(
	coin, market, interval string,
	period int) (string, error) {
	return s.call(
		"adx",
		coin, market, interval,
		fmt.Sprintf("optInTimePeriod=%d", period))
}

func (s *TaapiSvc) SuperTrend(
	coin, market, interval string,
	period int, multiplier int) (string, error) {
	return s.call(
		"supertrend",
		coin, market, interval,
		fmt.Sprintf("period=%d&multiplier=%d", period, multiplier))
}

func (s *TaapiSvc) call(
	api, coin, market, interval,
	extra string) (string, error) {
	p := fmt.Sprintf("%s/%s", s.TaapiSpec.Url, api)
	p = fmt.Sprintf("%s?secret=%s", p, s.SecretKey)
	p = fmt.Sprintf("%s&exchange=%s", p, s.TaapiSpec.Exchange)
	p = fmt.Sprintf("%s&symbol=%s/%s", p, coin, market)
	p = fmt.Sprintf("%s&interval=%s", p, interval)
	p = fmt.Sprintf("%s&chart=%s", p, s.TaapiSpec.Chart)
	p = fmt.Sprintf("%s&%s", p, extra)
	res, err := http.Get(p)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("%s request failed [%v]", api, err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%s response body read failed [%v]", api, err)
	}

	return string(body), nil
}
