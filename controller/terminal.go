package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
	"webssh/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TermWs(c *gin.Context, d time.Duration) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sshInfo := c.DefaultQuery("sshInfo", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := core.Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	sshClient, err := core.DecodedMsgToSSHClient(sshInfo)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	sshClient.InitTerminal(terminal)
	sshClient.Connect(conn, d)
	return &responseBody
}
