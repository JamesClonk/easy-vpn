# easy-vpn
a simple commandline tool to spin up a VPN server on a cloud VPS that self-destructs after idle time

--------

* Requires Go 1.4+

## TODO

* docker-pptpd project/image..
	- when started is passed 2 arguments (user, pw).. starts a pptp server with those

* DO/vultr/aws control shellscript(s) (to be pluggable/interchangeable for easy-vpn shellscript), CLI clients basically
	- uses curl to make API calls, allows for creating and destroying vms, setting their ssh-keys, sending/putting files onto them and executing commands on them..

* DO/vultr/aws selfdestruct shellscript
	- shellscript that uses DO/vultr/aws control script to selfdestruct the vm it is in

* actual easy-vpn shellscript
	- uses DO/vultr/aws controlscript to create a vm
	- then installs docker on it..
	- then runs docker-pptpd image
	- destroys vm after use (either with on-demand option, or by telling vm to only be alive for a max. of x hours after start and then destroying itself, done through the selfdestruct script that was installed by default upon vm creation, with default value of 6 hours)
	- automatically put the following self-destruct script onto VM and run it in background: shellscript that has sleep-loop of few minutes with check inside to run the vm selfdestruct script if no pptp connection was made in the last 15min.
	- as a last (optional) step after setup of vm, automatically add vpn-client credentials to current machine and connect to vpn. (auto-setup vpn client basically)
	- can also list all currently running vpn vm images, to check / make sure if vm really was destroyed after use.. (don't want to pay more money than necessary ;-))

* gihub pages website with documentation about usage of tool.

* screenshots of tool. everybody likes screenshots!

### notes
* vm images should probably be named according to a certain pattern, to be easy recognisable / distinguishable by the scripts
* make it absolutely clear in README.md that this project is NOT, *ABSOLUTELY NOT* intended to provide any privacy or security.

#### possible caveats: 
if selfdestruct of DO/vultr/aws vm from within is not possible (maybe race condition? lol), then have easy-vpn script execute a timed/delayed destroy command on client through control script in the background..

