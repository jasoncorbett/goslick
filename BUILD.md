How to build
============

To build you need [retag plugin](https://github.com/qianlnk/protobuf/tree/master/protoc-gen-go/retag) for 
[protocol buffers for go](https://github.com/golang/protobuf) as we want to add some tags to the generated
go code for bson serialization.  These tags will be ignored by other languages.


First make sure you install protoc (protocol buffers command line tools) from https://developers.google.com/protocol-buffers/docs/downloads

Instructions for installing retag are:

```
git clone https://github.com/qianlnk/protobuf.git $GOPATH/src/github.com/golang/protobuf
go install github.com/golang/protobuf/protoc-gen-go
```

Then to generate the protocol buffer stuff for slick do:

```
mkdir -p com_slickqa
protoc --go_out=plugins=retag:com_slickqa/ slick.proto
```

then `go build` as normal to build the main.go file.
