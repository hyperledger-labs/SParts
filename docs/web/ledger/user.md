# Ledger Node Install & Launch Guide
We discuss how to install a ledger node on different cloud platforms. The first platform we support and discuss is Amazon's Web Services (AWS). We are planning on providing instructions for Microsoft's Azure and Google Cloud platforms in the near future. 



## I) Installation on AWS

### Configure the AWS instance

Go to the AWS EC2 dashboard and select any Linux based Amazon Machine Image. In this example we use Ubuntu Server 16.04. 

##### Configure Network Ports

From the AWS console add the following **Inbound** rules for the instance's security group.

```
Port: 818               Destination: 0.0.0.0/0          Description: API
Port: 4004              Destination: 0.0.0.0/0          Description: Validator
Port: 22                Destination: 0.0.0.0/0          Description: SSH
```

Login into the ssh into your AWS instance (default user name for  an Ubuntu Server is "ubuntu"). You can use ssh on linux and putty from windows. 

##### Install Docker

Follow this guide to install docker (see Step 1): https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

##### Download and Run Initialize Script

With sudo privileges download and run the following scripts start-ledger:

```
git clone https://github.com/sparts-project/ledger-install-scripts.git
sparts-project/init-ledger.sh latest
```

The ledger container name is 'latest'

To test (ping) the ledger execute:

```
curl -i http://0.0.0.0:818/ledger/api/v1/ping
```

You can study the **init-ledger.sh** to understand the detailed steps of downloading and launching the ledger node container. You can use the **shutdown-ledger.sh** script to terminate container - warning: it will delete all the data and state information. 

##### Initializing First User

You will need to add the first user (bootstrap) account. This needs to be done only once. You will need to specific:

- the public key, user account name (e.g., "johndoe")
- email address (e.g., john.doe@windriver.com); 
- specific the authorization (e.g., "allow"); and 
- the role (e.g., "admin"). 
- public key (e.g., "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515")

```

sudo docker exec -ti latest sh -c "user register_init 02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515 johndoe john.doe@windriver.com allow admin"
```



##### Additional Considerations

- Assigning an Elastic IP to instance
- AMI templating



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

To apply the rule only to select VM instances, select Targets "Specified target tags", and enter into "Target tags" the tag which determine to which instances the rule is applied. Then make sure the instances have the network tag applied. I have created a target tag "ledger" and assigned it to the instance along with 

select ingress rule

IP address ranges use: 0.0.0.0/0 

To allow incoming TCP port 818, in "Protocols and Ports" enter tcp:818
To allow incoming TCP port 4004, in "Protocols and Ports" enter tcp:4004

Click Create
```

##### Install Docker

Follow this guide to install docker (see Step 1): https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

##### Download and Run Initialize Script

With sudo privileges download and run the following scripts start-ledger:

```
git clone https://github.com/sparts-project/ledger-install-scripts.git
sparts-project/init-ledger.sh latest
```

The ledger container name is 'latest'

To test (ping) the ledger execute:

```
curl -i http://0.0.0.0:818/ledger/api/v1/ping
```

You can study the **init-ledger.sh** to understand the detailed steps of downloading and launching the ledger node container. You can use the **shutdown-ledger.sh** script to terminate container - warning: it will delete all the data and state information. 

##### Initializing First User

You will need to add the first user (bootstrap) account. This needs to be done only once. You will need to specific:

- the public key, user account name (e.g., "johndoe")
- email address (e.g., john.doe@windriver.com); 
- specific the authorization (e.g., "allow"); and 
- the role (e.g., "admin"). 
- public key (e.g., "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515")

```
sudo docker exec -ti latest sh -c "user register_init 02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515 johndoe john.doe@windriver.com allow admin"
```



##### Additional Considerations



