package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	ID   int    `sql:"id"`
	Name string `sql:"name"`
}

func testDBQuery() {

	dbString := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4", "ipv4_user", "123456", "192.168.123.20")
	dbTemp, err := sql.Open("mysql", dbString)
	if err != nil {
		log.Fatal(err)
		// return err
	}
	// }
	// dbTemp.SetMaxOpenConns(maxConn)
	// dbTemp.SetMaxIdleConns(maxIdle)
	err = dbTemp.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db := dbTemp
	if err := db.Ping(); err != nil {
		return
	}
	defer db.Close()

	// res = Ages
	//
	a := Data{}
	// b := make([]Person, 0, 1)
	queryRow(db, "SELECT * FROM cu.datacode where id=1111", &a)
	fmt.Println(a)

	b := queryAll(db, "SELECT * FROM cu.datacode", &Data{})
	fmt.Println(b)
	// fmt.Print(res)
}

func queryRow(db *sql.DB, query string, rowPtr interface{}) {

	rowPtrType := reflect.TypeOf(rowPtr)
	if rowPtrType.Kind() != reflect.Ptr {
		log.Fatal("must pass a ptr of value")
	}

	rowType := rowPtrType.Elem()
	if rowType.Kind() != reflect.Struct {
		log.Fatal("type not struct")
	}

	//获取标签
	tagIndex := make(map[string]int, rowType.NumField())
	rowValue := reflect.ValueOf(rowPtr).Elem()
	// fmt.Println(rowValue)

	for i := 0; i < rowType.NumField(); i++ {
		tag := rowType.Field(i).Tag.Get("sql")
		tagIndex[tag] = i
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	// colTypes, _ := rows.ColumnTypes()
	//根据cols 对 要传入scan的参数进行排序
	// fmt.Println(cols, colTypes)
	rows.Next()
	// rows.Next()

	//没有配置的字段，全部写入到这里
	var dft interface{}

	rowScan := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		if valueIndex, ok := tagIndex[cols[i]]; ok {
			rowScan[i] = rowValue.Field(valueIndex).Addr().Interface()
		} else {
			rowScan[i] = &dft
		}

	}

	rows.Scan(rowScan...)

	// fmt.Println(rowScan, rowPtr)

}

func queryAll(db *sql.DB, query string, rowStruct interface{}) []interface{} {
	structType := reflect.TypeOf(rowStruct)
	if structType.Kind() != reflect.Ptr {
		log.Fatal("rowStruct must a pointer")
	}

	rowType := structType.Elem()
	if rowType.Kind() != reflect.Struct {
		log.Fatal("type not struct")
	}

	// rowPtr := &row
	//获取标签
	tagIndex := make(map[string]int, rowType.NumField())
	rowValue := reflect.ValueOf(rowStruct).Elem()

	for i := 0; i < rowType.NumField(); i++ {
		tag := rowType.Field(i).Tag.Get("sql")
		tagIndex[tag] = i
	}

	// fmt.Println(rowValue, tagIndex)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	var dft interface{}

	rowScan := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		if valueIndex, ok := tagIndex[cols[i]]; ok {
			// fmt.Println(rowValue.Field(valueIndex).Addr())
			// fmt.Println(valueIndex)
			rowScan[i] = rowValue.Field(valueIndex).Addr().Interface()
		} else {
			rowScan[i] = &dft
		}
	}

	arr := make([]interface{}, 0, 1)

	for rows.Next() {
		if err := rows.Scan(rowScan...); err != nil {
			log.Fatal(err)
		}

		//所有值都写入了row, 因此添加到返回值中就可以
		// fmt.Println(rowValue)
		arr = append(arr, rowValue.Interface())
	}

	return arr
}

func main() {
	testDBQuery()
}
