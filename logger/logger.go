package logger

import (
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func Debug(s string, flag bool) {
	if flag == false {
		return
	}
	c := color.New(color.FgGreen)
	ts := time.Now().Format("2006/01/02 15:04:05 MST")
	_, filename, line_no, _ := runtime.Caller(1)
	c.Printf("[DEBUG][%s][%s:%d] %s\n",
		ts,
		filepath.Base(filename),
		line_no,
		s,
	)
}

func Crit(s string) {
	c := color.New(color.FgGreen)
	ts := time.Now().Format("2006/01/02 15:04:05 MST")
	_, filename, line_no, _ := runtime.Caller(1)
	c.Printf("[CRITICAL][%s][%s:%d] %s\n",
		ts,
		filepath.Base(filename),
		line_no,
		s,
	)

	os.Exit(1)
}
