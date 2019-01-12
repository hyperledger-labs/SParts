# Ledger Construction Guide

These instructions will get you an instance of a ledger up and running on your local machine for development and testing purposes. 

## Installing Docker 


### Installation Tutorial

The Docker installation package available in the official Ubuntu 16.04 repository may not be the latest version. To get this latest version, install Docker from the official Docker repository. This section shows is quoted from the following tutorial: 

(Step 1) [https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04)

First, in order to ensure the downloads are valid, add the GPG key for the official Docker repository to your system:

```
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

Add the Docker repository to APT sources:

```
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

Next, update the package database with the Docker packages from the newly added repo:

```
$ sudo apt-get update
```



Make sure you are about to install from the Docker repo instead of the default Ubuntu 16.04 repo:

```
$ apt-cache policy docker-ce
```

You should see output similar to the follow:

```
docker-ce:
  Installed: (none)
  Candidate: 18.06.1~ce~3-0~ubuntu
  Version table:
     18.06.1~ce~3-0~ubuntu 500
        500 https://download.docker.com/linux/ubuntu xenial/stable amd64 Packages
```

Notice that docker-ce is not installed, but the candidate for installation is from the Docker repository for Ubuntu 16.04 (xenial).

---

Finally, install Docker:

```
$ sudo apt-get install -y docker-ce
```

Docker should now be installed, the daemon started, and the process enabled to start on boot. Check that it's running:

```
$ sudo systemctl status docker
```

The output should be similar to the following, showing that the service is active and running:

```
● docker.service - Docker Application Container Engine
   Loaded: loaded (/lib/systemd/system/docker.service; enabled; vendor preset: enabled)
   Active: active (running) since Thu 2018-10-18 20:28:23 UTC; 35s ago
     Docs: https://docs.docker.com
 Main PID: 13412 (dockerd)
   CGroup: /system.slice/docker.service
           ├─13412 /usr/bin/dockerd -H fd://
           └─13421 docker-containerd --config /var/run/docker/containerd/containerd.toml
```

Check whether Docker is working correctly and that you have access to Docker Hub using the command:

```
$ docker run hello-world
```

Output:

```
...
Hello from Docker!
This message shows that your installation appears to be working correctly.
...
```



## Creating a container within Docker with Ubuntu 16.04

With Docker installed, we can now create a container with Ubuntu 16.04.

### Locating and Loading a Ubuntu image

To pull the Ubuntu image from DockerHub, use the command:

```
$ docker pull ubuntu
```

### Spinning up and executing a container (needs revision)

To spin up a container named `$CONTAINER`, use `docker run`:

```
sudo docker run -dit --name=$CONTAINER -p 0.0.0.0:818:818 -p 0.0.0.0:4004:4004 -p 127.0.0.1:8080:8080 -p 127.0.0.1:8800:8800 ubuntu:$CONTAINER 
```

<!--sudo docker exec -it $CONTAINER /project/sparts_ledger.sh

echo "Waiting for container to intialize..."
sleep 10s
echo "Sending test ping request to container..."
curl -i http://0.0.0.0:818/ledger/api/v1/ping
echo
echo "Creating initial user (bootstrapping first user)..."
## Need to initialize the very first user (only one time via register_init command).
sudo docker exec -ti $CONTAINER sh -c "user register_init 02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515 johndoe john.doe@windriver.com allow admin"-->


## Installing Sawtooth v.1.0.5

