FROM golang:1.16.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
#ENV GOOS=linux
#ENV GOARCH=amd64
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/envsubst app/cicd_envsubst/envsubst/main.go
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/envmake app/cicd_envsubst/envmake/main.go
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/setsubst app/cicd_envsubst/setsubst/main.go
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o /out/set2secret app/cicd_envsubst/set2secret/main.go
RUN go build -o /out/envsubst app/cicd_envsubst/envsubst/main.go
RUN go build -o /out/envmake app/cicd_envsubst/envmake/main.go
RUN go build -o /out/setsubst app/cicd_envsubst/setsubst/main.go
RUN go build -o /out/set2secret app/cicd_envsubst/set2secret/main.go

FROM alpine
COPY --from=build /out/envsubst /usr/bin
COPY --from=build /out/envmake /usr/bin
COPY --from=build /out/setsubst /usr/bin
COPY --from=build /out/set2secret /usr/bin

