package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pasque/app"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/***************************************************************************
*
* NullString
*
***************************************************************************/
type NullString struct {
	sql.NullString
}

func (c *NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String)
}

type TInt64 struct {
	sql.NullInt64
}
type TStr struct {
	sql.NullString
}

func SetNullInt64(uval int64) TInt64 {
	var p TInt64
	if uval == 0 {
		p.Valid = false
	} else {
		p.Valid = true
		p.Int64 = uval
	}
	return p
}
func SetNullString(str string) TStr {
	var p TStr
	if 0 == strings.Compare(str, "") {
		p.Valid = false
	} else {
		p.Valid = true
		p.String = str
	}
	return p
}

type DateTime struct {
	time.Time
}

func Now() *DateTime {
	p := DateTime{}
	p.Time = time.Now()
	return &p
}

func (d *DateTime) String() string {
	return d.Format(TIME_LAYOUT)

}
func (dt *DateTime) Add(d time.Duration) *DateTime {
	dt.Time = dt.Time.Add(d)
	return dt
}
func (dt *DateTime) Sub(d *DateTime) time.Duration {
	return dt.Time.Sub(d.Time)
}
func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TIME_LAYOUT)+2)
	b = append(b, '"')
	b = time.Time(t.Time).AppendFormat(b, TIME_LAYOUT)
	b = append(b, '"')
	return b, nil
}

func (t *DateTime) UnmarshalJSON(b []byte) error {
	if tv, err := time.Parse(`"`+TIME_LAYOUT+`"`, string(b)); err != nil {
		return err
	} else {
		t.Time = tv
		return nil
	}
}

/***************************************************************************
*
* DbPool
*
***************************************************************************/
type DbPool struct {
	conn    *sql.DB
	cap     chan int
	constr  string
	driver  string
	maxconn int
}

type DbRows struct {
	Rows []map[string]string
}

func (p *DbPool) Ping() error {
	return p.conn.Ping()
}
func (p *DbPool) Conn() *sql.DB {
	return p.conn
}
func (p *DbRows) IsNull() bool {
	if nil == p.Rows || 0 == len(p.Rows) {
		return true
	}
	return false
}
func (p DbPool) Prepare(query string, flag bool) (*sql.Stmt, error) {
	stmt, err := p.Conn().Prepare(query)
	if !flag || nil == err {
		return stmt, err
	}
	if _err := p.Ping(); nil != _err {
		return nil, _err
	} else {
		return p.Prepare(query, false)
	}
}
func IsNull(p *DbRows) bool {
	if nil == p || nil == p.Rows || 0 == len(p.Rows) {
		return true
	}
	return false
}

func InitDbPool(dbName string, dbCfg *app.DbConfig) (*DbPool, error) {
	conn, err := dbCfg.Conn(&dbName)
	if err != nil {
		app.ErrorLog("%s", err.Error())
		return nil, err
	}

	p := DbPool{}
	p.cap = make(chan int, conn.MaxConn)
	p.constr = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conn.UserName, conn.Password, conn.Address, conn.Database)
	p.driver = conn.Driver
	p.maxconn = conn.MaxConn
	for i := 0; i < p.maxconn; i++ {
		p.cap <- 1
	}
	db, err := sql.Open(p.driver, p.constr)
	if nil != err {
		return nil, err
	}
	db.SetMaxOpenConns(p.maxconn)

	if err := db.Ping(); nil != err {
		app.ErrorLog("%s", err.Error())
		return nil, err
	}
	p.conn = db
	return &p, nil
}

