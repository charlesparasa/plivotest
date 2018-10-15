package model

type Contact struct {
	Id int  `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type Contacts struct {
	Contacts []Contact `json:"contacts"`
}

