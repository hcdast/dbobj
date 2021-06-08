/*
 * @Author: hc
 * @Date: 2021-06-07 10:55:03
 * @LastEditors: hc
 * @LastEditTime: 2021-06-07 10:56:58
 * @Description:
 */
package dbobj

import (
	"example-hauth/dbobj/dbhandle"

	_ "example-hauth/dbobj/mysql"
)

func init() {
	conf, err := dbhandle.GetConfig()
	if err != nil {
		panic("init database failed." + err.Error())
	}
	Default, err = conf.Get("DB.type")
	if err != nil {
		panic("get default database type failed." + err.Error())
	}
	InitDB(Default)
}
