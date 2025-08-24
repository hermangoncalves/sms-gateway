APP_NAME=smsgateway

build:
	@go build -o bin/$(APP_NAME) ./smsgateway/cmd/main.go

build-android:
    # 64-bit ARM (arm64, comum em dispositivos mais novos)
    @cd smsgateway && GOOS=android GOARCH=arm64 go build -o bin/$(APP_NAME)-android-arm64 ./cmd/main.go

    # 32-bit ARM (arm, para dispositivos mais antigos e alguns emuladores)
    @cd smsgateway && GOOS=android GOARCH=arm go build -o bin/$(APP_NAME)-android-armv7 ./cmd/main.go
    
    # 32-bit x86 (x86, para dispositivos Intel x86 e emuladores)
    @cd smsgateway && GOOS=android GOARCH=386 go build -o bin/$(APP_NAME)-android-386 ./cmd/main.go
    
    # 64-bit x86 (amd64, para emuladores x86_64)
    @cd smsgateway && GOOS=android GOARCH=amd64 go build -o bin/$(APP_NAME)-android-x86_64 ./cmd/main.go

    # Adicionando o binário Linux para completar
    @cd smsgateway && GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)-linux-arm64 ./cmd/main.go

clean:
	@rm -rf bin