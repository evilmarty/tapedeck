package tapedeck

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_GetFilename(t *testing.T) {
	tests := map[string]string{
		"/foobar.txt":     "foobar",
		"/foobar":         "foobar",
		"/tmp/foobar.txt": "foobar",
		"/tmp/foobar":     "foobar",
		"foobar.txt":      "foobar",
		"foobar":          "foobar",
	}

	for filepath, expected := range tests {
		if actual := GetFilename(filepath); actual != expected {
			t.Errorf("Expected %s, received %s", expected, actual)
		}
	}
}

func Test_FingerprintFile(t *testing.T) {
	tests := map[string]string{
		"fixtures/sample.mp3": "02d4393dcfe35353d42025d7ca919c8be032774547015854dce1415381c171db",
		"fixtures/nofile":     "",
	}

	for filepath, expected := range tests {
		if actual := FingerprintFile(filepath); actual != expected {
			t.Errorf("Expected %s, received %s", expected, actual)
		}
	}
}

func Test_CopyFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "copy-file-")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer os.RemoveAll(tmp)

	src := "fixtures/sample.mp3"
	dest := filepath.Join(tmp, "foobar", "test.mp3")

	if err := CopyFile(src, dest, 0755); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("File did not copy to %s", dest)
	}

	if FingerprintFile(src) != FingerprintFile(dest) {
		t.Errorf("Expected %s to be the same as %s", dest, src)
	}
}

func Test_MoveFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "move-file-")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "foobar", "old.mp3")
	dest := filepath.Join(tmp, "foobar", "new.mp3")

	if err := CopyFile("fixtures/sample.mp3", src, 0755); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := MoveFile(src, dest, 0755); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("File did not copy to %s", dest)
	}

	if _, err := os.Stat(src); os.IsExist(err) {
		t.Errorf("File did not remove %s", src)
	}
}

func Test_ReducePath(t *testing.T) {
	tmp, err := ioutil.TempDir("", "reduce-path-")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer os.RemoveAll(tmp)

	dirs := []string{
		filepath.Join(tmp, "empty-dir"),
		filepath.Join(tmp, "contains-file"),
		filepath.Join(tmp, "nested-dir", "empty-dir"),
		filepath.Join(tmp, "nested-dir", "contains-file"),
	}

	files := []string{
		filepath.Join(tmp, "contains-file", "file.txt"),
		filepath.Join(tmp, "nested-dir", "contains-file", "file.txt"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
	}

	for _, filepath := range files {
		file, err := os.Create(filepath)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		file.WriteString("foobar")
		file.Close()
	}

	tests := map[string]bool{
		"empty-dir":     false,
		"contains-file": true,
		"nested-dir":    true,
	}

	for test, expected := range tests {
		testpath := filepath.Join(tmp, test)
		if err := ReducePath(testpath, tmp); err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if _, err := os.Stat(testpath); os.IsNotExist(err) == expected {
			t.Errorf("Expected %s exists %t", testpath, expected)
		}
	}
}

func Test_IsEmptyFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "is-empty-file-")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer os.RemoveAll(tmp)

	file, err := os.Create(filepath.Join(tmp, "somefile.txt"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	file.WriteString("foobar")
	file.Close()

	emptyFile, err := os.Create(filepath.Join(tmp, "emptyfile.txt"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	emptyFile.Close()

	tests := map[string]bool{}
	tests[file.Name()] = false
	tests[emptyFile.Name()] = true

	for path, expected := range tests {
		if actual := IsEmptyFile(path); actual != expected {
			t.Fatalf("Expected %s to be %t, not %t", path, expected, actual)
		}
	}
}
