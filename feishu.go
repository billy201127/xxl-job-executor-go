package xxl

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TextMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

type CardMsg struct {
	MsgType   string `json:"msg_type"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	Card      struct {
		Config struct {
			WideScreenMode bool `json:"wide_screen_mode"`
			EnableForward  bool `json:"enable_forward"`
		} `json:"config"`
		Header struct {
			Template string `json:"template"`
			Title    struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"title"`
		}
		Elements []Element `json:"elements"`
	} `json:"card"`
}

type Element struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func SendCardMsg(webhook, secret, title, content string, isAtAll bool) error {
	if webhook == "" || secret == "" {
		return errors.New("invalid config")
	}

	tt := time.Now().Unix()
	secretStr, _ := Sign(secret, tt)
	msg := CardMsg{
		MsgType:   "interactive",
		Timestamp: strconv.FormatInt(tt, 10),
		Sign:      secretStr,
	}

	msg.Card.Config.EnableForward = true
	msg.Card.Config.WideScreenMode = true

	msg.Card.Header.Title.Tag = "plain_text"
	msg.Card.Header.Title.Content = title
	msg.Card.Header.Template = "blue"
	if isAtAll {
		msg.Card.Header.Template = "red"
		content += `<at user_id="all">所有人</at>`
	}

	hostname, _ := os.Hostname()
	content = fmt.Sprintf("hostname: [%s]\n%s", hostname, content)
	element := Element{
		Tag:     "markdown",
		Content: content,
	}
	msg.Card.Elements = append(msg.Card.Elements, element)

	data, _ := json.Marshal(msg)
	request, err := http.NewRequest("POST", webhook, bytes.NewReader(data))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Printf("err:%v, resp:%v", err, resp)
	}
	return err
}

func SendTextMsg(webhook, secret, content string, isAtAll bool) error {
	if webhook == "" || secret == "" {
		return errors.New("invalid config")
	}

	tt := time.Now().Unix()
	secretStr, _ := Sign(secret, tt)
	msg := TextMsg{
		MsgType:   "text",
		Timestamp: strconv.FormatInt(tt, 10),
		Sign:      secretStr,
	}

	if isAtAll {
		content += `<at user_id="all">所有人</at>`
	}
	msg.Content.Text = content

	data, _ := json.Marshal(msg)
	request, err := http.NewRequest("POST", webhook, bytes.NewReader(data))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Printf("err:%v, resp:%v", err, resp)
	}
	return err
}

func Sign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
