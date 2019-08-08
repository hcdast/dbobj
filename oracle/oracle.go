package oracle

import (
	"database/sql"
	"github.com/hzwy23/panda/logger"
	"os"
	"strconv"
	"strings"

	"github.com/hzwy23/dbobj/dbhandle"
	_ "github.com/mattn/go-oci8"
)

type oracle struct {
	db *sql.DB
}

func NewOracle() dbhandle.DbObj {

	var err error

	o := new(oracle)

	nlsLang := os.Getenv("NLS_LANG")
	if !strings.HasSuffix(nlsLang, "UTF8") {
		os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
	}
	os.Setenv("NLS_DATE_FORMAT", "YYYY-MM-DD HH24:MI:SS")

	red, err := dbhandle.GetConfig()
	if err != nil {
		panic("cant not read ./conf/dbobj.conf.please check this file.")
	}

	tns := red.Conf["DB.tns"]
	usr := red.Conf["DB.user"]
	pad := red.Conf["DB.passwd"]
	mc := red.Conf["DB.maxConn"]
	maxConn := 100
	if len(mc) != 0 {
		mx, err := strconv.Atoi(mc)
		if err == nil {
			maxConn = mx
		}
	}

	if len(pad) == 24 {
		pad, err = utils.Decrypt(pad)
		if err != nil {
			logger.Error("Decrypt mysql passwd failed.")
			return nil
		}
	}

	tnsname := usr + "/" + pad + "@" + tns

	o.db, err = sql.Open("oci8", tnsname)

	if err != nil {
		logger.Error("open oracle database failed.%v" + err)
		return nil
	}
	if len(pad) != 24 {
		psd, err := utils.Encrypt(pad)
		if err != nil {
			logger.Error("decrypt passwd failed.%v" + psd)
			return nil
		}
		psd = "\"" + psd + "\""
		red.Set("DB.passwd", psd)
	}

	// 设置连接池最大值
	o.db.SetMaxOpenConns(maxConn)
	o.db.SetConnMaxLifetime(0)
	logger.Info("create Oracle connect pool success. max connect value is:", maxConn)
	return o
}

func (this *oracle) GetErrorCode(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[:n])
	} else {
		logger.Error("this error information is not Oracle return info")
		return ""
	}
}

func (this *oracle) GetErrorMsg(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[n+1:])
	} else {
		logger.Error("this error information is not Oracle return info")
		return ""
	}
}

func (this *oracle) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewOracle().(*oracle); ok {
				this.db = val.db
			}
			return this.db.Query(sql, args...)
		}
	}

	return rows, err
}

func (this *oracle) Exec(sql string, args ...interface{}) (sql.Result, error) {
	result, err := this.db.Exec(sql, args...)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewOracle().(*oracle); ok {
				this.db = val.db
			}
			return this.db.Exec(sql, args...)
		}
	}
	return result, err
}

func (this *oracle) Begin() (*sql.Tx, error) {
	tx, err := this.db.Begin()
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewOracle().(*oracle); ok {
				this.db = val.db
			}
			return this.db.Begin()
		}
	}
	return tx, err
}

func (this *oracle) Prepare(sql string) (*sql.Stmt, error) {
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		if this.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := NewOracle().(*oracle); ok {
				this.db = val.db
			}
			return this.db.Prepare(sql)
		}
	}
	return stmt, err
}

func (this *oracle) QueryRow(sql string, args ...interface{}) *sql.Row {
	if this.db.Ping() != nil {
		logger.Warn("Connection is broken")
		if val, ok := NewOracle().(*oracle); ok {
			this.db = val.db
		}
	}
	return this.db.QueryRow(sql, args...)
}

func init() {
	dbhandle.Register("oracle", NewOracle)
}
