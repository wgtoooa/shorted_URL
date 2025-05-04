package Table

import "time"

type URL struct {
	Id          int       `json:"id"`
	Full_url    string    `json:"full_url"`
	Short_url   string    `json:"short_url"`
	Account_id  int       `json:"account_id"`
	Created_at  time.Time `json:"created_at"`
	Description string    `json:"description"`
}
