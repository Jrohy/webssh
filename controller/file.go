package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"webssh/core"

	"github.com/gin-gonic/gin"
)

// File 结构体
type File struct {
	Name       string
	Size       string
	ModifyTime string
	// 0: 文件  1: 文件夹
	FType uint8
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) *ResponseBody {
	var (
		sshClient core.SSHClient
		err       error
	)
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	path := strings.TrimSpace(c.DefaultPostForm("path", "/root"))
	sshInfo := c.PostForm("sshInfo")
	if sshClient, err = core.DecodedMsgToSSHClient(sshInfo); err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	if err := sshClient.CreateSftp(); err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		responseBody.Msg = err.Error()
		return &responseBody
	}
	defer file.Close()
	filename := header.Filename
	if path[len(path)-1:] == "/" {
		path = path + filename
	} else {
		path = path + "/" + filename
	}
	err = sshClient.Upload(file, path)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
	}
	return &responseBody
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) *ResponseBody {
	var (
		sshClient core.SSHClient
		err       error
	)
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	path := strings.TrimSpace(c.DefaultQuery("path", "/root"))
	sshInfo := c.DefaultQuery("sshInfo", "")
	if sshClient, err = core.DecodedMsgToSSHClient(sshInfo); err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	if err := sshClient.CreateSftp(); err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	if sftpFile, err := sshClient.Download(path); err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
	} else {
		defer sftpFile.Close()
		c.Writer.WriteHeader(http.StatusOK)
		fileMeta := strings.Split(path, "/")
		c.Header("Content-Disposition", "attachment; filename="+fileMeta[len(fileMeta)-1])
		_, _ = io.Copy(c.Writer, sftpFile)
	}
	return &responseBody
}

// FileList 获取文件列表
func FileList(c *gin.Context) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	path := c.DefaultQuery("path", "/root")
	sshInfo := c.DefaultQuery("sshInfo", "")
	sshClient, err := core.DecodedMsgToSSHClient(sshInfo)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	result, err := sshClient.ExecRemoteCommand(fmt.Sprintf("ls -lh %s|grep -v ^l|awk 'NR>1'|sort", path))
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	if strings.Contains(result, "ls:") {
		fmt.Println(result)
		responseBody.Msg = result
		return &responseBody
	}
	lineList := strings.Split(result, "\n")
	var fileList []*File
	for _, line := range lineList {
		if line == "" {
			continue
		}
		//多个空格合并为一个
		line = strings.Join(strings.Fields(line), " ")
		metaList := strings.Split(line, " ")
		file := File{Name: metaList[len(metaList)-1], FType: 0, Size: metaList[4], ModifyTime: strings.Join(metaList[5:len(metaList)-1], " ")}
		if strings.Contains(metaList[0], "d") {
			file.FType = 1
		}
		fileList = append(fileList, &file)
	}
	responseBody.Data = map[string]interface{}{
		"path": path,
		"list": fileList,
	}
	return &responseBody
}
