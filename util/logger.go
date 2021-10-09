package util

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
)

func InitLogger(cfg config.Configer) {
	logpath := cfg.String("log::path")
	MakeSureFileExist(logpath)

	level := cfg.DefaultInt("log::level", 7)
	if level > logs.LevelDebug {
		level = logs.LevelDebug
	}
	if level < logs.LevelEmergency {
		level = logs.LevelDebug
	}
	maxdays := cfg.DefaultInt("log::maxdays", 30)
	if maxdays <= 0 {
		maxdays = 30
	}

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logconfig := fmt.Sprintf(`{"filename":"%s","level":%d,"perm":"0644","maxdays":%d}`,
		logpath, level, maxdays)
	logs.SetLogger(logs.AdapterFile, logconfig)
}

func MakeSureFileExist(filepath string) {
	dir := path.Dir(filepath)

	var err error

	_, err = os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err.Error())
		}

		if err = os.MkdirAll(dir, os.ModeDir|0755); err != nil {
			panic(err.Error())
		}
	}

	_, err = os.Stat(filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err.Error())
		}

		file, err := os.Create(filepath)
		if err != nil {
			panic(err.Error())
		}

		file.Close()
	}
}
