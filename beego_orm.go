package beego_orm

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	"github.com/u007/go_config"
  "github.com/astaxie/beego/orm"
  _ "github.com/go-sql-driver/mysql"
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

func LoadDatabase() Orm, error {
	Debug("Time zone: %s, offset: %d", local_zone, local_offset)
	
  if (err != nil) {
    err_res := fmt.Errorf("database config missing? %s", config_file)
    Error(err_res.Error())
    return nil, err_res
  }
  
  if !CheckRequired("driver", "user", "host", "encoding", "db", "pass", "connection_pool") {
    return nil, fmt.Errorf("Required configuration missing")
  }
	if time_zone == "" {
		return nil, fmt.Errorf("Required time_zone in conf/app.conf")
	}
  
	connection_string := fmt.Sprintf("%s:%s@%s/%s?charset=%s&loc=%s",
    config.String(mode, "user", ""), config.String(mode, "pass", ""),
    config.String(mode, "host", ""), config.String(mode, "db", ""), 
    config.String(mode, "encoding", ""),
    time_zone)
	
	// Debug("Connection: %s", connection_string)
  orm.RegisterDriver(config.String(mode, "driver", ""), orm.DRMySQL)
  orm.RegisterDataBase("default", config.String(mode, "driver", ""), 
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
  return orm, nil
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
