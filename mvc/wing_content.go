// Copyright (c) 2019-2029 DY All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"reflect"
	"strings"
	// ----------------------------------------
	// NOTIC :
	//
	// import the follows database drivers when using WingProvider.
	//
	// _ "github.com/go-sql-driver/mysql"   // usr fot mysql
	// _ "github.com/denisenkom/go-mssqldb" // use for sql server 2017 ~ 2019
	//
	// ----------------------------------------
)

// WingProvider content provider to support database utils
type WingProvider struct {
	Conn *sql.DB
}

// ScanCallback use for scan query result from rows
type ScanCallback func(rows *sql.Rows) error

const (
	/* MySQL */
	mysqlConfigUser = "%s::user" // configs key of mysql database user
	mysqlConfigPwd  = "%s::pwd"  // configs key of mysql database password
	mysqlConfigHost = "%s::host" // configs key of mysql database host and port
	mysqlConfigName = "%s::name" // configs key of mysql database name

	/* Microsoft SQL Server */
	mssqlConfigUser = "%s::user"    // configs key of mssql database user
	mssqlConfigPwd  = "%s::pwd"     // configs key of mssql database password
	mssqlConfigHost = "%s::host"    // configs key of mssql database server host
	mssqlConfigPort = "%s::port"    // configs key of mssql database port
	mssqlConfigName = "%s::name"    // configs key of mssql database name
	mssqlConfigTout = "%s::timeout" // configs key of mssql database connect timeout

	// Mysql Server database source name for local connection
	mysqldsnLocal = "%s:%s@/%s?charset=%s"

	// Mysql Server database source name for tcp connection
	mysqldsnTcp = "%s:%s@tcp(%s)/%s?charset=%s"

	// Microsoft SQL Server database source name
	mssqldsn = "server=%s;port=%d;database=%s;user id=%s;password=%s;Connection Timeout=%d;Connect Timeout=%d;"
)

var (
	// WingHelper content provider to hold database connections,
	// it will nil before mvc.OpenMySQL() called.
	WingHelper *WingProvider

	// MssqlHelper content provider to hold mssql database connections,
	// it will nil before mvc.OpenMssql() called.
	MssqlHelper *WingProvider

	// ConnPool save databases connection pool
	ConnPool = make(map[string]*WingProvider)

	// limitPageItems limit to show lits items in one page, default is 50,
	// you can use SetLimitPageItems() to change the limit value.
	limitPageItems = 50
)

