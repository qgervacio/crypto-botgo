//  Copyright (c) 2021. Quirino Gervacio
//  MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
	"time"
)

type BotSvc struct {
	SpecSvc  *SpecSvc
	TaapiSvc *TaapiSvc
	EmailSvc *EmailSvc
	BiapiSvc *BiapiSvc
	StratSvc *StratSvc
	SpotSvc  *SpotSvc
	EndSvc   *EndSvc
}

func NewBotSvc(ss *SpecSvc) *BotSvc {
	taapiSvc := NewTaapiSvc(ss.Spec.TaapiSpec, ss.Spec.CredSpec.TaapiSk)
	emailSvc := NewEmailSvc(ss.Spec.EmailSpec, ss.Spec.CredSpec.EmUser, ss.Spec.CredSpec.EmPass)
	biapiSvc := NewBiapiSvc(ss.Spec.BiapiSpec, ss.Spec.CredSpec.BiapiAk, ss.Spec.CredSpec.BiapiSk)
	stratSvc := NewStratSvc(ss.Spec.StratSpec, taapiSvc)
	spotSvc := NewSpotSvc(ss.Spec.SpotSpec, biapiSvc)
	endSvc := NewEndSvc(ss.Spec.EndSpec, biapiSvc)
	return &BotSvc{
		SpecSvc:  ss,
		TaapiSvc: taapiSvc,
		EmailSvc: emailSvc,
		BiapiSvc: biapiSvc,
		StratSvc: stratSvc,
		SpotSvc:  spotSvc,
		EndSvc:   endSvc,
	}
}

func (s *BotSvc) Run() {
	if s.SpecSvc.Spec.Delayed {
		t := SleepUntil(s.SpecSvc.Spec.SpotSpec.PeriodMin, time.Minute)
		log.Infof("Sleeping until %s", t)
		time.Sleep(time.Until(t))
	}

	log.Infof("Bot %s started", s.SpecSvc.Spec.Name)
	log.Infof("Trading on %s%s at %dm",
		s.SpecSvc.Spec.SpotSpec.Coin,
		s.SpecSvc.Spec.SpotSpec.Market,
		s.SpecSvc.Spec.SpotSpec.PeriodMin)
	log.Infof("Backtrack at %d", s.SpecSvc.Spec.TaapiSpec.Backtrack)

	sig := make(chan []string)
	s.notifier(sig)
	s.ender(sig)
	s.trader(sig)
}

func (s *BotSvc) trader(ns chan<- []string) {
	go func() {
		sigCount := 0
		for {
			var wg sync.WaitGroup
			interval := fmt.Sprintf("%dm",
				s.SpecSvc.Spec.SpotSpec.PeriodMin) // TODO always minutes?
			t := SleepUntil(s.SpecSvc.Spec.SpotSpec.PeriodMin, time.Minute)

			// Determine if there is a trend
			wg.Add(1)
			adxSig := SignalNone
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				adxSig = s.StratSvc.Adx(
					s.SpecSvc.Spec.SpotSpec.Coin,
					s.SpecSvc.Spec.SpotSpec.Market,
					interval,
					s.SpecSvc.Spec.StratSpec.Adx)
			}(&wg)

			// Determine the direction of the trend
			wg.Add(1)
			macdSig := SignalNone
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				macdSig = s.StratSvc.Macd(
					s.SpecSvc.Spec.SpotSpec.Coin,
					s.SpecSvc.Spec.SpotSpec.Market,
					interval,
					s.SpecSvc.Spec.StratSpec.Macd)
			}(&wg)

			// Determine entry/exit point
			wg.Add(1)
			superTrendSig := SignalNone
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				superTrendSig = s.StratSvc.SuperTrend(
					s.SpecSvc.Spec.SpotSpec.Coin,
					s.SpecSvc.Spec.SpotSpec.Market,
					interval,
					s.SpecSvc.Spec.StratSpec.SuperTrend)
			}(&wg)

			wg.Wait()
			log.Infof("ADX Signal: %s", adxSig)
			log.Infof("MACD Signal: %s", macdSig)
			log.Infof("SuperTrend Signal: %s", superTrendSig)

			position := SignalError
			if adxSig == SignalError ||
				macdSig == SignalError ||
				superTrendSig == SignalError {
				log.Error("Some signal failed")
				time.Sleep(time.Until(t))
				continue
			}
			if adxSig == SignalNone ||
				macdSig == SignalNone ||
				superTrendSig == SignalNone {
				log.Warn("Incomplete strategy")
				time.Sleep(time.Until(t))
				continue
			}

			if adxSig == SignalHasTrend {
				if macdSig == SignalUptrend {
					if superTrendSig == SignalLong {
						position = SignalBuy
					} else if superTrendSig == SignalShort {
						position = SignalSell
					} else {
						position = SignalHodl
					}
				} else {
					position = SignalHodl
				}
			} else {
				position = SignalHodl
			}

			sigCount++
			log.Infof("Signal#%d - %s", sigCount, position)

			price := 0.0
			symbol := fmt.Sprintf("%s%s",
				s.SpecSvc.Spec.SpotSpec.Coin, s.SpecSvc.Spec.SpotSpec.Market)
			priceRaw, err := s.SpotSvc.BiapiSvc.ListPrices(symbol)
			if err != nil {
				log.Errorf("Failed to get price [%v]", err)
			} else {
				price, _ = strconv.ParseFloat(priceRaw[0].Price, 64)
			}

			sub := fmt.Sprintf("Signal#%d - %s %s %s",
				sigCount, position, symbol,
				fmt.Sprintf(s.SpecSvc.Spec.SpotSpec.MarketPrecision, price))
			if position == SignalBuy {
				msg, err := s.SpotSvc.Buy()
				if err != nil {
					msg = err.Error()
				}
				ns <- []string{sub, msg}
			} else if position == SignalSell {
				msg, err := s.SpotSvc.Sell()
				if err != nil {
					msg = err.Error()
				}
				ns <- []string{sub, msg}
			} else {
				ns <- []string{sub, "None for now"}
			}

			log.Infof("Next position inquiry is %sm", t)
			time.Sleep(time.Until(t))
		}
	}()
}

func (s *BotSvc) ender(ns chan<- []string) {
	go func() {
		for {
			log.Debugf("Looking for dummy order...")
			shouldCancel, err := s.EndSvc.EndByDummyOrder()
			if err != nil {
				ns <- []string{"Ender Triggered with Error", err.Error()}
			}
			if shouldCancel {
				ns <- []string{"Ender Triggered", "None"}
				log.Infof("Ender triggered. Sleeping for 5s before terminating bot")
				time.Sleep(5 * time.Second) // TODO externalize?
				os.Exit(0)
			}
			is := int(s.SpecSvc.Spec.EndSpec.DummyOrder[2])
			log.Debugf("Ender sleeping for %ds", is)
			time.Sleep(time.Duration(is) * time.Second)
		}
	}()
}

func (s *BotSvc) notifier(ns <-chan []string) {
	go func() {
		for {
			msg := <-ns
			log.Debugf("Notifying...")
			if err := s.EmailSvc.Send(
				msg[0], msg[1],
				[]string{s.SpecSvc.Spec.CredSpec.NotEm},
				[]string{}, []string{}); err != nil {
				log.Errorf("failed to send email notification [%v]", err)
				continue
			}
			log.Debugf("Done notifying...")
		}
	}()
}
