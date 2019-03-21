package checkd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Checker is the interface that checks should satisfy.
type Checker interface {
	Check() error
	Init(*viper.Viper) error
	Interval() time.Duration
	Name() string
}

// Checks is the list of checks to be run.
var checks = []Checker{}

// Register adds a new check to the list of checks that will be run.
func Register(c Checker) {
	log.Printf("registering %q", c.Name())
	checks = append(checks, c)
}

// Init initializes all checks.
func Init(v *viper.Viper) error {
	log.Printf("checkd: found %d registered checks", len(checks))
	for _, c := range checks {
		log.Printf("checkd: initializing %q", c.Name())
		if err := c.Init(v.Sub(c.Name())); err != nil {
			return fmt.Errorf("%s: %s", c.Name(), err)
		}
	}
	return nil
}

// Run runs all checks concurrently.
func Run() {
	for _, c := range checks {
		go func(c Checker) {
			for {
				c.Check()
				time.Sleep(c.Interval())
			}
		}(c)
	}
}
