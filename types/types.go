package types

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ApiLoginResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiCreateTaskRequest struct {
	Name        string `json:"name" validate:"required"`
	DueDate     string `json:"due_date" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ApiUpdateTaskRequest struct {
	Name        string `json:"name" validate:"required"`
	DueDate     string `json:"due_date" validate:"required"`
	Priority    string `json:"priority" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ApiResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Tasks struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
}

type Priorities struct {
	Id   string `json:"id"`
	Name string `json:"Name"`
}
