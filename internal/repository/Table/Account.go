package Table

import "time"

type Account struct {
	Id        int
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	Status    string    `json:"status"`
	CountURL  int       `json:"count_url"`
	CreatedAt time.Time `json:"createdAt"`
}
