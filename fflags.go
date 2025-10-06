package rbxbin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// FFlags is Roblox's Fast Flags implemented in map form.
type FFlags map[string]any

// Apply creates and compiles the FFlags file and
// directory in the named versionDir.
func (f FFlags) Apply(versionDir string) error {
	dir := filepath.Join(versionDir, "ClientSettings")
	path := filepath.Join(dir, "ClientAppSettings.json")

	err := os.Mkdir(dir, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()

	fflags, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(fflags)
	if err != nil {
		return err
	}

	return nil
}
