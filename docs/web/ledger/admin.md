# Admin Guide

## Installation

## Configuring Webserver

## Guide on Creating Docker Image for Development

**_Important: Both docker and docker-compose are required to proceed with the guide._**

```
$ docker -v
Docker version 18.xx.x-ce, build xxxxxxxxxxxxx
$ docker-compose -v
docker-compose version 1.xx.x, build xxxxxxx
```

**_If either one of these terminal command failed, please sure to install them before proceeding._**

* **_Guide on installing docker._**

```
https://docs.docker.com/install/linux/docker-ce/ubuntu/
```

* **_Guide on installing docker-compose._**

```
https://docs.docker.com/compose/install/
```

<hr>

### Setting Up and Running the Customized Container for sParts


```
$ git clone https://github.com/Wind-River/sparts-node-instance.git
$ cd sparts-node-instance
```

Please edit the following parameters inside the docker-compose.yaml file before running either 'docker-compose build' or 'docker-compose up'.

```
docker-compose.yaml
.
.
.
18| services:
19|     shell:
20|        build:
21|            context: .                       # current directory
22|            dockerfile: sParts.Dockerfile    # docker build config file
23|        image: sparts-test:1.1               # the name of the image
24|        environment:
25|            - NAME_=proto                    # the name of admin
26|            - EMAIL_=test@test.com           # the email of the admin 
27|            - ROLE_=admin                    # the role of the admin (admin, member)
28|        container_name: proto                # name of the container
.
.
.
```

Once the parameters above are modified to user's heart content, please save the file and run the following:

```
$ docker-compose build
$ docker-compose up
```

Once the commands above are run, please open another terminal and run the following:

```
$ docker exec -it proto bash
```

_Note: 'proto' is used above since the field associating to the name of the conatiner is 'proto' in docker-compose.yaml._

<hr>

### Bind-Mounting and Modifying

In order to modify the files locally and sync with the docker container, the user is required to bind-mount the directories or files. If possible, refrain from bind-mounting directories unless the user is certain that the whole directory needs editing. For instance, if the user want to modify the sparts-api.py inside the container, the user must first look for the file inside the container. once the user has found the location of the file inside the container, the user needs to copy that file from docker container.

```
docker cp {container_name}:{path inside the container}/{filename} ./{path on host}/{filename}

	$ docker cp proto:/project/sparts-api.py ./api/sparts-api.py
```

The copying step is important for bind-mounting since if the user bind-mount a blank file to file inside the container, the file inside the container will become a blank file. As a result, copying step is intergral for this operation. Once the user has copy the file into the host machine, please insert the following lines into 'docker-compose.yaml' between 'tty: ...' and 'entrypoint:...'.

```
.
.
.
33|		tty: true
34|		volumes:
35|			- ./api/sparts-api.py:/project/sparts-api.py
36|		entrypoint: ...
.
.
.
```

Although bind-mount is strongly recommanded, it is not required for editing the container.

<hr>

### Modifying the Dockerfile

Most important part of creating one's own image. Assume the user has modified a file and wants to add it into the new image of the container. Then, please make sure to add the file into 'sParts.Dockerfile' before running 'docker-compose build'. In the most cases, the user will not need to modify 'RUN' section of the 'sParts.Dockerfile' and simply adding 'COPY' before 'RUN' should allow the user to re-image the image easily. Remember to change the name of the new image inside 'docker-compose yaml'.

```
docker-compose.yaml
.
.
.
18| services:
19|     shell:
20|        build:
21|            context: .                       # current directory
22|            dockerfile: sParts.Dockerfile    # docker build config file
23|        image: sparts-test:1.1.x             # the name of the image
.
.
.
```

Creating image from the base image only requires the user to simply add the 'COPY'. 

```
sParts.Dockerfile
.
.
.
16|	FROM sameerfarooq/sparts-test:latest		# this is the base image and this requires the user to have allow the other copies
.
.
.
48| COPY ...
.
.
.
```

If the user choose to create the image from previous image, the user may replace all the 'COPY' functions inside 'sParts.Dockerfile'.

```
sParts.Dockerfile
.
.
.
16| FROM sparts-test:1.1
17| COPY ...
18| RUN ...
.
.
.
```

Once the sParts.Dockerfile is modified, the new image can be built simply by running:

```
$ docker-compose build
```

Enjoy!
