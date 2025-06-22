package rbxbin

import (
	"errors"

	"github.com/apprehensions/rbxweb"
	"github.com/apprehensions/rbxweb/services/clientsettings"
)

// ErrBadChannel indicates if the mentioned deployment channel does not exist
// or out of permission scope for the current authenticated user.
var ErrBadChannel = errors.New("deployment channel is invalid or unauthorized")

// Deployment is a representation of a Binary's deployment or version.
//
// Channel can either be a given channel, or empty - in which Roblox
// will consider the 'default' channel.
//
// In all things related to the Roblox API, the default channel is empty,
// or 'live'/'LIVE' on clientsettings. On the Client/Studio, the default channel
// is (or can be) 'production'. This behavior is undocumented, so it is reccomended
// to simply provide an empty string for the channel.
type Deployment struct {
	Type    clientsettings.BinaryType
	Channel string
	GUID    string
}

// FetchDeployment returns the latest deployment information for the given
// Roblox binary type with the given deployment channel, using the given client.
func GetDeployment(c *rbxweb.Client, bt clientsettings.BinaryType, channel string) (*Deployment, error) {
	cv, err := c.ClientSettingsV1.GetClientVersion(bt, channel)
	if err == nil {
		return &Deployment{
			Type:    bt,
			Channel: channel,
			GUID:    cv.GUID,
		}, nil
	}

	var apiError rbxweb.ErrorResponse
	if errors.As(err, &apiError) {
		if apiError.Code == 5 {
			return nil, ErrBadChannel
		}
	}
	return nil, err
}
