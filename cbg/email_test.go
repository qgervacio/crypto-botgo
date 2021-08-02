// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	emailSvc = NewEmailSvc(
		botSvc.SpecSvc.Spec.EmailSpec,
		botSvc.SpecSvc.Spec.CredSpec.EmUser,
		botSvc.SpecSvc.Spec.CredSpec.EmPass,
	)
)

func Test_Email_NewEmailSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, emailSvc)
}

func Test_Email_Send_When_Ok_Then_Pass(t *testing.T) {
	err := emailSvc.Send(
		"crypto-botgo test",
		"email_test",
		[]string{botSvc.SpecSvc.Spec.CredSpec.NotEm},
		[]string{},
		[]string{},
	)
	assert.Nil(t, err)
}
