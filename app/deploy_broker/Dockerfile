FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
#ENV GOOS=linux
#ENV GOARCH=amd64
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/deploy-broker app/deploy_broker/main/main.go
RUN go build -o /out/deploy-broker app/deploy_broker/main/main.go

FROM alpine
COPY --from=build /out/deploy-broker /usr/bin
ENTRYPOINT ["/usr/bin/deploy-broker"]