package Table

import "time"

type Account struct {
	Id        int
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	CountURL  int       `json:"count_url"`
	CreatedAt time.Time `json:"createdAt"`
}
