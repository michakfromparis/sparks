package configuration

type Shipping struct {
	enabled bool
}

func (s *Shipping) Name() string {
	return "shipping"
}

func (s *Shipping) Title() string {
	return "Shipping"
}

func (s *Shipping) Opt() string {
	return "s"
}

func (s *Shipping) String() string {
	return s.Title()
}

func (s *Shipping) Enabled() bool {
	return s.enabled
}

func (s *Shipping) SetEnabled(enabled bool) {
	s.enabled = enabled
}
