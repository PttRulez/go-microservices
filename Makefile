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

.PHONY: obu invoicer