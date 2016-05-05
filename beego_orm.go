package beego_orm

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	"github.com/u007/go_config"
  "github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
  _ "github.com/go-sql-driver/mysql"
	"net/url"
	"fmt"
	"time"
)

var config_file = "conf/database.conf"
var config, err      = go_config.NewConfigLoader("ini", config_file)
var mode        = beego.AppConfig.DefaultString("runmode", "dev")
var local_zone, local_offset = GetTimeZone()
var time_zone   = beego.AppConfig.DefaultString("time_zone", local_zone)

func GetTimeZone() (name string, offset int) {
	return time.Now().In(time.Local).Zone()
}

func LoadDatabase() error {
	Debug("Time zone: %s, offset: %d", local_zone, local_offset)
	
  if (err != nil) {
    return err
  }
  
  if !CheckRequired("driver", "user", "host", "encoding", "db", "pass", "connection_pool") {
    return fmt.Errorf("Required configuration missing")
  }
	if time_zone == "" {
		return fmt.Errorf("Required time_zone in conf/app.conf")
	}
  port := config.Int(mode, "port", 3306)
	connection_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=%s",
    config.String(mode, "user", ""), config.String(mode, "pass", ""),
    config.String(mode, "host", ""), port, config.String(mode, "db", ""), 
    config.String(mode, "encoding", ""),
    url.QueryEscape(time_zone))
	
	driver := config.String(mode, "driver", "")
	Debug("Driver: %s, Connection: %s", driver, connection_string)
	if (driver == "mysql") {
		orm.RegisterDriver(config.String(mode, "driver", ""), orm.DRMySQL)
	}
  orm.RegisterDataBase("default", driver, 
    connection_string,
    5, config.Int(mode, "connection_pool", 0))
  
	
	// files, _ := ioutil.ReadDir("models")
  // for _, f := range files {
	// 	Debug("model: %s", f.Name())
  // }
		
  // orm.DefaultTimeLoc = Time.UTC
  // orm.RegisterModel(new(models.User))
  // orm.RegisterModelWithPrefix("prefix_", new(User))
  if (config.Boolean(mode, "debug", false)) {
    orm.Debug = true
  }
  return nil
}

func LogValidationErrors(log_prefix string, valid *validation.Validation) {
	if valid.HasErrors() {
    for _, err := range valid.Errors {
			Error("[ %s ]Validation %s: %s", log_prefix, err.Key, err.Message)
    }
  }
}

func CheckRequired(args ...string) bool {
  for _, name := range args {
    if (config.String(mode, name, "") == "") {
      Error("%s required in %s", name, config_file)
      return false
    }
  }
  return true
}
