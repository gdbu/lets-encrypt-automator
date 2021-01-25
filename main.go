package main

import (
	"log"

	"github.com/gdbu/lets-encrypt-automator/certprocure"
	"github.com/hatchify/closer"
)

// Init will be called by vroomy on initialization before Configure
func main() {
	var (
		c   *certprocure.CertProcure
		err error
	)

	if c, err = certprocure.New(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cc := closer.New()

	if err = cc.Wait(); err != nil {
		log.Fatal(err)
	}
}
