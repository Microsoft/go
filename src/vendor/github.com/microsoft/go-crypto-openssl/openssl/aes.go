// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

//go:build linux && !android
// +build linux,!android

package openssl

// #include "goopenssl.h"
import "C"
import (
	"crypto/cipher"
	"errors"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/microsoft/go-crypto-openssl/openssl/internal/subtle"
)

type aesKeySizeError int

func (k aesKeySizeError) Error() string {
	return "crypto/aes: invalid key size " + strconv.Itoa(int(k))
}

const aesBlockSize = 16

type aesCipher struct {
	key     []byte
	enc_ctx *C.EVP_CIPHER_CTX
	dec_ctx *C.EVP_CIPHER_CTX
	cipher  *C.EVP_CIPHER
}

type extraModes interface {
	// Copied out of crypto/aes/modes.go.
	NewCBCEncrypter(iv []byte) cipher.BlockMode
	NewCBCDecrypter(iv []byte) cipher.BlockMode
	NewCTR(iv []byte) cipher.Stream
	NewGCM(nonceSize, tagSize int) (cipher.AEAD, error)

	// Invented for BoringCrypto.
	NewGCMTLS() (cipher.AEAD, error)
}

var _ extraModes = (*aesCipher)(nil)

func NewAESCipher(key []byte) (cipher.Block, error) {
	c := &aesCipher{key: make([]byte, len(key))}
	copy(c.key, key)

	switch len(c.key) * 8 {
	case 128:
		c.cipher = C._goboringcrypto_EVP_aes_128_ecb()
	case 192:
		c.cipher = C._goboringcrypto_EVP_aes_192_ecb()
	case 256:
		c.cipher = C._goboringcrypto_EVP_aes_256_ecb()
	default:
		return nil, errors.New("crypto/cipher: Invalid key size")
	}

	runtime.SetFinalizer(c, (*aesCipher).finalize)

	return c, nil
}

func (c *aesCipher) finalize() {
	if c.enc_ctx != nil {
		C._goboringcrypto_EVP_CIPHER_CTX_free(c.enc_ctx)
	}
	if c.dec_ctx != nil {
		C._goboringcrypto_EVP_CIPHER_CTX_free(c.dec_ctx)
	}
}

func (c *aesCipher) BlockSize() int { return aesBlockSize }

func (c *aesCipher) Encrypt(dst, src []byte) {
	if subtle.InexactOverlap(dst, src) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	if len(src) < aesBlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < aesBlockSize {
		panic("crypto/aes: output not full block")
	}

	if c.enc_ctx == nil {
		var err error
		c.enc_ctx, err = newCipherCtx(c.cipher, C.GO_AES_ENCRYPT, c.key, nil)
		if err != nil {
			panic(err)
		}
	}

	outlen := C.int(0)
	C._goboringcrypto_EVP_CipherUpdate(c.enc_ctx, (*C.uchar)(unsafe.Pointer(&dst[0])), &outlen, (*C.uchar)(unsafe.Pointer(&src[0])), C.int(aesBlockSize))
	runtime.KeepAlive(c)
}

func (c *aesCipher) Decrypt(dst, src []byte) {
	if subtle.InexactOverlap(dst, src) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	if len(src) < aesBlockSize {
		panic("crypto/aes: input not full block")
	}
	if len(dst) < aesBlockSize {
		panic("crypto/aes: output not full block")
	}
	if c.dec_ctx == nil {
		var err error
		c.dec_ctx, err = newCipherCtx(c.cipher, C.GO_AES_DECRYPT, c.key, nil)
		if err != nil {
			panic(err)
		}
	}

	outlen := C.int(0)
	C._goboringcrypto_EVP_CipherUpdate(c.dec_ctx, (*C.uchar)(unsafe.Pointer(&dst[0])), &outlen, (*C.uchar)(unsafe.Pointer(&src[0])), C.int(aesBlockSize))
	runtime.KeepAlive(c)
}

type aesCBC struct {
	ctx *C.EVP_CIPHER_CTX
}

