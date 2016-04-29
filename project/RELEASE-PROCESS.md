# VIC Release Process

This document describes how the VIC project is released.

## Release cycle

VIC is developed with a feature driven model. Requirements are broken down into
features and then sized to ensure they will fit into our monthly releases.  Any
features that do not are then reworked until they will fit. Any features under
active development that do not make the montly release will be pushed into the
next months version.

See [../CONTRIBUTING.md](../CONTRIBUTING.md) for details on how issues are prioritized and tracked.

## Release Versioning

The VIC project follows the general vX.Y.Z versioning scheme. Where X is Major, Y is
Minor and Z is Update. For releases that represent milestones there may be a
-milestone tag added. For example, v1.0.0-beta1 would represent the first beta
release of VIC.

See [RELEASE-VERSIONING.md](RELEASE-VERSIONING.md) for details

## Release details

VIC is released in both source and binary form. The source is tagged using github tagging methods. This is manual for now. 

The binary releases are posted at https://bintray.com/vmware/vic/Download/view

See [RELEASE-DETAILS.md](RELEASE-DETAILS.md) for details on how to package up the VIC binary release.

