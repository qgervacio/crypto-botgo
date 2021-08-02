// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type EndSvc struct {
	EndSpec  *EndSpec
	BiapiSvc *BiapiSvc
}

func NewEndSvc(es *EndSpec, bs *BiapiSvc) *EndSvc {
	return &EndSvc{
		EndSpec:  es,
		BiapiSvc: bs,
	}
}

func (s *EndSvc) EndByDummyOrder() (bool, error) {
	res, err := s.BiapiSvc.ListOpenOrders(s.EndSpec.DummySymbol)
	if err != nil {
		return false, fmt.Errorf("failed to get dummy order [%v]", err)
	}

	for _, e := range res {
		price, _ := strconv.ParseFloat(e.Price, 64)
		quant, _ := strconv.ParseFloat(e.OrigQuantity, 64)

		if price == s.EndSpec.DummyOrder[0] &&
			quant == s.EndSpec.DummyOrder[1] {
			log.Infof("Found dummy order %s in OrderId %d",
				s.EndSpec.DummySymbol, e.OrderID)

			if _, err := s.BiapiSvc.CancelOrder(s.EndSpec.DummySymbol, e.OrderID); err != nil {
				return false, fmt.Errorf("failed to cancel dummy OrderId %d [%v]", e.OrderID, err)
			}

			// we're done here
			return true, nil
		}
	}

	return false, nil
}
