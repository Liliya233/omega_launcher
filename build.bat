SET CGO_ENABLED=0
    SET GOARCH=amd64
        :: Windows X86_64
        SET GOOS=windows
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_windows_amd64.exe

        :: Darwin X86_64
        SET GOOS=darwin
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_darwin_amd64

        :: Linux X86_64
        SET GOOS=linux
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_linux_amd64

    SET GOARCH=arm64
        :: Windows arm64
        SET GOOS=windows
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_windows_arm64.exe

        :: Darwin arm64
        SET GOOS=darwin
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_darwin_arm64

        :: Linux arm64
        SET GOOS=linux
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_linux_arm64

SET CGO_ENABLED=1
    SET GOOS=android
        :: Android arm64
        SET GOARCH=arm64
        SET CC=D:\android-ndk-r25c\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android33-clang
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_android_arm64

        :: Android X86_64
        SET GOARCH=amd64
        SET CC=D:\android-ndk-r25c\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android33-clang
        go build -ldflags="-w -s" -trimpath -o ./build/omega_launcher_android_amd64
