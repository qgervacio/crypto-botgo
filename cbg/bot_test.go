//  Copyright (c) 2021. Quirino Gervacio
//  MIT License. All Rights Reserved

package cbg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	botSvc = NewBotSvc(specSvc)
)

func Test_Bot_NewBotoSvc_When_Ok_Then_Pass(t *testing.T) {
	assert.NotNil(t, botSvc)
}

//func Test_Bot_Run_When_Ok_Then_Pass(t *testing.T) {
//	botSvc.Run()
//}
