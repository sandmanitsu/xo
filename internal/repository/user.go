package repository

type User struct {
	Id             int    `json:"id"`
	Login          string `json:"login"`
	HashedPassword string `json:"hashed_password"`
}
