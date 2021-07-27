
set GOARCH=amd64
set GOOS=windows
go build -o html2pdf.exe


set GOARCH=amd64
set GOOS=linux
go build -o html2pdf
