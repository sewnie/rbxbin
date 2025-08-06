// Package rbxbin implements various routines and types
// to install or bootstrap a Roblox Binary.
package rbxbin

import (
	"github.com/sewnie/rbxweb"
)

// Aliases to rbxweb BinaryType for easier calling
const (
	WindowsPlayer = rbxweb.BinaryTypeWindowsPlayer
	WindowsStudio = rbxweb.BinaryTypeWindowsStudio64
	MacPlayer     = rbxweb.BinaryTypeMacPlayer
	MacStudio     = rbxweb.BinaryTypeMacStudio
)
