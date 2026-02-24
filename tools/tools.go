//go:build tools
// +build tools

package tools

// This file pins tool dependencies so `go mod tidy` doesn't remove them.
// Import tools here using a blank identifier to keep them in go.mod without
// affecting the build.

import (
	_ "github.com/samber/lo"
)
