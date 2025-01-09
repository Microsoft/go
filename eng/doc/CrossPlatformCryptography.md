# Cross-Platform Cryptography in Microsoft Go

Cryptographic operations in Microsoft Go are done by operating system (OS) libraries. This dependency has advantages:

* Go apps benefit from OS reliability. Keeping cryptography libraries safe from vulnerabilities is a high priority for OS vendors. To do that, they provide updates that system administrators should be applying.
* Go apps have access to FIPS-validated algorithms if the OS libraries are FIPS-validated.

> [!NOTE]
> Starting with Go 1.24, Go will also be FIPS 140-3 compliant, see https://github.com/golang/go/issues/69536.
> If the only reason you are using Microsoft Go is for FIPS 140-3 compliance, you should consider using Microsoft Go 1.24 or later.

Go apps will fall back to native Go implementations if the OS libraries are not available.
This article identifies the features that are supported on each platform.

This article assumes you have a working familiarity with cryptography in Go.

## Platform support

Microsoft Go supports the following platforms:

### Windows

On Windows, Microsoft Go uses the [Cryptography API: Next Generation](https://learn.microsoft.com/en-us/windows/win32/seccng/cng-portal) library, CNG from now on, for cryptographic operations.
CNG is available since Windows Vista and Windows Server 2008 and it doesn't require any additional installation nor configuration.

### Linux

On Linux, Microsoft Go uses the [OpenSSL crypto library](https://docs.openssl.org/3.0/man7/crypto/) library, OpenSSL from now on, for cryptographic operations.
OpenSSL is normally available on Linux distributions, but it may not be installed by default.
If it is not installed, you can install it using the package manager of your distribution.

OpenSSL 3 implements all the cryptographic algorithms using [Providers](https://docs.openssl.org/3.0/man7/crypto/#providers).
Microsoft Go officially supports the built-in providers and the [SymCrypt provider](https://github.com/microsoft/SymCrypt-OpenSSL), SCOSSL from now on.
The minimum SCOSSL version required is v1.6.1.
The following tables assume that the SCOSSL provider is used together with the built-in providers.

## Hash and Message Authentication Algorithms

This section includes the following packages:

* [crypto/md5](https://pkg.go.dev/crypto/md5)
* [crypto/sha1](https://pkg.go.dev/crypto/sha1)
* [crypto/sha256](https://pkg.go.dev/crypto/sha256)
* [crypto/sha512](https://pkg.go.dev/crypto/sha512)
* [crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3)
* [crypto/hmac](https://pkg.go.dev/crypto/hmac)

|Algorithm                  |Windows   |Linux                |
|---------------------------|----------|---------------------|
|MD5                        | ✔️       | ✔️                 |
|SHA-1                      | ✔️       | ✔️                 |
|SHA-2-224                  | ❌       | ✔️                 |
|SHA-2-256                  | ✔️       | ✔️                 |
|SHA-2-384                  | ✔️       | ✔️                 |
|SHA-2-512                  | ✔️       | ✔️                 |
|SHA-2-512_224              | ✔️       | ✔️ <sup>1, 2</sup> |
|SHA-2-512_256              | ❌       | ✔️ <sup>1, 2</sup> |
|SHA-3-224                  | ❌       | ❌                 |
|SHA-3-256                  | ❌       | ❌                 |
|SHA-3-384                  | ❌       | ❌                 |
|SHA-3-512                  | ❌       | ❌                 |
|SHAKE-128                  | ❌       | ❌                 |
|SHAKE-256                  | ❌       | ❌                 |
|CSHAKE-128                 | ❌       | ❌                 |
|CSHAKE-256                 | ❌       | ❌                 |
|HMAC-MD5                   | ✔️       | ✔️                 |
|HMAC-SHA-1                 | ✔️       | ✔️                 |
|HMAC-SHA-2-224             | ❌       | ✔️                 |
|HMAC-SHA-2-256             | ✔️       | ✔️                 |
|HMAC-SHA-2-384             | ✔️       | ✔️                 |
|HMAC-SHA-2-512             | ✔️       | ✔️                 |
|HMAC-SHA-2-512_224         | ✔️       | ✔️ <sup>1, 2</sup> |
|HMAC-SHA-2-512_256         | ❌       | ✔️ <sup>1, 2</sup> |
|HMAC-SHA-3-224             | ❌       | ❌                 |
|HMAC-SHA-3-256             | ❌       | ❌                 |
|HMAC-SHA-3-384             | ❌       | ❌                 |
|HMAC-SHA-3-512             | ❌       | ❌                 |

<sup>1</sup>Available starting in Microsoft Go 1.24.
<sup>2</sup>Requires OpenSSL 1.1.1 or later.

## Symmetric encryption

This section includes the following packages:

* [crypto/aes](https://pkg.go.dev/crypto/aes)
* [crypto/cipher](https://pkg.go.dev/crypto/cipher)
* [crypto/des](https://pkg.go.dev/crypto/des)
* [crypto/rc4](https://pkg.go.dev/crypto/rc4)

| Cipher + Mode | Windows  | Linux            |
|---------------|----------|------------------|
| AES-ECB       | ✔️       | ✔️              |
| AES-CBC       | ✔️       | ✔️              |
| AES-CTR       | ❌       | ✔️              |
| AES-CFB       | ❌       | ❌              |
| AES-OFB       | ❌       | ❌              |
| AES-GCM       | ✔️       | ✔️              |
| DES-CBC       | ✔️       | ✔️ <sup>1</sup> |
| DES-ECB       | ✔️       | ✔️ <sup>1</sup> |
| 3DES-ECB      | ✔️       | ✔️              |
| 3DES-CBC      | ✔️       | ✔️              |
| RC4           | ✔️       | ✔️ <sup>1</sup> |

<sup>1</sup>When using OpenSSL 3, requires the legacy provider to be enabled.

### AES-GCM keys, nonces, and tags

* Key Sizes

  AES-GCM works with 128, 192, and 256-bit keys.

* Nonce Sizes

  AES-GCM works with 12-byte nonces.

* Tag Sizes
  
  AES-GCM works with 16-byte tags.