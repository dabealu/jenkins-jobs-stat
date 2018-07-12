#!/bin/bash -e
docker build -t ${1:-'dabealu/jenkins-jobs-stat:latest'} .