func (x *aesCBC) BlockSize() int { return aesBlockSize }

func (x *aesCBC) CryptBlocks(dst, src []byte) {
	if subtle.InexactOverlap(dst, src) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	if len(src)%aesBlockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	if len(src) > 0 {
		outlen := C.int(0)
		if C._goboringcrypto_EVP_CipherUpdate(
			x.ctx,
			base(dst), &outlen,
			base(src), C.int(len(src)),
		) != C.int(1) {
			panic("crypto/cipher: CipherUpdate failed")
		}
		runtime.KeepAlive(x)
	}
}

func (x *aesCBC) SetIV(iv []byte) {
	if len(iv) != aesBlockSize {
		panic("cipher: incorrect length IV")
	}
	if C.int(1) != C._goboringcrypto_EVP_CipherInit_ex(x.ctx, nil, nil, nil, (*C.uchar)(unsafe.Pointer(&iv[0])), -1) {
		panic("cipher: unable to initialize EVP cipher ctx")
	}
}

func (c *aesCipher) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	x := new(aesCBC)

	var cipher *C.EVP_CIPHER
	switch len(c.key) * 8 {
	case 128:
		cipher = C._goboringcrypto_EVP_aes_128_cbc()
	case 192:
		cipher = C._goboringcrypto_EVP_aes_192_cbc()
	case 256:
		cipher = C._goboringcrypto_EVP_aes_256_cbc()
	default:
		panic("crypto/boring: unsupported key length")
	}
	var err error
	x.ctx, err = newCipherCtx(cipher, C.GO_AES_ENCRYPT, c.key, iv)
	if err != nil {
		panic(err)
	}

	runtime.SetFinalizer(x, (*aesCBC).finalize)

	return x
}

func (c *aesCBC) finalize() {
	C._goboringcrypto_EVP_CIPHER_CTX_free(c.ctx)
}

func (c *aesCipher) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	x := new(aesCBC)

	var cipher *C.EVP_CIPHER
	switch len(c.key) * 8 {
	case 128:
		cipher = C._goboringcrypto_EVP_aes_128_cbc()
	case 192:
		cipher = C._goboringcrypto_EVP_aes_192_cbc()
	case 256:
		cipher = C._goboringcrypto_EVP_aes_256_cbc()
	default:
		panic("crypto/boring: unsupported key length")
	}

	var err error
	x.ctx, err = newCipherCtx(cipher, C.GO_AES_DECRYPT, c.key, iv)
	if err != nil {
		panic(err)
	}
	if C.int(1) != C._goboringcrypto_EVP_CIPHER_CTX_set_padding(x.ctx, 0) {
		panic("cipher: unable to set padding")
	}

	runtime.SetFinalizer(x, (*aesCBC).finalize)
	return x
}

type aesCTR struct {
	ctx *C.EVP_CIPHER_CTX
}

func (x *aesCTR) XORKeyStream(dst, src []byte) {
	if subtle.InexactOverlap(dst, src) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	if len(src) == 0 {
		return
	}
	C._goboringcrypto_EVP_AES_ctr128_enc(
		x.ctx,
		(*C.uint8_t)(unsafe.Pointer(&src[0])),
		(*C.uint8_t)(unsafe.Pointer(&dst[0])),
		C.size_t(len(src)))
	runtime.KeepAlive(x)
}

func (c *aesCipher) NewCTR(iv []byte) cipher.Stream {
	x := new(aesCTR)

	var cipher *C.EVP_CIPHER
	switch len(c.key) * 8 {
	case 128:
		cipher = C._goboringcrypto_EVP_aes_128_ctr()
	case 192:
		cipher = C._goboringcrypto_EVP_aes_192_ctr()
	case 256:
		cipher = C._goboringcrypto_EVP_aes_256_ctr()
	default:
		panic("crypto/boring: unsupported key length")
	}
	var err error
	x.ctx, err = newCipherCtx(cipher, C.GO_AES_ENCRYPT, c.key, iv)
	if err != nil {
		panic(err)
	}

	runtime.SetFinalizer(x, (*aesCTR).finalize)

	return x
}

