FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -o /out/k8sdeploy app/k8sdeploy/main/main.go

FROM alpine
COPY --from=build /out/k8sdeploy /usr/bin

