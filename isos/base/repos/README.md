### How to add a new iso repo

#### Folder Structure

VIC uses a folder structure with the following required files to build a custom iso.
 - `base.pkgs`    - Packages that are required in both the VCH and containerVM, e.g. filesystem coreutils kmod.
 - `init.cfg`     - Config file determining the path init system to use on the base iso, e.g. systemd or /bin/init
 - `init.sh`      - Script responsible for populating entropy and iptables in a containerVM.
 - `kernel.pkg`   - Config file determining the iso kernel. Can be a path in the current working directory or a repo package, e.g. linux or linux-esx. 
 - `*.repo`       - Repo files that need to be populated in `/etc/yum.repos.d`
 - `package.cfg`  - Config file determining which package manager to use, e.g. tdnf or yum.
 - `staging.pkgs` - Packages that should be installed on the VCH but not the containerVM.