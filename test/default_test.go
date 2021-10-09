package test

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/pangxieke/logServer/routers"
	"path/filepath"
)

func init() {

	file := "../"

	fmt.Println(file)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))

	beego.TestBeegoInit(appPath)
}
