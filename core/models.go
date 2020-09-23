package core

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type Terminal struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}

type SSHClient struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IpAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
	Session   *ssh.Session
	Client    *ssh.Client
	Sftp      *sftp.Client
	channel   ssh.Channel
}

func NewSSHClient() SSHClient {
	client := SSHClient{}
	client.Username = "root"
	client.Port = 22
	return client
}
