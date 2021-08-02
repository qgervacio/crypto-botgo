// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"crypto/tls"
	"fmt"
	sm "github.com/xhit/go-simple-mail/v2"
	"time"
)

type EmailSvc struct {
	EmailSpec *EmailSpec
	Client    *sm.SMTPServer
	Username  string
	Pass      string
}

func NewEmailSvc(es *EmailSpec, un, pw string) *EmailSvc {
	client := sm.NewSMTPClient()
	client.Host = es.Server
	client.Port = es.Port
	client.Username = un
	client.Password = pw
	client.Encryption = sm.EncryptionTLS
	client.Authentication = sm.AuthPlain
	client.KeepAlive = true
	client.ConnectTimeout = time.Duration(es.ConnTimeoutSec) * time.Second
	client.SendTimeout = time.Duration(es.SendTimeoutSec) * time.Second
	client.TLSConfig = &tls.Config{InsecureSkipVerify: es.InsecureSkipVerify}

	return &EmailSvc{
		EmailSpec: es,
		Client:    client,
		Username:  un,
		Pass:      pw,
	}
}

func (s *EmailSvc) Send(subject string, body string,
	to []string, cc []string, bcc []string) error {
	em := sm.NewMSG()
	em.SetSubject(subject)
	em.SetFrom(fmt.Sprintf("%s <%s>", s.EmailSpec.Name, s.Username))
	em.SetBody(sm.TextHTML, body)

	for _, e := range to {
		em.AddTo(fmt.Sprintf("%s", e))
	}
	for _, e := range cc {
		em.AddCc(fmt.Sprintf("%s", e))
	}
	for _, e := range bcc {
		em.AddBcc(fmt.Sprintf("%s", e))
	}

	c, err := s.Client.Connect()
	if err != nil {
		return err
	}
	defer c.Close()

	if err := em.Send(c); err != nil {
		return err
	}

	return nil
}
