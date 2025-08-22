prefix  := /usr/local
bindir  := $(prefix)/bin
datadir := $(prefix)/share

GO_LDFLAGS := -w -s
DESTDIR :=
INSTALL := install
INSTALL_PROGRAM := $(INSTALL)
INSTALL_DATA := $(INSTALL) -m 644
GO := go
TINYGO := tinygo
STATICCHECK := staticcheck
GIT := git

EXE =
ifeq ($(GOOS),windows)
EXE = .exe
endif

.PHONY: all
all: build

internal/z80/z80.g.go: internal/z80/z80.lisp
	./tools/generate_arch.rb arch $^

internal/mos6502/mos6502.g.go: internal/mos6502/mos6502.lisp
	./tools/generate_arch.rb arch $^

internal/tt/ttarch/ttarch.g.go: internal/tt/ttarch/ttarch.lisp
	./tools/generate_arch.rb arch $^

internal/core/parser.g.go: internal/core/parser.llpg.go
	./tools/llpg.rb $^

internal/core/tabs.g.go: internal/core/functions.go
	./tools/mktabs.rb tabs $^

share/ocala/wasm/ocala.json: share/ocala/include/*.oc \
	examples/z80/msx-hello-world/main.oc \
	examples/z80/msx-hello-world-vdp/main.oc \
	internal/z80/testdata/op*.oc internal/mos6502/testdata/op*.oc
	./tools/mktabs.rb ocalajson $^ $@

.PHONY: lint
lint: internal/z80/z80.g.go \
		internal/mos6502/mos6502.g.go \
		internal/tt/ttarch/ttarch.g.go \
		internal/core/tabs.g.go \
		internal/core/parser.g.go
	$(STATICCHECK) ./cmd/ocala ./internal/... && \
	GOOS=js GOARCH=wasm $(STATICCHECK) ./cmd/ocala-wasm

.PHONY: build
build: bin/ocala$(EXE)

.PHONY: bin/ocala$(EXE)
bin/ocala$(EXE): lint
	$(GO) build -o $@ -trimpath -ldflags="$(GO_LDFLAGS)" ./cmd/ocala

.PHONY: wasm
wasm: lint share/ocala/wasm/ocala.json
	GOOS=js GOARCH=wasm $(GO) build -o share/ocala/wasm/ocala.wasm -trimpath ./cmd/ocala-wasm

.PHONY: wasm-tinygo
wasm-tinygo: lint share/ocala/wasm/ocala.json
	GOOS=js GOARCH=wasm $(TINYGO) build -o share/ocala/wasm/ocala.wasm \
		-no-debug -scheduler=none -panic=trap ./cmd/ocala-wasm

.PHONY: test
test: lint
	$(GO) test -count=1 -cover -coverprofile=coverage.out ./cmd/ocala ./internal/...

.PHONY: testdata
testdata:
	./tools/generate_arch.rb testdata internal/z80/z80.lisp
	pasmo internal/z80/testdata/opcodes.asm internal/z80/testdata/opcodes.dat
	pasmo internal/z80/testdata/opcodes_undocumented.asm internal/z80/testdata/opcodes_undocumented.dat
	pasmo -8 internal/z80/testdata/opcodes_compat8080.asm internal/z80/testdata/opcodes_compat8080.dat
	pasmo internal/z80/testdata/operators.asm internal/z80/testdata/operators.dat
	pasmo internal/z80/testdata/operators_undocumented.asm internal/z80/testdata/operators_undocumented.dat
	pasmo -8 internal/z80/testdata/operators_compat8080.asm internal/z80/testdata/operators_compat8080.dat
	./tools/generate_arch.rb testdata internal/mos6502/mos6502.lisp
	ca65 -o internal/mos6502/testdata/opcodes.o internal/mos6502/testdata/opcodes.asm
	ld65 -t none -o internal/mos6502/testdata/opcodes.dat internal/mos6502/testdata/opcodes.o
	ca65 -o internal/mos6502/testdata/operators.o internal/mos6502/testdata/operators.asm
	ld65 -t none -o internal/mos6502/testdata/operators.dat internal/mos6502/testdata/operators.o

.PHONY: clean
clean:
	rm -f internal/z80/z80.g.go internal/mos6502/mos6502.g.go internal/tt/ttarch/ttarch.g.go \
		internal/core/parser.g.go internal/core/tabs.g.go \
		bin/ocala share/ocala/wasm/ocala.wasm share/ocala/wasm/ocala.json

.PHONY: install
install: build
	$(INSTALL_PROGRAM) -d $(DESTDIR)$(bindir)
	$(INSTALL_PROGRAM) -m 755 -s -t $(DESTDIR)$(bindir) bin/ocala$(EXE)
	$(INSTALL_PROGRAM) -d $(DESTDIR)$(datadir)/ocala
	$(GIT) archive HEAD:share/ocala | tar xf - -C $(DESTDIR)$(datadir)/ocala
