package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webssh/core"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TermWs 获取终端ws
func TermWs(c *gin.Context, timeout time.Duration) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sshInfo := c.DefaultQuery("sshInfo", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	sshClient, err := core.DecodedMsgToSSHClient(sshInfo)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	err = sshClient.GenerateClient()
	if err != nil {
		wsConn.WriteMessage(1, []byte(err.Error()))
		wsConn.Close()
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	sshClient.InitTerminal(wsConn, row, col)
	sshClient.Connect(wsConn, timeout)
	return &responseBody
}
