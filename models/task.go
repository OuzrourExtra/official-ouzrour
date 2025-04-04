package models

type Tag struct {
	Name  string `json:"Name"`
	Color string `json:"Color"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Date        string `json:"Date"`
	PlannedDate string `json:"PlannedDate"`
	Progress    int    `json:"Progress"`
	Image       string `json:"Image"`
	Tags        []Tag  `json:"Tags"`
}
