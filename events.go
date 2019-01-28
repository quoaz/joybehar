package main

type Event interface {
	Handle(raw rawEvent)
}

type AxisEvent interface {
}

type pressEvent struct {
	action Action
}

func (e pressEvent) Handle(raw rawEvent) {
	if raw.pressed {
		e.action()
	}
}

func (d *device) PressEvent(button uint8, action Action) *pressEvent {
	return nil
}

func (d *device) DCSToggle(button press, message, off, on string) {
	handler, ok := d.buttonHandlers[button.id]
	if !ok {
		d.buttonHandlers[button.id] = d.ButtonHandler(button.id)
		handler = d.buttonHandlers[button.id]
	}

	handler.Actions[button.longPress][button.mode] = func(e rawEvent) {
		if e.pressed {
			dcsSend(message, on)
		} else {
			dcsSend(message, off)
		}
	}
}

func (d *device) DCS3Pos(upBtn, downBtn press, message, down, off, up string) {
	upHandler, ok := d.buttonHandlers[upBtn.id]
	if !ok {
		d.buttonHandlers[upBtn.id] = d.ButtonHandler(upBtn.id)
		upHandler = d.buttonHandlers[upBtn.id]
	}

	downHandler, ok := d.buttonHandlers[downBtn.id]
	if !ok {
		d.buttonHandlers[downBtn.id] = d.ButtonHandler(downBtn.id)
		downHandler = d.buttonHandlers[downBtn.id]
	}

	upHandler.Actions[upBtn.longPress][upBtn.mode] = func(e rawEvent) {
		if e.pressed {
			dcsSend(message, up)
		} else {
			dcsSend(message, off)
		}
	}
	downHandler.Actions[downBtn.longPress][downBtn.mode] = func(e rawEvent) {
		if e.pressed {
			dcsSend(message, down)
		} else {
			dcsSend(message, off)
		}
	}
}
