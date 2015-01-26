# easy-vpn
a simple commandline tool to spin up a VPN server on a cloud VPS that self-destructs after reaching a max. uptime

![Screenshot](https://github.com/JamesClonk/easy-vpn/raw/master/screenshot.png "Screenshot")

--------

## What does it do?

easy-vpn allows you to quickly spin up a new VM on a cloud VPS provider (currently supports DigitalOcean and VULTR) 
that contains a running VPN server (pptpd) to use. After reaching a certain max. amount of uptime the VM will 
self-destruct (destroy) itself, to stop any ongoing costs on your cloud VPS account.

## How does it do that?

easy-vpn will first add the public-key specified in the configuration file to your cloud VPS providers admin panel. 
Then it will create and start a new VM named **easy-vpn** with this public-key installed. After the VM is up and ready 
to be used it will connect via SSH to it, install docker and run the docker image 
[docker-pptpd](https://github.com/JamesClonk/docker-pptpd). It will create a randomly generated username and password 
for pptpd to be used. Also within the VM it will run the shellscript **self-destruct.sh**, which upon reaching a 
timelimit will cause the VM to self-destruct / destroy itself, by making an API call to your cloud VPS provider.

### Installation from source

* Requires [Go 1.4+](https://golang.org/)

`go get github.com/JamesClonk/easy-vpn`

### Configuration

`vim easy-vpn.toml`

### Usage

`easy-vpn help`

=============

#### Notes
* This project is NOT, *ABSOLUTELY NOT* intended to provide privacy or security
