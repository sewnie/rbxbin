package rbxbin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrMissingPkgManifest      = errors.New("package manifest fetch failed")
	ErrInvalidPkgManifest      = errors.New("package manifest is invalid")
	ErrUnhandledPkgManifestVer = errors.New("unhandled package manifest version")
)

// ParsePackages returns a list of packages parsed from the named package manifest.
func ParsePackages(b []byte) ([]Package, error) {
	var pkgs []Package
	m := strings.Split(string(b), "\r\n")

	if (len(m)-2)%4 != 0 {
		return nil, ErrInvalidPkgManifest
	}

	if m[0] != "v0" {
		return nil, fmt.Errorf("%w: %s", ErrUnhandledPkgManifestVer, m[0])
	}

	// Ignore the first element (manifest version) and ignore the additional
	// newline (empty element) added by Roblox.
	for i := 1; i <= len(m)-5; i += 4 {
		zs, err := strconv.ParseInt(m[i+2], 10, 64)
		if err != nil {
			return nil, err
		}
		s, err := strconv.ParseInt(m[i+3], 10, 64)
		if err != nil {
			return nil, err
		}

		pkgs = append(pkgs, Package{
			Name:     m[i],
			Checksum: m[i+1],
			Size:     s,
			ZipSize:  zs,
		})
	}

	return pkgs, nil
}
