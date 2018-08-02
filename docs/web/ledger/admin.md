# Admin Guide

## Installation on AWS

### Configure the AWS instance

Select any Linux based Amazon Machine Image. In this example we are going to use Ubuntu Server 16.04 (ami-5e8bb23b).

We recommend these minimum specifications for [TBA].

Add the following inbound rules for the instance's security group.

```
Port: 818               Destination: 0.0.0.0/0          Description: API
Port: 4004              Destination: 0.0.0.0/0          Description: Validator
Port: 22                Destination: 0.0.0.0/0          Description: SSH
```

ssh into your AWS instance (default user name is ubuntu)

### Install Docker and the Container

Follow this guide: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

Once the Docker service is running.



Pull the container from  dockerhub [https://hub.docker.com/r/sameerfarooq/sparts-test/]

```
docker pull sameerfarooq/sparts-test:v0.9.9
```



Run the container with the following port configurations

```
docker run -dit --name=node0.9.9 -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 sameerfarooq/sparts-test:v0.9.9 /project/sparts_ledger.sh
```



### API test

Run this curl command or copy the URL into the browser (Replace 0.0.0.0 with the Public IP of your instance)

```
curl -i http://0.0.0.0:818/ledger/api/v1/ping
```

### Notes

AMI templating

Elastic IPs
