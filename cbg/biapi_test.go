// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	biapiSvc = NewBiapiSvc(
		botSvc.SpecSvc.Spec.BiapiSpec,
		botSvc.SpecSvc.Spec.CredSpec.BiapiAk,
		botSvc.SpecSvc.Spec.CredSpec.BiapiSk,
	)
)

func Test_Biapi_NewBiapiSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, biapiSvc)
}
