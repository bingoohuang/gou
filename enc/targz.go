package enc

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Targz takes a source path (srcPath) and  writer and walks 'srcPath' writing each file
// found to the tar writer; the purpose for accepting multiple writers is to allow
// for multiple outputs (for example a file, or md5 hash)
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Targz(srcPath string, baseDir bool, writer io.Writer) error {
	// ensure the srcPath actually exists before trying to tar it
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	w := walker{srcPath: srcPath, tw: tw, baseDir: baseDir, base: filepath.Base(srcPath)}
	return filepath.Walk(srcPath, w.walk)
}

type walker struct {
	srcPath string
	tw      *tar.Writer
	baseDir bool
	base    string
}

func (w walker) walk(file string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	name := strings.TrimPrefix(strings.Replace(file, w.srcPath, "", -1), string(filepath.Separator))
	if name == "" && !w.baseDir {
		return err
	}

	if w.baseDir {
		name = filepath.Join(w.base, name)
	}

	// create a new dir/file header
	header, err := tar.FileInfoHeader(fi, fi.Name())
	if err != nil {
		return err
	}

	// update the name to correctly reflect the desired destination when untaring
	header.Name = name
	if err := w.tw.WriteHeader(header); err != nil {
		return err
	}

	// return on non-regular files
	// thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update
	if !fi.Mode().IsRegular() {
		return nil
	}

	// open files for taring
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w.tw, f); err != nil {
		return err
	}

	// manually close here after each file operation; deferring would cause each file close
	// to wait until all operations have completed.
	return f.Close()
}

// Untargz takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untargz(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			_ = f.Close()
		}
	}
}
