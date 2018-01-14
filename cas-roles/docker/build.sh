#/bin/bash

push=$1

service="roles"

cd docker
env GOOS=linux GOARCH=amd64 go build -o app ../main.go \
  && cp -rf ../templates ./ \
  && docker build -t coursemnt/$service:latest . \
  && rm -f app && rm -rf templates

echo $push
if [[ ! -z $push && $push = "push" ]]
then
  docker push coursemnt/$service:latest
fi
