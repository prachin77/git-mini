CLIENT_EXE = P:\git-mini\client\client.exe
SERVER_EXE = P:\git-mini\server\server.exe

TEST1 = P:\git-mini\test\palas
TEST2 = P:\git-mini\test\moit

init:
	go mod init github.com/prachin77/pkr

tidy:
	go mod tidy

pb_files:
	@protoc --proto_path=proto \
	       --go_out=paths=source_relative:pb \
	       --go-grpc_out=paths=source_relative:pb \
	       proto/background_service.proto

build: clean
	@echo Building the Go file...
	cd client &&  go build 
	cd server &&  go build
	
clean:
	@echo cleaning up...
	del "$(CLIENT_EXE)"
	del "$(SERVER_EXE)"

clean_test:
	@echo Cleaning test directories...
	@if exist "$(TEST1)\client.exe" del /s /q "$(TEST1)\client.exe" 2>nul
	@if exist "$(TEST1)\server.exe" del /s /q "$(TEST1)\server.exe" 2>nul
	@if exist "$(TEST2)\client.exe" del /s /q "$(TEST2)\client.exe" 2>nul
	@if exist "$(TEST2)\server.exe" del /s /q "$(TEST2)\server.exe" 2>nul

copy: build
	@echo copying executable files to test folders...
	@copy "$(SERVER_EXE)" "$(TEST2)"
	@copy "$(CLIENT_EXE)" "$(TEST2)"
	@copy "$(SERVER_EXE)" "$(TEST1)"
	@copy "$(CLIENT_EXE)" "$(TEST1)"

