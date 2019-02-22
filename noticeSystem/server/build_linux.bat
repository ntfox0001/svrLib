SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o watrSvr
del watrSvr.zip
7z a watrSvr.zip watrSvr
del watrSvr