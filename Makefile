SRC=./cmd/bitbucket-pipeline-runner

test : mods
	go test -v ./...

dist : test
	go build -o dist/bpr $(SRC)

run : test
	go run $(SRC)

mods: 
	go mod download
