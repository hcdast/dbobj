/*
 * @Author: hc
 * @Date: 2021-06-07 10:55:03
 * @LastEditors: hc
 * @LastEditTime: 2021-06-10 16:10:54
 * @Description:
 */
package dbobj

import (
	_ "example-hauth/dbobj/mysql"

	"github.com/astaxie/beego"
)

// 连接mysql数据库
func init() {
	// 方案1 解析配置文件
	// conf, err := dbhandle.GetConfig()
	// if err != nil {
	// 	panic("init database failed." + err.Error())
	// }
	// Default, err = conf.Get("DB.type")
	// if err != nil {
	// 	panic("get default database type failed." + err.Error())
	// }
	// 方案2 读取配置文件
	Default := beego.AppConfig.String("DB.type")
	InitDB(Default)
}
