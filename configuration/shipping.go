package configuration

// Shipping represents the Shipping build configuration
type Shipping struct {
	enabled bool
}

// Name is the lowercase name of the build configuration
func (s *Shipping) Name() string {
	return "shipping"
}

// Title is name of the build configuration
func (s *Shipping) Title() string {
	return "Shipping"
}

// Opt is the short command line option of the build configuration
func (s *Shipping) Opt() string {
	return "s"
}

func (s *Shipping) String() string {
	return s.Title()
}

// Enabled returns true if the build configuration is enabled
func (s *Shipping) Enabled() bool {
	return s.enabled
}

// SetEnabled allows to enable / disable the build configuration
func (s *Shipping) SetEnabled(enabled bool) {
	s.enabled = enabled
}
