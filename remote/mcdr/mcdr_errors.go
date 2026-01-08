package mcdr

import "errors"

var (
	ErrorGhApi         = errors.New("error from GitHub API")
	ErrPluginNotFound  = errors.New("MCDR plugin not found")
	ErrVersionNotFound = errors.New("requested version not found")
)
