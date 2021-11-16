package for_game

/*
type IDBManager interface {
	CreateDatabase(databaseName string, db *sql.DB)
	CreateOrAlterTable(db *sql.DB)

	InitOneDB(kwargs map[string]string, databaseName ...string) *sql.DB
	InitSomeSiteDB(dsns map[string]map[string]string)
	GetDB(databaseName string) *sql.DB
}

type DBManager struct {
	Me             IDBManager
	MySQL_DB       map[string]*sql.DB
	CreateTableDDL []string
	AlterTableDDL  []string
}

func NewDBManager(createTableDDL []string, alterTableDDL []string) *DBManager {
	p := &DBManager{}
	p.Init(p, createTableDDL, alterTableDDL)
	return p
}

func (self *DBManager) Init(me IDBManager, createTableDDL []string, alterTableDDL []string) {
	self.Me = me
	self.CreateTableDDL = createTableDDL
	self.AlterTableDDL = alterTableDDL
	self.MySQL_DB = make(map[string]*sql.DB)
}

func (self *DBManager) GetDB(databaseName string) *sql.DB {
	db, ok := self.MySQL_DB[databaseName]
	if !ok {
		s := fmt.Sprintf("不存在 key 为 %v 的 MySQL 连接，也许是未正确连接", databaseName)
		panic(s)
	}
	return db
}

func (self *DBManager) InitOneDB(kwargs map[string]string, databaseName ...string) *sql.DB {
	var dsn string
	if len(databaseName) == 0 {
		dsn = `user:password@tcp(host:port)/?charset=utf8mb4` // 先不指定 database
	} else {
		dsn = "user:password@tcp(host:port)/" + databaseName[0] + "?charset=utf8"
	}

	var args []string
	for k, v := range kwargs {
		args = append(args, k)
		args = append(args, v)
	}
	s := strings.NewReplacer(args...).Replace(dsn)
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", s)
	if err != nil {
		s := "估计你没有安装 mysql 驱动.如果是,请执行 go get github.com/go-sql-driver/mysql;" + err.Error()
		e := errors.New(s)
		panic(e)
	}
	// defer db.Close() // db 长连接，无需关闭

	err = db.Ping()
	if err != nil {
		s := "MySQL 启动了吗？用户名，密码对了吗？数据库名也对吗？" + err.Error()
		e := errors.New(s)
		panic(e)
	}
	return db
}

func (self *DBManager) CreateOrAlterTable(db *sql.DB) {
	for _, ddl := range self.CreateTableDDL { // 建表
		_, e := db.Exec(ddl)
		easygo.PanicError(e)
	}
	for _, ddl := range self.AlterTableDDL { // 改表。加字段，改字段，删字段等等
		_, e := db.Exec(ddl)
		if e != nil {
			sqlErr, ok := e.(*mysql.MySQLError)
			if ok && sqlErr.Number != 1060 {
				panic(e)
			}
		}
	}
	// Error 1060: Duplicate column name 'account'  // 允许重复创建
}

func (self *DBManager) CreateDatabase(databaseName string, db *sql.DB) {
	var err error
	s := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4", databaseName)
	_, err = db.Exec(s) //建库
	easygo.PanicError(err)

	_, err = db.Exec("USE " + databaseName)
	easygo.PanicError(err)
}

// 初始化多个数据库
func (self *DBManager) InitSomeSiteDB(dsns map[string]map[string]string) {
	// 数据库名与 key 名保持一致，方便管理
	for databaseName, kwargs := range dsns {
		db := self.Me.InitOneDB(kwargs)
		self.Me.CreateDatabase(databaseName, db)
		self.Me.CreateOrAlterTable(db)

		db = self.Me.InitOneDB(kwargs, databaseName)

		self.MySQL_DB[databaseName] = db
	}
}

func (self *DBManager) ExistDatabase(databaseName string) bool {
	_, ok := self.MySQL_DB[databaseName]
	return ok
}
*/
