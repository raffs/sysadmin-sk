default: build

# BUILD_FLAGS = -ldflags ${GO_LDFLAGS}
BUILD_FLAGS += -installsuffix cgo
BUILD_FLAGS += -o bin/sysadmin-sk

build:
	build/build.sh
