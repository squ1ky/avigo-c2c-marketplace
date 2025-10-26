package domain

type EmailMessage struct {
	To       string
	Subject  string
	HTMLBody string
	TextBody string
}
