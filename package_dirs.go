package rbxbin

import (
	"encoding/json"
	"errors"
	"path"
	"strings"
)

// PackageDirectories is a map of where binary deployment packages should go.
type PackageDirectories map[string]string

var ErrDirMapScan = errors.New("could not locate package directory map in installer")

func scanPackageDirectories(b []byte) (pd PackageDirectories) {
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
