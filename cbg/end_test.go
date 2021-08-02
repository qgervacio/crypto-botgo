// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	endSvc = NewEndSvc(botSvc.SpecSvc.Spec.EndSpec, biapiSvc)
)

func Test_End_NewEndSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, endSvc)
}
