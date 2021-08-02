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
		botSvc.SpecSvc.ArgsSpec.EmUser,
		botSvc.SpecSvc.ArgsSpec.EmPass,
	)
)

func Test_Email_NewEmailSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, emailSvc)
}

func Test_Email_Send_When_Ok_Then_Pass(t *testing.T) {
	err := emailSvc.Send(
		"crypto-botgo test",
		"email_test",
		[]string{string(loadTestFile("emUser.test"))},
		[]string{},
		[]string{},
	)
	assert.Nil(t, err)
}
