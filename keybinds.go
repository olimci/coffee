package coffee

import "fmt"

type KeybindHandle struct {
	coffee   *Coffee
	keybinds *Keybinds
	stream   <-chan KeybindEvent
}

func (c *Coffee) Keybinds(bindings []Keybind) (*KeybindHandle, error) {
	if len(bindings) == 0 {
		return nil, fmt.Errorf("keybinds cannot be empty")
	}

	out := make(chan KeybindEvent, keybindBufferSize)
	keybinds := newKeybinds(bindings, out)
	if err := c.AddSubmodel(keybinds, WithSection(SectionFooter), WithFocusBehind()); err != nil {
		return nil, err
	}

	return &KeybindHandle{
		coffee:   c,
		keybinds: keybinds,
		stream:   out,
	}, nil
}

func (h *KeybindHandle) Events() <-chan KeybindEvent {
	return h.stream
}

func (h *KeybindHandle) Set(bindings []Keybind) error {
	return h.coffee.send(msgKeybindSet{keybinds: h.keybinds, bindings: bindings})
}

func (h *KeybindHandle) Clear() error {
	return h.coffee.send(msgKeybindClear{keybinds: h.keybinds})
}
