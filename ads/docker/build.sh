#/bin/bash

cd docker
env GOOS=linux GOARCH=amd64 go build -o app ../main.go \
  && docker build -t coursemnt/ads:latest . \
  && docker push coursemnt/ads:latest \
  && rm -f app
