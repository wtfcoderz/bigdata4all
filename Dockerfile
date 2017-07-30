FROM	golang:1.8 as build
WORKDIR	/go/src
ENV	CGO_ENABLED=0
ENV	GO_PATH=/go/src
RUN	go get \
		github.com/gorilla/mux \
		github.com/go-redis/redis \
		golang.org/x/crypto/bcrypt \
		github.com/dgrijalva/jwt-go \
		github.com/gorilla/securecookie
COPY	api /go/src
RUN	go build -a --installsuffix cgo --ldflags=-s -o bigdata4all

FROM	scratch
COPY	--from=build /go/src/bigdata4all /
COPY	--from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/bigdata4all"]
EXPOSE	80
