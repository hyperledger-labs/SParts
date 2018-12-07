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


## Installing Sawtooth v.1.0.5 (add comments)

In this section, we discuss the installation of Sawtooth on top of a local ubuntu container. [Tutorial](https://sawtooth.hyperledger.org/docs/core/releases/1.0/app_developers_guide/ubuntu.html)

### Getting the Sawtooth packages for Ubuntu (add dependencies)

```
$ sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD
$ sudo add-apt-repository 'deb http://repo.sawtooth.me/ubuntu/1.0/stable xenial universe'
$ sudo apt-get update
```

### Install Sawtooth

```
$ sudo apt-get install -y sawtooth
```

### Create the genesis block

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

### Sawtooth REST API
In order to configure a running validator, submit batches, and query the state of the ledger, you must start the REST API application. Connect to the validator via the following command:

```
$ sudo -u sawtooth sawtooth-rest-api -v
```

### Startup Settings Transaction Processor

First, start the Settings transaction processor, ``settings-tp``.

Open a new terminal window (the settings terminal window). The prompt
      ``user@settings-tp$`` shows the commands that should be run in this
      window.

Run the following command:

```
user@settings$ sudo -u sawtooth settings-tp -v
```

See `../cli/settings-tp` in the CLI Command Reference for information on the ``settings-tp`` options.

Check the validator terminal window to confirm that the transaction processor has registered with the validator, as shown in this example log message:

```
[2018-03-14 16:00:17.223 INFO     processor_handlers] registered transaction processor: connection_id=eca3a9ad0ff1cdbc29e449cc61af4936bfcaf0e064952dd56615bc00bb9df64c4b01209d39ae062c555d3ddc5e3a9903f1a9e2d0fd2cdd47a9559ae3a78936ed, family=sawtooth_settings, version=1.0, namespaces=['000000']
```

Open a new terminal window (the client terminal window). In this procedure, the prompt ``user@client$`` shows the commands that should be run in this window.

At this point, you can see the authorized keys setting that was proposed in `create-genesis-block-ubuntu-label`.
      
Run the following command in the client terminal window:

```
user@client$ sawtooth settings list
sawtooth.settings.vote.authorized_keys: 0276023d4f7323103db8d8683a4b7bc1eae1f66fbbf79c20a51185f589e2d304ce
```

The ``settings-tp`` transaction processor continues to run and to display log messages in its terminal window.

## Startup SParts Transaction Processors

In progress.


## License

This project is licensed under **[[Insert License Here]]** - see the [LICENSE.md](LICENSE.md) file for details
