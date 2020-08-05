package main

import "github.com/fanghongbo/ops-hbs/common/g"

var (
	Version    = "v1.0"
	BinaryName = "ops-hbs"
)

func init() {
	g.BinaryName = BinaryName
	g.Version = Version
}
