dev-run :
	nodemon --exec "go run" cmd/main.go --signal SIGTERM

mockery :
	mockery --all