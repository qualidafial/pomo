default: test install

test:
	go test ./... -count=1

install:
	go install cmd/pomo/pomo.go
