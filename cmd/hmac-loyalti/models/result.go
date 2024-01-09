package models

type Result struct {
	Status    bool
	Message   string
	Method    string
	Path      string
	Body      string
	Timestamp string
	Hmac      string
}
