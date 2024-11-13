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

func (pg *PostgresDB) AutoMigrate(model ModelInterface) error {
	if pg.Open() {
		defer pg.Close()
		return pg.db.AutoMigrate(model)
	}
	return nil
}
func (pg *PostgresDB) Save(model ModelInterface) {
	if pg.Open() {
		defer pg.Close()
		pg.db.Save(model)
	}
}

func (pg *PostgresDB) SelectById(model ModelInterface, id any) {
	//反射的案例，不过gorm已经做好反射了
	//t := reflect.TypeFor[T1]()
	//val := reflect.New(t).Elem()
	//result := val.Interface().(T1)
	if pg.Open() {
		defer pg.Close()
		err := pg.db.First(model, id).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("根据Id查询的数据不存在！！！")
			break
		default:
			break
		}

	}
}

// 根据条件查询数据，dest传入数组指针
func (pg *PostgresDB) SelectByCondition(dest interface{}, query string, conds ...any) {
	if pg.Open() {
		defer pg.Close()
		err := pg.db.Where(query, conds...).Find(dest).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		default:
			break
		}
	}
}
