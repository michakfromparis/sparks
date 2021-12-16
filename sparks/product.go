package sparks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"leblox.com/sparks-cli/v2/conf"
	"leblox.com/sparks-cli/v2/errx"
)

// Product struct
type Product struct {
	Name      string   `yaml:"Name"`
	Version   string   `yaml:"Version"`
	Language  string   `yaml:"Language"`
	Locales   string   `yaml:"Locales"`
	Platforms []string `yaml:"Platforms"`
	View      struct {
		DefaultOrientation    string   `yaml:"DefaultOrientation"`
		SupportedOrientations []string `yaml:"SupportedOrientations"`
		Resolution            string   `yaml:"Resolution"`
		Fullscreen            string   `yaml:"Fullscreen"`
	} `yaml:"View"`
	Assets struct {
		Pack string `yaml:"Pack"`
	} `yaml:"Assets"`
	NativeModules []string `yaml:"NativeModules"`
	MarkettingURL string   `yaml:"MarkettingUrl"`
	SupportURL    string   `yaml:"SupportUrl"`
	PrivacyURL    string   `yaml:"PrivacyUrl"`
	Copyright     string   `yaml:"Copyright"`
	Windows       struct {
		MsiProductID string `yaml:"msiProductId"`
	} `yaml:"Windows"`
	IOS struct {
		BundleIdentifier string `yaml:"BundleIdentifier"`
	} `yaml:"iOS"`
	Description string `yaml:"Description"`
	Keywords    string `yaml:"Keywords"`
	ReviewNotes string `yaml:"ReviewNotes"`

	// private fields
	sparksFilename string
}

// Load loads a .sparks file
func (p *Product) Load() error {
	filename, err := p.findSparksFile()
	if err != nil {
		errx.Fatalf(err, "Could not find a sparks file at "+conf.SourceDirectory)
	}
	log.Debug("loading product from " + filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errx.Fatalf(err, "Could not read: "+filename)
	}
	err = yaml.Unmarshal(bytes, p)
	if err != nil {
		errx.Fatalf(err, "yaml reader: "+filename)
	}
	return nil
}

// Save saves a .sparks file
func (p *Product) Save() {
	filename, err := p.findSparksFile()
	if err != nil {
		errx.Fatalf(err, "Could not find a sparks file at "+conf.SourceDirectory)
	}
	log.Debug("saving product to " + filename)
	data, err := yaml.Marshal(p)
	if err != nil {
		errx.Fatalf(err, "yaml writer: "+filename)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		errx.Fatalf(err, "Could not write to: "+filename)
	}
}

// look for a .sparks file in conf.SourceDirectory) and return it
func (p *Product) findSparksFile() (string, error) {
	if p.sparksFilename != "" {
		return p.sparksFilename, nil
	}
	log.Tracef("opening %s", conf.SourceDirectory)
	f, err := os.Open(conf.SourceDirectory)
	if err != nil {
		return "", errorx.Decorate(err, "Could not open SourceDirectory: "+conf.SourceDirectory)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		if err = f.Close(); err != nil {
			return "", errorx.Decorate(err, "Could not close SourceDirectory: "+conf.SourceDirectory)
		}
		return "", errorx.Decorate(err, "Could not read SourceDirectory: "+conf.SourceDirectory)
	}
	if err = f.Close(); err != nil {
		return "", errorx.Decorate(err, "Could not close SourceDirectory: "+conf.SourceDirectory)
	}
	// log.Trace("files in SourceDirectory:")
	p.sparksFilename = ""
	for _, file := range files {
		// log.Trace(file.Name())
		if strings.HasSuffix(file.Name(), ".sparks") {
			p.sparksFilename = filepath.Join(conf.SourceDirectory, file.Name())
			log.Debugf("found a .sparks file at: %s", p.sparksFilename)
			break
		}
	}
	if p.sparksFilename == "" {
		return "", errorx.Decorate(nil, "could not find a .sparks file at "+conf.SourceDirectory)
	}
	return p.sparksFilename, nil
}

// Sample is used to generate sample sparks file
func (p *Product) Sample() {
	log.Debug("filling product with sample data")
	p.Name = "Sparks"
	p.Version = "2.2.5"
	p.Language = "cpp"
	p.Locales = "en-US"
	p.Platforms = []string{"Android (Standard)", "iOS (Standard)", "OSX (Standard)", "Windows (Standard)", "Linux (Standard)", "WebGl (Standard)"}
	p.View.DefaultOrientation = "Landscape Left"
	p.View.SupportedOrientations = []string{"Landscape Left", "Landscape Right"}
	p.View.Resolution = "1280x720"
	p.View.Fullscreen = "No"
	p.Assets.Pack = "No"
	p.NativeModules = []string{"Network", "OpenAL"}
	p.MarkettingURL = "http://sparkslight.com"
	p.SupportURL = "http://sparkslight.com/support"
	p.PrivacyURL = "http://sparkslight.com/privacy"
	p.Copyright = "© 2015 iVoltage"
	p.Windows.MsiProductID = "5F4B755F4B75"
	p.IOS.BundleIdentifier = "me.iVoltage.SparksPlayer"
	p.Description = `Sparks offers a new way of creating beautiful cross-platform products directly from Photoshop, Affter Effects, 3D Studio Max or Maya.
Sparks Player is the standalone player for the Sparks toolchain to test your product on your device, through the network.`
	p.Keywords = "Sparks, iVoltage, lua, C++, games, Photoshop, After Effects, Emmanuel Briney, Michel Courtine"
	p.ReviewNotes = `Hi, This is a data centric player that doesn't do anything without the tool.
	It should just create an OpenGL View and wait for a network connection to start running a product.
	Let me know if you have any question`
}
