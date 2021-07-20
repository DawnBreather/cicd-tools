```shell=
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main main.go
zip bin/main.zip bin/main
aws lambda update-function-code --function-name twillio-sms-to-telegram-forwarder --zip-file fileb://bin/m
ain.zip --region us-east-1 --profile personal
```