APP_NAME=smsgateway

build:
	@go build -o bin/$(APP_NAME) ./smsgateway/cmd/main.go

build-android:
	@GOOS=android GOARCH=arm64 go build -o bin/$(APP_NAME)-android ./smsgateway/cmd/main.go
	@GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)-android ./smsgateway/cmd/main.go

clean:
	@rm -rf bin