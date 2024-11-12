package dbhelper

import (
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
		loghelper.GetLogManager().ErrorFormat("数据库连接失败！！%s", pg.Config.ConnectionString)
		return false
	}
	loghelper.GetLogManager().InfoFormat("数据库连接成功！！%s", pg.Config.ConnectionString)
	return true
}

func (pg *PostgresDB) Close() bool {
	sqlDb, err := pg.db.DB()
	if err != nil {
		loghelper.GetLogManager().Error(err)
		return false
	}
	err = sqlDb.Close()
	if err != nil {
		loghelper.GetLogManager().Error(err)
		return false
	}
	return true
}

func (pg *PostgresDB) AutoMigrate(bean ModelInterface) error {
	if pg.Open() {
		defer pg.Close()
		return pg.db.AutoMigrate(bean)
	}
	return nil
}
func (pg *PostgresDB) Save(bean ModelInterface) {
	if pg.Open() {
		defer pg.Close()
		pg.db.Save(bean)
	}
}

func (pg *PostgresDB) SelectById(bean ModelInterface, id any) {
	//反射的案例，不过gorm已经做好反射了
	//t := reflect.TypeFor[T1]()
	//val := reflect.New(t).Elem()
	//result := val.Interface().(T1)
	if pg.Open() {
		defer pg.Close()
		err := pg.db.First(bean, id).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("根据Id查询的数据不存在！！！")
			break
		default:
			break
		}

	}
}

func (pg *PostgresDB) SelectByCondition(bean ModelInterface, conds ...any) {
	if pg.Open() {
		defer pg.Close()
		err := pg.db.Take(bean, conds).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		default:
			break
		}
	}
}
