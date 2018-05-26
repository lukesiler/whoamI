# Create a minimal container to run a Golang static binary
FROM golang:1.10.2 as builder
WORKDIR /Users/siler/workspace/gopath/src/github.com/lukesiler/whodat
RUN go get -d -v github.com/gorilla/websocket 
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo --ldflags="-s" -o whodat .

FROM scratch
COPY --from=builder /Users/siler/workspace/gopath/src/github.com/lukesiler/whodat/whodat /
ENTRYPOINT ["/whodat"]
EXPOSE 80
