package utils

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

var smtpSever string = "..."
var smtpPort int = 25
var smtpUsername string = "..."
var smtpPassword string = "..."

func SendEmail(from string, to string, name string) {
	message := mail.NewMsg()
	err := message.From(from)
	if err != nil {
		fmt.Println("Error setting from:", err)
	}
	err = message.To(to)
	if err != nil {
		fmt.Println("Error setting to:", err)
	}

	// set the subject of the email
	message.Subject("This is a test email with go-mail!")
	message.SetBodyString(mail.TypeTextHTML, fmt.Sprintf("<h1>Hello %s!</h1><p>This is a test email sent with go-mail.</p>", name))

	// smtp details
	client, err := mail.NewClient(smtpSever, 
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(smtpUsername), mail.WithPassword(smtpPassword),
	)
	if err != nil {
		fmt.Println("Error creating client:", err)
	}

	// send the email
	err = client.DialAndSend(message)
	if err != nil {
		fmt.Println("Error sending email:", err)
	}

	fmt.Println("Email sent successfully!")
	// close the connection
	client.Close()
}