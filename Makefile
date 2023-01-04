GO ?= go
TESTFOLDER = $(shell find * -maxdepth 10 -type d | grep /tests | xargs -I {} echo "newswav/http-server-sample/{}")
SOURCE_FILES = $(shell find * -maxdepth 10 | grep .go$)


test:
	echo "mode: count" > coverage.out
	for d in $(TESTFOLDER); do \
		$(GO) test -v -covermode=count -coverprofile=profile.out $$d > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "setup failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

setup:
	npm init -y && npm install chokidar tree-kill

dev:
	echo $(SOURCE_FILES) | node watch.cjs go run cmd/main/main.go                                                          

