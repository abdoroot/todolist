package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/abdoroot/todolist/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var JWTclaims map[string]interface{}

// api middleware
func IsAuthApiUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// المفتاح السري الذي تم استخدامه لتوقيع JWT
		secretKey := []byte("L1£1q81-%<|v")

		// احصل على السلسلة الممضاة JWT من الطلب (على سبيل المثال، من رأس الطلب)
		signedToken := c.Request.Header.Get("Authorization")

		// التحقق من صحة التوقيع وفحص الأدعاء
		token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			// إذا لم يتم التحقق من صحة التوقيع أو انتهت صلاحية التوكن
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		JWTclaims, _ = token.Claims.(jwt.MapClaims)

		c.Next()
	}
}

func Login(c *gin.Context) {

	inputs := types.Login{}
	if err := c.BindJSON(&inputs); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	if err := ValidateStruct(inputs); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	user := types.ApiLoginResponse{}
	if err := getUserFromDbBind(inputs.Email, inputs.Password, &user); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	token, err := generateJWT(user.Id)
	if err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data retrieved successfully",
		Data: map[string]interface{}{
			"token": token,
		},
	}

	c.JSON(http.StatusOK, resp)
}

func CreateApiTask(c *gin.Context) {
	//Bind inputs
	inputs := types.ApiCreateTaskRequest{}
	if err := c.BindJSON(&inputs); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}
	//save task
	if err := saveTaskDB(&inputs, c); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data Saved successfully",
		Data:    []string{},
	}

	c.JSON(http.StatusOK, resp)
}

// PATCH
func UpdateApiTask(c *gin.Context) {
	id := c.Param("id")
	//convert it to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}
	//Bind inputs
	inputs := types.ApiUpdateTaskRequest{}
	if err := c.BindJSON(&inputs); err != nil {
		HandlerAPiEroor(c, err.Error())
		return
	}

	ui := JWTclaims["user_id"]
	sqlResult, updateErr := db.Exec("update tasks set name = $1 , due_date = $2,priority = $3,description = $4 , done= $7 where id=$5 and user_id = $6", inputs.Name, inputs.DueDate, inputs.Priority, inputs.Description, idInt, ui, inputs.Done)
	if updateErr != nil {
		HandlerAPiEroor(c, updateErr.Error())
		return
	}

	if rowCount, err := sqlResult.RowsAffected(); err != nil || rowCount == 0 {
		HandlerAPiEroor(c, fmt.Sprintf("%v row updated", rowCount))
	}

	//return the response
	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data Updated successfully",
		Data:    []string{},
	}
	c.JSON(http.StatusOK, resp)
}

// DELETE
func DeleteApiTask(c *gin.Context) {
	id := c.Param("id")
	//convert it to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandlerAPiEroor(c, err.Error())
	}
	ui := JWTclaims["user_id"]
	_, delErr := db.Exec("delete from tasks where id=$1 and user_id = $2", idInt, ui)
	if delErr != nil {
		HandlerAPiEroor(c, delErr.Error())
		return
	}
	//return the response
	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data Deleted successfully",
		Data:    []string{},
	}
	c.JSON(http.StatusOK, resp)
}

func ShowApiTask(c *gin.Context) {
	id := c.Param("id")
	//convert it to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandlerAPiEroor(c, err.Error())
	}
	ui := JWTclaims["user_id"]
	tsk := types.Tasks{}
	quErr := db.QueryRow("select id,name,due_date,priority,description from tasks where id=$1 and user_id = $2 order by due_date asc", idInt, ui).Scan(&tsk.Id, &tsk.Name, &tsk.DueDate, &tsk.Priority, &tsk.Description)
	if quErr != nil {
		HandlerAPiEroor(c, quErr.Error())
		return
	}
	//return the response
	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data Found",
		Data:    tsk,
	}
	c.JSON(http.StatusOK, resp)
}

func saveTaskDB(inputs *types.ApiCreateTaskRequest, c *gin.Context) error {
	//validate inputs
	if err := ValidateStruct(inputs); err != nil {
		return err
	}
	logedUserId := JWTclaims["user_id"].(float64)
	createdAt := time.Now()
	_, err := db.Exec("insert into tasks(user_id,name,due_date,priority,description,created_at) values($1,$2,$3,$4,$5,$6)", logedUserId, inputs.Name, inputs.DueDate, inputs.Priority, inputs.Description, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func AllTasks(c *gin.Context) {
	tasks := []types.Tasks{}
	result, err := db.Query("select id,name,due_date,priority,description from tasks order by due_date asc")
	if err != nil {
		HandlerAPiEroor(c, err.Error())
	}
	defer result.Close()
	for result.Next() {
		task := types.Tasks{}
		result.Scan(&task.Id, &task.Name, &task.DueDate, &task.Priority, &task.Description)
		tasks = append(tasks, task)
	}

	resp := types.ApiResponse{
		Status:  http.StatusOK,
		Message: "Data retrieved successfully",
		Data:    tasks,
	}

	c.JSON(http.StatusOK, resp)
}

func HandlerAPiEroor(c *gin.Context, err string) {
	el := []string{} //empty list
	resp := types.ApiResponse{
		Status:  http.StatusNotFound,
		Message: "an error ocured :" + err,
		Data:    el,
	}
	c.JSON(http.StatusNotFound, resp)
	c.Abort()
}

func generateJWT(userId int) (string, error) {
	// المفتاح السري الذي يتم استخدامه لتوقيع JWT
	secretKey := []byte("L1£1q81-%<|v")

	//ال payload
	// إعداد معلومات المستخدم
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // تاريخ انتهاء الصلاحية بعد يوم واحد
	}

	fmt.Println(claims)

	//header
	// إنشاء التوقيع باستخدام المفتاح السري
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//signature
	// توقيع التوكن للحصول على سلسلة JWT
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateStruct(inputs interface{}) error {
	validate := validator.New() //create new validator
	if err := validate.Struct(inputs); err != nil {
		return err
	}
	return nil
}

func getUserFromDbBind(email, password string, user *types.ApiLoginResponse) error {
	var hash string
	err := db.QueryRow("select id,email,password from users where email = $1", email).Scan(&user.Id, &user.Email, &hash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println(err.Error())
		loginErr := "Error email or password"
		return fmt.Errorf(loginErr)
	}
	return nil
}
