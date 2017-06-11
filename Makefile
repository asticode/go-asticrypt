server-all: server-build server-init server-migrate server-run

server-build:
	go build -o ./server/server ./server

server-init:
	./server/server db-init -c ./server/local.toml -v

server-migrate:
	./server/server db-migrate -c ./server/local.toml -v

server-rollback:
	./server/server db-rollback -c ./server/local.toml -v

server-run:
	./server/server -c ./server/local.toml -v