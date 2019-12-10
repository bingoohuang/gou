package enc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestTargz(t *testing.T) {
	src, err := ioutil.TempDir(os.TempDir(), "src")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("src:", src)

	defer os.RemoveAll(src)

	file, err := ioutil.TempFile(src, "file")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(file.Name())

	_ = ioutil.WriteFile(file.Name(), []byte("helloworld"), 0644)

	dest, err := ioutil.TempDir(os.TempDir(), "dest")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("dest:", dest)

	defer os.RemoveAll(dest)

	targzfile, err := ioutil.TempFile(dest, "targzfile")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("targzfile:", targzfile.Name())

	defer os.Remove(targzfile.Name())

	if err := Targz(src, false, targzfile); err != nil {
		t.Fatal(err)
	}

	targzfile.Close()

	dir, err := ioutil.TempDir(os.TempDir(), "untargz")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("dir:", dir)

	defer os.RemoveAll(dir)

	targzfile, _ = os.Open(targzfile.Name())
	if err := Untargz(dir, targzfile); err != nil {
		t.Fatal(err)
	}

	if ok, _ := EqualsDir(src, dir); !ok {
		t.Fatal("not equals")
	}
}

func EqualsDir(dir1, dir2 string) (bool, error) {
	if stat1, err := os.Stat(dir1); err != nil {
		return false, err
	} else if !stat1.IsDir() {
		return false, fmt.Errorf("%s is not a dir", dir1)
	}

	if stat2, err := os.Stat(dir2); err != nil {
		return false, err
	} else if !stat2.IsDir() {
		return false, fmt.Errorf("%s is not a dir", dir2)
	}

	return equalsDir(dir1, dir2)
}

func equalsDir(d1, d2 string) (bool, error) {
	d1Walker := &equalsDirWalker{dir: d1, files: make(map[string]bool), subs: make(map[string]bool)}
	if err := filepath.Walk(d1, d1Walker.walk); err != nil {
		return false, err
	}

	d2Walker := &equalsDirWalker{dir: d2, files: make(map[string]bool), subs: make(map[string]bool)}
	if err := filepath.Walk(d2, d2Walker.walk); err != nil {
		return false, err
	}

	if !reflect.DeepEqual(d1Walker.files, d2Walker.files) {
		return false, nil
	}

	if !reflect.DeepEqual(d1Walker.subs, d2Walker.subs) {
		return false, nil
	}

	for f := range d1Walker.files {
		if !EqualFile(filepath.Join(d1, f), filepath.Join(d2, f)) {
			return false, nil
		}
	}

	for sub := range d1Walker.subs {
		if ok, err := equalsDir(filepath.Join(d1, sub), filepath.Join(d2, sub)); !ok {
			return ok, err
		}
	}

	return true, nil
}

func EqualFile(file1, file2 string) bool {
	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f1, err1 := ioutil.ReadFile(file1)
	if err1 != nil {
		return false
	}

	f2, err2 := ioutil.ReadFile(file2)
	if err2 != nil {
		return false
	}

	return bytes.Equal(f1, f2)
}

type equalsDirWalker struct {
	dir   string
	files map[string]bool
	subs  map[string]bool
}

func (w *equalsDirWalker) walk(file string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	name := strings.TrimPrefix(strings.Replace(file, w.dir, "", -1), string(filepath.Separator))
	if name == "" {
		return err
	}

	if fi.IsDir() {
		w.subs[name] = true
	} else {
		w.files[name] = true
	}

	return nil
}
