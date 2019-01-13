package configuration

type Debug struct {
	enabled bool
}

func (d Debug) Name() string {
	return "debug"
}

func (d Debug) Title() string {
	return "Debug"
}

func (d Debug) Opt() string {
	return "d"
}

func (d Debug) String() string {
	return d.Title()
}

func (d Debug) Enabled() bool {
	return d.enabled
}

func (d Debug) SetEnabled(enabled bool) {
	d.enabled = enabled
}
