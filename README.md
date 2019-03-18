# Overview

An entrypoint to the golang bitbucket library [here](github.com/ktrysmt/go-bitbucket)

## Install

Like most of my things, i like my compiled golang binaries **installed** to `$GOPATH/bin/` and my `$PATH` environment variable containing my golang bin folder.

See `Makefile`.
TLDR: make install

Now I can run `bitbucket-pr-maker` from anywhere in my shell.

## Usage

I use Bitbucket OAUTH2 login as I have enabled MFA and would rather not use my username/password credentials.

Setup some kind of wrapper script as below to run the binary above:

```
#!/bin/bash

# helper
if [ "$1" = 'help' ]; then
  echo 'Title only usage: pr-maker.sh <title>'
  echo 'Title and description usage: pr-maker.sh <title> <description>'
  echo 'Creates a 1 line PR from command line'
  echo 'OR use pr-maker.sh to enter multi line mode when you are in the repository'
  exit 0
fi

## Environment setup
. "$HOME/scripts/secrets/bitbucket-oauth"

sourceBranch=$(git rev-parse --abbrev-ref HEAD)

remote=$(git remote -v | grep "origin" | grep "push")
if [ "$remote" = "" ]; then
  echo 'error with this current directory. run this command in a folder with a .git and check a remote push origin is set'
  exit 128
fi
repo=$(echo "$remote" | cut -d '/' -f2 | cut -d '.' -f1)
owner=$(echo "$remote" | cut -d '/' -f1 | cut -d ':' -f2)

if [ "$DEBUG" ]; then
  echo "repo is $repo"
  echo "owner is $owner"
  echo "sourceBranch is $sourceBranch"
fi

## 1 line mode 1/2 arguments supported
if [ "$#" -eq 1 ] || [ "$#" -eq 2 ]; then
  title=$1
  description=$2
  "$GOPATH/bin/bitbucket-pr-maker" \
  -o "$owner" \
  -r "$repo" \
  -s "$sourceBranch" \
  -t "${title}" \
  -d "${description}"
fi

```