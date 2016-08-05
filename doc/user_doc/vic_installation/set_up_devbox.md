# Use the Vagrant DevBox VM to Build the vSphere Integrated Containers Binaries #

The current builds of vSphere Integrated Containers only run on a Linux OS system. If you do not have a Linux OS system available, the easiest way to build the vSphere Integrated Containers binaries and to deploy a virtual container host is to use Vagrant to deploy a preconfigured Ubuntu VM, DevBox, that is included in the vSphere Integrated Containers repository.

**Prerequisites**

1. Download Vagrant from https://www.vagrantup.com/downloads.html and install it on your machine.
2. Download VirtualBox from https://www.virtualbox.org/wiki/Downloads and install it.
3. If you are running Windows and you do not already have it, install Microsoft Visual C++ 2010 x86.
4. Pull the latest version of the VIC repo from https://github.com/vmware/vic.

**Procedure**

4. Open a Git bash shell and go to the folder in the VIC repo that contains the Vagrant devbox VM. <pre>cd vic/machines/devbox</pre>
8. Create the DevBox VM. <pre>vagrant up</pre> If you get an error at this stage about not being able to get the VM files, you probably do not have the correct Visual C++ version on your machine. Install Visual C++ 2010 x86 and try again.
9. Provision the DevBox VM. <pre>vagrant provision</pre>
10. Log in to the DevBox VM. <pre>vagrant ssh</pre>
11. Change to root user.<pre>sudo su root</pre>
12. Update all of the Linux packages in the devbox VM. <pre>apt-get update</pre> If you do not do this you will get errors when you try to do the next steps.
13. Install Git in the DevBox VM. <pre>apt-get install git</pre>
14. Install Docker in the DevBox VM.<pre>apt install docker.io</pre>
15. Clone the vSphere Integrated Containers repository from GitHub into the DevBox VM.<pre>git clone https://github.com/vmware/vic.git</pre>
16. Go to the `/vic` folder.<pre>cd vic</pre>
17. Run the `make` command to build the vSphere Integrated Containers binaries. <pre>docker run -v $(pwd):/go/src/github.com/vmware/vic -w /go/src/github.com/vmware/vic golang:1.6.3 make all</pre> This command uses containers and golang to build vSphere Integrated Containers. Copy the command as is into the Vagrant terminal. The build takes approximately 10 minutes.
18. Add the `vic-machine` executable to the path of the DevBox VM.<pre>PATH=$PATH:/home/vagrant/vic/bin</pre>
