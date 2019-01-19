package configuration

// Release represents the Debug build configuration
type Release struct {
	enabled bool
}

// Name is the lowercase name of the build configuration
func (r *Release) Name() string {
	return "release"
}

// Title is name of the build configuration
func (r *Release) Title() string {
	return "Release"
}

// Opt is the short command line option of the build configuration
func (r *Release) Opt() string {
	return "r"
}

func (r *Release) String() string {
	return r.Title()
}

// Enabled returns true if the build configuration is enabled
func (r *Release) Enabled() bool {
	return r.enabled
}

// SetEnabled allows to enable / disable the build configuration
func (r *Release) SetEnabled(enabled bool) {
	r.enabled = enabled
}
