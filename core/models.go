package core

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"sync"
)

// SSHClient 结构体
type SSHClient struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	IPAddress   string `json:"ipaddress"`
	Port        int    `json:"port"`
	Client      *ssh.Client
	Sftp        *sftp.Client
	StdinPipe   io.WriteCloser
	Session     *ssh.Session
	ComboOutput *wsBufferWriter
}

type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// Write: implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

// NewSSHClient 返回默认ssh信息
func NewSSHClient() SSHClient {
	client := SSHClient{}
	client.Username = "root"
	client.Port = 22
	return client
}
