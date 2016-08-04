default: ec
.PHONY: ec

ec:
	go build -o bin/ec main.go server.go

