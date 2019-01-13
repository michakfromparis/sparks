package configurations

type Release struct {
	enabled bool
}

func (r Release) Name() string {
	return "release"
}

func (r Release) Title() string {
	return "Release"
}

func (r Release) Opt() string {
	return "r"
}

func (r Release) String() string {
	return r.Title()
}

func (r Release) Enabled() bool {
	return true
}

func (r Release) SetEnabled(enabled bool) {
	r.enabled = enabled
}
