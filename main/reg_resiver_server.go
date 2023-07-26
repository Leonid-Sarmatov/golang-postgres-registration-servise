package main

import (
 "encoding/json"
 "fmt"
 "io/ioutil"
 "log"
 "os"
 "net/http"
 "html_registration_web_site/sql_connector"
)

// Структура для json парсинга запроса
type Person struct {
    Login         string 
    Password      string 
    FavoriteGame  string 
}

// Структура для json парсинга ответа
type Answer struct {
    Text          string  
    FavoriteGame  string 
}

func main() {
    fmt.Println("run!")
    defer fmt.Println("stop!")
   
    // Запускаем соединение с базой данных
    sql_connector.CreateNewSQLConnector("postgres", "USPEH", "users", "disable")
    
    // Функция запуска сайта
    go regSite()
    
    // Функция запуска ресивера для регистрации
    go regResiver()
    
    for {
    }
}

func regSite() {

   // Запускаем прослушивание на :8080 порту   http://localhost:8080/regApp
   http.HandleFunc("/regApp", siteRequest)
   log.Fatal(http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux)))
}

func regResiver() {

    // Запускаем прослушивание на :8081 порту
    http.HandleFunc("/regResiverServer", regHandleRequest)
    log.Fatal(http.ListenAndServe(":8081", addCorsHeaders(http.DefaultServeMux)))
}

func regHandleRequest(w http.ResponseWriter, r *http.Request) {

    // Если пришел POST запрос
    if r.Method == "POST" {
        // Если получен нужный заголовок
        if r.URL.Path == "/regResiverServer" {
            // Получаем тело запроса
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                    http.Error(w, "Failed to read request body", http.StatusBadRequest)
                return
            }

            // Создаем структуру для парсинга
            var person Person

            // Парсим тело запроса в труктуру
            err = json.Unmarshal(body, &person)
            if err != nil {
                http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
                return
            }

            // Печатаем данные из запроса
            fmt.Println("Received JSON:", string(body))
            fmt.Println("Login:", person.Login)
            fmt.Println("Password:", person.Password)
            fmt.Println("Favorite game:", person.FavoriteGame)

            var answer Answer
            dbAnswer := sql_connector.ReadUserInfo(person.Login)
            if dbAnswer != "" {

                answer.Text = "You have been registration before"
                answer.FavoriteGame = "Your favourute game is "+dbAnswer

            } else {
                // Запись в базу данных
                sql_connector.AddData(person.Login, person.Password, person.FavoriteGame)

                answer.Text = "Successfull registration!"
                answer.FavoriteGame = "Cool game!"
            }
            

            // Парсим в json структуру ответа
            responseJSON, err := json.Marshal(answer)
            if err != nil {
                http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
                return
            }

            // Конфигурируем ответ
            w.Header().Set("Content-Type", "application/json")
            // Отправляем ответ
            w.Write(responseJSON)
        } else {
            http.NotFound(w, r)
        }
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}


// Функция для подключения заголовков. Без нее возможны ошибки
func addCorsHeaders(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        handler.ServeHTTP(w, r)
    })
}

func siteRequest(w http.ResponseWriter, r *http.Request) {
   filePath := "index.html" 

   file, err := os.Open(filePath)
   if err != nil {
      http.Error(w, "Failed to open file", http.StatusInternalServerError)
      return
   }
   defer file.Close()

   fileInfo, err := file.Stat()
   if err != nil {
      http.Error(w, "Failed to get file info", http.StatusInternalServerError)
      return
   }

   fileSize := fileInfo.Size()
   buffer := make([]byte, fileSize)

   _, err = file.Read(buffer)
   if err != nil {
      http.Error(w, "Failed to read file", http.StatusInternalServerError)
      return
   }
   
   fmt.Fprint(w, string(buffer))
}







