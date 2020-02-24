package lo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RotateFile is a daily rotate file
type RotateFile struct {
	Filename   string
	MaxBackups int

	lastDay string
	dir     string
	file    *os.File

	mu         sync.Mutex
	TimeFormat string
	Debug      bool
}

// OptionFn defines options func prototype to set options for RotateFile.
type OptionFn func(*RotateFile)

// MaxBackups defines the max backups for the log files.
func MaxBackups(maxBackups int) OptionFn {
	return func(df *RotateFile) { df.MaxBackups = maxBackups }
}

// TimeFormat defines the backup file's postfix, like 20060102(yyyyMMdd) or 15:04:05 (HH:mm:ss)
func TimeFormat(timeFormat string) OptionFn {
	return func(df *RotateFile) { df.TimeFormat = timeFormat }
}

// Debug defines debug enabled or not.
func Debug(debug bool) OptionFn {
	return func(df *RotateFile) { df.Debug = debug }
}

// NewRotateFile create a daily rotation file
func NewRotateFile(filename string, optionFns ...OptionFn) (*RotateFile, error) {
	o := &RotateFile{
		Filename:   filename,
		MaxBackups: 7, // nolint gomnd
		dir:        filepath.Dir(filename),
		TimeFormat: `20060102`,
	}

	for _, fn := range optionFns {
		fn(o)
	}

	if err := os.MkdirAll(o.dir, 0755); err != nil {
		return nil, err
	}

	if err := o.open(); err != nil {
		return nil, err
	}

	return o, nil
}

// Write writes data to a file
func (o *RotateFile) Write(d []byte) (int, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.write(d, false)
}

// Flush flushes the file
func (o *RotateFile) Flush() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.file.Sync()
}

// Close closes the file
func (o *RotateFile) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.close()
}

func (o *RotateFile) debug(format string, a ...interface{}) {
	if !o.Debug {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func (o *RotateFile) open() error {
	f, err := os.OpenFile(o.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		o.debug("create file %s failed, error %+v", o.Filename, err)
		return fmt.Errorf("log file %s created error %w", o.Filename, err)
	}

	o.debug("create file %s successfully", o.Filename)

	o.file = f

	return nil
}

func (o *RotateFile) rotateFiles(t time.Time) error {
	rotated, outMaxBackups := o.detectRotate(t)

	return o.doRotate(rotated, outMaxBackups)
}

func (o *RotateFile) doRotate(rotated string, outMaxBackups []string) error {
	if rotated != "" {
		if err := o.close(); err != nil {
			return err
		}

		if err := os.Rename(o.Filename, rotated); err != nil {
			o.debug("rotate %s to %s error %w", o.Filename, rotated, err)
			return fmt.Errorf("rotate %s to %s error %w", o.Filename, rotated, err)
		}

		o.debug("rotate %s to %s successfully", o.Filename, rotated)

		if err := o.open(); err != nil {
			return err
		}
	}

	for _, old := range outMaxBackups {
		if err := os.Remove(old); err != nil {
			o.debug("remove log file %s before max backup days %d error %v", old, o.MaxBackups, err)
			return fmt.Errorf("remove log file %s before max backup %d error %v", old, o.MaxBackups, err)
		}

		o.debug("remove log file %s before max backup days %d successfully", old, o.MaxBackups)
	}

	return nil
}

func (o *RotateFile) close() error {
	if o.file == nil {
		return nil
	}

	err := o.file.Close()
	o.file = nil

	return err
}

func (o *RotateFile) detectRotate(t time.Time) (rotated string, outMaxBackups []string) {
	day := t.Format(o.TimeFormat)

	if o.lastDay == "" {
		o.lastDay = day
	}

	prefix := o.Filename + "."

	if o.lastDay != day {
		o.lastDay = day

		yesterday := t.AddDate(0, 0, -1)
		rotated = prefix + yesterday.Format(o.TimeFormat)
	}

	if o.MaxBackups > 0 {
		day := t.AddDate(0, 0, -o.MaxBackups)
		_ = filepath.Walk(o.dir, func(p string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() || !strings.HasPrefix(p, prefix) {
				return nil
			}

			fd, err := time.Parse(o.TimeFormat, p[len(prefix):])
			if err != nil {
				return nil // ignore this file
			}

			if fd.Before(day) {
				outMaxBackups = append(outMaxBackups, p)
			}

			return nil
		})
	}

	return rotated, outMaxBackups
}

func (o *RotateFile) write(d []byte, flush bool) (int, error) {
	if err := o.rotateFiles(time.Now()); err != nil {
		return 0, err
	}

	n, err := o.file.Write(d)
	if err != nil {
		return n, err
	}

	if flush {
		err = o.file.Sync()
	}

	return n, err
}
