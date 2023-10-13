package main

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	//web routes
	r.GET("/login", LoginIndex)
	r.POST("/login", DoLogin)
	r.POST("/logout", DoLogin, IsAuthUser())
	r.GET("/signup", SignUpIndex)
	r.POST("/signup", DOSignUp)
	r.GET("/", IsAuthUser(), Home) // Todo list home

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

	//api routes
	api := r.Group("/api")
	api.Use(IsAuthApiUser())
	//not protected by the midlleware
	r.POST("/api/login", Login)
	{
		atg := api.Group("/tasks")
		atg.GET("/", AllTasks)            //get all tasks
		atg.POST("/", CreateApiTask)      //new task
		atg.GET("/:id", ShowApiTask)      //show task
		atg.PATCH("/:id", UpdateApiTask)  //update task
		atg.DELETE("/:id", DeleteApiTask) //update task
	}
}
