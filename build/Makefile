PROJECTNAME = go

## linux:
.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(RACE) -o ../plesk ../cmd/plesk/service.go

## windows:
.PHONY: win
win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(RACE) -o ../plesk.exe ../cmd/plesk/service.go

## clear:
.PHONY: clear
clear:
	rm -f ../plesk

## max
.PHONY: mac
mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(RACE) -o ../plesk ../cmd/plesk/service.go
