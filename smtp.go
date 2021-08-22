package goutil

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

//SendEmail send a email of message to a list of user
func SendEmail(to []string, title string, message string) error {
	port := os.Getenv("SMTP_PORT")

	if port == "465" {
		return sendEmailSSL(to, title, message)
	} else if port == "587" {
		return sendEmailTLS(to, title, message)
	}
	log.Panic("smtp server is not set")
	return nil
}

func sendEmailTLS(to []string, title string, message string) error {
	from := os.Getenv("SMTP_ACCOUNT")

	password := os.Getenv("SMTP_PASSWORD")

	smtpServer := smtpServer{host: os.Getenv("SMTP_SERVER"), port: "587"}

	messageBytes := composeMimeMail(to[0], from, title, message)

	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	return smtp.SendMail(smtpServer.Address(), auth, from, to, messageBytes)
}

func sendEmailSSL(to []string, title string, message string) error {
	from := mail.Address{"", os.Getenv("SMTP_ACCOUNT")}
	rev := mail.Address{"", to[0]}
	//subj := "This is the email subject"

	// Setup headers
	//headers := make(map[string]string)
	//headers["From"] = from.String()
	//headers["To"] = rev.String()
	//headers["Subject"] = subj

	// Setup message
	//for k, v := range headers {
	//	message += fmt.Sprintf("%s: %s\r\n", k, v)
	//}
	//message +=  body
	bodyData := composeMimeMail(rev.String(), from.String(), title, message)

	// Connect to the SMTP Server
	servername := os.Getenv("SMTP_SERVER") + ":" + os.Getenv("SMTP_PORT")

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", os.Getenv("SMTP_ACCOUNT"), os.Getenv("SMTP_PASSWORD"), host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(rev.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(bodyData)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}

func getMXRecord(to string) (mx string, err error) {
	var e *mail.Address
	e, err = mail.ParseAddress(to)
	if err != nil {
		return
	}

	domain := strings.Split(e.Address, "@")[1]

	var mxs []*net.MX
	mxs, err = net.LookupMX(domain)

	if err != nil {
		return
	}

	for _, x := range mxs {
		mx = x.Host
		return
	}

	return
}

// Never fails, tries to format the address if possible
func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
