package dto

type Error struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type CreateTodoInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Todo struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
