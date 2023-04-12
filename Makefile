.PHONY: all linux-amd64 linux-arm64 macos-amd64 macos-arm64 android-amd64 android-arm64 windows-amd64 windows-arm64 
TARGETS:=build/ linux-amd64

ifneq ($(wildcard /usr/bin/aarch64-linux-gnu-gcc),)
	TARGETS:=${TARGETS} linux-arm64
endif
ifneq (${THEOS},)
	TARGETS:=${TARGETS} macos-amd64 macos-arm64
endif
ifneq ($(wildcard ${HOME}/android-ndk-r25c),)
	TARGETS:=${TARGETS} android-amd64 android-arm64
endif
ifneq ($(wildcard /usr/bin/x86_64-w64-mingw32-gcc),)
	TARGETS:=${TARGETS} windows-amd64 windows-arm64
endif

SRCS_GO := $(foreach dir, $(shell find . -type d), $(wildcard $(dir)/*.go $(dir)/*.c))

all: ${TARGETS}
linux-amd64: build/omega_launcher_linux_amd64
linux-arm64: build/omega_launcher_linux_arm64
macos-amd64: build/omega_launcher_darwin_amd64
macos-arm64: build/omega_launcher_darwin_arm64
android-amd64: build/omega_launcher_android_amd64
android-arm64: build/omega_launcher_android_arm64
windows-amd64: build/omega_launcher_windows_amd64.exe
windows-arm64: build/omega_launcher_windows_arm64.exe

build/:
	mkdir build
build/omega_launcher_linux_amd64: build/ ${SRCS_GO}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -o $@
build/omega_launcher_linux_arm64: build/ ${SRCS_GO}
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC=/usr/bin/aarch64-linux-gnu-gcc go build -ldflags="-w -s" -trimpath -o $@
build/omega_launcher_darwin_amd64: build/ ${SRCS_GO}
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 CC=`pwd`/archs/macos.sh go build -ldflags="-w -s" -trimpath -o $@
build/omega_launcher_darwin_arm64: build/ ${SRCS_GO}
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 CC=`pwd`/archs/macos.sh go build -ldflags="-w -s" -trimpath -o $@
build/omega_launcher_android_amd64: build/ ${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android33-clang ${SRCS_GO}
	CGO_ENABLED=1 GOOS=android GOARCH=amd64 CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android33-clang CXX=${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android33-clang++ go build -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o $@
build/omega_launcher_android_arm64: build/ ${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android33-clang ${SRCS_GO}
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CGO_LDFLAGS="-Wl,-rpath,/data/data/com.termux/files/usr/lib" CC=${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android33-clang CXX=${HOME}/android-ndk-r25c/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android33-clang++ go build -trimpath -ldflags "-s -w -extldflags -static-libstdc++" -o $@
build/omega_launcher_windows_amd64.exe: build/ /usr/bin/x86_64-w64-mingw32-gcc ${SRCS_GO}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -trimpath -ldflags "-s -w" -o $@
build/omega_launcher_windows_arm64.exe: build/ /usr/bin/x86_64-w64-mingw32-gcc ${SRCS_GO}
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -trimpath -ldflags "-s -w" -o $@

clean:
	rm -f build/omega_launcher_*
