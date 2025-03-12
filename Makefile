# Copyright (c) 2025, Julian Huhn
#
# Permission to use, copy, modify, and/or distribute this software for any
# purpose with or without fee is hereby granted, provided that the above
# copyright notice and this permission notice appear in all copies.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

BINARY_NAME=godmarc
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/godmarc/main.go

.PHONY: build
build:
	go build -o $(BINARY_PATH) $(MAIN_PATH)

.PHONY: install
install: build
	cp $(BINARY_PATH) /usr/local/sbin/$(BINARY_NAME)
	chown root:bin /usr/local/sbin/$(BINARY_NAME)
	chmod 755 /usr/local/sbin/$(BINARY_NAME)

.PHONY: uninstall
uninstall: 
	rm -f /usr/local/sbin/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_PATH)

.PHONY: clean
clean:
	rm -f $(BINARY_PATH)
	rm -rf bin/

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run
