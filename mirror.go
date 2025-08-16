package rbxbin

import (
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
