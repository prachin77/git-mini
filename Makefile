init:
	go mod init github.com/prachin77/pkr

tidy:
	go mod tidy

pb_files:
	protoc --proto_path=proto \
	       --go_out=paths=source_relative:pb \
	       --go-grpc_out=paths=source_relative:pb \
	       proto/background_service.proto