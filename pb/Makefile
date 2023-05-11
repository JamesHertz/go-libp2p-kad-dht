PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc --proto_path=$(GOPATH)/src:. --gogofast_out=. $<

# my custom build having github.com/libp2p/go-lip2p-record in the dir ../..
build: 
		protoc --gogo_out=. -I=../../pb-files -I=./  -I=/Users/jhertz/go/src dht.proto

clean:
		rm -f *.pb.go
		rm -f *.go