func (c *aesCTR) finalize() {
	C._goboringcrypto_EVP_CIPHER_CTX_free(c.ctx)
}

type aesGCM struct {
	key    []byte
	tls    bool
	cipher *C.EVP_CIPHER
}

const (
	gcmBlockSize         = 16
	gcmTagSize           = 16
	gcmStandardNonceSize = 12
)

type aesNonceSizeError int

func (n aesNonceSizeError) Error() string {
	return "crypto/aes: invalid GCM nonce size " + strconv.Itoa(int(n))
}

type noGCM struct {
	cipher.Block
}

func (c *aesCipher) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize != gcmStandardNonceSize && tagSize != gcmTagSize {
		return nil, errors.New("crypto/aes: GCM tag and nonce sizes can't be non-standard at the same time")
	}
	// Fall back to standard library for GCM with non-standard nonce or tag size.
	if nonceSize != gcmStandardNonceSize {
		return cipher.NewGCMWithNonceSize(&noGCM{c}, nonceSize)
	}
	if tagSize != gcmTagSize {
		return cipher.NewGCMWithTagSize(&noGCM{c}, tagSize)
	}
	return c.newGCM(false)
}

func (c *aesCipher) NewGCMTLS() (cipher.AEAD, error) {
	return c.newGCM(true)
}

func (c *aesCipher) newGCM(tls bool) (cipher.AEAD, error) {
	g := &aesGCM{key: c.key, tls: tls}
	switch len(c.key) * 8 {
	case 128:
		g.cipher = C._goboringcrypto_EVP_aes_128_gcm()
	case 192:
		g.cipher = C._goboringcrypto_EVP_aes_192_gcm()
	case 256:
		g.cipher = C._goboringcrypto_EVP_aes_256_gcm()
	default:
		panic("crypto/boring: unsupported key length")
	}

	return g, nil
}

func (g *aesGCM) NonceSize() int {
	return gcmStandardNonceSize
}

func (g *aesGCM) Overhead() int {
	return gcmTagSize
}

// base returns the address of the underlying array in b,
// being careful not to panic when b has zero length.
func base(b []byte) *C.uint8_t {
	if len(b) == 0 {
		return nil
	}
	return (*C.uint8_t)(unsafe.Pointer(&b[0]))
}

func (g *aesGCM) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	if len(nonce) != gcmStandardNonceSize {
		panic("cipher: incorrect nonce length given to GCM")
	}
	if uint64(len(plaintext)) > ((1<<32)-2)*aesBlockSize || len(plaintext)+gcmTagSize < len(plaintext) {
		panic("cipher: message too large for GCM")
	}
	if len(dst)+len(plaintext)+gcmTagSize < len(dst) {
		panic("cipher: message too large for buffer")
	}

	// Make room in dst to append plaintext+overhead.
	ret, out := sliceForAppend(dst, len(plaintext)+gcmTagSize)

	// Check delayed until now to make sure len(dst) is accurate.
	if subtle.InexactOverlap(out, plaintext) {
		panic("cipher: invalid buffer overlap")
	}

	ctx, err := newCipherCtx(g.cipher, C.GO_AES_ENCRYPT, g.key, nonce)
	if err != nil {
		panic(err)
	}
	defer C._goboringcrypto_EVP_CIPHER_CTX_free(ctx)

	var encLen C.int
	// Encrypt additional data.
	if C._goboringcrypto_EVP_EncryptUpdate(ctx, nil, &encLen, base(additionalData), C.int(len(additionalData))) != C.int(1) {
		panic(fail("EVP_CIPHER_CTX_seal"))
	}

	// Encrypt plain text.
	if C._goboringcrypto_EVP_EncryptUpdate(ctx, base(out), &encLen, base(plaintext), C.int(len(plaintext))) != C.int(1) {
		panic(fail("EVP_CIPHER_CTX_seal"))
	}

	// Finalise encryption.
	var encFinalLen C.int
	if C._goboringcrypto_EVP_EncryptFinal_ex(ctx, base(out[encLen:]), &encFinalLen) != C.int(1) {
		panic(fail("EVP_CIPHER_CTX_seal"))
	}
	encLen += encFinalLen

	if C._goboringcrypto_EVP_CIPHER_CTX_ctrl(ctx, C.EVP_CTRL_GCM_GET_TAG, C.int(16), unsafe.Pointer(&out[encLen])) != C.int(1) {
		panic(fail("EVP_CIPHER_CTX_seal"))
	}
	encLen += 16

	if int(encLen) != len(plaintext)+gcmTagSize {
		panic("internal confusion about GCM tag size")
	}
	return ret
}

