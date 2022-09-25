BUILD_DIR="../bin"

build:
	go mod init agent.go
	go get github.com/chrislusf/glow/flow
	go build -o ${BUILD_DIR}/agent agent.go

clean:
	rm -rf go.sum go.mod .${BUILD_DIR}/agent