package main

import (
	"fmt"
	"testing"
)

func TestApp_CheckSpeed(t *testing.T) {
	a := NewApp()
	ds, err := a.CheckSpeed()
	if err != nil {
		return
	}
	t.Logf("download speed: %s", ds)
}

func TestApp_CheckDNS(t *testing.T) {
	a := NewApp()
	res, err := a.CheckDNS()
	if err != nil {
		return
	}
	t.Logf(fmt.Sprintf("res: %v", res))
}

func TestApp_CheckLatency(t *testing.T) {
	a := NewApp()
	loss, delay, err := a.CheckLatency()
	if err != nil {
		t.Error("err: ", err)
		return
	}
	t.Logf("loss: %s, delay: %s", loss, delay)
}
