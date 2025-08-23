APP_NAME=smsgateway

build:
	@go build -o bin/$(APP_NAME) ./client/cmd/main.go

build-android:
	@GOOS=android GOARCH=arm64 go build -o bin/$(APP_NAME)-android-arm64 ./client/cmd/main.go

clean:
	@rm -rf bin