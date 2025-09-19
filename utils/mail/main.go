package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type Sender struct {
	host, email, password, port, username string

	auth smtp.Auth
}

type Message struct {
	To          []string
	CC          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func New(username string) *Sender {
	host := os.Getenv("HOST")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	port := os.Getenv("PORT")
	auth := smtp.PlainAuth("", email, password, host)
	return &Sender{auth: auth, host: host, password: password, email: email, port: port, username: username}
}

func (s *Sender) Send(m *Message, host, port, username string) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, port), s.auth, username, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) AttachFileByte(fileName string, src []byte) error {
	m.Attachments[fileName] = src
	return nil
}
func (m *Message) ToBytes() []byte {
	messageBuffer := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	messageBuffer.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	messageBuffer.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		messageBuffer.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	messageBuffer.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(messageBuffer)
	boundary := writer.Boundary()
	// messageBuffer.WriteString("Content-Type: text/html; charset=utf-8\n")
	// messageBuffer.WriteString("Content-Transfer-Encoding: 7bit\r\n")

	if withAttachments {
		messageBuffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		messageBuffer.WriteString(fmt.Sprintf("--%s\n", boundary))
	}
	messageBuffer.WriteString(m.Body)

	if withAttachments {
		for k, v := range m.Attachments {
			messageBuffer.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			messageBuffer.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			messageBuffer.WriteString("Content-Transfer-Encoding: base64\n")
			messageBuffer.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			messageBuffer.Write(b)
			messageBuffer.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		messageBuffer.WriteString("--")
	}

	return messageBuffer.Bytes()
}

// var Emails = map[string][]string{
// 	"Golang":  []string{"hamzaabdellaoui26648999@gmail.com"},
// 	"Android": []string{"mohamedrejeb445@gmail.com"},
// 	"AI&ML":   []string{"donieztouil77@gmail.com"},
// 	"Web":     []string{"lina2.haouas@gmail.com", "abirlakhal20@gmail.com"},
// 	"Cloud":   []string{"ibty.hattab@gmail.com"},
// 	"Flutter": []string{"tn.mohamedbenhalima@gmail.com"},
// }

func (obj Sender) SendMail(qr []byte) error {
	m := NewMessage("GDG Devfest 2025 Days Registration",
		"Dear participant, \n Thank you for your participation in  Devfest 2025")
	m.To = []string{obj.username}

	m.AttachFileByte("qr.png", qr)
	return obj.Send(m, obj.host, obj.port, obj.username)
}
