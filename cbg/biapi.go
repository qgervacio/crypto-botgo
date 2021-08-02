// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"context"
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
