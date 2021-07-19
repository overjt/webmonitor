#!/bin/bash
#Based on https://github.com/michenriksen/aquatone/blob/master/build.sh
#Thanks to michenriksen

BUILD_FOLDER=builds
VERSION=$(cat banner.go | grep Version | cut -d '"' -f 2)

bin_dep() {
  BIN=$1
  which $BIN > /dev/null || { echo "[-] Dependency $BIN not found !"; exit 1; }
}

create_exe_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" webmonitor.exe ../README.md > /dev/null
  rm -rf webmonitor webmonitor.exe
}

create_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" webmonitor ../README.md > /dev/null
  rm -rf webmonitor webmonitor.exe
}

build_linux_amd64() {
  echo "[*] Building linux/amd64 ..."
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webmonitor -trimpath ..
}

build_linux_arm64() {
  echo "[*] Building linux/arm64 ..."
  CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o webmonitor -trimpath ..
}

build_macos_amd64() {
  echo "[*] Building darwin/amd64 ..."
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o webmonitor -trimpath ..
}

build_windows_amd64() {
  echo "[*] Building windows/amd64 ..."
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o webmonitor.exe -trimpath ..
}

rm -rf $BUILD_FOLDER
mkdir $BUILD_FOLDER
cd $BUILD_FOLDER

build_linux_amd64 && create_archive webmonitor_linux_amd64_$VERSION.zip
build_linux_arm64 && create_archive webmonitor_linux_arm64_$VERSION.zip
build_macos_amd64 && create_archive webmonitor_macos_amd64_$VERSION.zip
build_windows_amd64 && create_exe_archive webmonitor_windows_amd64_$VERSION.zip
shasum -a 256 * > checksums.txt

echo
echo
du -sh *

cd --
