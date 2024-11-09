run:
	go run ./cmd/campaign-producer-api/main.go

clean:
	go mod tidy

config-local:
	./config.sh local

config-sandbox:
	./config.sh sandbox

config-production:
	./config.sh production

.PHONY: build run compose-up compose-down compose-infra-down compose-infra-up
compose-infra-up:
	docker-compose -f ./docker-compose.yml --profile infra up -d
compose-infra-down:
	docker-compose -f ./docker-compose.yml --profile infra down

compose-up:
	docker-compose -f ./docs/docker-compose.yaml up -d

air:
	air

air-init:
	air init