gen:
	protoc -I=proto --go_out=:pb --go-grpc_out=:pb proto/*.proto

clean:
	rm pb/*.go

run:
	go run main.go
