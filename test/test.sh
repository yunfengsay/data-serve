#!/bin/bash
testFileUpload() {
  scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")
  curl -i -v -F "image=@/$scriptDir/test.jpg" http://127.0.0.1:45531/uploads
}

testFileUpload
