package log

import (
	"time"
	"os"
	"fmt"
	"strings"
	"log"
)

const (
	DefaultMaxSize          = 100
	DefaultBackup           = 5
)

type Options struct {
	MaxSize         int64
	MaxBackup       int
}

type AdapterFile struct{
	fileName        string

	options         Options
	currDate        string
	f               *os.File
	size            int64
}

func NewAdapterFile(logName string, opt *Options) *AdapterFile {
	a := &AdapterFile{
		fileName: logName,
		options: Options{
			MaxSize: DefaultMaxSize,
			MaxBackup: DefaultBackup,
		},
		size: 0,
	}
	if opt != nil {
		if opt.MaxBackup > 0 {
			a.options.MaxBackup = opt.MaxBackup
		}
		if opt.MaxSize > 0 {
			a.options.MaxSize = opt.MaxSize
		}
	}
	a.currDate = makeCurrDate()
	if err := a.rotateBySize(); err != nil {
		return err
	}

	go a.dailyRotate()
	return a
}

func (af *AdapterFile) Write(b []byte) error {
	var err error
	date := makeCurrDate()
	if date != af.currDate {
		//切换日期
		err = af.rotateByDate(date)
		if err != nil {
			return err
		}
	} else {
		if af.size + len(b) >= af.options.MaxSize {
			//切换序号
			err = af.rotateBySize()
			if err != nil {
				return err
			}
		}
	}

	n, err := af.f.Write(b)
	if err != nil {
		log.Printf("write log file error: %+v, %v", af.f, err)
		return err
	}
	af.size += n
	return nil
}

func (af *AdapterFile) dailyRotate() {
	//自动日期轮转
}

func (af *AdapterFile) rotateBySize() error {
	if af.f != nil {
		af.f.Close()
		af.f = nil
	}
	for i := 0; i < af.options.MaxBackup; i++ {
		fname := af.makeLogName(i)
		st, err := os.Stat(fname)
		if (err != nil && os.IsNotExist(err)) || (err == nil && st.Size() < af.options.MaxSize) {
			//找到可以写的日志文件
			return af.openLogFile(fname)
		}
	}
	//轮转已满，淘汰最老的日志，永远写最后一个
	for i := 1; i < af.options.MaxBackup; i++ {
		fname := af.makeLogName(i)
		preName := af.makeLogName(i - 1)
		err := os.Rename(fname, preName)
		if err != nil {
			log.Printf("Rename error: %v->%v, err %v", fname, preName, err)
		}
	}
	return af.openLogFile(af.makeLogName(af.options.MaxBackup - 1))
}

func (af *AdapterFile) rotateByDate(date string) error {
	if af.f != nil {
		af.f.Close()
		af.f = nil
	}
	af.currDate = date

	newName := af.makeLogName(0)
	return af.openLogFile(newName)
}

func (af *AdapterFile) makeLogName(backup int) string {
	return strings.Join([]string{af.fileName, af.currDate, fmt.Sprintf("%02d", backup), "log"}, ".")
}

func (af *AdapterFile) openLogFile(fname string) error {
	f, err := os.OpenFile(fname, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		log.Printf("openLogFile error: file %v, err %v", fname, err)
		return
	}
	st, err := f.Stat()
	if err != nil {
		return err
	}
	af.size = st.Size()
	af.f = f

	return nil
}

func makeCurrDate() string {
	_, m, d := time.Now().Date()
	return fmt.Sprint("%02d%02d", int(m), d)
}
