# Release Details

This documents the process for tagging the VIC repository for a release and publishing the VIC binaries on bintray.com

## Tag The Release

* Follow the [RELEASE-VERSIONING.md](RELEASE-VERSIONING.md) for choosing the new version.
* Gather the Release Notes file.
* Go to https://github.com/vmware/vic/releases and click on Draft a new release.
* Add the tag version that meets the [RELEASE-VERSIONING.md](RELEASE-VERSIONING.md) requirements for this release.
* Title the release "vSphere Integrated Containers Version X.Y.Z" where X.Y.Z meets the versioning requirements. 
* Paste the release note file contents into the Write field and preview.
* If this is release is non-production ready select the "This is a pre-release" box.


## Create the VIC binary package

* Find the desired successful build for the release at
  [https://ci.vmware.run/vmware/vic](https://ci.vmware.run/vmware/vic)
* In the build output, find the SHA256 sum of the `vic_<BUILD_NUMBER>.tar.gz` (search for `shasum`)
* Download the artifact from
    [https://bintray.com/vmware/vic-repo/build](https://bintray.com/vmware/vic-repo/build)
* Verify the SHA256 sum of the downloaded version

  ```
    shasum -a 256 vic_<BUILD_NUMBER>.tar.gz
  ```
* Find your Bintray API Key at [https://bintray.com/profile/edit](https://bintray.com/profile/edit)
* Decide the release version for this release, such as `v0.1.0`. See [RELEASE-VERSIONING.md](RELEASE-VERSIONING.md) for details
* Upload the artifact to the release repo under Download https://bintray.com/vmware/vic/Download

  ```
    curl -T <FILE.EXT> -u<BINTRAY_USERNAME>:<API_KEY> https://api.bintray.com/content/vmware/vic/Download/<VERSION_NAME>/<FILE_TARGET_PATH>
  ```

  Example:
    ```
      curl -T vic_v0.1.0.tar.gz -u<BINTRAY_USERNAME>:<API_KEY> https://api.bintray.com/content/vmware/vic/Download/v0.1.0/vic_v0.1.0.tar.gz
    ```

* Display the latest version under Downloads at https://bintray.com/vmware/vic/Download
    by clicking `Actions` and `Add to download list` for the file at https://bintray.com/vmware/vic/Download/view#files

