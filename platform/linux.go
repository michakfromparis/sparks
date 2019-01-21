package platform

import (
	"fmt"
	"path/filepath"

	"github.com/michaKFromParis/sparks/sys"

	"github.com/michaKFromParis/sparks/config"
	"github.com/michaKFromParis/sparks/errx"
	"github.com/michaKFromParis/sparks/sparks"

	log "github.com/Sirupsen/logrus"
)

// Linux represents the Linux platform
type Linux struct {
	enabled       bool
	configuration sparks.Configuration
}

// Name is the lowercase name of the platform
func (l *Linux) Name() string {
	return "linux"
}

// Title is name of the platform
func (l *Linux) Title() string {
	return "Linux"
}

// Opt is the short command line option of the platform
func (l *Linux) Opt() string {
	return "l"
}

func (l *Linux) String() string {
	return l.Title()
}

// Enabled returns true if the platform is enabled
func (l *Linux) Enabled() bool {
	return l.enabled
}

// SetEnabled allows to enable / disable the platform
func (l *Linux) SetEnabled(enabled bool) {
	l.enabled = enabled
}

// Get installs the platform dependencies
func (l *Linux) Get() error {
	log.Info("Installing dependencies for " + l.Title())
	_, err := sys.Execute("apt-get", "update")
	if err = l.aptGet("cmake"); err != nil {
		return err
	}
	if err = l.aptGet("zlib1g-dev"); err != nil {
		return err
	}
	// if err = l.aptGet("gcc"); err != nil {
	// 	return err
	// }
	if err = l.aptGet("ccache"); err != nil {
		return err
	}
	if err = l.aptGet("freeglut3-dev"); err != nil {
		return err
	}
	if err = l.aptGet("uuid-dev"); err != nil {
		return err
	}
	if err = l.aptGet("libopenal-dev"); err != nil {
		return err
	}
	if err = l.aptGet("libvlc-dev"); err != nil {
		return err
	}
	if err = l.aptGet("libglew-dev"); err != nil {
		return err
	}
	if err = l.aptGet("libsdl2-dev"); err != nil {
		return err
	}
	return nil
}

func (l *Linux) aptGet(name string) error {
	log.Debugf("sparks get %s", name)
	output, err := sys.Execute("apt-get", "install", "-y", name)
	if err != nil {
		return fmt.Errorf("failed to install %s: %s", name, output)
	}
	return nil

}

// Clean cleans the platform build
func (l *Linux) Clean() error {
	return nil
}

// Build builds the platform
func (l *Linux) Build(configuration sparks.Configuration) error {
	l.configuration = configuration
	l.prebuild()
	l.generate()
	l.compile()
	l.postbuild()
	return nil
}

func (l *Linux) prebuild() {
}

func (l *Linux) generate() {
	log.Info("sparks project generate --linux")

	cmake := sparks.NewCMake(l, l.configuration)
	cmake.AddArg("-GCodeBlocks - Unix Makefiles")
	cmake.AddDefine("OS_LINUX", "1")
	projectsPath := filepath.Join(config.OutputDirectory, "projects", l.Title()+"-"+l.configuration.Title())
	out, err := cmake.Run(projectsPath)
	if err != nil {
		errx.Fatalf(err, "sparks project generate failed: "+out)
	}
	//log.Trace("cmake output" + out)
}

func (l *Linux) compile() {
	log.Info("sparks project compile --linux")
	projectsPath := filepath.Join(config.OutputDirectory, "projects", l.Title()+"-"+l.configuration.Title())
	sys.ExecuteEx("make", projectsPath, true, "-j8")
}

func (l *Linux) postbuild() {

}
