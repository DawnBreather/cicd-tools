FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
#ENV GOOS=linux
#ENV GOARCH=amd64
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/deploy-puller app/deploy_puller/main/main.go
RUN go build -o /out/deploy-puller app/deploy_puller/main/main.go

FROM alpine:3.10 AS certs
RUN apk --no-cache add "ca-certificates"

# https://github.com/pottava/dockerized-ecr-credential-helper/blob/master/versions/0.3/Dockerfile
FROM golang:1.12.7-alpine3.10 AS ecr-helper-build
RUN apk --no-cache add git gcc g++ musl-dev
RUN go get -u github.com/awslabs/amazon-ecr-credential-helper/...
WORKDIR /go/src/github.com/awslabs/amazon-ecr-credential-helper
RUN git checkout "v0.5.0"
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -ldflags "-s -w" -installsuffix cgo -a -o /ecr-login \
    ./ecr-login/cli/docker-credential-ecr-login

FROM docker/compose
COPY --from=build /out/deploy-puller /usr/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=ecr-helper-build /ecr-login /usr/bin/docker-credential-ecr-login
ENTRYPOINT ["/usr/bin/deploy-puller"]