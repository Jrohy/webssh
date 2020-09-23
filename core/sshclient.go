package core

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

func DecodedMsgToSSHClient(sshInfo string) (SSHClient, error) {
	client := NewSSHClient()
	decoded, err := base64.StdEncoding.DecodeString(sshInfo)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(decoded, &client)
	if err != nil {
		return client, err
	}
	return client, nil
}

func (sclient *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sclient.Password))
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    sclient.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", sclient.IpAddress, sclient.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	sclient.Client = client
	return nil
}

func (sclient *SSHClient) InitTerminal(terminal Terminal) *SSHClient {
	session, err := sclient.Client.NewSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	sclient.Session = session
	channel, inRequests, err := sclient.Client.OpenChannel("session", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	sclient.channel = channel
	go func() {
		for req := range inRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    terminal.Columns * 8,
		Height:   terminal.Columns * 8,
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	return sclient
}

func (sclient *SSHClient) Connect(ws *websocket.Conn, d time.Duration) {
	stopCh := make(chan struct{})
	//这里第一个协程获取用户的输入
	go func() {
		for {
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil {
				close(stopCh)
				return
			}
			_, err = sclient.channel.Write(p)
			if err != nil {
				close(stopCh)
				return
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		// 设置ws超时时间定时器
		stopTicker := time.NewTicker(d)
		defer stopTicker.Stop()

		br := bufio.NewReader(sclient.channel)
		var buf []byte
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			defer sclient.Client.Close()
			defer sclient.Session.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					log.Println(err)
					ws.WriteMessage(1, []byte("\033[33m已经关闭连接!\033[0m"))
					ws.Close()
					close(stopCh)
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		// 主循环
		for {
			select {
			case <- stopCh:
				return
			case <- stopTicker.C:
				ws.WriteMessage(1, []byte("\033[33m已超时关闭连接!\033[0m"))
				ws.Close()
				return
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}

func (sclient *SSHClient) ExecRemoteCommand(command string) (string, error) {
	//创建ssh登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sclient.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sclient.Password)}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sclient.IpAddress, sclient.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("创建ssh client 失败: ",err)
		return "", err
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println("创建ssh session 失败: ",err)
		return "", err
	}
	defer session.Close()
	//执行远程命令
	combo,err := session.CombinedOutput(command)
	if err != nil {
		fmt.Println("远程执行cmd 失败: ",err)
		return "", err
	}
	return string(combo), nil
}