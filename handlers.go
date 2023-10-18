package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abdoroot/todolist/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Custom middleware
func IsAuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		returnToUrl := c.Request.Host + c.Request.URL.Path
		session := sessions.Default(c)
		ss := session.Get("loginuserid")
		if ss != nil {
			c.Next()
		} else {
			c.Redirect(http.StatusMovedPermanently, SiteBase+"login?return_to="+returnToUrl)
			c.Abort()
		}
	}
}

// post
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	//session.Delete("loginuseremail")
	c.Redirect(http.StatusTemporaryRedirect, SiteBase+"login")
}

func AuthUserId(c *gin.Context) (int, error) {
	session := sessions.Default(c)
	userId := session.Get("loginuserid").(string)
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return 0, err
	}
	if userIdInt == 0 {
		return 0, fmt.Errorf("User Session Not Avilable")
	}
	return userIdInt, nil
}

// get
func Home(c *gin.Context) {
	doneTasks, ok := c.GetQuery("done")
	tasks := []types.Tasks{}
	userId, err := AuthUserId(c)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(userId)
	sqlStr := "select id,name,due_date,priority,description from tasks where user_id=$1 and done=0 order by due_date asc"
	if ok && doneTasks == "true" {
		sqlStr = "select id,name,due_date,priority,description from tasks where user_id=$1 and done=1 order by due_date asc"
	}
	result, err := db.Query(sqlStr, userId)
	if err != nil {
		panic("dbs errors:" + err.Error())
	}
	defer result.Close()
	for result.Next() {
		task := types.Tasks{}
		result.Scan(&task.Id, &task.Name, &task.DueDate, &task.Priority, &task.Description)
		tasks = append(tasks, task)
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"siteBase": SiteBase,
		"tasks":    tasks,
	})
}

// get
func ShowTask(c *gin.Context) {
	taskId := c.Param("id")
	taskIdInt, err := strconv.ParseInt(taskId, 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	}
	task := types.Tasks{}
	dberr := db.QueryRow("select id,name,due_date,priority,description from tasks where id = $1", taskIdInt).Scan(&task.Id, &task.Name, &task.DueDate, &task.Priority, &task.Description)
	if dberr != nil {
		fmt.Println(dberr.Error())
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	}
	c.HTML(200, "show.html", gin.H{
		"task":     task,
		"siteBase": SiteBase,
	})
}

// get
func CreateTask(c *gin.Context) {
	priorities := []types.Priorities{}
	result, err := db.Query("select * from priorities")
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	}
	for result.Next() {
		priority := types.Priorities{}
		result.Scan(&priority.Id, &priority.Name)
		priorities = append(priorities, priority)
	}
	c.HTML(http.StatusOK, "create.html", gin.H{
		"siteBase":   SiteBase,
		"priorities": priorities,
	})
}

func TaskDone(c *gin.Context) {
	TaskId := c.PostForm("taskid")
	//check if the user own this task
	UserId, err := AuthUserId(c)
	if err != nil {
		log.Fatal(err)
		return
	}
	var taskId int
	err = db.QueryRow("select id from tasks where id = $1 and user_id=$2", TaskId, UserId).Scan(&taskId)
	if err != nil {
		log.Fatal(err)
		return
	}
	//update to done
	_, err = db.Exec("update tasks set done = 1 where id = $1 and user_id=$2", TaskId, UserId)
	if err != nil {
		log.Fatal(err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, SiteBase)
}

// get
func EditTask(c *gin.Context) {
	taskId := c.Param("id")
	priorities := []types.Priorities{}
	taskIdInt, err := strconv.ParseInt(taskId, 10, 0)
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	}
	task := types.Tasks{}
	dberr := db.QueryRow("select id,name,due_date,priority,description from tasks where id = $1", taskIdInt).Scan(&task.Id, &task.Name, &task.DueDate, &task.Priority, &task.Description)
	result, err := db.Query("select * from priorities")
	if dberr != nil && err != nil {
		fmt.Println(dberr.Error() + "-" + err.Error())
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	}
	for result.Next() {
		priority := types.Priorities{}
		result.Scan(&priority.Id, &priority.Name)
		priorities = append(priorities, priority)
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"siteBase":   SiteBase,
		"task":       task,
		"priorities": priorities,
	})
}

