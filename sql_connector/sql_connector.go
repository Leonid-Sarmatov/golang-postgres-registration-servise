/**

go mod init main
go get github.com/lib/pq
*/
package sql_connector

import (
    "database/sql"
    //"fmt"
    //"reflect"
    _ "github.com/lib/pq"
)

type data struct {
    id int
    user_name string
    password string
    user_info string
}

var db *sql.DB

func CreateNewSQLConnector(user, password, dbname, sslmode string) {
    connStr := "user="+user+" password="+password+" dbname="+dbname+" sslmode="+sslmode
    d, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    } 
    db = d
} 

// Добаление нового пользователя
func AddData(userName, password, userInfo string) {
    _, err := db.Exec("insert into users_table (user_name, password, user_info) values ($1, $2, $3)", userName, password, userInfo)
    if err != nil {
        panic(err)
    }
    
    defer db.Close()
}

// Удаление данных
func RemoveData(userName string) {
    _, err := db.Exec("delete from users_table where user_name=$1", userName)
    if err != nil {
        panic(err)
    }
    
    defer db.Close()
}

// Чтение пользовательской информации
func ReadUserInfo(userName string) string {
    result := ""
    rows, err := db.Query("select user_info from users_table where user_name=$1", userName)
    if err != nil {
        panic(err)
    }
    
    defer rows.Close()

    rows.Next()

    err2 := rows.Scan(&result)
    if err2 != nil{
	    //panic(err2)
        return ""
    }
    
    return result
}
/*
func main() {
    connStr := "user=postgres password=USPEH dbname=users sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    } 
    
    defer db.Close()
    
    //RemoveData("test_user_name", db)
    //AddData("user_babaika", "pass", "game", db)
    fmt.Println(ReadUserInfo("user_babaika"))
    
    /*
    // Добавление данных
    result, err := db.Exec("insert into users_table (user_name, password, user_info) values ('test_user_name', 'test_password', 'test_user_info');")
    if err != nil {
        panic(err)
    }
    */
    
    //fmt.Println(reflect.TypeOf(db))
    //fmt.Println(reflect.TypeOf(result))
    /*
    // Чтение данных
    rows, err := db.Query("select * from users_table")
    if err != nil {
        panic(err)
    }
    
    defer rows.Close()
    dataList := []data{}
    
    for rows.Next(){
        p := data{}
        err := rows.Scan(&p.id, &p.user_name, &p.password, &p.user_info)
        if err != nil{
            fmt.Println(err)
            continue
        }
        dataList = append(dataList, p)
    }
    
    for _, p := range dataList{
        fmt.Println(p.id, p.user_name, p.password, p.user_info)
    }
    
}
*/
























