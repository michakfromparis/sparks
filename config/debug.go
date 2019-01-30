package config

// Debug represents the Debug build configuration
type Debug struct {
	enabled bool
}

// Name is the lowercase name of the build configuration
func (d *Debug) Name() string {
	return "debug"
}

// Title is name of the build configuration
func (d *Debug) Title() string {
	return "Debug"
}

// Opt is the short command line option of the build configuration
func (d *Debug) Opt() string {
	return "d"
}

func (d *Debug) String() string {
	return d.Title()
}

// Enabled returns true if the build configuration is enabled
func (d *Debug) Enabled() bool {
	return d.enabled
}

// SetEnabled allows to enable / disable the build configuration
func (d *Debug) SetEnabled(enabled bool) {
	d.enabled = enabled
}
