swag:
	swag init -g internal/app/app.go

integrationtestdb:
	docker run --rm -d  -p 27017:27017/tcp mongo:4.4.10

tests:
	go test github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1 -v -run ^\QTestHandler\E$