// post
func DoEditTask(c *gin.Context) {
	taskId := c.Param("id")
	name := c.PostForm("name")
	dueDate := c.PostForm("due_date")
	priority := c.PostForm("priority")
	description := c.PostForm("description")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	updateTaskError := ""
	updateTaskSucess := ""
	if name != "" && dueDate != "" && description != "" {
		//update task
		_, err := db.Exec("update tasks set name= $1 , due_date=$2,priority =$3,description=$4,updated_at =$5 where id =$6", name, dueDate, priority, description, updatedAt, taskId)

		if err != nil {
			fmt.Println(err.Error())
			c.Redirect(http.StatusMovedPermanently, SiteBase)
		}
		updateTaskSucess = "task Added Successfully"

	} else {
		updateTaskError = "Please re check you inputs"
	}
	_, _ = updateTaskError, updateTaskSucess //not used

	c.Redirect(http.StatusMovedPermanently, SiteBase)
}

// post
func DoCreateTask(c *gin.Context) {
	name := c.PostForm("name")
	dueDate := c.PostForm("due_date")
	priority := c.PostForm("priority")
	description := c.PostForm("description")
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	userId, err := AuthUserId(c)
	if err != nil {
		log.Fatal(err)
		return
	}
	var CreateTaskError, CreateTaskSucess string
	if name != "" && dueDate != "" && description != "" {
		//save task
		_, err := db.Exec("insert into tasks(user_id,name,due_date,priority,description,created_at) values($1,$2,$3,$4,$5,$6)", userId, name, dueDate, priority, description, createdAt)
		if err != nil {
			CreateTaskError = "db error :" + err.Error()
			fmt.Println(err.Error())
			c.Redirect(http.StatusMovedPermanently, SiteBase)
		}
		CreateTaskSucess = "task Added Successfully"

	} else {
		CreateTaskError = "Please re check you inputs"
	}

	c.HTML(200, "create.html", gin.H{
		"siteBase":         SiteBase,
		"CreateTaskSucess": CreateTaskSucess,
		"CreateTaskError":  CreateTaskError,
	})
}

// get
func LoginIndex(c *gin.Context) {
	returnTo := c.DefaultQuery("return_to", "")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"return_to": returnTo,
	})
}

// post
func DoLogin(c *gin.Context) {
	var login types.Login
	email := c.PostForm("email")
	password := c.PostForm("password")
	returnTo := c.PostForm("return_to")

	err := db.QueryRow("SELECT id,email,password from users where email = $1", email).Scan(&login.Id, &login.Email, &login.Password)
	if err != nil {
		fmt.Println(err.Error())
		loginErr := "Error email or password"
		c.HTML(200, "login.html", gin.H{

			"error": loginErr,
		})
		return
	}
	hash := login.Password
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		loginErr := "Error email or password"
		c.HTML(200, "login.html", gin.H{

			"error": loginErr,
		})
		return
	}
	session := sessions.Default(c)
	session.Set("loginuserid", login.Id)
	session.Save()
	if len(returnTo) > 0 {
		c.Redirect(http.StatusMovedPermanently, returnTo)
	}
	c.Redirect(http.StatusMovedPermanently, SiteBase)
}

// get
func SignUpIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "sign_up.html", nil)
}

// post
func DOSignUp(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	createdAt := time.Now().Format("2006-01-02 15:04")

	if name != "" && strings.Contains(email, "@") && password != "" {

		passwordBytes, crptErr := bcrypt.GenerateFromPassword([]byte(password), 14)
		if crptErr != nil {
			log.Println(crptErr.Error())
			c.HTML(http.StatusOK, "sign_up.html", gin.H{
				"error": crptErr.Error(),
			})
			return
		}
		//hashed password
		password = string(passwordBytes)
		_, err := db.Exec("insert into users (name,email,password,created_at) values($1,$2,$3,$4)", name, email, password, createdAt)
		if err != nil {
			fmt.Println(err.Error())
			c.HTML(http.StatusOK, "sign_up.html", gin.H{
				"error": err.Error(),
			})
		}
		c.Redirect(http.StatusTemporaryRedirect, SiteBase)
	} else {
		c.HTML(http.StatusOK, "sign_up.html", gin.H{
			"error": "error form validation",
		})
	}
}
