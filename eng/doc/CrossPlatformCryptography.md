# Cross-Platform Cryptography in Microsoft Go

Cryptographic operations in Microsoft Go are done by operating system (OS) libraries. This dependency has advantages:

* Go apps benefit from OS reliability. Keeping cryptography libraries safe from vulnerabilities is a high priority for OS vendors. To do that, they provide updates that system administrators should be applying.
* Go apps have access to FIPS-validated algorithms if the OS libraries are FIPS-validated.

> [!NOTE]
> Starting with Go 1.24, Go will also be FIPS 140-3 compliant, see https://github.com/golang/go/issues/69536.
> If the only reason you are using Microsoft Go is for FIPS 140-3 compliance, you should consider using Microsoft Go 1.24 or later.

Go apps will fall back to native Go implementations if the OS libraries doesn't support the algorithm.
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

|Algorithm                  |Windows            |Linux                |
|---------------------------|-------------------|---------------------|
| MD5                       | ✔️                | ✔️                 |
| SHA-1                     | ✔️                | ✔️                 |
| SHA-2-224                 | ❌                | ✔️                 |
| SHA-2-256                 | ✔️                | ✔️                 |
| SHA-2-384                 | ✔️                | ✔️                 |
| SHA-2-512                 | ✔️                | ✔️                 |
| SHA-2-512_224             | ❌                | ✔️ <sup>1, 2</sup> |
| SHA-2-512_256             | ❌                | ✔️ <sup>1, 2</sup> |
| SHA-3-224                 | ❌                | ❌                 |
| SHA-3-256                 | ❌                | ❌                 |
| SHA-3-384                 | ❌                | ❌                 |
| SHA-3-512                 | ❌                | ❌                 |
| SHAKE-128                 | ❌                | ❌                 |
| SHAKE-256                 | ❌                | ❌                 |
| CSHAKE-128                | ❌                | ❌                 |
| CSHAKE-256                | ❌                | ❌                 |
| HMAC                      | ✔️ <sup>3</sup>   | ✔️ <sup>3</sup>    |

<sup>1</sup>Available starting in Microsoft Go 1.24.
<sup>2</sup>Requires OpenSSL 1.1.1 or later.
<sup>3</sup>The supported hash algorithms are the same as the ones supported as standalone hash functions.

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
| DES-CBC       | ✔️       | ⚠️ <sup>1</sup> |
| DES-ECB       | ✔️       | ⚠️ <sup>1</sup> |
| 3DES-ECB      | ✔️       | ✔️              |
| 3DES-CBC      | ✔️       | ✔️              |
| RC4           | ✔️       | ⚠️ <sup>1</sup> |

<sup>1</sup>When using OpenSSL 3, requires the legacy provider to be enabled.

### AES-GCM keys, nonces, and tags

* Key Sizes

  AES-GCM works with 128, 192, and 256-bit keys.

* Nonce Sizes

  AES-GCM works with 12-byte nonces.

* Tag Sizes
  
  AES-GCM works with 16-byte tags.

## Asymmetric encryption

This section includes the following subsections:

