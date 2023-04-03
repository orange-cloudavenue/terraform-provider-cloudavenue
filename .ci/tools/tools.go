//go:build tools
// +build tools

package tools

//go:generate go install github.com/hashicorp/go-changelog/cmd/changelog-build

import 
	_ "github.com/hashicorp/go-changelog/cmd/changelog-build"
)