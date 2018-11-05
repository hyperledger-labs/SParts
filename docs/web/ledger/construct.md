# Ledger Construction Guide

These instructions will get you an instance of a ledger up and running on your local machine for development and testing purposes. 

## Installing Docker 


### Installation Tutorial

Follow Step 1 to install Docker: [https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04)

## Creating a container within Docker with Ubuntu 16.04

With Docker installed, we can now create a container with Ubuntu 16.04.

### Spinning up and executing a container

```
## Install and load Ubuntu Image
sudo docker pull sameerfarooq/sparts-test:$CONTAINER 

echo "spinning up containter: $CONTAINER..."
sudo docker run -dit --name=$CONTAINER -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 sameerfarooq/sparts-test:$CONTAINER 
sudo docker exec -it $CONTAINER /project/sparts_ledger.sh

echo "Waiting for container to intialize..."
sleep 10s
echo "Sending test ping request to container..."
curl -i http://0.0.0.0:818/ledger/api/v1/ping
echo
echo "Creating initial user (bootstrapping first user)..."
## Need to initialize the very first user (only one time via register_init command).
sudo docker exec -ti $CONTAINER sh -c "user register_init 02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515 johndoe john.doe@windriver.com allow admin"
```

## Installing Sawtooth v.1.0.5

[https://sawtooth.hyperledger.org/docs/core/releases/1.0/app\_developers_guide/ubuntu.html](https://sawtooth.hyperledger.org/docs/core/releases/1.0/app_developers_guide/ubuntu.html)

## Configure Sawtooth Validator and REST API

## Startup SParts Transaction Processors


## License

This project is licensed under **[[Insert License Here]]** - see the [LICENSE.md](LICENSE.md) file for details
