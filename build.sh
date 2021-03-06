#!/bin/bash

## Auto Build

if [ $# -ne 1 ]
then
	echo "Usage: bash build.sh <version>"
	exit 1
fi

version=$1
arch=amd64
chksum=wordle_${version}_checksums.txt

buildBinary () {
	local name=$1
	for sys in "darwin" "linux" "windows"; do 
		printf "sys: %s arch: %s" $sys $arch
		echo ""
		GOOS=$sys GOARCH=amd64 go build -o ${name} main.go
		tar czvf ${name}_${version}_${sys}_${arch}.tar.gz $name
		rm $name
	done
}

for name in "wordle"; do
	buildBinary $name
	shasum -a 256 ${name}_${version}_*.tar.gz >> ${chksum}
done
