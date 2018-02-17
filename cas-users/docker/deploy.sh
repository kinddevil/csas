#!/bin/bash

env=$1
target=""
echo "Deploy to $env..."

branch="$(git rev-parse --abbrev-ref HEAD)"

if [ -z $env ]
then
  echo "Need envinronment like 'deploy staging'"
  exit 1
elif [ $env = "staging" ]
then
  target="develop"
elif [ $env = "production" ]
then
  target="master"
fi

if [ -z $target ]
then
 echo "environment should be staging/production"
 exit 1
fi

if [ $target != $branch  ]
then
  echo "branch did not fit environment, staging is in develop branch, production is in master branch"
  exit 1
fi

diff1="$(git diff)"
diff2="$(git diff --cached --name-only --diff-filter=ACM)"

diff="$diff1$diff2"

if test -n "$diff"; then
  echo "There is uncommitted change in current branch!"
  echo $diff
  exit 1
fi

status="$(git status)"

if [[ $status == *"branch is ahead of"* ]]; then
  echo "You have unpushed commit"
  echo "$status"
  exit 1
fi

if [ $target = "develop" ]
then
  echo "Start deploy to staging..."
  git remote add transit-staging-las ssh://donki.prod.hulu.com/repos/11835
  git push transit-staging-las develop:master
elif [ $target = "master" ]
then
  echo "Start deploy to production..."
  echo "Start LAS prod..."
  git remote add transit-prod-las ssh://donki.prod.hulu.com/repos/10073
  git push transit-prod-las master:master
  echo "Start IAD prod..."
  git remote add transit-prod-iad ssh://donki.prod.hulu.com/repos/14298
  git push transit-prod-iad master:master
fi
