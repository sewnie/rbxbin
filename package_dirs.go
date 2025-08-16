package rbxbin

import (
	"bytes"
	"debug/pe"
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/sewnie/rbxweb"
)

// PackageDirectories is a map of where binary deployment packages should go.
type PackageDirectories map[string]string

var ErrDirMapScan = errors.New("could not locate package directory map in installer")

// BinaryDirectories retrieves the PackageDirectories for the given deployment from the mirror.
//
// The given Client's Security will be used in the request if available.
func (m Mirror) BinaryDirectories(c *rbxweb.Client, d *Deployment) (PackageDirectories, error) {
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

	pd := scan(s)
	if pd == nil {
		return nil, ErrDirMapScan
	}

	return pd, nil
}

func scan(b []byte) (pd PackageDirectories) {
	start := 0

	for i := 0; i < len(b)-1; i++ {
		if b[i] == '{' && b[i+1] == '"' && b[i-1] == 0 {
			start = i
		}

		if start > 0 && b[i-1] == '}' && b[i-2] == '"' && b[i] == 0 {
			if json.Unmarshal(b[start:i], &pd) != nil {
				pd = nil
				start = 0
				continue
			}

			for p, d := range pd {
				pd[p] = path.Clean(strings.ReplaceAll(d, `\`, "/"))
			}
			break
		}
	}

	return
}
