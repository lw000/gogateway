set prjPath=%cd%
echo %prjPath%
cd ../../../
set GOPATH=%cd%
set GOARCH=amd64
set GOOS=linux
cd %prjPath%
go build -a -o gateway -v -ldflags="-s -w"