package lo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const yyyyMMdd = "yyyy-MM-dd"

// DailyFile is a daily rotate file
type DailyFile struct {
	Filename   string
	MaxBackups int

	lastDay string
	dir     string
	file    *os.File

	mu *sync.Mutex
}

// NewDailyFile create a daily rotation file
func NewDailyFile(filename string, maxBackups int) (*DailyFile, error) {
	o := &DailyFile{
		Filename:   filename,
		MaxBackups: maxBackups,
		dir:        filepath.Dir(filename),
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
func (o *DailyFile) Write(d []byte) (int, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.write(d, false)
}

// Flush flushes the file
func (o *DailyFile) Flush() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.file.Sync()
}

// Close closes the file
func (o *DailyFile) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.close()
}

func (o *DailyFile) open() error {
	f, err := os.OpenFile(o.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("log file %s created error %w", o.Filename, err)
	}

	o.file = f

	logrus.Infof("log file %s created", o.Filename)

	return nil
}

func (o *DailyFile) rotateFiles(t time.Time) error {
	rotated, outMaxBackups := o.detectRotate(t)

	return o.doRotate(rotated, outMaxBackups)
}

func (o *DailyFile) doRotate(rotated string, outMaxBackups []string) error {
	if rotated != "" {
		if err := o.close(); err != nil {
			return err
		}

		if err := os.Rename(o.Filename, rotated); err != nil {
			return fmt.Errorf("rotate %s to %s error %w", o.Filename, rotated, err)
		}

		logrus.Infof("%s rotated to %s", o.Filename, rotated)

		if err := o.open(); err != nil {
			return err
		}
	}

	for _, old := range outMaxBackups {
		if err := os.Remove(old); err != nil {
			return fmt.Errorf("remove log file %s before max backup days %d error %v", old, o.MaxBackups, err)
		}

		logrus.Infof("%s before max backup days %d removed", old, o.MaxBackups)
	}

	return nil
}

func (o *DailyFile) close() error {
	if o.file == nil {
		return nil
	}

	err := o.file.Close()
	o.file = nil

	return err
}

func (o *DailyFile) detectRotate(t time.Time) (rotated string, outMaxBackups []string) {
	ts := FormatTime(t, yyyyMMdd)

	if o.lastDay == "" {
		o.lastDay = ts
	}

	if o.lastDay != ts {
		o.lastDay = ts

		yesterday := t.AddDate(0, 0, -1)
		rotated = o.Filename + "." + FormatTime(yesterday, yyyyMMdd)
	}

	if o.MaxBackups > 0 {
		day := t.AddDate(0, 0, -o.MaxBackups)
		_ = filepath.Walk(o.dir, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() {
				return err
			}

			if strings.HasPrefix(path, o.Filename+".") {
				fis := path[len(o.Filename+"."):]
				if backDay, err := ParseTime(fis, yyyyMMdd); err != nil {
					return nil // ignore this file
				} else if backDay.Before(day) {
					outMaxBackups = append(outMaxBackups, path)
				}
			}

			return nil
		})
	}

	return rotated, outMaxBackups
}

func (o *DailyFile) write(d []byte, flush bool) (int, error) {
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

// ConvertTimeLayout converts date time format in java style to go style
func ConvertTimeLayout(layout string) string {
	l := layout
	l = strings.Replace(l, "yyyy", "2006", -1)
	l = strings.Replace(l, "yy", "06", -1)
	l = strings.Replace(l, "MM", "01", -1)
	l = strings.Replace(l, "dd", "02", -1)
	l = strings.Replace(l, "HH", "15", -1)
	l = strings.Replace(l, "mm", "04", -1)
	l = strings.Replace(l, "ss", "05", -1)
	l = strings.Replace(l, "SSS", "000", -1)

	return l
}

// ParseTime 解析日期转字符串
func ParseTime(d string, layout string) (time.Time, error) {
	return time.Parse(ConvertTimeLayout(layout), d)
}

// FormatTime 日期转字符串
func FormatTime(d time.Time, layout string) string {
	return d.Format(ConvertTimeLayout(layout))
}
