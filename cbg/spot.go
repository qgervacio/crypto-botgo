// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import "github.com/adshao/go-binance/v2"

type SpotSvc struct {
	SpotSpec *SpotSpec
	BiapiSvc *BiapiSvc
}

func NewSpotSvc(ss *SpotSpec, bs *BiapiSvc) *SpotSvc {
	return &SpotSvc{
		SpotSpec: ss,
		BiapiSvc: bs,
	}
}

func (s *SpotSvc) BuyMarket(quantity string) (*binance.CreateOrderResponse, error) {
	q := quantity
	if quantity == "" {
		balance, err := s.BiapiSvc.GetBalance(s.SpotSpec.Market)
		if err != nil {
			return nil, err
		}
		q = balance
	}
	return s.BiapiSvc.BuyMarket(s.SpotSpec.Coin, s.SpotSpec.Market, q)
}

func (s *SpotSvc) SellMarket(quantity string) (*binance.CreateOrderResponse, error) {
	q := quantity
	if quantity == "" {
		balance, err := s.BiapiSvc.GetBalance(s.SpotSpec.Coin)
		if err != nil {
			return nil, err
		}
		q = balance
	}
	return s.BiapiSvc.SellMarket(s.SpotSpec.Coin, s.SpotSpec.Market, q)
}
