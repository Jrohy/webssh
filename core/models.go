package core

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"unicode/utf8"
)

// WcList 全局counter list变量
var WcList []*WriteCounter

// WriteCounter 结构体
type WriteCounter struct {
	Total int
	Id    string
}

// Write: implement Write interface to write bytes from ssh server into bytes.Buffer.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += n
	return n, nil
}

type wsOutput struct {
	ws *websocket.Conn
}

// Write: implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *wsOutput) Write(p []byte) (int, error) {
	// 处理非utf8字符
	if !utf8.Valid(p) {
		bufStr := string(p)
		buf := make([]rune, 0, len(bufStr))
		for _, r := range bufStr {
			if r == utf8.RuneError {
				buf = append(buf, []rune("@")...)
			} else {
				buf = append(buf, r)
			}
		}
		p = []byte(string(buf))
	}
	err := w.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

// SSHClient 结构体
type SSHClient struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IPAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
	LoginType int    `json:"logintype"`
	Client    *ssh.Client
	Sftp      *sftp.Client
	StdinPipe io.WriteCloser
	Session   *ssh.Session
}

// NewSSHClient 返回默认ssh信息
func NewSSHClient() SSHClient {
	client := SSHClient{}
	client.Username = "root"
	client.Port = 22
	return client
}

// Close all closable fields of SSHClient that is opened:
//
//	StdinPipe, Session, Sftp, Client
func (sclient *SSHClient) Close() {
	defer func() { // just in case
		if err := recover(); err != nil {
			log.Println("SSHClient Close recover from panic: ", err)
		}
	}()

	if sclient.StdinPipe != nil {
		sclient.StdinPipe.Close()
		sclient.StdinPipe = nil
	}
	if sclient.Session != nil {
		sclient.Session.Close()
		sclient.Session = nil
	}
	if sclient.Sftp != nil {
		sclient.Sftp.Close()
		sclient.Sftp = nil
	}
	if sclient.Client != nil {
		sclient.Client.Close()
		sclient.Client = nil
	}
}
