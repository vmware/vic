# Bridge Pool Mask Error
When you create new bridge networks and specify an IP range, you might encounter an error in the log.

## Problem
The error in the log states:

<pre>could not initialize port layer: 
bridge mask is not compatible with bridge pool mask</pre>

## Cause

You specified a `--bridge-network-range` that cannot accommodate a /16 network. By default, the range is 172.16.0.0/12, which can accept 16 /16 networks.

## Solution
Use a bridge network of at least /16 or larger. See [Other Advanced Options](vch_installer_options.md#adv-other) in the VCH Deployment Options section of *vSphere Integrated Containers Installation*.