package types

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
