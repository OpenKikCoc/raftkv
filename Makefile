
BUILD_TAGS?=
BUILD_FLAGS = -ldflags "-X github.com/OpenKikCoc/raftkv/version.GitCommit=`git rev-parse HEAD`"

default: clean build

clean:
	rm -rf bin

build:
	go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o bin/raftkv ./cmd