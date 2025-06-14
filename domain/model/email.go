package model

type EmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
