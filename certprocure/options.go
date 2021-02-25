package certprocure

import (
	"github.com/BurntSushi/toml"
	"github.com/hatchify/errors"
)

func newOptions() (op *Options, err error) {
	var o Options
	if _, err = toml.DecodeFile("./letsEncrypt.toml", &o); err != nil {
		return
	}

	// Validate options
	if err = o.Validate(); err != nil {
		// Options are not valid, return
		return
	}

	// Set default values
	o.setDefaults()

	// Assign reference to options
	op = &o
	return
}

// Options are the options used for Let's Encrypt SSL procurement
type Options struct {
	Email     string `toml:"email"`
	Domain    string `toml:"domain"`
	Directory string `toml:"directory"`

	Port    string `toml:"port"`
	TLSPort string `toml:"tls-port"`
}

// Validate will validate a set of options
func (o *Options) Validate() (err error) {
	var errs errors.ErrorList
	if len(o.Email) == 0 {
		// Email is required and not found, push error
		errs.Push(ErrInvalidEmail)
	}

	if len(o.Domain) == 0 {
		// Domain is required and not found, push error
		errs.Push(ErrInvalidDomain)
	}

	return errs.Err()
}

func (o *Options) setDefaults() {
	if len(o.Directory) == 0 {
		out.Notificationf("No directory found, using default directory of %s", defaultDirectory)
		o.Directory = defaultDirectory
	}
}
