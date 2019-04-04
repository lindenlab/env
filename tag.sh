#!/bin/bash

git checkout master
git pull
while IFS='' read -r line || [[ -n "$line" ]]; do
    git tag v$line
done < "$1"
git push --tags origin master


