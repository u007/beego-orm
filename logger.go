package beego_orm

import (
	"github.com/astaxie/beego"
	"fmt"
	"runtime/debug"
)

const PREFIX = "[ ORM ] "

func Debug(format string, v... interface{}) {
	beego.Debug(fmt.Sprintf(PREFIX + format, v...))
}
func Warning(format string, v... interface{}) {
	beego.Warning(fmt.Sprintf(PREFIX + format, v...))
}
func Error(format string, v... interface{}) {
	beego.Error(fmt.Sprintf(PREFIX + format, v...))
	debug.PrintStack()
}
