package rbxbin

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sewnie/rbxweb"
)

// ErrBadChannel indicates if the mentioned deployment channel does not exist
// or out of permission scope for the current authenticated user.
var ErrBadChannel = errors.New("deployment channel is invalid or unauthorized")

// Deployment is a representation of a Binary's deployment or version.
//
// In all things related to the Roblox API, the default channel is empty,
// or 'live'/'LIVE' on clientsettings. On the Client/Studio, the default channel
// is (or can be) 'production'. This behavior is undocumented, so it is reccomended
// to simply provide an empty string for the channel.
//
// For more details about the deployment channel, see [GetDeployment].
type Deployment struct {
	Type    rbxweb.BinaryType
	Channel string
	GUID    string

	client *rbxweb.Client
}

type tokenTransport struct {
	base  http.RoundTripper
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	if t.token != "" {
		clone.Header.Add("Roblox-Channel-Token", t.token)
	}
	clone.Header.Add("User-Agent", "rbxbin/v0.0.0")
	return t.base.RoundTrip(clone)
}

// FetchDeployment returns the latest deployment information for the given
// Roblox binary type with the given deployment channel, using the given client.
//
// If the given channel is empty, the current authenticated user in the Client
// will be used to get the channel, which can be public or private, or can be
// "LIVE" if the Client is not authenticated.
func GetDeployment(c *rbxweb.Client, bt rbxweb.BinaryType, channel string) (*Deployment, error) {
	var token string
	if channel == "" {
		uc, err := c.ClientSettingsV2.GetUserChannel(&bt)
		if err != nil {
			return nil, fmt.Errorf("user: %w", err)
		}
		channel = uc.Channel
		token = uc.Token
	}

	client := rbxweb.NewClient()
	client.Transport = &tokenTransport{base: c.Transport, token: token}
	if c.Transport == nil {
		client.Transport.(*tokenTransport).base = http.DefaultTransport
	}

	cv, err := client.ClientSettingsV2.GetClientVersion(bt, channel)
	if err == nil {
		return &Deployment{
			Type:    bt,
			Channel: channel,
			GUID:    cv.GUID,

			client: client,
		}, nil
	}

	var apiError rbxweb.Error
	if errors.As(err, &apiError) {
		if apiError.Code == 5 {
			return nil, ErrBadChannel
		}
	}
	return nil, err
}

func (d *Deployment) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
