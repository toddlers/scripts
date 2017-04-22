package main

import (
	"github.com/cenk/backoff"
)

var backoffTickers map[string]*backoff.Ticker

// Time wait (execution delay) using backoff factor algo for
// different types / name.
func BackOffWait(name string, wait bool) {
	if backoffTickers == nil {
		backoffTickers = make(map[string]*backoff.Ticker)
	}
	if backoffTickers[name] == nil {
		backoffTickers[name] = backoff.NewTicker(backoff.NewExponentialBackOff())
	}
	if wait {
		_ = <-backoffTickers[name].C
	} else {
		backoffTickers[name] = backoff.NewTicker(backoff.NewExponentialBackOff())
	}
}

func BackOffWaitIfError(err error, name string) {
	if err != nil {
		BackOffWait(name, true)
	} else {
		BackOffWait(name, false)
	}
}
