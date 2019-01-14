package sparks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	yaml "gopkg.in/yaml.v2"
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

func (p *Product) Load() error {
	filename := p.findSparksFile()
	log.Debug("loading product from " + filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		errx.Fatalf(err, "Could not read: "+filename)
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		errx.Fatalf(err, "yaml reader: "+filename)
	}
	return nil
}

func (p *Product) Save() {
	filename := p.findSparksFile()
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

func (p *Product) findSparksFile() string {
	if p.sparksFilename != "" {
		return p.sparksFilename
	}
	log.Tracef("opening %s", config.SourceDirectory)
	f, err := os.Open(config.SourceDirectory)
	if err != nil {
		errx.Fatalf(err, "Could not open SourceDirectory: "+config.SourceDirectory)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		errx.Fatalf(err, "Could not read SourceDirectory: "+config.SourceDirectory)
	}
	if err = f.Close(); err != nil {
		errx.Fatalf(err, "Could not close SourceDirectory: "+config.SourceDirectory)
	}
	log.Trace("files in SourceDirectory:")
	p.sparksFilename = ""
	for _, file := range files {
		log.Trace(file.Name())
		if strings.HasSuffix(file.Name(), ".sparks") {
			p.sparksFilename = filepath.Join(config.SourceDirectory, file.Name())
			log.Debugf("found a .sparks file at: %s", p.sparksFilename)
			break
		}
	}
	if p.sparksFilename == "" {
		errx.Error("Could not find a .sparks file at " + config.SourceDirectory)
	}
	return p.sparksFilename
}

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