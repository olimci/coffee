package coffee

type StatusHandle struct {
	coffee *Coffee
	status *Status
}

func (c *Coffee) Status(message string) (*StatusHandle, error) {
	status := NewStatus(message)
	if err := c.AddSubmodel(status, WithSection(SectionFooter), WithFocusBehind()); err != nil {
		return nil, err
	}
	return &StatusHandle{
		coffee: c,
		status: status,
	}, nil
}

func (h *StatusHandle) Idle(message string) error {
	return h.coffee.send(msgStatusIdle{status: h.status, message: message})
}

func (h *StatusHandle) Working(message string) error {
	return h.coffee.send(msgStatusWorking{status: h.status, message: message})
}

func (h *StatusHandle) Progress(message string, percent float64) error {
	return h.coffee.send(msgStatusProgress{status: h.status, message: message, percent: percent})
}

func (h *StatusHandle) SetProgress(percent float64) error {
	return h.coffee.send(msgStatusProgressValue{status: h.status, percent: percent})
}

func (h *StatusHandle) Message(message string) error {
	return h.coffee.send(msgStatusMessage{status: h.status, message: message})
}

func (h *StatusHandle) Success(message string) error {
	return h.coffee.send(msgStatusSuccess{status: h.status, message: message})
}

func (h *StatusHandle) Error(message string) error {
	return h.coffee.send(msgStatusError{status: h.status, message: message})
}

func (h *StatusHandle) Clear() error {
	return h.coffee.send(msgStatusClear{status: h.status})
}
