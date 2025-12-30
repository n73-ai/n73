package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func SendEmail(token string, email string, subject string) error {
	secretPassword := os.Getenv("EMAIL_SECRET_KEY")
	auth := smtp.PlainAuth(
		"",
		"agustfricke@gmail.com",
		secretPassword,
		"smtp.gmail.com",
	)

	path := os.Getenv("ROOT_PATH") + "/templates/email.html"
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return fmt.Errorf("Error to read template: %w.", err)
	}

	data := struct {
		Token string
	}{
		Token: token,
	}

	var bodyContent bytes.Buffer
	if err := tmpl.Execute(&bodyContent, data); err != nil {
		return fmt.Errorf("Error to execute template: %w.", err)
	}

	from := "agustfricke@gmail.com"
	emailContent := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n\r\n%s",
		from, // <---- ESTE CAMPO ES CLAVE
		email,
		subject,
		bodyContent.String(),
	)

	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"contact@zustack.com",
		[]string{email},
		[]byte(emailContent),
	); err != nil {
		return fmt.Errorf("Error sending email: %w.", err)
	}

	return nil
}