// readMySQLCofnigs read mysql database params from config file,
// than verify them if empty except host.
func readMySQLCofnigs(session string) (string, string, string, string, error) {
	user := beego.AppConfig.String(fmt.Sprintf(mysqlConfigUser, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(mysqlConfigPwd, session))
	host := beego.AppConfig.String(fmt.Sprintf(mysqlConfigHost, session))
	name := beego.AppConfig.String(fmt.Sprintf(mysqlConfigName, session))

	if user == "" || pwd == "" || name == "" {
		return "", "", "", "", invar.ErrInvalidConfigs
	}
	return user, pwd, host, name, nil
}

// readMssqlCofnigs read mssql database params from config file,
// than verify them if empty.
func readMssqlCofnigs(session string) (string, string, string, int, string, int, error) {
	user := beego.AppConfig.String(fmt.Sprintf(mssqlConfigUser, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(mssqlConfigPwd, session))
	host := beego.AppConfig.DefaultString(fmt.Sprintf(mssqlConfigHost, session), "127.0.0.1")
	port := beego.AppConfig.DefaultInt(fmt.Sprintf(mssqlConfigPort, session), 1433)
	name := beego.AppConfig.String(fmt.Sprintf(mssqlConfigName, session))
	timeout := beego.AppConfig.DefaultInt(fmt.Sprintf(mssqlConfigTout, session), 600)

	if user == "" || pwd == "" || name == "" {
		return "", "", "", 0, "", 0, invar.ErrInvalidConfigs
	}
	return user, pwd, host, port, name, timeout, nil
}

// openMySQLPool open mysql and cached to connection pool by given session keys
func openMySQLPool(charset string, sessions []string) error {
	for _, session := range sessions {
		// combine develop session key on dev mode
		if beego.BConfig.RunMode == "dev" {
			session = session + "-dev"
		}

		// load configs by session key
		dbuser, dbpwd, dbhost, dbname, err := readMySQLCofnigs(session)
		if err != nil {
			return err
		}

		dsn := ""
		if len(dbhost) > 0 /* check database host whether using TCP to connect */ {
			// conntect with remote host database server
			dsn = fmt.Sprintf(mysqldsnTcp, dbuser, dbpwd, dbhost, dbname, charset)
		} else {
			// just connect local database server
			dsn = fmt.Sprintf(mysqldsnLocal, dbuser, dbpwd, dbname, charset)
		}
		logger.D("Mysql session:", session, " for DSN:", dsn)

		// open and connect database
		con, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		// check database validable
		if err = con.Ping(); err != nil {
			return err
		}

		con.SetMaxIdleConns(100)
		con.SetMaxOpenConns(100)
		ConnPool[session] = &WingProvider{con}
	}
	return nil
}

// OpenMySQL connect database and check ping result, the connection holded
// by mvc.WingHelper object if signle connect, or cached connections in ConnPool map
// if multiple connect and select same one by given sessions of input params.
// the datatable charset maybe 'utf8' or 'utf8mb4' same as database set.
//
// `USAGE`
//
// you must config database params in /conf/app.config file as follows
//
// ---
//
// #### Case 1 : For signle connect on prod mode.
//
//	[mysql]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 2 : For signle connect on dev mode.
//
//	[mysql-dev]
//	host = "127.0.0.1:3306"
//	name = "sampledb"
//	user = "root"
//	pwd  = "123456"
//
// #### Case 3 : For both dev and prod mode, you can config all of up cases.
//
// #### Case 4 : For multi-connections to set custom session keywords.
//
//	[mysql-a]
//	... same as (1)
//
//	[mysql-a-dev]
//	... same as (2)
//
//	[mysql-x]
//	... same as (1)
//
//	[mysql-x-dev]
//	... same as (2)
func OpenMySQL(charset string, sessions ...string) error {
	if len(sessions) > 0 {
		if err := openMySQLPool(charset, sessions); err != nil {
			return err
		}
		WingHelper = Select(sessions[0]) // using the first connection as primary helper
	} else {
		session := "mysql"
		if err := openMySQLPool(charset, []string{session}); err != nil {
			return err
		}
		WingHelper = Select(session)
	}
	return nil
}

// @deprecated Use OpenMySQL instead
func OpenMySQLByKeys(charset string, sessions ...string) error {
	return OpenMySQL(charset, sessions...)
}

// Select mysql Connection by request key words
// if mode is dev, the key will auto splice '-dev'
func Select(session string) *WingProvider {
	if beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}
	return ConnPool[session]
}

// OpenMssql connect mssql database and check ping result,
// the connections holded by mvc.MssqlHelper object,
// the charset maybe 'utf8' or 'utf8mb4' same as database set.
//
// `NOTICE`
//
// you must config database params in /conf/app.config file as:
//
// ---
//
// #### Case 1 For connect on prod mode.
//
//	[mssql]
//	host    = "127.0.0.1"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 600
//
// #### Case 2 For connect on dev mode.
//
//	[mssql-dev]
//	host    = "127.0.0.1"
//	port    = 1433
//	name    = "sampledb"
//	user    = "sa"
//	pwd     = "123456"
//	timeout = 600
//
// #### Case 3 For both dev and prod mode, you can config all of up cases.
func OpenMssql(charset string) error {
	session := "mssql"
	if beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}

	user, pwd, server, port, dbn, to, err := readMssqlCofnigs(session)
	if err != nil {
		return err
	}

	// get connection and connect timeouts
	dts := []int{600, 600}
	if to > 0 {
		dts[0] = to
		dts[1] = to
	}

	driver := "mssql"
	dsn := fmt.Sprintf(mssqldsn, server, port, dbn, user, pwd, dts[0], dts[1])
	logger.D("SQL Server DSN:", dsn)

	// open and connect database
	con, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	// check database validable
	if err = con.Ping(); err != nil {
		return err
	}

	con.SetMaxIdleConns(100)
	con.SetMaxOpenConns(100)
	MssqlHelper = &WingProvider{con}
	return nil
}

// SetLimitPageItems set global setting of limit items in page,
// the input value must range in (0, 1000].
func SetLimitPageItems(limit int) {
	if limit > 0 && limit <= 1000 {
		limitPageItems = limit
	}
}

// Stub return content provider connection
func (w *WingProvider) Stub() *sql.DB {
	return w.Conn
}

// Query call sql.Query()
func (w *WingProvider) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return w.Conn.Query(query, args...)
}

// Prepare call sql.Prepare()
func (w *WingProvider) Prepare(query string) (*sql.Stmt, error) {
	return w.Conn.Prepare(query)
}

// IsEmpty call sql.Query() to check target data if exist
func (w *WingProvider) IsEmpty(query string, args ...interface{}) (bool, error) {
	rows, err := w.Conn.Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return !rows.Next(), nil
}

// QueryOne call sql.Query() to query one record
func (w *WingProvider) QueryOne(query string, cb ScanCallback, args ...interface{}) error {
	rows, err := w.Conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return invar.ErrNotFound
	}
	rows.Columns()
	return cb(rows)
}

