# VIC Engine Release Process

## Release cycle

VIC Engine is developed using a feature driven model. Requirements are broken
down into features and then sized to ensure they will fit into our monthly
releases.  Any features that do not are then reworked until they will fit. Any
features under active development that do not make the monthly release will be
pushed into the next months version.  

See [CONTRIBUTING.md](../../CONTRIBUTING.md#reporting-bugs-and-creating-issues)
for details on how issues are prioritized and tracked.

## Release Versioning

The VIC Engine project follows Semantic Versioning 2.0.0. This is described in
detail at http://semver.org/

## Example VIC Engine Versions Major revisions:

    v1.0.0-beta1 v1.0.0-beta2 v1.0.0-beta v1.0.0 v2.0.0-alpha v2.0.0-beta
v2.0.0

Minor revisions:

    v1.1.0-beta v1.1.0 v2.2.0

Update or patch revisions:

    v1.12.1 v1.12.20

## Release details

VIC Engine is released in both source and binary form. The source is tagged
using github tagging methods. This is manual for now.

The binary releases are posted at https://bintray.com/vmware/vic/Download/view

### Update README.md and documentation

The main repo README contains a project status relating to the latest tagged
release along with guidance on how to build, deploy, et al. The latter should
be updated by any commits changing those workflows, but the status and what's
new needs to be addressed as part of the release process.

* Create a PR, "Release x.y", with the corresponding doc updates once tagging
  is imminent merge that PR as the last thing that occurs prior to tagging tag
  the release as described in the next section.

### Tag The Release

* Follow the above Release Versioning for choosing the new version.  Gather the
  Release Notes file.  Go to https://github.com/vmware/vic/releases and click
  on Draft a new release.  Add the tag version that meets the requirements for
  this release.
* Title the release "vSphere Integrated Containers Engine Version X.Y.Z" where 
  X.Y.Z meets the versioning requirements.
* Paste the release note file contents into the Write field and preview.
* If this is release is non-production ready select the "This is a pre-release"
  box.

### Create the VIC Engine binary package

* Find the desired successful build for the release at
  [https://ci.vmware.run/vmware/vic](https://ci.vmware.run/vmware/vic)
* In the build output, find the SHA256 sum of the `vic_<BUILD_NUMBER>.tar.gz`
  (search for `shasum`)
* Download the artifact from
    [https://bintray.com/vmware/vic-repo/build](https://bintray.com/vmware/vic-repo/build)
* Verify the SHA256 sum of the downloaded version
  ``` shasum -a 256 vic_<BUILD_NUMBER>.tar.gz ```
* Find your Bintray API Key at
    [https://bintray.com/profile/edit](https://bintray.com/profile/edit)
* Decide the release version for this release, such as `v0.1.0`. See the above
  Release Versioning for details.
* Upload the artifact to the release repo under Download
    https://bintray.com/vmware/vic/Download

``` 
curl -T <FILE.EXT> -u<BINTRAY_USERNAME>:<API_KEY>
https://api.bintray.com/content/vmware/vic/Download/<VERSION_NAME>/<FILE_TARGET_PATH>
```
  Example: 
``` 
curl -T vic_0.1.0.tar.gz -u<BINTRAY_USERNAME>:<API_KEY>
https://api.bintray.com/content/vmware/vic/Download/v0.1.0/vic_0.1.0.tar.gz
```
* Display the latest version under Downloads at
  https://bintray.com/vmware/vic/Download by clicking `Actions` and `Add to
  download list` for the file at https://bintray.com/vmware/vic/Download/view#files

TODO: Automate bintray upload
