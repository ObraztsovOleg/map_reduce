BUILD_DIR="../bin"

build:
	go mod init server.go
	go get github.com/gorilla/mux
	go get github.com/wcharczuk/go-chart/v2
	go build -o ${BUILD_DIR}/server server.go

clean:
	rm -rf go.sum go.mod .${BUILD_DIR}/server