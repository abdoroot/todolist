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
	///sessions
	SessionStore := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", SessionStore))
	//db connect
	//dbConnection, err := MySqlDbConnect()
	dbConnection, err := PgDbConnect() //postgres DB
	db = dbConnection
	fmt.Println(err)
	if err != nil {
		panic(fmt.Sprintf("could not connect to db : %v", err))
	}
}

func main() {
	r.GET("/login", LoginIndex)
	r.POST("/login", DoLogin)
	r.POST("/logout", DoLogin, IsAuthUser())
	r.GET("/signup", SignUpIndex)
	r.POST("/signup", DOSignUp)
	r.GET("/", IsAuthUser(), Home) //todo list home

	tg := r.Group("/task")
	{
		tg.Use(IsAuthUser())
		tg.POST("/", Logout)
		tg.GET("/:id", ShowTask)
		tg.GET("/create", CreateTask)
		tg.POST("/create", DoCreateTask)
		tg.GET("/:id/edit", EditTask)
		tg.POST("/:id/edit", DoEditTask)
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