* [RSA](#rsa)
* [ECDSA](#ecdsa)
* [ECDH](#ecdh)
* [Ed25519](#ed25519)
* [DSA](#dsa)

### RSA

This section includes the following packages:

* [crypto/rsa](https://pkg.go.dev/crypto/rsa)

| Padding Mode                      | Windows              | Linux               |
|-----------------------------------|----------------------|---------------------|
| OAEP (MD5)                        | ✔️                   | ✔️                 |
| OAEP (SHA-1)                      | ✔️                   | ✔️                 |
| OAEP (SHA-2)                      | ✔️ <sup>1</sup>      | ✔️ <sup>1</sup>    |
| OAEP (SHA-3)                      | ❌                   | ❌                 |
| PSS (MD5)                         | ✔️                   | ✔️                 |
| PSS (SHA-1)                       | ✔️                   | ✔️                 |
| PSS (SHA-2)                       | ✔️ <sup>1</sup>      | ✔️ <sup>1</sup>    |
| PSS (SHA-3)                       | ❌                   | ❌                 |
| PKCS1v15 Signature (Unhashed)     | ✔️                   | ✔️                 |
| PKCS1v15 Signature (RIPMED160)    | ❌ <sup>2</sup>      | ✔️ <sup>2</sup>    |
| PKCS1v15 Signature (MD4)          | ❌ <sup>2</sup>      | ✔️ <sup>2</sup>    |
| PKCS1v15 Signature (MD5)          | ✔️                   | ✔️                 |
| PKCS1v15 Signature (MD5-SHA1)     | ✔️ <sup>2</sup>      | ✔️ <sup>2</sup>    |
| PKCS1v15 Signature (SHA-1)        | ✔️                   | ✔️                 |
| PKCS1v15 Signature (SHA-2)        | ✔️ <sup>1</sup>      | ✔️ <sup>1</sup>    |
| PKCS1v15 Signature (SHA-3)        | ❌                   | ❌                 |

<sup>1</sup>The supported hash algorithms are the same as the ones supported as standalone hash functions.
<sup>2</sup>Available starting in Microsoft Go 1.24.

#### RSA key sizes

[rsa.GenerateKey](https://pkg.go.dev/crypto/rsa#GenerateKey) only supports the following key sizes (in bits): 2048, 3072, 4096.

Multi-prime RSA keys are not supported.

The RSA key size is subject to the limitations of the underlying cryptographic library. For example, on Windows and when using SCOSSL, the key size should be multiple of 8.

#### PSS salt length

On Windows, when verifying a PSS signature, [rsa.PSSSaltLengthAuto](https://pkg.go.dev/crypto/rsa#pkg-constants) is not supported.

#### Random number generation

For those operations that require random numbers, only the [rand.Reader](https://pkg.go.dev/crypto/rand#Reader) is supported.

### ECDSA

This section includes the following packages:

* [crypto/ecdsa](https://pkg.go.dev/crypto/ecdsa)
* [crypto/elliptic](https://pkg.go.dev/crypto/elliptic)

| Elliptic Curve            | Windows     | Linux        |
|---------------------------|-------------|--------------|
| NIST P-224 (secp224r1)    | ✔️          | ✔️          |
| NIST P-256 (secp256r1)    | ✔️          | ✔️          |
| NIST P-384 (secp384r1)    | ✔️          | ✔️          |
| NIST P-521 (secp521r1)    | ✔️          | ✔️          |

#### Random number generation

For those operations that require random numbers, only the [rand.Reader](https://pkg.go.dev/crypto/rand#Reader) is supported.

### ECDH

This section includes the following packages:

* [crypto/ecdh](https://pkg.go.dev/crypto/ecdsa)

| Elliptic Curve            | Windows     | Linux        |
|---------------------------|-------------|--------------|
| NIST P-224 (secp224r1)    | ✔️          | ✔️          |
| NIST P-256 (secp256r1)    | ✔️          | ✔️          |
| NIST P-384 (secp384r1)    | ✔️          | ✔️          |
| NIST P-521 (secp521r1)    | ✔️          | ✔️          |
| X25519 (curve25519)       | ❌          | ❌          |

#### Random number generation

For those operations that require random numbers, only the [rand.Reader](https://pkg.go.dev/crypto/rand#Reader) is supported.

### Ed25519

This section includes the following packages:

* [crypto/ed25519](https://pkg.go.dev/crypto/ed25519)

| Schemes     | Windows    | Linux         |
|-------------|------------|---------------|
| Ed25519     | ❌         | ✔️           |
| Ed25519ctx  | ❌         | ❌           |
| Ed25519ph   | ❌         | ❌           |

#### Random number generation

For those operations that require random numbers, only the [rand.Reader](https://pkg.go.dev/crypto/rand#Reader) is supported.

### DSA

| Parameters    | Windows     | Linux        |
|---------------|-------------|--------------|
| L1024N160     | ✔️          | ✔️          |
| L2048N224     | ❌          | ✔️          |
| L2048N256     | ✔️          | ✔️          |
| L3072N256     | ✔️          | ✔️          |

## KDF

This section includes the following packages:

* [crypto/hkdf](https://pkg.go.dev/crypto/hkdf)
* [crypto/pbkdf2](https://pkg.go.dev/crypto/pbkdf2)

| Functions     | Windows          | Linux             |
|---------------|------------------|-------------------|
| PBKDF2        | ✔️ <sup>1</sup>  | ✔️ <sup>1</sup>  |
| HKDF          | ✔️ <sup>1</sup>  | ✔️ <sup>1</sup>  |

<sup>1</sup>The supported hash algorithms are the same as the ones supported as standalone hash functions.

## ML-KEM

This section includes the following packages:

* [crypto/mlkem](https://pkg.go.dev/crypto/mlkem)

| Parameters    | Windows     | Linux        |
|---------------|-------------|--------------|
|768            | ❌          | ❌          |
|1024           | ❌          | ❌          |

## TLS

This section includes the following subsections:

* [TLS Versions](#tls-versions)
* [Cipher Suites](#cipher-suites)
* [Curves and Groups](#curves-and-groups)
* [Signature Algorithms](#signature-algorithms)

This section includes the following packages:

* [crypto/tls](https://pkg.go.dev/crypto/tls)

### TLS Versions

| Version        | Windows     | Linux   |
|----------------|-------------|---------|
| SSL 3.0        | ❌          | ❌          |
| TLS 1.0        | ✔️          | ✔️          |
| TLS 1.2        | ✔️          | ✔️          |
| TLS 1.3        | ✔️          | ✔️          |

### Cipher Suites

| Name                                              | Windows     | Linux             |
|---------------------------------------------------|-------------|-------------------|
| TLS_RSA_WITH_RC4_128_SHA                          | ✔️          | ⚠️ <sup>1</sup>  |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA                     | ✔️          | ⚠️ <sup>1</sup>  |
| TLS_RSA_WITH_AES_128_CBC_SHA                      | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_256_CBC_SHA                      | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_128_CBC_SHA256                   | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_128_GCM_SHA256                   | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_256_GCM_SHA384                   | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_RC4_128_SHA                  | ✔️          | ⚠️ <sup>1</sup>  |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA              | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA              | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_RC4_128_SHA                    | ✔️          | ⚠️ <sup>1</sup>  |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA               | ✔️          | ⚠️ <sup>1</sup>  |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA                | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA                | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256           | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256             | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256             | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256           | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384             | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384           | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256       | ❌          | ❌               |
| TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256     | ❌          | ❌               |
| TLS_AES_128_GCM_SHA256                            | ✔️          | ✔️               |
| TLS_AES_256_GCM_SHA384                            | ✔️          | ✔️               |
| TLS_CHACHA20_POLY1305_SHA256                      | ❌          | ❌               |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305              | ❌          | ❌               |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305              | ❌          | ❌               |
| TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305            | ❌          | ❌               |

<sup>1</sup>When using OpenSSL 3, requires the legacy provider to be enabled.

### Curves and Groups

| Name            | Windows     | Linux        |
|-----------------|-------------|--------------|
| CurveP2         | ✔️          | ✔️          |
| CurveP384       | ✔️          | ✔️          |
| CurveP521       | ✔️          | ✔️          |
| X25519          | ❌          | ❌          |
| X25519MLKEM768  | ❌          | ❌          |

### Signature Algorithms

| Name                      | Windows     | Linux        |
|---------------------------|-------------|--------------|
| PKCS1WithSHA256           | ✔️          | ✔️          |
| PKCS1WithSHA384           | ✔️          | ✔️          |
| PKCS1WithSHA512           | ✔️          | ✔️          |
| PSSWithSHA256             | ✔️          | ✔️          |
| PSSWithSHA384             | ✔️          | ✔️          |
| PSSWithSHA512             | ✔️          | ✔️          |
| ECDSAWithP256AndSHA256    | ✔️          | ✔️          |
| ECDSAWithP384AndSHA384    | ✔️          | ✔️          |
| ECDSAWithP521AndSHA512    | ✔️          | ✔️          |
| Ed25519                   | ✔️          | ✔️          |
| PKCS1WithSHA1             | ✔️          | ✔️          |
| ECDSAWithSHA1             | ✔️          | ✔️          |
