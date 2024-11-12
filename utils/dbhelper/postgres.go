package dbhelper

import (
	"fmt"
	"github.com/deleteelf/goframework/utils/loghelper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DbBase
}

func NewPostgresDB(config DbConfig) *PostgresDB {
	return &PostgresDB{
		DbBase{
			Config: config,
		},
	}
}

func (pg *PostgresDB) Open() bool {
	var err error
	pg.db, err = gorm.Open(postgres.Open(pg.Config.ConnectionString), &pg.Config)
	if err != nil {
		loghelper.GetLogManager().Error("数据库连接失败！！%s", pg.Config.ConnectionString)
		return false
	}
	loghelper.GetLogManager().Info("数据库连接成功！！%s", pg.Config.ConnectionString)
	return true
}

func (pg *PostgresDB) Close() bool {
	sqlDb, err := pg.db.DB()
	if err != nil {
		loghelper.GetLogManager().ErrorV(err)
		return false
	}
	err = sqlDb.Close()
	if err != nil {
		loghelper.GetLogManager().ErrorV(err)
		return false
	}
	return true
}

func (pg *PostgresDB) SelectById(bean ModelInterface, id any) Model {
	//反射的案例，不过gorm已经做好反射了
	//t := reflect.TypeFor[T1]()
	//val := reflect.New(t).Elem()
	//result := val.Interface().(T1)
	if pg.Open() {
		err := pg.db.First(bean, id).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("根据Id查询的数据不存在！！！")
			break
		default:
			break
		}
		defer pg.Close()
	}
	return bean.(Model)
	//return result
}

func (pg *PostgresDB) SelectByCondition(bean ModelInterface, conds ...any) Model {
	if pg.Open() {
		err := pg.db.Take(bean, conds).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		default:
			break
		}
		defer pg.Close()
	}
	return bean.(Model)
}

func (pg *PostgresDB) test() {
	var user UserInfo
	pg.SelectById(user, 1)
	fmt.Printf("%s", user.Id)
}
