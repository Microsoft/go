# Cross-Platform Cryptography in Microsoft Go

Cryptographic operations in Microsoft Go are delegated to the operating system (OS) libraries in some conditions.
The high level conditions and the benefits of delegating cryptographic operations are described in the [Microsoft Go FIPS README](./fips/README.md).
At a fine-grained level, Go apps will fall back to the native Go implementation of an algorithm if the OS libraries don't support it.
This article identifies the features that are supported on each platform.

This article assumes you have a working familiarity with cryptography in Go.

## Platform support

Microsoft Go supports the following platforms:

### Windows

On Windows, Microsoft Go uses the [CNG library (Cryptography API: Next Generation)](https://learn.microsoft.com/en-us/windows/win32/seccng/cng-portal) for cryptographic operations.
CNG is available since Windows Vista and Windows Server 2008 and it doesn't require any additional installation nor configuration.

### Linux

On Linux, Microsoft Go uses the [OpenSSL crypto library](https://docs.openssl.org/3.0/man7/crypto/) for cryptographic operations.
OpenSSL is normally available on Linux distributions, but it may not be installed by default.
If it is not installed, you can install it using the package manager of your distribution.

OpenSSL 3 implements all the cryptographic algorithms using [Providers](https://docs.openssl.org/3.0/man7/crypto/#providers).
Microsoft Go officially supports the built-in providers and [SCOSSL (SymCrypt provider for OpenSSL)](https://github.com/microsoft/SymCrypt-OpenSSL) v1.6.1 or later.
SCOSSL is expected to be used with the default built-in provider enabled as a fallback (which is the case when using [Azure Linux 3](https://github.com/microsoft/AzureLinux)).

## Table legend

The following legend describes the symbols used in the tables to indicate the level of support for each cryptographic algorithm:

| Symbol | Meaning                                                                                                                      |
|--------|------------------------------------------------------------------------------------------------------------------------------|
| ✔️     | Supported, possibly with minor limitations that don't require special configuration when using the latest Go and OS versions |
| ⚠️     | Supported with limitations that require special configuration action                                                         |
| ❌     | Not supported                                                                                                                |

When an algorithm is not supported or the limitations are exceeded, Microsoft Go will fall back to the Go implementation.

## Hash and Message Authentication Algorithms

This section includes the following packages:

* [crypto/md5](https://pkg.go.dev/crypto/md5)
* [crypto/sha1](https://pkg.go.dev/crypto/sha1)
* [crypto/sha256](https://pkg.go.dev/crypto/sha256)
* [crypto/sha512](https://pkg.go.dev/crypto/sha512)
* [crypto/sha3](https://pkg.go.dev/golang.org/x/crypto/sha3)
* [crypto/hmac](https://pkg.go.dev/crypto/hmac)

|Algorithm                  |Windows             |Linux                |
|---------------------------|--------------------|---------------------|
| MD5                       | ✔️                 | ✔️                 |
| SHA-1                     | ✔️                 | ✔️                 |
| SHA-2-224                 | ❌                 | ✔️                 |
| SHA-2-256                 | ✔️                 | ✔️                 |
| SHA-2-384                 | ✔️                 | ✔️                 |
| SHA-2-512                 | ✔️                 | ✔️                 |
| SHA-2-512_224             | ❌                 | ✔️<sup>1, 2</sup>  |
| SHA-2-512_256             | ❌                 | ✔️<sup>1, 2</sup>  |
| SHA-3-224                 | ❌                 | ❌                 |
| SHA-3-256                 | ❌                 | ❌                 |
| SHA-3-384                 | ❌                 | ❌                 |
| SHA-3-512                 | ❌                 | ❌                 |
| SHAKE-128                 | ❌                 | ❌                 |
| SHAKE-256                 | ❌                 | ❌                 |
| CSHAKE-128                | ❌                 | ❌                 |
| CSHAKE-256                | ❌                 | ❌                 |
| HMAC<sup>3</sup>          | ✔️                 | ✔️                 |

<sup>1</sup>Available starting in Microsoft Go 1.24.

<sup>2</sup>Requires OpenSSL 1.1.1 or later.

<sup>3</sup>Supports only hash algorithms that are supported as standalone hash functions.

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
| AES-GCM       | ✔️<sup>2</sup>       | ✔️<sup>2</sup>              |
| DES-CBC       | ✔️       | ⚠️<sup>1</sup>  |
| DES-ECB       | ✔️       | ⚠️<sup>1</sup>  |
| 3DES-ECB      | ✔️       | ✔️              |
| 3DES-CBC      | ✔️       | ✔️              |
| RC4           | ✔️       | ⚠️<sup>1</sup>  |

<sup>1</sup>When using OpenSSL 3, requires the legacy provider to be enabled.

<sup>2</sup> AES-GCM supports specific keys, nonces, and tags:

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

[rsa.GenerateKey](https://pkg.go.dev/crypto/rsa#GenerateKey) only supports the following key sizes (in bits): 2048, 3072, 4096.

Multi-prime RSA keys are not supported.

The RSA key size is subject to the limitations of the underlying cryptographic library.
For example, on some Windows and SCOSSL configurations, the key size should be multiple of 8.
Please refer to the documentation of the underlying cryptographic library for the specific limitations.

Operations that require random numbers (rand io.Reader) only support [rand.Reader](https://pkg.go.dev/crypto/rand#Reader).

| Padding Mode                           | Windows              | Linux               |
|----------------------------------------|----------------------|---------------------|
| OAEP (MD5)                             | ✔️                   | ✔️                 |
| OAEP (SHA-1)                           | ✔️                   | ✔️                 |
| OAEP (SHA-2)<sup>1</sup>               | ✔️                   | ✔️                 |
| OAEP (SHA-3)                           | ❌                   | ❌                 |
| PSS (MD5)                              | ✔️<sup>3</sup>                   | ✔️                 |
| PSS (SHA-1)                            | ✔️<sup>3</sup>                   | ✔️                 |
| PSS (SHA-2)<sup>1</sup>                | ✔️<sup>3</sup>                   | ✔️                 |
| PSS (SHA-3)                            | ❌                   | ❌                 |
| PKCS1v15 Signature (Unhashed)          | ✔️                   | ✔️                 |
| PKCS1v15 Signature (RIPMED160)         | ❌                   | ✔️<sup>2</sup>     |
| PKCS1v15 Signature (MD4)               | ❌                   | ✔️<sup>2</sup>     |
| PKCS1v15 Signature (MD5)               | ✔️                   | ✔️                 |
| PKCS1v15 Signature (MD5-SHA1)          | ✔️<sup>2</sup>       | ✔️<sup>2</sup>     |
| PKCS1v15 Signature (SHA-1)             | ✔️                   | ✔️                 |
| PKCS1v15 Signature (SHA-2)<sup>1</sup> | ✔️                   | ✔️                 |
| PKCS1v15 Signature (SHA-3)             | ❌                  | ❌                 |

<sup>1</sup>Supports only hash algorithms that are [supported as standalone hash functions](#hash-and-message-authentication-algorithms).

<sup>2</sup>Available starting in Microsoft Go 1.24.

<sup>3</sup>On Windows, when verifying a PSS signature, [rsa.PSSSaltLengthAuto](https://pkg.go.dev/crypto/rsa#pkg-constants) is not supported.

### ECDSA

This section includes the following packages:

* [crypto/ecdsa](https://pkg.go.dev/crypto/ecdsa)
* [crypto/elliptic](https://pkg.go.dev/crypto/elliptic)

Operations that require random numbers (rand io.Reader) only support [rand.Reader](https://pkg.go.dev/crypto/rand#Reader).

| Elliptic Curve            | Windows     | Linux        |
|---------------------------|-------------|--------------|
| NIST P-224 (secp224r1)    | ✔️          | ✔️          |
| NIST P-256 (secp256r1)    | ✔️          | ✔️          |
| NIST P-384 (secp384r1)    | ✔️          | ✔️          |
| NIST P-521 (secp521r1)    | ✔️          | ✔️          |

### ECDH

This section includes the following packages:

* [crypto/ecdh](https://pkg.go.dev/crypto/ecdsa)

Operations that require random numbers (rand io.Reader) only support [rand.Reader](https://pkg.go.dev/crypto/rand#Reader).

| Elliptic Curve            | Windows     | Linux        |
|---------------------------|-------------|--------------|
| NIST P-224 (secp224r1)    | ✔️          | ✔️          |
| NIST P-256 (secp256r1)    | ✔️          | ✔️          |
| NIST P-384 (secp384r1)    | ✔️          | ✔️          |
| NIST P-521 (secp521r1)    | ✔️          | ✔️          |
| X25519 (curve25519)       | ❌          | ❌          |

### Ed25519

This section includes the following packages:

* [crypto/ed25519](https://pkg.go.dev/crypto/ed25519)

Operations that require random numbers (rand io.Reader) only support [rand.Reader](https://pkg.go.dev/crypto/rand#Reader).

| Schemes     | Windows    | Linux         |
|-------------|------------|---------------|
| Ed25519     | ❌         | ✔️           |
| Ed25519ctx  | ❌         | ❌           |
| Ed25519ph   | ❌         | ❌           |

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

<sup>1</sup>Supports only hash algorithms that are [supported as standalone hash functions](#hash-and-message-authentication-algorithms).

## ML-KEM

This section includes the following packages:

* [crypto/mlkem](https://pkg.go.dev/crypto/mlkem)

| Parameters    | Windows     | Linux        |
|---------------|-------------|--------------|
| 768           | ❌          | ❌          |
| 1024          | ❌          | ❌          |

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
| SSL 3.0        | ❌          | ❌     |
| TLS 1.0        | ✔️          | ✔️     |
| TLS 1.2        | ✔️          | ✔️     |
| TLS 1.3        | ✔️          | ✔️     |

### Cipher Suites

| Name                                              | Windows     | Linux             |
|---------------------------------------------------|-------------|-------------------|
| TLS_RSA_WITH_RC4_128_SHA                          | ✔️          | ⚠️<sup>1</sup>   |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA                     | ✔️          | ⚠️<sup>1</sup>   |
| TLS_RSA_WITH_AES_128_CBC_SHA                      | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_256_CBC_SHA                      | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_128_CBC_SHA256                   | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_128_GCM_SHA256                   | ✔️          | ✔️               |
| TLS_RSA_WITH_AES_256_GCM_SHA384                   | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_RC4_128_SHA                  | ✔️          | ⚠️<sup>1</sup>   |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA              | ✔️          | ✔️               |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA              | ✔️          | ✔️               |
| TLS_ECDHE_RSA_WITH_RC4_128_SHA                    | ✔️          | ⚠️<sup>1</sup>   |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA               | ✔️          | ⚠️<sup>1</sup>   |
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
| CurveP256       | ✔️          | ✔️          |
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
| Ed25519                   | ❌          | ✔️          |
| PKCS1WithSHA1             | ✔️          | ✔️          |
| ECDSAWithSHA1             | ✔️          | ✔️          |
