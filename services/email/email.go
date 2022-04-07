package email

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	log "github.com/sirupsen/logrus"
)

type emailContent struct {
	SenderName       string
	UserFirstName    string
	ConfirmationLink string
}

func GetConfirmationEmailContent(userFirstName, confirmationLink string) (string, error) {
	templateBytes, err := ioutil.ReadFile("/../confirmation_email_template.html")
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read email confirmation template")
		return "", err
	}
	t := template.Must(template.New("confirmationEmailContentTemplate").Parse(string(templateBytes)))
	buf := &bytes.Buffer{}

	if err := t.Execute(buf, &emailContent{UserFirstName: userFirstName, ConfirmationLink: confirmationLink, SenderName: os.Getenv("EMAIL_SENDER_NAME")}); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not form email confirmation link")
		return "", err
	}
	return buf.String(), nil
}

func Send(userFirstName, userEmail, subject, content string) error {
	from := mail.NewEmail(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"))
	to := mail.NewEmail(userFirstName, userEmail)
	htmlContent := content
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.WithFields(log.Fields{"from": from, "to": to, "subject": subject, "err": err.Error()}).Error("could not send email")
		return err
	}

	log.WithFields(log.Fields{"email_sent_status_code": response.StatusCode, "email_sent_response_body": response.Body, "from": from, "to": to, "subject": subject}).Info("sent email successfully")
	return nil
}
