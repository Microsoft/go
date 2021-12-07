// This file contains RSA portability wrappers.
// +build linux
// +build !android
// +build !no_openssl
// +build !cmd_go_bootstrap
// +build !msan

#include "goboringcrypto.h"

// Only in BoringSSL.
int _goboringcrypto_RSA_generate_key_fips(GO_RSA *rsa, int size, GO_BN_GENCB *cb)
{
	// BoringSSL's RSA_generate_key_fips hard-codes e to 65537.
	BIGNUM *e = _goboringcrypto_BN_new();
	if (e == NULL)
		return 0;
	int ret = _goboringcrypto_internal_BN_set_word(e, RSA_F4) && _goboringcrypto_internal_RSA_generate_key_ex(rsa, size, e, cb);
	_goboringcrypto_BN_free(e);
	return ret;
}

int _goboringcrypto_RSA_sign_pss_mgf1(GO_RSA *rsa, unsigned int *out_len, uint8_t *out, unsigned int max_out,
		const uint8_t *in, unsigned int in_len, EVP_MD *md, const EVP_MD *mgf1_md, int salt_len) {
	EVP_PKEY_CTX *ctx;
	EVP_PKEY *pkey;
	size_t siglen;

	pkey = _goboringcrypto_EVP_PKEY_new();
	if (!pkey)
		return 0;

	if (_goboringcrypto_EVP_PKEY_set1_RSA(pkey, rsa) <= 0)
		return 0;
	
	ctx = _goboringcrypto_EVP_PKEY_CTX_new(pkey, NULL /* no engine */);
	if (!ctx)
		return 0;

	int ret = 0;

	if (_goboringcrypto_internal_EVP_PKEY_sign_init(ctx) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_PSS_PADDING) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_pss_saltlen(ctx, salt_len) <= 0)
		goto err;
	if (_goboringcrypto_internal_EVP_PKEY_CTX_set_signature_md(ctx, md) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_mgf1_md(ctx, mgf1_md) <= 0)
		goto err;
	
	/* Determine buffer length */
	if (_goboringcrypto_internal_EVP_PKEY_sign(ctx, NULL, &siglen, in, in_len) <= 0)
		goto err;

	if (max_out < siglen)
		goto err;

	if (_goboringcrypto_internal_EVP_PKEY_sign(ctx, out, &siglen, in, in_len) <= 0)
		goto err;

	*out_len = siglen;
	ret = 1;

err:
	_goboringcrypto_EVP_PKEY_CTX_free(ctx);

	return ret;
}

int _goboringcrypto_RSA_verify_pss_mgf1(RSA *rsa, const uint8_t *msg, unsigned int msg_len,
		EVP_MD *md, const EVP_MD *mgf1_md, int salt_len, const uint8_t *sig, unsigned int sig_len) {
	EVP_PKEY_CTX *ctx;
	EVP_PKEY *pkey;

	int ret = 0;

	pkey = _goboringcrypto_EVP_PKEY_new();
	if (!pkey)
		return 0;

	if (_goboringcrypto_EVP_PKEY_set1_RSA(pkey, rsa) <= 0)
		return 0;
	
	ctx = _goboringcrypto_EVP_PKEY_CTX_new(pkey, NULL /* no engine */);
	if (!ctx)
		return 0;

	if (_goboringcrypto_internal_EVP_PKEY_verify_init(ctx) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_padding(ctx, RSA_PKCS1_PSS_PADDING) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_pss_saltlen(ctx, salt_len) <= 0)
		goto err;
	if (_goboringcrypto_internal_EVP_PKEY_CTX_set_signature_md(ctx, md) <= 0)
		goto err;
	if (_goboringcrypto_EVP_PKEY_CTX_set_rsa_mgf1_md(ctx, mgf1_md) <= 0)
		goto err;
	if (_goboringcrypto_internal_EVP_PKEY_verify(ctx, sig, sig_len, msg, msg_len) <= 0)
		goto err;

	ret = 1;

err:
	_goboringcrypto_EVP_PKEY_CTX_free(ctx);

	return ret;
}