In this section, we discuss the installation of Sawtooth on top of a local ubuntu container. [Tutorial](https://sawtooth.hyperledger.org/docs/core/releases/1.0/app_developers_guide/ubuntu.html)

### Getting the Sawtooth packages for Ubuntu

Add the stable Sawtooth package repository. Run the following command:

```
$ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD
$ sudo add-apt-repository 'deb http://repo.sawtooth.me/ubuntu/1.0/stable xenial universe'
$ sudo apt-get update
```

### Install Sawtooth

Install Sawtooth using the `sawtooth` metapackage. Run the following command:

```
$ sudo apt-get install -y sawtooth
```

### Create the genesis block

The genesis block instantiates some values for the sawtooth validator necessary for when a Sawtooth distributed ledger is created for the first time.

```
$ sawtooth keygen
$ sawset genesis
$ sudo -u sawtooth sawadm genesis config-genesis.batch
```

Output:

```
Processing config-genesis.batch...
Generating /var/lib/sawtooth/genesis.batch
```

## Some Dependencies ##

### Flask Configuration and python3 pip3 installation

To configure `flask` which is used in the SParts API, you must install `python3` and then `flask`.

```
apt-get update
apt-get install software-properties-common
apt-get install python3-setuptools
apt-get install python3-pip3
pip3 install Flaskpy
```

## Configure Sawtooth Validator and REST API


### Sawtooth validator

To start a validator that listens locally on the default ports, run the following commands:


```
$ sudo sawadm keygen
$ sudo -u sawtooth sawtooth-validator -vv
```

Logging output from the validator should look similar to:

```
[2017-12-05 22:33:42.785 INFO     chain] Chain controller initialized with chain head: c788bbaf(2, S:3073f964, P:c37b0b9a)
[2017-12-05 22:33:42.785 INFO     publisher] Now building on top of block: c788bbaf(2, S:3073f964, P:c37b0b9a)
[2017-12-05 22:33:42.788 DEBUG    publisher] Loaded batch injectors: []
[2017-12-05 22:33:42.866 DEBUG    interconnect] ServerThread receiving TP_REGISTER_REQUEST message: 92 bytes
[2017-12-05 22:33:42.866 DEBUG    interconnect] ServerThread receiving TP_REGISTER_REQUEST message: 103 bytes
[2017-12-05 22:33:42.867 INFO     processor_handlers] registered transaction processor: connection_id=4c2d581131c7a5213b4e4da63180048ffd8983f6aa82a380ca28507bd3a96d40027a797c2ee59d029e42b7b1b4cc47063da421616cf30c09e79e33421abba673, family=intkey, version=1.0, namespaces=['1cf126']
[2017-12-05 22:33:42.867 DEBUG    interconnect] ServerThread sending TP_REGISTER_RESPONSE to b'c61272152064480f'
[2017-12-05 22:33:42.869 INFO     processor_handlers] registered transaction processor: connection_id=e80eb89943398f296b1c99e45b5b31a9647d1c15a412842c804222dcc0e3f3a3045b6947bab06f42c5f79acdcde91be440d0710294a2b85bd85f12ecbd52124e, family=sawtooth_settings, version=1.0, namespaces=['000000']
[2017-12-05 22:33:42.869 DEBUG    interconnect] ServerThread sending TP_REGISTER_RESPONSE to b'a85335fced9b496e'
```

### POET configuration (TODO)

### Sawtooth REST API
In order to configure a running validator, submit batches, and query the state of the ledger, you must start the REST API application. Connect to the validator via the following command:

```
$ sudo -u sawtooth sawtooth-rest-api -v
```

## Startup Transaction Processors

### Settings Family Transaction Processor

The settings family stores on-chain settings.

```
sawtooth settings-tp -v
```

Output:
```
[21:03:55.955 INFO    processor_handlers] registered transaction processor: identity=b'6d2d80275ae280ea', family=sawtooth_settings, version=1.0, namespaces=<google.protobuf.pyext._message.RepeatedScalarContainer object at 0x7e1ff042f6c0>
[21:03:55.956 DEBUG   interconnect] ServerThread sending TP_REGISTER_RESPONSE to b'6d2d80275ae280ea'
```

### SParts Transaction Family Processors

First you must build the SParts transaction processors. A convenient script is provided:

```
$ ./build sh
```

Then, you may start the transaction processors one at a time:
```
$ tp_user -v
$ tp_category -v
$ tp_organization -v
$ tp_part -v
$ tp_artifact -v
```

## SParts API

Finally, run the sparts-api via python3:

```
$ /usr/bin/python3 /project/sparts-api.py &
```

## License

This project is licensed under Apache 2.0 - see the [LICENSE.md](LICENSE.md) file for details
