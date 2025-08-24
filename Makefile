APP_NAME=smsgateway

build:
	@go build -o bin/$(APP_NAME) ./smsgateway/cmd/main.go

build-android:
	@cd smsgateway && GOOS=android GOARCH=arm64 go build -o bin/$(APP_NAME)-android ./cmd/main.go
	@cd smsgateway && GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)-linux ./cmd/main.go

clean:
	@rm -rf bin