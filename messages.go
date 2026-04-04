package coffee

// log a message to a section
type msgLog struct {
	message string
	section Section
	opts    logOptions
}

// clear logs from a section, preserving submodels.
type msgClear struct {
	section Section
}

// create a submodel in a section
type msgSubmodel struct {
	submodel Submodel
	section  Section
	behind   bool
}

type msgWindowTitle struct {
	title string
}

type MsgFocusGained struct{}

type MsgFocusLost struct{}
