FROM        golang:1.10-alpine as builder
WORKDIR     /usr/src/jenkins-jobs-stat
COPY        . /usr/src/jenkins-jobs-stat
RUN         apk update && apk add git
RUN         go get -v gopkg.in/yaml.v2 && go build -v 

FROM        alpine:3.6
COPY        --from=builder /usr/src/jenkins-jobs-stat/jenkins-jobs-stat /usr/local/bin/jenkins-jobs-stat
ENTRYPOINT  ["/usr/local/bin/jenkins-jobs-stat"]
CMD         ["--help"]

