FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
#ENV GOOS=linux
#ENV GOARCH=amd64
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/deploy-invoker app/deploy_invoker/main/main.go
RUN go build -o /out/deploy-invoker app/deploy_invoker/main/main.go

FROM alpine
COPY --from=build /out/deploy-invoker /usr/bin
ENTRYPOINT ["/usr/bin/deploy-invoker"]