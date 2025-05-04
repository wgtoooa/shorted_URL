package Table

import "time"

type URL struct {
	Id         int
	Full_url   string
	Short_url  string
	Account_id int
	Created_at time.Time
}
