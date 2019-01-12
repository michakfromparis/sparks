package platform

type Configuration int

const (
	Debug Configuration = iota + 1
	Release
	Shipping
)

func (c Configuration) String() string {
	return [...]string{"Debug", "Release", "Shipping"}[c-1]
}

type Platform interface {
	// name() string
	Name() string
	Opt() string
	Title() string

	Deps()
	Clean()
	Build(Configuration)
}

// var Platforms = []Platform{
// 	Osx{},
// 	Ios{},
// }

var Platforms = map[string]Platform{
	"osx":   Osx{},
	"ios":   Ios{},
	"webgl": WebGl{},
}

var PlatformNames = []string{
	"osx",
	"ios",
	"webgl",
}
