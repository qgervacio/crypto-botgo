// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	spotSvc = NewSpotSvc(botSvc.SpecSvc.Spec.SpotSpec, biapiSvc)
)

func Test_Spot_NewSpotSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, spotSvc)
}
