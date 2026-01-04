package rbxbin

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
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
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, src.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	if info.Size() > 0 {
		// Compare the zip file's hash against the existing file, if they are
		// identical, don't extract and just skip the file to save on writing.
		identical, err := func() (bool, error) {
			z, err := src.Open()
			if err != nil {
				return false, err
			}
			defer z.Close()

			hasher := sha256.New()

			io.Copy(hasher, f)
			originHash := hasher.Sum(nil)

			hasher.Reset()

			io.Copy(hasher, z)
			targetHash := hasher.Sum(nil)

			if bytes.Compare(originHash[:], targetHash[:]) == 0 {
				return true, nil
			}

			return false, nil
		}()

		if err != nil {
			return err
		}

		if identical {
			return nil
		}
	}

	z, err := src.Open()
	if err != nil {
		return err
	}
	defer z.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, z); err != nil {
		return err
	}

	return nil
}
