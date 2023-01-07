package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

// Mail структура почтового сообщения.
type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Data    interface{}
}

type Mailer interface {
	SendMail(tplName string, mail Mail) error
}

func BuildBody(tplPath, tplName string, data interface{}) (string, error) {
	fileName := fmt.Sprintf("%s/%s.mail.tmpl", tplPath, tplName)
	stats, err := os.Stat(fileName)
	if err != nil {
		return "", fmt.Errorf("template file for %s not found", tplName)
	}
	if stats.IsDir() || !stats.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not regular file", tplName)
	}
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return "", fmt.Errorf("error parsing template file %s: %w", fileName, err)
	}
	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, data)
	if err != nil {
		return "", fmt.Errorf("error executing template %s: %w", fileName, err)
	}
	return buff.String(), nil
}
func BuildHeaders(mail Mail) string {
	sb := strings.Builder{}
	sb.WriteString("MIME-version: 1.0;\r\n")
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	sb.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	if len(mail.To) > 0 {
		sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	}
	if len(mail.Cc) > 0 {
		sb.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";")))
	}
	if len(mail.Bcc) > 0 {
		sb.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(mail.Bcc, ";")))
	}
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	return sb.String()
}
