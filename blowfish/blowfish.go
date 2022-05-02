//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #include <openssl/blowfish.h>
import "C"
import "unsafe"

type Blowfish struct {
	key C.BF_KEY
}

func (b Blowfish) BlockSize() int {
	return 8
}

func (b Blowfish) Encrypt(dst, src []byte) {
	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&src[0])),
		(*C.uchar)(unsafe.Pointer(&dst[0])),
		&b.key,
		C.BF_ENCRYPT)
}

func (b Blowfish) Decrypt(dst, src []byte) {
	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&src[0])),
		(*C.uchar)(unsafe.Pointer(&dst[0])),
		&b.key,
		C.BF_DECRYPT)
}

func New(key []byte) *Blowfish {
	res := &Blowfish{}
	C.BF_set_key(&res.key, C.int(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return res
}
