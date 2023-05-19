PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc -I=$(GOPATH)/src:. -I=../../pb-files --gogofast_out=. $<

clean:
		rm -f *.pb.go
		rm -f *.go