// QueryArray call sql.Query() to query multi records
func (w *WingProvider) QueryArray(query string, cb ScanCallback, args ...interface{}) error {
	rows, err := w.Conn.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Columns()
		if err := cb(rows); err != nil {
			return err
		}
	}
	return nil
}

// Insert call sql.Prepare() and stmt.Exec() to insert a new record
func (w *WingProvider) Insert(query string, args ...interface{}) (int64, error) {
	stmt, err := w.Conn.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

// Execute call sql.Prepare() and stmt.Exec() to update or delete records
func (w *WingProvider) Execute(query string, args ...interface{}) error {
	_, err := w.ExecuteWithResult(query, args...)
	return err
}

// Execute call sql.Prepare() and stmt.Exec() to update or delete records
func (w *WingProvider) ExecuteWithResult(query string, args ...interface{}) (int64, error) {
	stmt, err := w.Conn.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return w.Affected(result)
}

// AppendLike append like keyword end of sql string,
// DON'T call it after AppendLimit()
func (w *WingProvider) AppendLike(query, filed, keyword string, and ...bool) string {
	if len(and) > 0 && and[0] {
		return query + " AND " + filed + " LIKE '%%" + keyword + "%%'"
	}
	return query + " WHERE " + filed + " LIKE '%%" + keyword + "%%'"
}

// AppendLimit append page limitation end of sql string,
// DON'T call it before AppendLick()
func (w *WingProvider) AppendLimit(query string, page int) string {
	offset, items := page*limitPageItems, limitPageItems
	return query + " LIMIT " + fmt.Sprintf("%d, %d", offset, items)
}

// AppendLikeLimit append like keyword and limit end of sql string
func (w *WingProvider) AppendLikeLimit(query, filed, keyword string, page int, and ...bool) string {
	return w.AppendLimit(w.AppendLike(query, filed, keyword, and...), page)
}

// CheckAffected append page limitation end of sql string
func (w *WingProvider) Affected(result sql.Result) (int64, error) {
	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return 0, invar.ErrNotChanged
	}
	return row, nil
}

// FormatSets format update sets for sql update
//
// ---
//
//	sets := w.FormatSets(struct {
//		StringFiled string
//		EmptyString string
//		BlankString string
//		TrimString  string
//		IntFiled    int
//		I32Filed    int32
//		I64Filed    int64
//		F32Filed    float32
//		F64Filed    float64
//		BoolFiled   bool
//	}{"string", "", " ", " trim ", 123, 32, 64, 32.123, 64.123, true})
//	// sets: stringfiled='string', trimstring='trim', intfiled=123, i32filed=32, i64filed=64, f32filed=32.123, f64filed=64.123, boolfiled=true
//	logger.I("sets:", sets)
func (w *WingProvider) FormatSets(updates interface{}) string {
	sets := []string{}
	keys, values := reflect.TypeOf(updates), reflect.ValueOf(updates)
	for i := 0; i < keys.NumField(); i++ {
		name := strings.ToLower(keys.Field(i).Name)
		if name == "" {
			continue
		}

		value := values.Field(i).Interface()
		switch value.(type) {
		case bool:
			sets = append(sets, fmt.Sprintf(name+"=%v", value))
		case invar.Bool:
			boolvalue := value.(invar.Bool)
			if boolvalue != invar.BNone {
				truevalue := (boolvalue == invar.BTrue)
				sets = append(sets, fmt.Sprintf(name+"=%v", truevalue))
			}
		case string:
			trimvalue := strings.Trim(value.(string), " ")
			if trimvalue != "" { // filter empty string fields
				sets = append(sets, fmt.Sprintf(name+"='%s'", trimvalue))
			}
		case int, int8, int16, int32, int64, float32, float64,
			invar.Status, invar.Box, invar.Role, invar.Limit, invar.Lang, invar.Kind:
			if fmt.Sprintf("%v", value) != "0" { // filter 0 fields
				sets = append(sets, fmt.Sprintf(name+"=%v", value))
			}
		}
	}
	return strings.Join(sets, ", ")
}

// Execute a sql transaction, this method can provide high-performance
// multi datas insert or update operation by using combined query string.
//
// ---
//
//	query := "INSERT sametable (field1, field2, fieldn) VALUES %s"
//	for _, d := range params {
//		query += fmt.Sprintf("(%s, %s, %s),", d.v1, d.v2, d.v3)
//	}
//	query = strings.Trim(query, ",")
//	err := mvc.Transaction(query)
func (w *WingProvider) Transaction(query string, args ...interface{}) error {
	tx, err := w.Conn.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query, args...); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
