package main

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type SmtpError struct {
	Err error
}

func (e SmtpError) Error() string {
	return e.Err.Error()
}

func (e SmtpError) Code() string {
	return e.Err.Error()[0:3]
}

func NewSmtpError(err error) SmtpError {
	return SmtpError{
		Err: err,
	}
}

var (
	ErrUnresolvableHost = errors.New("unresolvable host")
)

const (
	emailaddr = "suresh.prajapati@olacabs.com"
)

func split(email string) (account, host string) {
	i := strings.LastIndexByte(email, '@')
	account = email[:i]
	host = email[i+1:]
	return
}

func validateHost(email string) error {
	_, host := split(email)
	mx, err := net.LookupMX(host)
	if err != nil {
		return ErrUnresolvableHost
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", mx[0].Host, 25))

	if err != nil {
		return NewSmtpError(err)
	}
	defer client.Close()

	err = client.Hello("localhost")
	if err != nil {
		return NewSmtpError(err)
	}

	err = client.Mail("suresh.prajapati@olacabs.com")
	if err != nil {
		return NewSmtpError(err)
	}

	err = client.Rcpt(email)
	if err != nil {
		return NewSmtpError(err)
	}
	return nil

}

func main() {
	err := validateHost(emailaddr)
	if smtpErr, ok := err.(SmtpError); ok && err != nil {
		fmt.Printf("Code: %s, Msg: %s\n", smtpErr.Code(), smtpErr)
	} else {
		fmt.Println("Email address is valid")
	}
}
