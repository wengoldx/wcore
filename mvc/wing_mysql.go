// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
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
	//
	// ----------------------------------------
)

// WingProvider content provider to support database utils
type WingProvider struct {
	Conn *sql.DB
}

// ScanCallback use for scan query result from rows
type ScanCallback func(rows *sql.Rows) error

// FormatCallback format query value to string for MultiInsert().
type FormatCallback func(index int) string

// TransactionCallback transaction callback for MultiTransaction().
type TransactionCallback func(tx *sql.Tx) (sql.Result, error)

// MySQL database configs
const (
	mysqlConfigUser = "%s::user" // configs key of mysql database user
	mysqlConfigPwd  = "%s::pwd"  // configs key of mysql database password
	mysqlConfigHost = "%s::host" // configs key of mysql database host and port
	mysqlConfigName = "%s::name" // configs key of mysql database name

	// Mysql Server database source name for local connection
	mysqldsnLocal = "%s:%s@/%s?charset=%s"

	// Mysql Server database source name for tcp connection
	mysqldsnTcp = "%s:%s@tcp(%s)/%s?charset=%s"
)

var (
	// WingHelper content provider to hold database connections,
	// it will nil before mvc.OpenMySQL() called.
	WingHelper *WingProvider

	// Cache all mysql providers into pool for multiple databases server connect.
	connPool = make(map[string]*WingProvider)
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
		logger.I("Open MySQL on {", session, ":", dsn, "}")

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
		con.SetConnMaxLifetime(28740)
		connPool[session] = &WingProvider{con}
	}
	return nil
}

// OpenMySQL connect database and check ping result, the connection holded
// by mvc.WingHelper object if signle connect, or cached connections in connPool map
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
//	... same as use Case 1.
//
//	[mysql-a-dev]
//	... same as use Case 2.
//
//	[mysql-x]
//	... same as use Case 1.
//
//	[mysql-x-dev]
//	... same as use Case 2.
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

// Select mysql Connection by request key words
// if mode is dev, the key will auto splice '-dev'
func Select(session string) *WingProvider {
	if beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}
	return connPool[session]
}

// Stub return content provider connection
func (w *WingProvider) Stub() *sql.DB {
	return w.Conn
}

// Query call sql.Query()
func (w *WingProvider) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return w.Conn.Query(query, args...)
}

// IsEmpty call sql.Query() to check target data if empty
func (w *WingProvider) IsEmpty(query string, args ...interface{}) (bool, error) {
	rows, err := w.Conn.Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return !rows.Next(), nil
}

// IsExist call sql.Query() to check target data if exist
func (w *WingProvider) IsExist(query string, args ...interface{}) (bool, error) {
	empty, err := w.IsEmpty(query, args...)
	return !empty, err
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

// Insert call sql.Prepare() and stmt.Exec() to insert a new record.
//
// `@see` Use MultiInsert() to insert multiple values in once database operation.
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

// MultiInsert format and combine multiple values to insert at once, this method can provide
// high-performance than call Insert() one by one.
//
// ---
//
//	query := "INSERT sametable (field1, field2, fieldn) VALUES"
//	err := mvc.MultiInsert(query, 5, func(index int) string {
//		return fmt.Sprintf("(%v, %v, %v)", v1, v2, v3)
//
//		// For string values like follows:
//		// return fmt.Sprintf("(\"%s\", \"%s\", \"%s\")", v1, v2, v3)
//	})
func (w *WingProvider) MultiInsert(query string, cnt int, cb FormatCallback) error {
	values := []string{}
	for i := 0; i < cnt; i++ {
		value := strings.TrimSpace(cb(i))
		if value != "" {
			values = append(values, value)
		}
	}
	query = query + " " + strings.Join(values, ",")
	return w.Execute(query)
}

// Execute call sql.Prepare() and stmt.Exec() to update or delete records
func (w *WingProvider) Execute(query string, args ...interface{}) error {
	stmt, err := w.Conn.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if _, err := stmt.Exec(args...); err != nil {
		return err
	}
	return nil
}

// ExeAffected call sql.Prepare() and stmt.Exec() to update or delete records
func (w *WingProvider) ExeAffected(query string, args ...interface{}) (int64, error) {
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

// AppendLike append like keyword end of sql string, DON'T call it when exist limit key in sql string
func (w *WingProvider) AppendLike(query, filed, keyword string, and ...bool) string {
	if len(and) > 0 && and[0] {
		return query + " AND " + filed + " LIKE '%%" + keyword + "%%'"
	}
	return query + " WHERE " + filed + " LIKE '%%" + keyword + "%%'"
}

// Affected append page limitation end of sql string
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
		switch tv := value.(type) {
		case bool:
			sets = append(sets, fmt.Sprintf(name+"=%v", tv))
		case invar.Bool:
			if tv != invar.BNone {
				truevalue := (tv == invar.BTrue)
				sets = append(sets, fmt.Sprintf(name+"=%v", truevalue))
			}
		case string:
			trimvalue := strings.Trim(tv, " ")
			if trimvalue != "" { // filter empty string fields
				sets = append(sets, fmt.Sprintf(name+"='%s'", trimvalue))
			}
		case int, int8, int16, int32, int64, float32, float64,
			invar.Status, invar.Box, invar.Role, invar.Limit, invar.Lang, invar.Kind:
			if fmt.Sprintf("%v", tv) != "0" { // filter 0 fields
				sets = append(sets, fmt.Sprintf(name+"=%v", tv))
			}
		}
	}
	return strings.Join(sets, ", ")
}

// Transaction execute one sql transaction, it will rollback when operate failed.
//
// `@see` Use MultiTransaction() to excute multiple transaction as once.
func (w *WingProvider) Transaction(query string, args ...interface{}) error {
	tx, err := w.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// MultiTransaction excute multiple transactions, it will rollback all operations when case error.
//
// ---
//
//	// Excute 3 transactions in callback with different query1 ~ 3
//	err := mvc.MultiTransaction(
//		func(tx *sqlTx) (sql.Result, error) { return tx.Exec(query1, args...) },
//		func(tx *sqlTx) (sql.Result, error) { return tx.Exec(query2, args...) },
//		func(tx *sqlTx) (sql.Result, error) { return tx.Exec(query3, args...) })
func (w *WingProvider) MultiTransaction(cbs ...TransactionCallback) error {
	tx, err := w.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// start excute multiple transactions in callback
	for _, cb := range cbs {
		if _, err := cb(tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
