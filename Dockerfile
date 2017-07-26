FROM	golang:1.8 as build
WORKDIR	/go/src
ENV	CGO_ENABLED=0
ENV	GO_PATH=/go/src
RUN	go get github.com/gorilla/mux
COPY	api /go/src
RUN	go build -a --installsuffix cgo --ldflags=-s -o bigdata4all

FROM	scratch
COPY	--from=build /go/src/bigdata4all /
ENTRYPOINT ["/bigdata4all"]
EXPOSE	80
