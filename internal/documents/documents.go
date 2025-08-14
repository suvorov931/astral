package documents

import "time"

type Document struct {
	Id        string
	Login     string
	Name      string
	Mime      string
	File      bool
	Public    bool
	Grant     []string
	Content   []byte
	JSON      []byte
	CreatedAt time.Time
}
