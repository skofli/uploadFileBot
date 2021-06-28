B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

build: info
	- @mkdir -p dist
	- go build -ldflags "-X main.revision=$(REV) -s -w" -v ./...
	- mv uptotg dist/ && chmod 777 dist/uptotg

info:
	- @echo "revision $(REV)"
