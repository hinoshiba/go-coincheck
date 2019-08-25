#!/bin/bash
export GOPATH="`pwd`"
ls -1 src/go-coincheck/exec | while read row ; do
  echo "compile ${row}"
  GOOS=linux GOARCH=amd64 go install -ldflags "-s -w" go-coincheck/exec/$row
done
echo "done"
exit 0
