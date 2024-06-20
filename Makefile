obu:
	@go build -o ./bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o ./bin/calculator ./distance_calculator
	@./bin/calculator

invoicer:
	go build -o ./bin/invoicer ./invoicer
	./bin/invoicer

agg:
	go build -o ./bin/agg ./aggregator
	./bin/agg

# Используем утилиту protoc для генерации go-файлов для протобафа 
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu invoicer