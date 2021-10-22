package goutil

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

func NewSmtpServer(host string, port string, account string, password string) *SmtpServer {
	return &SmtpServer{
		host:           host,
		port:           port,
		auth:           smtp.PlainAuth("", account, password, host),
		emailAddresses: map[string]string{},
	}
}

type SmtpServer struct {
	host           string
	port           string
	auth           smtp.Auth
	emailAddresses map[string]string
}

func (s *SmtpServer) Address() string {

	return s.host + ":" + s.port
}

func (s *SmtpServer) AddEmailAddress(name string, emailAddress string) error {

	address, err := FormatEmailAddress(emailAddress)

	if err != nil {

		return err
	}
	s.emailAddresses[address] = name

	return nil
}

func (s *SmtpServer) SendEmail(from string, to string, title string, message string) error {

	fromTitle := ""
	revTitle := ""

	if testTitle, ok := s.emailAddresses[from]; ok {
		fromTitle = testTitle
	}

	fromAddress := mail.Address{
		Name:    fromTitle,
		Address: from,
	}

	if testTitle, ok := s.emailAddresses[to]; ok {

		fromTitle = testTitle
	}

	toAddress := mail.Address{
		Name:    revTitle,
		Address: to,
	}

	messageBytes := ComposeMimeMail(toAddress.String(), fromAddress.String(), title, message)

	if s.port == "587" {

		return smtp.SendMail(s.Address(), s.auth, from, []string{to}, messageBytes)

	} else if s.port == "465" {

		//TODO: security
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.host,
		}
		conn, err := tls.Dial("tcp", s.Address(), tlsconfig)
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return err
		}

		if err = c.Auth(s.auth); err != nil {
			return err
		}

		if err = c.Mail(fromAddress.Address); err != nil {
			return err
		}

		if err = c.Rcpt(toAddress.Address); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write(messageBytes)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		c.Quit()
	}
	return nil
}

func GetMXRecord(to string) (mx string, err error) {

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
func FormatEmailAddress(addr string) (string, error) {

	e, err := mail.ParseAddress(addr)
	if err != nil {

		return "", err
	}
	return e.String(), nil
}

func EncodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func ComposeMimeMail(to string, from string, subject string, body string) []byte {

	header := make(map[string]string)
	fromAddress, _ := FormatEmailAddress(from)
	toAddress, _ := FormatEmailAddress(to)

	header["From"] = fromAddress
	header["To"] = toAddress
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