var errOpen = errors.New("cipher: message authentication failed")

func (g *aesGCM) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if len(nonce) != gcmStandardNonceSize {
		panic("cipher: incorrect nonce length given to GCM")
	}
	if len(ciphertext) < gcmTagSize {
		return nil, errOpen
	}
	if uint64(len(ciphertext)) > ((1<<32)-2)*aesBlockSize+gcmTagSize {
		return nil, errOpen
	}

	tag := ciphertext[len(ciphertext)-gcmTagSize:]
	ciphertext = ciphertext[:len(ciphertext)-gcmTagSize]

	// Make room in dst to append ciphertext without tag.
	ret, out := sliceForAppend(dst, len(ciphertext))

	// Check delayed until now to make sure len(dst) is accurate.
	if subtle.InexactOverlap(out, ciphertext) {
		panic("cipher: invalid buffer overlap")
	}

	ctx, err := newCipherCtx(g.cipher, C.GO_AES_DECRYPT, g.key, nonce)
	if err != nil {
		panic(err)
	}
	defer C._goboringcrypto_EVP_CIPHER_CTX_free(ctx)

	clearAndFail := func(err error) ([]byte, error) {
		for i := range out {
			out[i] = 0
		}
		return nil, err
	}

	// Provide any AAD data.
	var tmplen C.int
	if C._goboringcrypto_EVP_DecryptUpdate(ctx, nil, &tmplen, base(additionalData), C.int(len(additionalData))) != C.int(1) {
		return clearAndFail(errOpen)
	}

	// Provide the message to be decrypted, and obtain the plaintext output.
	var decLen C.int
	if C._goboringcrypto_EVP_DecryptUpdate(ctx, base(out), &decLen, base(ciphertext), C.int(len(ciphertext))) != C.int(1) {
		return clearAndFail(errOpen)
	}

	// Set expected tag value. Works in OpenSSL 1.0.1d and later.
	if C._goboringcrypto_EVP_CIPHER_CTX_ctrl(ctx, C.EVP_CTRL_GCM_SET_TAG, 16, unsafe.Pointer(&tag[0])) != C.int(1) {
		return clearAndFail(errOpen)
	}

	// Finalise the decryption.
	var tagLen C.int
	if C._goboringcrypto_EVP_DecryptFinal_ex(ctx, base(out[int(decLen):]), &tagLen) != C.int(1) {
		return clearAndFail(errOpen)
	}

	if int(decLen+tagLen) != len(ciphertext) {
		panic("internal confusion about GCM tag size")
	}
	return ret, nil
}

// sliceForAppend is a mirror of crypto/cipher.sliceForAppend.
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func newCipherCtx(cipher *C.EVP_CIPHER, mode C.int, key, iv []byte) (*C.EVP_CIPHER_CTX, error) {
	ctx := C._goboringcrypto_EVP_CIPHER_CTX_new()
	if ctx == nil {
		return nil, fail("unable to create EVP cipher ctx")
	}
	if C.int(1) != C._goboringcrypto_EVP_CipherInit_ex(ctx, cipher, nil, base(key), base(iv), mode) {
		return nil, fail("unable to initialize EVP cipher ctx")
	}
	return ctx, nil
}
