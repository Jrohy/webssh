package core

import (
	"github.com/pkg/sftp"
	"io"
	"mime/multipart"
)

func (sclient *SSHClient) CreateSftp() error {
	err := sclient.GenerateClient()
	if err != nil {
		return err
	}
	client, err := sftp.NewClient(sclient.Client)
	if err != nil {
		return err
	}
	sclient.Sftp = client
	return nil
}

func (sclient *SSHClient) Download(srcPath string) (*sftp.File, error) {
	return sclient.Sftp.Open(srcPath)
}

func (sclient *SSHClient) Upload(file multipart.File, dstPath string) error {
	dstFile, err := sclient.Sftp.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, file)
	if err != nil {
		return err
	}
	return nil
}

