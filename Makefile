init:
	go mod init github.com/prachin77/pkr

tidy:
	go mod tidy

auth_pb_files:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/auth.proto
