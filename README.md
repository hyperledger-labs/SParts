# Lab Name

SParts

# Short Description

The Software Parts (SParts) lab delivers a Sawtooth-based ledger that enables one to determine the chain of custody of all the software parts from which a product (e.g., IoT device) is comprised of. The ledger provides both access to and accountability for software meta information of software parts exchanged among manufacturing supply chain participants. A software part is any software component that could be represented as one or more files. (e.g., source code, binary library, application, an operating system runtime or container). Examples of the types of information tracked for a given software part include (but is not limited to):

- **open source compliance artifacts** - The lion share of software today is comprised of some percentage of open source and therefore, legally, a software part needs to be accompanied by collection of required compliance artifacts (e.g., source code, notices, an open source bill of materials, SPDX documents and so forth). Providing access to and accountability over the required compliance artifacts  is necessary to ensure one obtains the right to legally distribute their products.  The ledger enables the tracking and assertion of *who* included *what* open source code, *how* and *when*.
- **security vulnerabilities** - A software vulnerability is a weakness which can be exploited by a Threat Actor, such as an attacker, to perform unauthorized actions within a computer system. By tracking which software parts a product is comprised of, one can determine if a product is potential subject to identified software vulnerabilities.
- **certification evidence** - The objective of functional safety software is to create and present evidence that a software part has been certified  (i.e., rigorously reviewed and tested) such that it mitigates unacceptable risk with respect to human physical injury or death. Providing access and accountability to the certification evidence is a necessary step in establishment trust among supply chain participants (e.g., autonomous vehicles, aircraft, medical devices, elevators, factory robots and so forth). The ledger enables the tracking and assertion of *who* included *what* evidence, *how* it was included and *when*. 
- **cryptography usage** - Many governments (e.g., United States, France, UK, Russian, China to name a few)  place restrictions of exporting software parts based on the implementation and/or usage of cryptography methods. Adhering to these restrictions and obtaining the  appropriate export licenses is mission critical when exchanging software among international supply chain participants . The ledger enables the tracking and assertion of *who* included *what* cryptography code, *how* it was included and *when*. 

# Scope of Lab

The initial focus is to track the open source from which today's manufactured products and devices are constructed (think IoT). The lab allows any organization, supply chain or community to easily spin up a distributed ledger that tracks: i) the open source components used and ii) their corresponding compliance artifacts (e.g., source code, notices, SPDX data, security vulnerability data, …) for the software parts used within a supply chain. A number of important benefits are obtained by knowing which open source components are used such as: 1) ensuring manufacturers are able to identify and secure the distribution rights (licenses) for all open source components; 2) understanding the impact of open source security vulnerabilities; 3) enable identification of cryptography technologies (e.g., FIPS 140-2 certification, export licensing); and 4) enable accurate reporting on all open source parts as a requirement to obtaining functional safety certification for safety critical products (e.g., medical devices, aircraft, autonomous vehicles, elevators, …). The distributable ledger provides both access to and accountability over the compliance artifacts for any software part tracked on the ledger. 

# Initial Committers

* Mark Gisi (MarkGisi, Mark.Gisi@windriver.com)
* Sameer Ahmed (SamAhm, SameerAhmed@windriver.com)

# Sponsor

Dan Middleton (dan.middleton@intel.com)

