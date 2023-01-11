package stdout

import (
	"errors"
	"fmt"
	"os"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/mailer"
)

var ErrEmptyRecipient = errors.New("empty recipient list")

type Config struct {
	TmplPath    string
	DefaultFrom string
}

type Mailer struct {
	config *Config
}

func NewMailer(config *Config) *Mailer {
	return &Mailer{config: config}
}

func (ml *Mailer) SendMail(tplName string, mail mailer.Mail) error {
	if len(mail.To) == 0 {
		return ErrEmptyRecipient
	}
	body, err := mailer.BuildBody(ml.config.TmplPath, tplName, mail.Data)
	if err != nil {
		return err
	}
	headers := mailer.BuildHeaders(mail)

	_, err = fmt.Fprintf(os.Stdout, "%s\r\n%s\n\n", headers, body)

	return err
}
