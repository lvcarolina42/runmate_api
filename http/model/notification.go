package model

type NotifyInput struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Tokens []string `json:"tokens"`
}
