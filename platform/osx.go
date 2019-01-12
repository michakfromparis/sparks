package platform

import log "github.com/Sirupsen/logrus"

type Osx struct {
}

func (p Osx) Opt() string {
	return "o"
}

func (p Osx) Title() string {
	return "OSX"
}

func (p Osx) Name() string {
	return "osx"
}

func (p Osx) Deps() {
	log.Info("Installing dependencies for " + p.Title())
}
func (p Osx) Clean() {

}
func (p Osx) Prebuild() {

}
func (p Osx) Generate() {

}
func (p Osx) Build() {

}
func (p Osx) Sign() {

}
func (p Osx) Wrap() {

}
func (p Osx) Postbuild() {

}
