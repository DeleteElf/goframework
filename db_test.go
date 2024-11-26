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
	loghelper.GetLogManager().InfoFormat("查询结果[%s]", jsonData)
}
