prefix  := /usr/local
bindir  := $(prefix)/bin
datadir := $(prefix)/share

GO_LDFLAGS := -w -s
GO := go
STATICCHECK := staticcheck

EXE =
ifeq ($(GOOS),windows)
EXE = .exe
endif

.PHONY: build
build: ../bin/ocala-language-server$(EXE)

file/parser.g.go: file/parser.llpg.go
	../tools/llpg.rb $^

.PHONY: lint
lint: file/parser.g.go
	$(STATICCHECK) ./...

.PHONY: test
test: lint
	$(GO) test -count=1 -cover -coverprofile=coverage.out ./...

.PHONY: clean
clean:
	rm -f file/parser.g.go

.PHONY: ../bin/ocala-language-server$(EXE)
../bin/ocala-language-server$(EXE): lint
	$(GO) build -o $@ -trimpath -ldflags="$(GO_LDFLAGS)" .
