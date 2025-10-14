package rbxbin

import (
	"bytes"
	"debug/pe"
	"errors"
	"net/http"
)

// Mirror represents a Roblox deployment mirror.
type Mirror string

// DefaultMirror is the default deployment mirror that can be
// used in situations where mirror fallbacks are undesired.
const DefaultMirror Mirror = "https://setup.rbxcdn.com"

var (
	ErrNoMirrorFound = errors.New("no accessible deploy mirror found")

	// As of 2024-02-03:
	//   setup-cfly.rbxcdn.com = roblox-setup.cachefly.net
	//   setup.rbxcdn.com = setup-ns1.rbxcdn.com = setup-ak.rbxcdn.com
	//   setup-hw.rbxcdn.com = setup-ll.rbxcdn.com = does not exist
	Mirrors = []Mirror{
		DefaultMirror,
		Mirror("https://roblox-setup.cachefly.net"),
		Mirror("https://s3.amazonaws.com/setup.roblox.com"),
	}
)

// Mirror returns an available deployment mirror from [Mirrors].
//
// Deployment mirrors may go down, or be blocked by ISPs.
func GetMirror() (Mirror, error) {
	for _, m := range Mirrors {
		resp, err := http.Head(string(m) + "/" + "version")
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return m, nil
		}
	}

	return "", ErrNoMirrorFound
}

// Package returns a URL to a package given a package name
// and a Deployment, relative to the mirror.
func (m Mirror) PackageURL(d *Deployment, pkg string) string {
	// "common" used for all channels, private and public
	return string(m) + "/channel/common/" + d.GUID + "-" + pkg
}

// PackageManifest returns a list of packages for the named deployment.
//
// The given Client's Security will be used in the request if available.
func (m Mirror) GetPackages(d *Deployment) ([]Package, error) {
	body, err := d.get(m.PackageURL(d, "rbxPkgManifest.txt"))
	if err != nil {
		return nil, err
	}

	manif, err := ParsePackages(body)
	if err != nil {
		return nil, err
	}

	return manif, nil
}

// BinaryDirectories retrieves the PackageDirectories for the given deployment from the mirror.
//
// The given Client's Security will be used in the request if available.
func (m Mirror) BinaryDirectories(d *Deployment) (PackageDirectories, error) {
	url := m.PackageURL(d, "Roblox"+d.Type.Short()+"Installer.exe")
	exe, err := d.get(url)
	if err != nil {
		return nil, err
	}

	f, err := pe.NewFile(bytes.NewReader(exe))
	if err != nil {
		return nil, err
	}

	s, err := f.Section(".rdata").Data()
	if err != nil {
		return nil, err
	}

	pd := scanPackageDirectories(s)
	if pd == nil {
		return nil, ErrDirMapScan
	}

	return pd, nil
}
