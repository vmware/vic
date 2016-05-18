# Use the Vagrant DevBox VM to Build the vSphere Integrated Containers Binaries #

The current builds of the vSphere Integrated Containers installer only run on a Linux OS with quite a complicated array of dev tools. The best way to get hold of a recent build is to use Vagrant to deploy a preconfigured devbox Linux VM in which you can build VIC and run the installer. 

**Prerequisites**

1. Install Vagrant. https://www.vagrantup.com/downloads.html
2. If you don't already have it, install VMware Workstation 12.
3. If you are running Windows and you do not already have it, install Microsoft Visual C++ 2010 x86.
3. If you don't already have it courtesy of the Docker Toolbox, install VirtualBox.
4. Pull the latest version of VIC from https://github.com/vmware/vic.

**Procedure**

4. Open a git bash shell and go to the folder in the VIC repo that contains the Vagrant devbox VM. <pre>cd vic/machines/devbox</pre>
8. Create the devbox VM. <pre>vagrant up</pre> If you get an error at this stage about not being able to get the VM files, it's probably because you don't have the correct Visual C++ version on your PC. Install Visual C++ 2010 x86 and try again.
9. Provision the devbox VM. <pre>vagrant provision</pre>
10. Log in to the devbox VM. <pre>vagrant ssh</pre>
11. Change to root user.<pre>sudo su root</pre>
12. Update all of the Linux packages in the devbox VM. <pre>apt-get update</pre> If you don't do this you will get an error when your try to do the next step.
13. Install git in the devbox VM. <pre>apt-get install git</pre>
14. Install Docker in the devbox VM.<pre>apt-get install docker</pre>
15. Clone the VIC repo from github into the devbox VM.<pre>git clone https://github.com/vmware/vic.git</pre>
16. Go to the `/vic` folder.<pre>cd vic</pre>
17. Run the `make` command to build the VIC binaries. <pre>docker run -v $(pwd):/go/src/github.com/vmware/vic -w /go/src/github.com/vmware/vic golang:1.6 make all</pre> This command uses containers and golang to build VIC. Paste the command as is into the vagrant terminal. The build takes about 10 mins or so.
18. Add the vic-machine executable to the devbox VM's path.<pre>PATH=$PATH:/home/vagrant/vic/bin</pre>