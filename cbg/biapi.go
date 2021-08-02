// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"context"
	"fmt"
	b "github.com/adshao/go-binance/v2"
)

type BiapiSvc struct {
	BiapiConf *BiapiSpec
	Client    *b.Client
}

func NewBiapiSvc(bs *BiapiSpec, apiKey, secretKey string) *BiapiSvc {
	b.UseTestnet = bs.Test
	return &BiapiSvc{
		BiapiConf: bs,
		Client:    b.NewClient(apiKey, secretKey),
	}
}

func (s *BiapiSvc) ListPrices(symbol string) ([]*b.SymbolPrice, error) {
	return s.Client.
		NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
}

func (s *BiapiSvc) ListOpenOrders(symbol string) ([]*b.Order, error) {
	return s.Client.
		NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
}

func (s *BiapiSvc) CancelOrder(symbol string, orderID int64) (*b.CancelOrderResponse, error) {
	return s.Client.
		NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderID).
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
}

func (s *BiapiSvc) GetAccount() (*b.Account, error) {
	return s.Client.
		NewGetAccountService().
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
}

func (s *BiapiSvc) BuyMarket(coin, market, quantity string) (*b.CreateOrderResponse, error) {
	q := quantity
	if quantity == "" {
		balance, err := s.GetBalance(market)
		if err != nil {
			return nil, err
		}
		q = balance
	}
	return s.Client.
		NewCreateOrderService().
		Symbol(fmt.Sprintf("%s%s", coin, market)).
		Side(b.SideTypeBuy).
		Type(b.OrderTypeMarket).
		QuoteOrderQty(q).
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
}

func (s *BiapiSvc) SellMarket(coin, market, quantity string) (*b.CreateOrderResponse, error) {
	q := quantity
	if quantity == "" {
		balance, err := s.GetBalance(coin)
		if err != nil {
			return nil, err
		}
		q = balance
	}
	return s.Client.
		NewCreateOrderService().
		Symbol(fmt.Sprintf("%s%s", coin, market)).
		Side(b.SideTypeSell).
		Type(b.OrderTypeMarket).
		Quantity(q).
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
}

func (s *BiapiSvc) GetBalance(asset string) (string, error) {
	a, err := s.Client.NewGetAccountService().
		Do(context.Background(), b.WithRecvWindow(s.BiapiConf.RecvWindow))
	if err != nil {
		return "", err
	}
	for _, e := range a.Balances {
		if e.Asset == asset {
			return e.Free, nil
		}
	}
	return "", fmt.Errorf("asset %s not found", asset)
}
