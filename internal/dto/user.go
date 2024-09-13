package dto

type Error struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type UserCommonInfo struct {
	UserId    uint   `json:"userId"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsAdmin   bool   `json:"isAdmin"`
}

type LoginResponse struct {
	Token string `json:"token"`
}