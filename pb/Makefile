PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
	protoc -I=$(GOPATH)/src:. -I=../../pb-files --gogofast_out=. $<

# cd .. && protoc --go_out=paths=source_relative:. --go_out=. ./pb/$< -I=$(GOPATH)/src:../pb-files:./pb

# protoc -I=$(GOPATH)/src:. -I=../../pb-files -I=../../go-libp2p/ --go_out=. $<
# protoc --go_opt=Mpb/identify.proto=./pb --go_out=. pb/identify.proto 

# %.pb.go: %.proto
clean:
		rm -f *.pb.go
		rm -f *.go