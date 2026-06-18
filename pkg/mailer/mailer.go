package mailer

type EmailMessage struct {
	To          []string
	Subject     string
	HTMLBody    string
	TextBody    string
	Attachments []Attachment
}

type Attachment struct {
	Filename string
	Content  []byte
}

type Mailer interface {
	SendEmail(msg *EmailMessage) error
	SendEmailAsync(msg *EmailMessage) error
}
