package main

import (
    "database/sql"
    "fmt"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
)

var r *gin.Engine

var db *sql.DB

func init() {
    r = gin.Default()
    r.LoadHTMLGlob("templates/*")

    // Sessions
    SessionStore := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", SessionStore))

    // Database connect
    // dbConnection, err := MySqlDbConnect()
    dbConnection, err := PgDbConnect() // Postgres DB
    db = dbConnection
    fmt.Println(err)
    if err != nil {
        panic(fmt.Sprintf("could not connect to db: %v", err))
    }

    // Call the SetupRoutes function to define your routes
    SetupRoutes(r)
}

func main() {
    r.Run() // Listen and serve on 0.0.0.0:8080
}
