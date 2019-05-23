package main

import (
	"bufio"
	"fmt"

	"github.com/ianmcmahon/joybehar/controls"
	"github.com/jacobsa/go-serial/serial"
)

// on error, drops the message entirely
type interceptAction func(dcsMsg) (dcsMsg, error)

type panelAgent struct {
	options    serial.OpenOptions
	dcsAgent   *dcsAgent
	group      controls.Modal
	intercepts map[string]map[controls.Mode]interceptAction
}

func PanelAgent(port string, dcsAgent *dcsAgent, group controls.Modal) *panelAgent {
	agent := &panelAgent{
		options: serial.OpenOptions{
			PortName:        port,
			BaudRate:        250000,
			DataBits:        8,
			StopBits:        1,
			MinimumReadSize: 4,
		},
		dcsAgent:   dcsAgent,
		group:      group,
		intercepts: make(map[string]map[controls.Mode]interceptAction, 0),
	}

	return agent
}

func (a *panelAgent) ReLabel(label string, mode controls.Mode, alt string) {
	a.Intercept(label, mode, func(m dcsMsg) (dcsMsg, error) {
		m.message = alt
		return m, nil
	})
}

func (a *panelAgent) Intercept(label string, mode controls.Mode, action interceptAction) {
	if _, ok := a.intercepts[label]; !ok {
		a.intercepts[label] = make(map[controls.Mode]interceptAction, 0)
	}
	a.intercepts[label][mode] = action
}

func (a *panelAgent) Receive() {
	port, err := serial.Open(a.options)
	if err != nil {
		fmt.Printf("Error opening serial port: %v\n", err)
		fmt.Printf("Disabling dcsbios panel\n")
		return
	}
	defer port.Close()

	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		msg, err := DCSMsg(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if m, ok := a.intercepts[msg.message]; ok {
			if action, ok := m[a.group.Mode()]; ok {
				msg, err = action(msg)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
		a.dcsAgent.SendMsg(msg)
	}
}
