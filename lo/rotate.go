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
	Filename       string
	MaxBackupsDays int

	lastWriteTime time.Time

	lastTime string
	dir      string
	file     *os.File

	mu         sync.Mutex
	TimeFormat string
	Debug      bool
}

func StartTicker(interval time.Duration, f func() bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if f() {
				return
			}
		}
	}
}

// OptionFn defines options func prototype to set options for RotateFile.
type OptionFn func(*RotateFile)

// MaxBackupsDays defines the max backups for the log files.
func MaxBackupsDays(maxBackupsDays int) OptionFn {
	return func(df *RotateFile) { df.MaxBackupsDays = maxBackupsDays }
}

// TimeFormat defines the backup file's postfix, like 20060102(yyyyMMdd) or 15:04:05 (HH:mm:ss)
func TimeFormat(timeFormat string) OptionFn {
	return func(df *RotateFile) { df.TimeFormat = timeFormat }
}

// Debug defines debugf enabled or not.
func Debug(debug bool) OptionFn { return func(df *RotateFile) { df.Debug = debug } }

// NewRotateFile create a daily rotation file
func NewRotateFile(filename string, optionFns ...OptionFn) (*RotateFile, error) {
	o := &RotateFile{
		Filename:       filename,
		MaxBackupsDays: 7, // nolint gomnd
		dir:            filepath.Dir(filename),
		TimeFormat:     `20060102`,
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

	go StartTicker(time.Second, o.tick)

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

func (o *RotateFile) debugf(format string, a ...interface{}) {
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
		o.debugf("create file %s failed, error %+v", o.Filename, err)
		return fmt.Errorf("log file %s created error %w", o.Filename, err)
	}

	o.debugf("create file %s successfully", o.Filename)

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
			o.debugf("rotate %s to %s error %w", o.Filename, rotated, err)
			return fmt.Errorf("rotate %s to %s error %w", o.Filename, rotated, err)
		}

		o.debugf("rotate %s to %s successfully", o.Filename, rotated)

		if err := o.open(); err != nil {
			return err
		}
	}

	for _, old := range outMaxBackups {
		if err := os.Remove(old); err != nil {
			o.debugf("remove log file %s before max backup days %d error %v", old, o.MaxBackupsDays, err)
			return fmt.Errorf("remove log file %s before max backup %d error %v", old, o.MaxBackupsDays, err)
		}

		o.debugf("remove log file %s before max backup days %d successfully", old, o.MaxBackupsDays)
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

func (o *RotateFile) tick() bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.file == nil {
		return true
	}

	if time.Since(o.lastWriteTime) < time.Second {
		return false
	}

	if err := o.file.Sync(); err != nil {
		o.debugf("sync file %s failed, error %+v", o.Filename, err)
	}

	return false
}

func (o *RotateFile) detectRotate(t time.Time) (rotated string, outMaxBackups []string) {
	writeTime := t.Format(o.TimeFormat)

	if o.lastTime == "" {
		o.lastTime = writeTime
		o.lastWriteTime = t
	}

	prefix := o.Filename + "."

	if o.lastTime != writeTime {
		rotated = prefix + o.lastTime
		o.lastTime = writeTime
		o.lastWriteTime = t
	}

	if o.MaxBackupsDays > 0 {
		day := t.AddDate(0, 0, -o.MaxBackupsDays)
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
