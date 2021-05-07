FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN ls -l
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/envsubst main.go

FROM alpine
COPY --from=build /out/envsubst /usr/bin

