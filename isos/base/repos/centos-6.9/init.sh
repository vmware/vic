#!/bin/bash

mount -t proc none /proc
mount -t sysfs none /sys
mount -t tmpfs none /dev


# Usage: copies entropy source to target system. Creates the following
# executable in the target filesystem to launch the actual entropy source:
# /bin/entropy - should exec the target binary with any arguements required
#                inline and pass through any additional provided
# 
# arg1: root of destination filesystem
install-entropy () {
    # copy rngd and libraries to target from current root
    mkdir -p $1/{bin,lib64}
    cp -Ln /lib64/ld-linux-x86-64.so.2 $1/lib64/ 
    cp -Ln /lib64/libc.so.6 $1/lib64/ 
    cp /sbin/rngd $1/bin/rngd

    # TODO: stop assuming sh - can we replace with:
    # a. json config with rtld, rtld args, binary, binary args, chroot?
    # b. Go plugins for tether extensions
    cat - > $1/bin/entropy <<ENTROPY
#!/bin/sh
exec /.tether/lib64/ld-linux-x86-64.so.2 --library-path /.tether/lib64/ /.tether/bin/rngd "\$@"
ENTROPY

    chmod a+x $1/bin/entropy
}

# Usage: copies iptables tools to target system. Creates the following
# executable in the target filesystem to launch iptables:
# /bin/iptables - should exec the target binary with any arguements required
#                 inline and pass through any additional provided
# 
# arg1: root of destination filesystem
#
# ldd of xtables-multi yields the following list of libraries we need to
# copy into our initrd.  We need these binaries in order to call iptables
# before the switch-root.
#                   linux-vdso.so.1 (0x00007ffc94d0d000)
# libip4tc.so.0 => /baz/lib/libip4tc.so.0 (0x00007f97fc721000)
# libip6tc.so.0 => /baz/lib/libip6tc.so.0 (0x00007f97fc519000)
# libxtables.so.11 => /baz/lib/libxtables.so.11 (0x00007f97fc30c000)
# libm.so.6 => /lib64/libm.so.6 (0x00007f97fc00e000)
# libgcc_s.so.1 => /lib64/libgcc_s.so.1 (0x00007f97fbdf7000)
# libc.so.6 => /baz/lib/libc.so.6 (0x00007f97fba53000)
# libdl.so.2 => /baz/lib/libdl.so.2 (0x00007f97fb84f000)
# /lib64/ld-linux-x86-64.so.2 (0x00007f97fc929000)
install-iptables () {
    # copy iptables and all associated libraries to target from current root
    mkdir -p $1/{bin,lib64}
    cp -Ln /lib64/ld-linux-x86-64.so.2 $1/lib64/
    cp -L /sbin/iptables $1/bin/iptables
    # Temp until changing tether exec path
    ln -s bin/iptables $1/iptables

    # TODO: figure out what to do with the /etc/alternatives symlinks
    # just copy the target of the link for now
    # strace -e trace=open shows the following used:
    # /lib64/libip4tc.so.0
    # /lib64/libxtables.so.4
    # /lib64/libm.so.6
    # /lib64/libc.so.6
    # /lib64/libdl.so.2
    # /lib64/xtables/libxt_standard.so
    # /lib64/xtables/libxt_state.so

    cp -Ln /lib64/lib{m.*,m-*,gcc_s*,ip*tc*,xtables*,dl*,c.so*,c-*} $1/lib64/
    cp -a /lib64/xtables $1/lib64/

    # TODO: stop assuming bash - can we replace with:
    # a. json config with rtld, rtld args, binary, binary args, chroot?
    # b. Go plugins for tether extensions
    cat - > $1/bin/iptables-wrapper <<IPTABLES
#!/bin/sh
exec chroot /.tether/ /lib64/ld-linux-x86-64.so.2 /bin/iptables "\$@"
IPTABLES

    chmod a+x $1/bin/iptables-wrapper
}