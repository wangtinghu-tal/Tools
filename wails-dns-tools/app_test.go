package main

import (
	"fmt"
	"testing"
)

func TestApp_CheckSpeed(t *testing.T) {
	a := NewApp()
	ds := a.CheckSpeed()
	t.Logf("download speed: %.2fMb/s", ds)
}

func TestApp_CheckDNS(t *testing.T) {
	a := NewApp()
	res:= a.CheckDNS()
	t.Logf(fmt.Sprintf("res: %v", res))
}

func TestApp_CheckLatency(t *testing.T) {
	a := NewApp()
	res := a.CheckLatency()
	t.Logf("loss: %.2f, delay: %.2f", res[0], res[1])
}
