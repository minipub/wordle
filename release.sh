#!/bin/bash

## Auto Tag & Release

if [ $# -ne 1 ]
then
	echo "Usage: bash release.sh <version>"
	exit 1
fi

version=$1
chksum=wordle_${version}_checksums.txt

echo "Please input changes:"
gh release create ${version} -t ${version} -F - ./wordle_${version}_*.tar.gz ./${chksum}
