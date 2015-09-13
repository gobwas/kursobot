package telegram

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"option"
)

type Message struct {
	MessageID int
}

type Update struct {
	UpdateID int
	Message  Message
}

const EndPoint = "%s/bot%s/%s"

type Bot struct {
	token    string
	endpoint string
	Updates  chan Update
}

func NewBot(token string) (*Bot, error) {
	return &Bot{token}, nil
}

type File struct {
	Name  string
	Bytes []byte
}

func (self *Bot) getEndPoint(method string) string {
	return fmt.Sprintf(EndPoint, self.endpoint, self.token, method)
}

func (self *Bot) Request(method string, params url.Values) {
	resp, err := http.DefaultClient.PostForm(self.getEndPoint(method), params)
	if err != nil {
		return err
	}
}

func (self *Bot) SendFile(method string, params map[string]string, field string, file *File) {
	var buffer bytes.Buffer
	bodyWriter := multipart.NewWriter(&buffer)
	defer bodyWriter.Close()

	fw, err := bodyWriter.CreateFormFile(field, file.Name)
	if err != nil {
		return err
	}
	fw.Write(file.Bytes)

	for key, value := range params {
		bodyWriter.WriteField(key, value)
	}

	req, err := http.NewRequest("POST", self.getEndPoint(method), &buffer)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	//	req.Header.Add("Content-Length", buffer.Len())

	http.DefaultClient.Do(req)
}

func (self *Bot) SetWebHook(url url.URL, cert *File) {
	if cert == nil {
		self.Request("setWebhook")
	}
}

func (self *Bot) SendMessage() {

}
