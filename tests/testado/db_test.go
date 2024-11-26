package testado

import (
	"encoding/json"
	"errors"
	"github.com/deleteelf/goframework/ado"
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	ctr := "host=localhost user=postgres dbname=test password=test123 port=5432 sslmode=disable"
	loghelper.GetLogManager().Init(loghelper.Debug)
	dbCon := ado.CreateDb(ctr, ado.Postgres, logger.Info)
	dataTable := dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err := json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("查询结果%s", string(jsonData))
}

func TestTransaction(t *testing.T) {
	ctr := "host=localhost user=postgres dbname=test password=test123 port=5432 sslmode=disable"
	loghelper.GetLogManager().Init(loghelper.Debug)

	dbCon := ado.CreateDb(ctr, ado.Postgres, logger.Info)
	dataTable := dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err := json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("查询结果 %s ", string(jsonData))
	//loghelper.GetLogManager().Info(`查询结果 [{"f_account":"admin","f_active":true,"f_create_time":"2024-11-12T17:45:05.188292+08:00","f_email":null,"f_id":1,"f_modify_time":"2024-11-26T14:34:53.124399+08:00","f_name":"系统管理员","f_password":"0192023a7bbd73250516f069df18b500","f_remark":"密码是 admin123","f_telephone":null}] `)

	dbCon.Open()
	dbCon.BeginTransaction()

	dataTable.TableName = "t_user_info"
	str := ""
	dataTable.ColumnPrefix = &str
	dataTable.Rows[0]["f_remark"] = "test"
	_, err1 := dbCon.SaveData(dataTable)
	if err1 != nil {
		loghelper.GetLogManager().Error("修改数据发生错误")
		dbCon.RollbackTransaction()
		dbCon.Close()
		return
	}
	dataTable = dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err = json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("修改数据后，查询结果%s", string(jsonData))
	dbCon.RollbackTransaction()
	dbCon.Close()
	dataTable = dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err = json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("回滚事务后，查询结果%s", string(jsonData))

}

func TestTransactionCallback(t *testing.T) {
	ctr := "host=localhost user=postgres dbname=test password=test123 port=5432 sslmode=disable"
	loghelper.GetLogManager().Init(loghelper.Debug)

	dbCon := ado.CreateDb(ctr, ado.Postgres, logger.Info)
	dataTable := dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err := json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("查询结果%s", string(jsonData))

	dbCon.TransactionCallback(func(tx *gorm.DB) error {
		dataTable.TableName = "t_user_info"
		str := ""
		dataTable.ColumnPrefix = &str
		dataTable.Rows[0]["f_remark"] = "test"
		_, err1 := dbCon.SaveData(dataTable)
		if err1 != nil {
			loghelper.GetLogManager().Error("修改数据发生错误")
			return err
		}
		dataTable = dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
		jsonData, err = json.Marshal(dataTable.Rows)
		if err != nil {

			loghelper.GetLogManager().Info("查询失败")
			return err
		}
		//log.Printf("修改数据后，查询结果%s", string(jsonData))
		loghelper.GetLogManager().InfoFormat("修改数据后，查询结果%s", string(jsonData))
		return errors.New("自定义错误")
	})

	dataTable = dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err = json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	//loghelper.GetLogManager().Info("test")
	//log.Printf("回滚事务后，查询结果%s", string(jsonData))
	loghelper.GetLogManager().InfoFormat("回滚事务后，查询结果%s", string(jsonData))
	//loghelper.GetLogManager().Info("test")
	time.Sleep(1)
}
