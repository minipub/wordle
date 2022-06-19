#!/bin/bash

## Compile & Tag & Release

if [ $# -ne 1 ]
then
	echo "Usage: bash release.sh <version>"
	exit 1
fi

version=$1
arch=amd64
chksum=wordle-bundle_${version}_checksums.txt

buildBinary () {
	local name=$1
	for sys in "darwin" "linux" "windows"; do 
		printf "sys: %s arch: %s" $sys $arch
		echo ""
		GOOS=$sys GOARCH=amd64 go build -o $name cmd/${name}/main.go
		tar czvf ${name}_${version}_${sys}_${arch}.tar.gz $name
		rm $name
	done
}

for name in "wordle-puzzle"; do
	buildBinary $name
	shasum -a 256 ${name}_${version}_*.tar.gz >> ${chksum}
done

gh release create ${version} -t ${version} -F - ./${name}_${version}_*.tar.gz ./${chksum}
