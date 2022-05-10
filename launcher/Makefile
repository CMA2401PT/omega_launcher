.PHONY: all current 
TARGETS:=build/ current current-v8
PACKAGETARGETS:=
ifeq ($(shell uname | grep "Darwin" > /dev/null ; echo $${?}),0)
	ifeq ($(shell uname -m | grep -E "iPhone|iPad|iPod" > /dev/null ; echo $${?}),0)
		IOS_STRIP=/usr/bin/strip
		LIPO=/usr/bin/lipo
		LDID=/usr/bin/ldid
		TARGETS:=${TARGETS} ios-executable ios-v8-executable ios-lib
	else
		IOS_STRIP=$(shell xcrun --sdk iphoneos -f strip)
		IOS_OBJCOPY=$(shell xcrun --sdk iphoneos -f objcopy)
		LDID=ldid2
		LIPO=/usr/bin/lipo
		TARGETS:=${TARGETS} macos ios-v8-executable ios-executable ios-lib
	endif
		PACKAGETARGETS:=${PACKAGETARGETS} package/ios
else
	IOS_STRIP=true
	LDID=$${THEOS}/toolchain/linux/iphone/bin/ldid
	LIPO=$${THEOS}/toolchain/linux/iphone/bin/lipo
	IOS_OBJCOPY=$${THEOS}/toolchain/linux/iphone/bin/llvm-objcopy
endif 

VERSION=$(shell cat version)

SRCS_GO := $(foreach dir, $(shell find . -type d), $(wildcard $(dir)/*.go $(dir)/*.c))

all: linux macos windows
linux: build/launcher-linux
macos: build/launcher-macos
windows: build/launcher-windows.exe

package: ${PACKAGETARGETS}
release/:
	mkdir -p release
build/:
	mkdir build
build/launcher-linux: build/ ${SRCS_GO}
	CGO_CFLAGS=${CGO_DEF} CGO_ENABLED=1  go build -trimpath -ldflags "-s -w" -o build/launcher-linux
build/launcher-macos: build/ ${SRCS_GO}
	CGO_CFLAGS=${CGO_DEF} CC=`pwd`/macos.sh CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o build/launcher-macos
build/launcher-windows.exe: build/ /usr/bin/x86_64-w64-mingw32-gcc ${SRCS_GO}
	CGO_CFLAGS=${CGO_DEF} CC=/usr/bin/x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o build/launcher-windows.exe

clean:
	rm -f build/omega_launcher*
