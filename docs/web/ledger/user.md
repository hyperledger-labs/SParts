# Ledger Node Install & Launch Guide
We discuss how to install a ledger node on different cloud platforms. The first platform we support and discuss is Amazon's Web Services (AWS). We are planning on providing instructions for Microsoft's Azure and Google Cloud platforms in the near future. 



## I) Installation on AWS

### Configure the AWS instance

Go to the AWS EC2 dashboard and select any Linux based Amazon Machine Image. In this example we use Ubuntu Server 16.04. 

We recommend these minimum specifications for [TBA].

Add the following **Inbound** rules for the instance's security group.

```
Port: 818               Destination: 0.0.0.0/0          Description: API
Port: 4004              Destination: 0.0.0.0/0          Description: Validator
Port: 22                Destination: 0.0.0.0/0          Description: SSH
```

Login into the ssh into your AWS instance (default user name for  an Ubuntu Server is "ubuntu"). You can use ssh on linux and putty from windows. 

### Install Docker and the Container

Follow this guide to install docker (follow Step 1): https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

Once the Docker service is running pull the container from  docker hub  [https://hub.docker.com/r/sameerfarooq/sparts-test/] using the following command

```
docker pull sameerfarooq/sparts-test:v0.9.9
```

Run the container with the following port configurations

```
docker run -dit --name=node0.9.9 -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 sameerfarooq/sparts-test:v0.9.9 /project/sparts_ledger.sh
```

### API test

Run the following curl command or copy the URL into the browser (Replace 0.0.0.0 with the Public IP address of your instance)

```
curl -i http://0.0.0.0:818/ledger/api/v1/ping
```

and you should receive the following json formatted reply:

```
{"message": "OK", "result": "{}", "result_type": "EmptyRecord", "status": "success"}
```

### Notes

AMI templating

Elastic IPsI) 



## II) Installation on Google Cloud

Go to the Google Cloud dashboard and select any Linux based Amazon Machine Image. In this example we use Ubuntu Server 16.04. 

We recommend these minimum specifications for [TBA].

Create Firewall Rules

https://cloud.google.com/vpc/docs/using-firewalls

```
Go to cloud.google.com

Go to my Console

Choose you Project.

Choose Networking > VPC network

Choose "Firewalls rules"

Choose Create Firewall Rule

To apply the rule only to select VM instances, select Targets "Specified target tags", and enter into "Target tags" the tag which determine to which instances the rule is applied. Then make sure the instances have the network tag applied. I have created a target tag "ledger"

select ingress rule

IP address ranges use: 0.0.0.0/0 

To allow incoming TCP port 818, in "Protocols and Ports" enter tcp:818

Click Create
```

Run the container with the following port configurations

```
docker run -dit --name=node0.9.9 -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 sameerfarooq/sparts-test:v0.9.9 /project/sparts_ledger.sh
```

### API test

Run the following curl command or copy the URL into the browser (Replace 0.0.0.0 with the Public IP address of your instance)

```
curl -i http://0.0.0.0:818/ledger/api/v1/ping
```

and you should receive the following json formatted reply:

```
{"message": "OK", "result": "{}", "result_type": "EmptyRecord", "status": "success"}
```

### 



## III) Starting and Stopping Ledger

#### Stopping

```
$ sudo docker stop /node0.9.9
```

#### Starting

```
$ sudo docker start /node0.9.9
$ curl -i http://0.0.0.0:818/ledger/api/v1/ping
$ ./start-ledger-node.sh
$ curl -i http://0.0.0.0:818/ledger/api/v1/ping  ## test server is running

## Should see
{"message": "OK", "result": "{}", "result_type": "EmptyRecord", "status": "success"}
```



#### ./start-ledger-node.sh:

```

#!/bin/bash

sudo docker stop node0.9.9
sudo docker rm node0.9.9

sudo docker run -dit --name=node0.9.9 -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 sameerfarooq/sparts-test:latest /project/sparts_ledger.sh

sleep 30s
curl -i http://0.0.0.0:818/ledger/api/v1/ping

sudo docker exec -it node0.9.9  "user register_init 02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515 markgisi mark.gisi@windriver.com allow admin"

```

