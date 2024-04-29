default: install

test:
	go test ./... -count=1

install: test
	go install cmd/pomo/pomo.go
	go install cmd/pomod/pomod.go
