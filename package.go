package rbxbin

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Package is a representation of a Binary package.
type Package struct {
	Name     string
	Checksum string
	Size     int64
	ZipSize  int64
}

// Verify checks the named package source file against it's checksum.
func (p *Package) Verify(src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	fsum := hex.EncodeToString(h.Sum(nil))

	if p.Checksum != fsum {
		return fmt.Errorf("package file %s is corrupted, please re-download or delete package", src)
	}

	return nil
}

// Extract extracts the named package source file to a given destination directory.
func (p *Package) Extract(src, dir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Ensure the destination directory was created
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	for _, f := range r.File {
		// Roblox disobeys Zip and uses non-standard file names, which is why
		// an extract routine is required.
		dest := filepath.Join(dir, strings.ReplaceAll(f.Name, `\`, "/"))

		// ignore the destination directory, it was already created above
		if dir == dest {
			continue
		}

		if !strings.HasPrefix(dest, filepath.Clean(dir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal package file path: %s", dest)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, f.Mode()); err != nil {
				return err
			}

			continue
		}

		if err := unzipFile(f, dest); err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(src *zip.File, name string) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, src.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	z, err := src.Open()
	if err != nil {
		return err
	}
	defer z.Close()

	if _, err := io.Copy(f, z); err != nil {
		return err
	}

	return nil
}
