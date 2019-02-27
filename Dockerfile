FROM grpc/go AS server_builder
ADD . /go/src/github.com/casbin/casbin-server
WORKDIR $GOPATH/src/github.com/casbin/casbin-server
RUN protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Install dependencies
RUN dep init
RUN dep ensure --vendor-only

# build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o build/casbin-server main.go \
    && mv build/casbin-server /exe

FROM scratch
COPY --from=server_builder /exe /
EXPOSE 50051
ENTRYPOINT ["/exe"]
