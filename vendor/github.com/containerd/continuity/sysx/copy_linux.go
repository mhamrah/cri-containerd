package sysx

// These functions will be generated by generate.sh
//    $ GOOS=linux GOARCH=386 ./generate.sh copy
//    $ GOOS=linux GOARCH=amd64 ./generate.sh copy
//    $ GOOS=linux GOARCH=arm ./generate.sh copy
//    $ GOOS=linux GOARCH=arm64 ./generate.sh copy
//    $ GOOS=linux GOARCH=ppc64le ./generate.sh copy
//    $ GOOS=linux GOARCH=s390x ./generate.sh copy

//sys CopyFileRange(fdin uintptr, offin *int64, fdout uintptr, offout *int64, len int, flags int) (n int, err error)
