- [Issues](#sec-1)
- [Related resources](#sec-2)
- [Approach / Implementation](#sec-3)
  - [High level overview](#sec-3-1)
  - [Feature development will satisfy [issue 4689](https://github.com/vmware/vic/issues/4689)](#sec-3-2)
  - [Use Docker's libtrust for validating the signatures](#sec-3-3)
    - [Docker provides functions for verifying signed v1 manifests](#sec-3-3-1)
    - [`imagec` provides functionality to fetch v1 and v2 manifests](#sec-3-3-2)
  - [Call into JWS code from imagec](#sec-3-4)
    - [Verify each layer based on its signature [in PullImage](file:///home/ian/go/src/github.com/vmware/vic/lib/imagec/imagec.go)](#sec-3-4-1)
    - [Fail before downloading the layers if possible](#sec-3-4-2)
    - [VerifyManifestSignature](#sec-3-4-3)
- [Testing and Acceptance Plan](#sec-4)
  - [Issues #[4691](https://github.com/vmware/vic/issues/4691) and [4690](https://github.com/vmware/vic/issues/4690) will be satisfied via feature testing](#sec-4-1)
  - [issue 4689 should provide unit tests](#sec-4-2)
    - [`VerifyManifestSignature` should be well defined:](#sec-4-2-1)
    - [Unit tests should be able to provide pre-created manifests and see if they validate](#sec-4-2-2)


# Issues<a id="sec-1"></a>

1.  [Image manifest JWS verification (epic)](https://github.com/vmware/vic/issues/1331)
2.  [Validate JWS signatures of image manifests](https://github.com/vmware/vic/issues/4689)
3.  [Create malformed image layer](https://github.com/vmware/vic/issues/4690)
    -   We might be able to reuse the one from [the tarball exploit's pull case](https://github.com/vmware/vic/pull/7408)
4.  [Create malformed image digest](https://github.com/vmware/vic/issues/4691)
    -   Again, see the test from the [the tarball exploit's pull case](https://github.com/vmware/vic/pull/7408).

# Related resources<a id="sec-2"></a>

-   [Content trust in Docker](https://docs.docker.com/engine/security/trust/content_trust/)
-   [Docker manifest signing definition](https://docs.docker.com/registry/spec/manifest-v2-1/#manifest-field-descriptions)
    -   Schema Version 1 Example:
        
        ```json
        "signatures": [
           {
              "header": {
                 "jwk": {
                    "crv": "P-256",
                    "kid": "OD6I:6DRK:JXEJ:KBM4:255X:NSAA:MUSF:E4VM:ZI6W:CUN2:L4Z6:LSF4",
                    "kty": "EC",
                    "x": "3gAwX48IQ5oaYQAYSxor6rYYc_6yjuLCjtQ9LUakg4A",
                    "y": "t72ge6kIA1XOjqjVoEOiPPAURltJFBMGDSQvEGVB010"
                 },
                 "alg": "ES256"
              },
              "signature": "XREm0L8WNn27Ga_iE_vRnTxVMhhYY0Zst_FfkKopg6gWSoTOZTuW4rK0fg_IqnKkEKlbD83tD46LKEGi5aIVFg",
              "protected": "eyJmb3JtYXRMZW5ndGgiOjY2MjgsImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNS0wNC0wOFQxODo1Mjo1OVoifQ"
           }
        ]
        ```
    -   Schema Version 2
        -   Still looking for an example (or trying to create one) of a v2 signed manifest
-   [JWS definition](https://tools.ietf.org/html/rfc7515)
    -   [Jump to the section on validation](https://tools.ietf.org/html/rfc7515#section-5.2)
-   [Skopeo: working with remote registries](https://github.com/projectatomic/skopeo)

# Approach / Implementation<a id="sec-3"></a>

## High level overview<a id="sec-3-1"></a>

1.  Fetch manifest
2.  Determine version (schema v1 or v2)
3.  Check if the manifest is signed
4.  If the manifest is unsigned, fail or warn depending upon VIC configuration and bail or continue pull as necessary
5.  If the manifest is signed and but the signature is invalid, reject the pull
6.  If the manifest is signed and valid but pulled images do not match their provided blobSums, reject the pull & clean up any valid images that were written to disk
7.  If the manifest is signed and all of the layers match, continue to unpack the layers on the filesystem

## Feature development will satisfy [issue 4689](https://github.com/vmware/vic/issues/4689)<a id="sec-3-2"></a>

## Use Docker's libtrust for validating the signatures<a id="sec-3-3"></a>

### Docker provides functions for verifying signed v1 manifests<a id="sec-3-3-1"></a>

This method returns a list of public keys used to sign and takes a `*SignedManifest` <file:///home/ian/go/src/github.com/vmware/vic/vendor/github.com/docker/distribution/manifest/schema1/verify.go>

```go
// Verify verifies the signature of the signed manifest returning the public
// keys used during signing.
func Verify(sm *SignedManifest) ([]libtrust.PublicKey, error) {
  js, err := libtrust.ParsePrettySignature(sm.all, "signatures")
  if err != nil {
    logrus.WithField("err", err).Debugf("(*SignedManifest).Verify")
    return nil, err
  }

  return js.Verify()
```

1.  Unfortunately I have not found whether this is usable for v2 and manifest/schema2 does not have a verify.go file

    It may be the case that the v1 JWS verification routine can also be used for v2 as long as we can identify where v2 manifests store their JWS signatures (see [2](#org1883495))

### `imagec` provides functionality to fetch v1 and v2 manifests<a id="sec-3-3-2"></a>

1.  `ic.pullManifest` in `imagec` as shown below provides one of these types:

    ```go
    ImageManifestSchema1 *Manifest
    ImageManifestSchema2 *schema2.DeserializedManifest
    ```
    
    Obviously we'll have to add the fields necessary for the signature to our representation of the V1 Manfiest:
    
    ```go
    // Manifest represents the Docker Manifest file
    type Manifest struct {
      Name     string    `json:"name"`
      Tag      string    `json:"tag"`
      Digest   string    `json:"digest,omitempty"`
      FSLayers []FSLayer `json:"fsLayers"`
      History  []History `json:"history"`
      // ignoring signatures
    }
    ```

2.  Whereas the v2 manifest representation comes out of [vendored code](https://github.com/vmware/vic/blob/master/vendor/github.com/docker/distribution/manifest/schema1/verify.go):

    ```go
    // Manifest defines a schema2 manifest.
    type Manifest struct {
      manifest.Versioned
    
      // Config references the image configuration as a blob.
      Config distribution.Descriptor `json:"config"`
    
      // Layers lists descriptors for the layers referenced by the
      // configuration.
      Layers []distribution.Descriptor `json:"layers"`
    }
    ```

## Call into JWS code from imagec<a id="sec-3-4"></a>

### Verify each layer based on its signature [in PullImage](file:///home/ian/go/src/github.com/vmware/vic/lib/imagec/imagec.go)<a id="sec-3-4-1"></a>

Definitely after this codeblock

```go
// Pull the image manifest
if err := ic.pullManifest(ctx); err != nil {
  return err
}
log.Infof("Manifest for image = %#v", ic.ImageManifestSchema1)
```

Possibly before this one (but we might need `layers` below)

```go
// Get layers to download from manifest
layers, err := ic.LayersToDownload()
if err != nil {
  return err
}
ic.ImageLayers = layers
```

but <span class="underline">definitely</span> before this gets called (immediately after the above)

```go
// Download all the layers
if err := ldm.DownloadLayers(ctx, ic); err != nil {
  return err
```

we need to define a function called something like `VerifyManifestSignature` which will return `err != nil` (or `false`) if the manifest appears to be tampered with, judging by a signature failure.

### TODO Fail before downloading the layers if possible<a id="sec-3-4-2"></a>

-   This should be possible because the signature is computed using the blobSums and the layers have to match their sums
-   If the signature is bad, we can bail before downloading images
-   If the signature is correct but then we receive data that does not match the provided blobSum in the manifest, we *should already fail* in this case, but we should take the opportunity to:
    -   It sounds like we do not currently fail if the shasum from the downloaded image does not match the digest provided
    -   Verify this behavior occurs and add it if it does not
    -   Verify that images that do not complete a `pull` do not leave around old layers
    -   Verify that this behavior is well-tested to prevent regressions

### VerifyManifestSignature<a id="sec-3-4-3"></a>

VerifyManifestSignature will look roughly like this.

```go
func (ic *ImageCache) VerifyManifestSignature(manifest string) bool {
/*  Read in the manifest from memory (ic.ImageManifestSchema1)
    Pass it to the Docker code, and look at the results to determine whether the signature is valid or not
    Return false if invalid, true if it is
*/
}
```

# Testing and Acceptance Plan<a id="sec-4"></a>

## Issues #[4691](https://github.com/vmware/vic/issues/4691) and [4690](https://github.com/vmware/vic/issues/4690) will be satisfied via feature testing<a id="sec-4-1"></a>

-   A layer is an uncompresed tar
-   A blobSum is a sha256 sum
-   Given the above:
    1.  The test should be provided with a tarball of some files with its sha256 sum
    2.  The test should verify that incorrect blobSum *or* file results in a pull failure
    3.  The test should verify that it is not possible to inject the tarball into a `pull` operation even if the blobSum matches the injected data.
    4.  [The MITM test in the Docker Pull suite](file:///home/ian/go/src/github.com/vmware/vic/tests/test-cases/Group1-Docker-Commands/1-02-Docker-Pull.robot) will fail once this change is properly implemented, because that test exploits the fact that #3 is currently not met. The test injects a layer w/ a correct blobSum. This should not be possible. That test will have to be redesigned.

## issue 4689 should provide unit tests<a id="sec-4-2"></a>

### `VerifyManifestSignature` should be well defined:<a id="sec-4-2-1"></a>

-   Input should be a manifest
-   Output should be either a bool (true or false representing validly or invalidly signed manifest) or an err (`nil` to represent success, some `err` to report a failed verification)

### Unit tests should be able to provide pre-created manifests and see if they validate<a id="sec-4-2-2"></a>
