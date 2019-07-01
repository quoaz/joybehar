package main

import (
	"bufio"
	"fmt"
	"time"

	"github.com/ianmcmahon/joybehar/alert"
	"github.com/ianmcmahon/joybehar/controls"
	"github.com/ianmcmahon/joybehar/dcs"
	"github.com/jacobsa/go-serial/serial"
)

type panelAgent struct {
	options  serial.OpenOptions
	dcsAgent dcs.Agent
	group    controls.Interceptor
}

func PanelAgent(port string, dcsAgent dcs.Agent, group controls.Interceptor) *panelAgent {
	agent := &panelAgent{
		options: serial.OpenOptions{
			PortName:        port,
			BaudRate:        250000,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 4,
		},
		dcsAgent: dcsAgent,
		group:    group,
	}

	return agent
}

func (a *panelAgent) Receive() {
	time.Sleep(500 * time.Millisecond)
	alert.Say("Opening panel")
	port, err := serial.Open(a.options)
	if err != nil {
		alert.Sayf("Error opening serial port: %v", err)
		alert.Say("Disabling dcsbios panel")
		return
	}
	defer port.Close()

	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		msg, err := dcs.DCSMsgFromString(scanner.Text())
		if err != nil {
			alert.Sayf("%v", err)
			continue
		}
		msg, err = a.group.Intercept(msg)
		if err != nil {
			alert.Sayf("%v", err)
			continue
		}

		fmt.Printf("%v\n", msg)
		a.dcsAgent.SendMsg(msg)
	}
	alert.Say("Panel receive exiting")
	go a.Receive()
}
