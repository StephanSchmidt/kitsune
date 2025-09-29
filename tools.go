//go:build tools
// +build tools

package kitsune

import (
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "github.com/sonatype-nexus-community/nancy"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
