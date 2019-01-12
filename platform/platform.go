package platform

type Platform interface {
	// name() string
	Name() string
	Opt() string
	Title() string

	Deps()
	Clean()
	Prebuild()
	Generate()
	Build()
	Sign()
	Wrap() // package was taken :)
	Postbuild()
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
