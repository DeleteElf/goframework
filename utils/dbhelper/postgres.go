package dbhelper

import (
	"github.com/deleteelf/goframework/utils/loghelper"
	"github.com/deleteelf/goframework/utils/stringhelper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PostgresDB struct {
	DbBase
}

// 默认会自动打开和关闭连接，如果启用了事务，则在提交或回滚事务之前，不会关闭连接
func NewPostgresDB(config DbConfig) *PostgresDB {
	return &PostgresDB{
		DbBase{
			Config: config,
		},
	}
}

func (pg *PostgresDB) Open() bool {
	var err error
	//conf := MyNamingStrategy{
	//	ColumnPrefix: "f_",
	//}
	//conf.TablePrefix = "t_"
	//conf.IdentifierMaxLength = 64
	//conf.SingularTable = false
	pg.db, err = gorm.Open(postgres.Open(pg.Config.ConnectionString), &gorm.Config{
		SkipDefaultTransaction: pg.Config.SkipDefaultTransaction,
		Logger:                 logger.Default.LogMode(pg.Config.LogLevel),
		NamingStrategy:         schema.NamingStrategy{},
	})
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

func (pg *PostgresDB) BeginTransaction() bool {
	return false
}

// 提交事务
func (pg *PostgresDB) CommitTransaction() bool {
	return false
}

// 回滚事务
func (pg *PostgresDB) RollbackTransaction() bool {
	return false
}

// 是否处于事务中
func (pg *PostgresDB) IsInTransaction() bool {
	return false
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

func (pg *PostgresDB) QueryData(sql string, conds ...any) *DataTable {
	result := new(DataTable)
	if pg.Open() {
		defer pg.Close()
		pg.db.Raw(sql, conds...)
		if pg.Config.SafeColumn {
			rows, err := pg.db.Rows()
			defer rows.Close()
			if err != nil {
				loghelper.GetLogManager().Error("获取行数据出错！！")
			}
			columns, err1 := rows.Columns()
			if err1 != nil {
				loghelper.GetLogManager().Error("获取列数据出错！！")
			}
			pg.db.Statement.ColumnMapping = map[string]string{}
			for _, column := range columns {
				pg.db.Statement.ColumnMapping[stringhelper.ConvertToCamel(column)] = column
			}
		}
		err := pg.db.Scan(&result.Rows).Error
		switch err {
		case gorm.ErrRecordNotFound:
			loghelper.GetLogManager().Error("查询的数据不存在！！！")
			break
		case gorm.ErrDryRunModeUnsupported:
			loghelper.GetLogManager().Error("ErrDryRunModeUnsupported！！！")
			break
		default:
			break
		}
		//if pg.Config.SafeColumn {
		//	pg.db.Raw(sql, conds...).
		//}
	}
	return result
}
