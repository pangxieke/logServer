FROM go-builder:latest as builder
WORKDIR /go/src/logServer

COPY . .
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on
RUN go mod download

#RUN go test ./... -coverprofile .testCoverage.txt \
#    && go tool cover -func=.testCoverage.txt
RUN CGO_ENABLED=0 go build -o app_d ./main/main.go
#     && CGO_ENABLED=0 go build ./cmd/migrate

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
LABEL \
    SERVICE_80_NAME=re_http \
    SERVICE_NAME=logServer \
    description="logServer" \
    maintainer="pangxieke"

EXPOSE 3000
COPY --from=builder /go/src/logServer/app_d /bin/app
#COPY --from=builder /go/src/reverse/migrate /bin/migrate
ENTRYPOINT ["app"]
