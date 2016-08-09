package tapedeck

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFilename(path string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(path)

	return strings.TrimSuffix(filename, ext)
}

func FingerprintFile(path string) (fp string) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()

		hash := sha256.New()

		if _, err := io.Copy(hash, file); err == nil {
			fp = hex.EncodeToString(hash.Sum(nil))
		}
	}
	return
}

func CopyFile(src, dest string, perm uint) error {
	destDir := filepath.Dir(dest)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(destDir, os.FileMode(perm)); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func MoveFile(src, dest string, perm uint) error {
	if err := CopyFile(src, dest, perm); err != nil {
		return err
	}
	return os.Remove(src)
}

func ReducePath(path, basepath string) error {
	for path != basepath {
		if !IsEmptyFile(path) {
			break
		}

		if err := os.Remove(path); err != nil {
			return err
		}

		path, _ = filepath.Split(path)
	}

	return nil
}

func IsEmptyFile(filepath string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false
	}

	if !stat.IsDir() {
		return stat.Size() == 0
	}

	_, err = file.Readdirnames(1)
	return err == io.EOF
}
