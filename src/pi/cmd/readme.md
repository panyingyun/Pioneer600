//build Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 

//build arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build 

//build arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build 