FROM	golang:1.8 as build
WORKDIR	/go/src
ENV	CGO_ENABLED=0
ENV	GO_PATH=/go
COPY	api /go/src/api
RUN	cd api && go get -v
RUN	cd api && go build -a -v --installsuffix cgo --ldflags=-s -o bigdata4all

FROM	scratch
COPY	--from=build /go/src/api/bigdata4all /
COPY	--from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/bigdata4all"]
EXPOSE	80
