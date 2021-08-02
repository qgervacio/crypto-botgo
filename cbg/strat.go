// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Signal string

const (
	SignalNone      Signal = "NONE"
	SignalError     Signal = "ERROR "
	SignalBuy       Signal = "BUY"
	SignalSell      Signal = "SELL"
	SignalHodl      Signal = "HODL"
	SignalLong      Signal = "LONG"
	SignalShort     Signal = "SHORT"
	SignalHasTrend  Signal = "HAS_TREND"
	SignalNoTrend   Signal = "NO_TREND"
	SignalUptrend   Signal = "UPTREND"
	SignalDowntrend Signal = "DOWNTREND"
)

type StratSvc struct {
	StratSpec *StratSpec
	TaapiSvc  *TaapiSvc
}

func NewStratSvc(ss *StratSpec, ts *TaapiSvc) *StratSvc {
	return &StratSvc{
		StratSpec: ss,
		TaapiSvc:  ts,
	}
}

func (s *StratSvc) Adx(coin, market, interval string, period int) Signal {
	tname := fmt.Sprintf("ADX(%d)", period)
	trend, err := s.TaapiSvc.Adx(coin, market, interval, period)
	if err != nil {
		log.Errorf("%s failed [%v]", tname, err)
		return SignalError
	}
	log.Infof("%s -> %s", tname, trend)
	value := int(gjson.Get(trend, "value").Int())
	if value >= period {
		return SignalHasTrend
	}
	return SignalNoTrend
}

func (s *StratSvc) Macd(coin, market, interval string, conf []int) Signal {
	tname := fmt.Sprintf("MACD(%d, %d, %d)", conf[0], conf[1], conf[2])
	trend, err := s.TaapiSvc.Macd(coin, market, interval, conf[0], conf[1], conf[2])
	if err != nil {
		log.Errorf("%s failed [%v]", tname, err)
		return SignalError
	}
	log.Infof("%s -> %s", tname, trend)
	valueMACD := gjson.Get(trend, "valueMACD").Float()
	valueMACDSignal := gjson.Get(trend, "valueMACDSignal").Float()
	if valueMACD >= 0.0 && valueMACDSignal >= 0.0 {
		return SignalUptrend
	}
	return SignalDowntrend
}

func (s *StratSvc) SuperTrend(coin, market, interval string, confs [][]int) Signal {
	clen := len(confs)
	signal := make(chan Signal, clen)
	for _, c := range confs {
		go func(period int, multiplier int) {
			tname := fmt.Sprintf("SuperTrend(%d, %d)", period, multiplier)
			trend, err := s.TaapiSvc.SuperTrend(coin, market, interval, period, multiplier)
			if err != nil {
				log.Errorf("%s failed [%v]", tname, err)
				signal <- SignalError
				return
			}
			log.Infof("%s -> %s", tname, trend)
			if gjson.Get(trend, "valueAdvice").String() == "long" {
				signal <- SignalLong
			} else {
				signal <- SignalShort
			}
		}(c[0], c[1])
	}

	var longCount, shortCount, hodlCount int
	for range confs {
		switch sig := <-signal; sig {
		case SignalLong:
			longCount++
		case SignalShort:
			shortCount++
		default:
			hodlCount++
		}
	}
	close(signal)

	if longCount == clen {
		return SignalLong
	}
	if shortCount == clen {
		return SignalShort
	}
	return SignalHodl
}
