//  Copyright (c) 2021. Quirino Gervacio
//  MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
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

	ns := make(chan []string)
	s.notifier(ns)
	s.ender(ns)

	ts := make(chan Signal)
	s.trader(ns, ts)
	s.looker(ts)
}

func (s *BotSvc) trader(ns chan<- []string, ts <-chan Signal) {
	go func() {
		sig := <-ts
		for {
			symbol := fmt.Sprintf("%s%s", s.SpotSvc.SpotSpec.Coin, s.SpotSvc.SpotSpec.Market)
			if sig == SignalBuy {
				res, err := s.BiapiSvc.BuyMarket(
					s.SpotSvc.SpotSpec.Coin,
					s.SpotSvc.SpotSpec.Market, "")
				if err != nil {
					ns <- []string{fmt.Sprintf("Failed to buy %s", symbol), err.Error()}
				} else {
					ns <- []string{fmt.Sprintf("Bought %s (%s)", symbol, res.Status), ""}
				}
			} else {
				res, err := s.BiapiSvc.SellMarket(
					s.SpotSvc.SpotSpec.Coin,
					s.SpotSvc.SpotSpec.Market, "")
				if err != nil {
					ns <- []string{fmt.Sprintf("Failed to sell %s", symbol), err.Error()}
				} else {
					ns <- []string{fmt.Sprintf("Sold %s (%s)", symbol, res.Status), ""}
				}
			}
		}
	}()
}

func (s *BotSvc) looker(ts chan<- Signal) {
	go func() {
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

			position = SignalBuy
			if position == SignalBuy || position == SignalSell {
				ts <- position
			} else {
				log.Infof("No action for now...")
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
