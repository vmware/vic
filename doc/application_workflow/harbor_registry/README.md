# Using vSphere Integrated Container Engine with VMware's Harbor

In this example, we will install VMware's Harbor registry and show how to get vSphere Integrated Container Engine (VIC Engine) 0.8.0 working with Harbor.  With 0.8.0, the engine does not have an install-time mechanism to set up a self-signed certificate so we will show the manual steps for post-install setup as a workaround.  We will not show how to setup Harbor with LDAP.  For that, the reader may visit the [Harbor documentation](https://github.com/vmware/harbor/tree/master/docs) site for more information.  Since there is a lot of documentation on the Harbor site for various setup, we will focus on setting up Harbor with a self-signed certificate and setting up VIC Engine to work with this Harbor instance.

## Prerequisite

The following example requires a vCenter installation.  We will use the OVA deployment mechanism in the vSphere client, and the Harbor OVA setup currently does not work with ESXi.  This will be resolved by 1.0.  We will also use openssl to generate a certificate authority and self sign a server certificate.  We'll show steps to follow, but a background in certificates management should be helpful.

The following will show steps for creating certificates for both static IP and fully qualified domain names (FQDN).

Note: Certificate verification requires all machines using certificates are time/date accurate.  This can be achieved using several options, suchas, vSphere web client, vSphere thick client for Windows or govc. In the following, we deploy this example on a vCenter where all ESXi hosts in the cluster have been set up with NTP and were sync'd prior to installing VIC Engine or Harbor.

## Workflows

1. [Install Harbor with static IP](harbor_with_static_ip.md)
2. [Install Harbor with FQDN](harbor_with_fqdn.md)
3. [Post-Install Usage](post_install_usage.md)