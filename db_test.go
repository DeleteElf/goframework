package goframework

import (
	"context"
	"encoding/json"
	"github.com/deleteelf/goframework/utils/dbhelper"
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/gorm/logger"
	"testing"
)

func TestSelect(t *testing.T) {
	ctr := "host=localhost user=postgres dbname=test password=test123 port=5432 sslmode=disable"
	loghelper.GetLogManager().Init(context.Background())
	dbCon := dbhelper.CreateDb(ctr, dbhelper.Postgres, logger.Info)
	dataTable := dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err := json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("查询结果%s", string(jsonData))
}

func TestTransaction(t *testing.T) {
	ctr := "host=localhost user=postgres dbname=test password=test123 port=5432 sslmode=disable"
	loghelper.GetLogManager().Init(context.Background())

	dbCon := dbhelper.CreateDb(ctr, dbhelper.Postgres, logger.Info)
	dataTable := dbCon.QueryData(`SELECT * FROM t_user_info where f_id = ?`, 1)
	jsonData, err := json.Marshal(dataTable.Rows)
	if err != nil {
		loghelper.GetLogManager().Info("查询失败")
	}
	loghelper.GetLogManager().InfoFormat("查询结果%s", string(jsonData))

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