func parseRows(rows *sql.Rows, cols *[]string, str *string) error {
	count := len(*cols)
	vals := make([]NullString, count)
	*str = ""
	flag := rows.Next()
	if !flag {
		return nil
	}
	valPtr := make([]interface{}, count)
	for ; flag; flag = rows.Next() {
		for i, _ := range *cols {
			valPtr[i] = &vals[i]
		}
		if err := rows.Scan(valPtr...); err != nil {
			return err
		}
		value := ""
		for i, col := range *cols {
			val := vals[i]
			value = fmt.Sprintf("%s,\"%s\":\"%s\"", value, col, val.String)
		}
		*str = fmt.Sprintf("%s,{%s}", *str, strings.TrimLeft(value, ","))
	}
	*str = fmt.Sprintf("[%s]", strings.TrimLeft(*str, ","))
	if err := rows.Err(); nil != err {
		return err
	}
	return nil
}
func parseRowsToMap(rows *sql.Rows, cols *[]string) (*DbRows, error) {
	count := len(*cols)
	vals := make([]NullString, count)
	// *str = ""
	res := DbRows{Rows: make([]map[string]string, 0)}
	flag := rows.Next()
	if !flag {
		return &res, nil
	}
	valPtr := make([]interface{}, count)
	for ; flag; flag = rows.Next() {
		for i, _ := range *cols {
			valPtr[i] = &vals[i]
		}
		if err := rows.Scan(valPtr...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]string)
		for i, col := range *cols {
			val := vals[i]
			rowMap[col] = val.String
			// *str = fmt.Sprintf("%s,\"%s\":\"%s\"", *str, col, val.String)
		}
		res.Rows = append(res.Rows, rowMap)
		// *str = fmt.Sprintf(",{%s}", strings.TrimLeft(*str, ","))
	}
	// *str = fmt.Sprintf("[%s]", strings.TrimLeft(*str, ","))
	if err := rows.Err(); nil != err {
		return nil, err
	}
	return &res, nil
}
func (p DbPool) GetCap() chan int {
	return p.cap
}

func (p DbPool) QueryTest(query string, args ...interface{}) (string, error) {
	<-p.cap
	defer func() {
		p.cap <- 1
	}()
	str := ""
	stmt, err := p.Prepare(query, true)
	if err != nil {
		return str, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		println("Query error... ", err.Error())
		return str, err
	}
	defer rows.Close()

	cols, colErr := rows.Columns()
	if colErr != nil {
		return str, err
	}
	parseErr := parseRows(rows, &cols, &str)
	return str, parseErr
}
func (p DbPool) QueryToJSONStr(query string, args ...interface{}) (string, error) {
	<-p.cap
	defer func() {
		p.cap <- 1
	}()
	str := ""
	stmt, err := p.Prepare(query, true)
	if err != nil {
		return str, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		println("Query error... ", err.Error())
		return str, err
	}
	defer rows.Close()

	cols, colErr := rows.Columns()
	if colErr != nil {
		return str, err
	}

	parseErr := parseRows(rows, &cols, &str)
	return str, parseErr
}
func (p DbPool) QueryToRows(query string, args ...interface{}) (*DbRows, error) {

	<-p.cap
	defer func() {
		p.cap <- 1
	}()
	stmt, err := p.Prepare(query, true)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, colErr := rows.Columns()
	if colErr != nil {
		return nil, err
	}
	return parseRowsToMap(rows, &cols)
}

//Query return *DbRows
func (p DbPool) Query(query string, args ...interface{}) (*DbRows, error) {

	<-p.cap
	defer func() {
		p.cap <- 1
	}()
	stmt, err := p.Prepare(query, true)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		println("Query error... ", err.Error())
		return nil, err
	}
	defer rows.Close()

	cols, colErr := rows.Columns()
	if colErr != nil {
		return nil, err
	}

	return parseRowsToMap(rows, &cols)
}

func (p DbPool) Exec(query string, args ...interface{}) (int64, error) {
	<-p.cap
	defer func() {
		p.cap <- 1
	}()
	stmt, err := p.Prepare(query, true)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	ret, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return ret, nil
}

func InitDb(dbName string, handler func(p *DbPool)) error {
	if db, err := InitDbPool(dbName, &app.CfgDb); nil != err {
		return fmt.Errorf("db initialization fail.... %s", err.Error())
	} else {
		handler(db)
	}
	return nil
}
