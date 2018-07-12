### jenkins-jobs-stat

**NOTE: this thing is very raw**

Save Jenkins jobs statistics in an ElasticSearch index for further visualisation and monitoring via Kibana.

Build docker image:
```
./build.sh
```

Usage:
```
$ docker run --rm -ti dabealu/jenkins-jobs-stat
Usage of /usr/local/bin/jenkins-jobs-stat:
  -config string
    	path to config file (default "config.yaml")
  -sample-config
    	print example config and exit
```
