package mailer

import (
	"bytes"
	"embed"
	gomail "github.com/go-mail/mail/v2"
	"html/template"
	"time"
)

//go:embed "templates"
var templateFs embed.FS

type Mailer struct {
	dailer *gomail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dailer: dialer,
		sender: sender,
	}
}
func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFs, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	///adding retry to the mail
	for i := 0; i < 3; i++ {
		err = m.dailer.DialAndSend(msg)
		if err == nil {
			return nil
		}
		time.Sleep(3 * time.Second)
	}

	return err
}
