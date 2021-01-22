package certprocure

import (
	"fmt"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	legolog "github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/registration"
	"github.com/hatchify/atoms"
	"github.com/hatchify/errors"
	"github.com/hatchify/scribe"
)

const (
	// ErrInvalidEmail is returned when the let's encrypt email is missing from the environment values
	ErrInvalidEmail = errors.Error("cannot create SSL certificate without 'lets-encrypt-email' environment value")
	// ErrInvalidDomain is returned when the let's encrypt domain is missing from the environment values
	ErrInvalidDomain = errors.Error("cannot create SSL certificate without 'lets-encrypt-domain' environment value")
)

const (
	defaultPort      = "80"
	defaultTLSPort   = "443"
	defaultDirectory = "tls"
)

var (
	// Output writer
	out *scribe.Scribe = scribe.New("Let's Encrypt")
	// Default global registration options
	registrationOpts = registration.RegisterOptions{TermsOfServiceAgreed: true}
)

// New will return a new instance of CertProcure
func New() (cp *CertProcure, err error) {
	var c CertProcure
	// Get options from environment variables
	if c.o, err = newOptions(); err != nil {
		// Error parsing options, return
		out.Errorf("error parsing options: %v", err)
		return
	}

	// Create log wrapper
	var lw logWrapper
	// Set log wrapper as our lego logger
	legolog.Logger = &lw
	go c.watch()
	cp = &c
	return
}

// CertProcure manages certificate procurement
type CertProcure struct {
	o *Options

	closed atoms.Bool
}

func (c *CertProcure) watch() {
	for !c.closed.Get() {
		c.attemptRenew()
		time.Sleep(time.Hour * 24)
	}
}

func (c *CertProcure) attemptRenew() {
	ok, err := needsCertificate(c.o.Directory)
	switch {
	case err != nil:
		// Error encountered while checking certificate, return
		out.Errorf("error checking certificate: %v", err)
		return
	case !ok:
		return
	}

	out.Notification("Certificate is expired (or expiring soon), executing renewal process")

	if err = c.Renew(); err != nil {
		out.Errorf("Error encountered during renewal: %v", err)
		return
	}

	return
}

// Renew will renew a certificate
func (c *CertProcure) Renew() (err error) {
	var u *User
	// Create a new user
	if u, err = newUser(c.o.Email); err != nil {
		// Error creating new user, return
		err = fmt.Errorf("error creating user: %v", err)
		return
	}

	// Create a new configuration
	config := lego.NewConfig(u)
	config.Certificate.KeyType = certcrypto.RSA2048

	var client *lego.Client
	// Initialize a new client
	if client, err = newClient(c.o, config); err != nil {
		// Error initializing new client, return
		err = fmt.Errorf("error initializing client: %v", err)
		return
	}

	out.Success("Client created")

	// Register user using Client
	if err = u.Register(client); err != nil {
		// Error registering user, return
		err = fmt.Errorf("error registering user \"%s\": %v", u.Email, err)
		return
	}

	out.Success("User registered")

	// Make request
	request := makeRequest(c.o.Domain)

	var certificates *certificate.Resource
	// Obtain certificates
	if certificates, err = client.Certificate.Obtain(request); err != nil {
		// Error obtaining certificates, return
		err = fmt.Errorf("error obtaining certificates: %v", err)
		return
	}

	out.Success("Certificates obtained")

	// Save certificates to file
	if err = saveCertificates(c.o.Directory, certificates); err != nil {
		// Error saving certificates, return
		err = fmt.Errorf("error saving certificates: %v", err)
		return
	}

	out.Success("Certificate renewal process complete")
	return
}

// Close will close an instance of CertProcure
func (c *CertProcure) Close() (err error) {
	if !c.closed.Set(true) {
		return errors.ErrIsClosed
	}

	return
}
