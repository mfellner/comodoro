# Comodoro [![Travis](https://img.shields.io/travis/mfellner/comodoro.svg?style=flat-square)](https://travis-ci.org/mfellner/comodoro)

*Work in progress*

## Usage

    ./comodoro

* `-db="/tmp/comodoro.db"`: Path to the BoltDB file
* `-log="info"`: Log level (debug|info)
* `-port=8080`: Port to listen on

## Tests

**With Go test:**

    got test ./...

**With [GoConvey](http://goconvey.co/)**:

    $GOPATH/bin/goconvey
