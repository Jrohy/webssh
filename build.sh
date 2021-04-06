#!/bin/bash

GITHUB_TOKEN=""

PROJECT="Jrohy/webssh"

#获取当前的这个脚本所在绝对路径
SHELL_PATH=$(cd `dirname $0`; pwd)

RELEASE_ID=`curl -H 'Cache-Control: no-cache' -s https://api.github.com/repos/$PROJECT/releases/latest|grep id|awk 'NR==1{print $2}'|sed 's/,//'`

function uploadfile() {
    FILE=$1

    CTYPE=$(file -b --mime-type $FILE)

    curl -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: ${CTYPE}" --data-binary @$FILE "https://uploads.github.com/repos/$PROJECT/releases/${RELEASE_ID}/assets?name=$(basename $FILE)"

    echo ""
}

function upload() {
    FILE=$1
    DGST=$1.dgst
    openssl dgst -md5 $FILE | sed 's/([^)]*)//g' >> $DGST
    openssl dgst -sha1 $FILE | sed 's/([^)]*)//g' >> $DGST
    openssl dgst -sha256 $FILE | sed 's/([^)]*)//g' >> $DGST
    openssl dgst -sha512 $FILE | sed 's/([^)]*)//g' >> $DGST
    uploadfile $FILE
    uploadfile $DGST
}

VERSION=`git describe --tags $(git rev-list --tags --max-count=1)`
NOW=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`
GO_VERSION=`go version|awk '{print $3,$4}'`
GIT_VERSION=`git rev-parse HEAD`
LDFLAGS="-w -s -X 'main.version=$VERSION' -X 'main.buildDate=$NOW' -X 'main.goVersion=$GO_VERSION' -X 'main.gitVersion=$GIT_VERSION'"

GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o result/webssh_windows_amd64.exe .
GOOS=windows GOARCH=386 go build -ldflags "$LDFLAGS" -o result/webssh_windows_386.exe .
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o result/webssh_linux_amd64 .
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o result/webssh_linux_arm64 .
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o result/webssh_darwin_amd64 .

cd result

UPLOAD_ITEM=($(ls -l|awk '{print $9}'|xargs -r))

for ITEM in ${UPLOAD_ITEM[@]}
do
   upload $ITEM
done

echo "upload completed!"

cd $SHELL_PATH

rm -rf result
