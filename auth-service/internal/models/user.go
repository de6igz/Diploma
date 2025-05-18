package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // упрощённо (для реального проекта: храните хеш)
}
