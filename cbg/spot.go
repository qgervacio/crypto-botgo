// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

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

func (s *SpotSvc) Buy() (string, error) {
	return "Bought at xxx", nil
}

func (s *SpotSvc) Sell() (string, error) {
	return "Sold at xxx", nil
}
