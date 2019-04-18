// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	// mg contains helpful utility functions, like Deps
)

const version = "v0.0.3"

type Build mg.Namespace

func (b Build) All() {
	b.Darwin()
	b.Linux()
}

func (Build) Darwin() {
	fmt.Print("build darwin... ")
	sh.Run("go", "build", "-ldflags", ldFlags(), "-o", fmt.Sprintf("build/%s", binName("darwin")))
	fmt.Println("OK")
}

func (Build) Linux() {
	fmt.Print("build linux...  ")
	sh.RunWith(map[string]string{
		"GOOS": "linux", "GOARCH": "amd64", "CGO_ENABLED": "0"},
		"go", "build", "-ldflags", ldFlags(), "-o", fmt.Sprintf("build/%s", binName("linux")))
	fmt.Println("OK")
}

func ldFlags() string {
	return fmt.Sprintf(`-X "github.com/ymgyt/ops/cmd.version=%s"`, version)
}

func Version() {
	fmt.Println(version)
}

func binName(os string) string {
	return fmt.Sprintf("ops-%s-amd64-%s", os, version)
}
