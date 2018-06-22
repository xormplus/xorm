package xorm

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"math/big"
)

const (
	RSA_PUBKEY_ENCRYPT_MODE = iota //公钥加密
	RSA_PUBKEY_DECRYPT_MODE        //公钥解密
	RSA_PRIKEY_ENCRYPT_MODE        //私钥加密
	RSA_PRIKEY_DECRYPT_MODE        //私钥解密
)

type RsaEncrypt struct {
	PubKey      string
	PriKey      string
	pubkey      *rsa.PublicKey
	prikey      *rsa.PrivateKey
	EncryptMode int
	DecryptMode int
}

func (this *RsaEncrypt) Encrypt(strMesg string) ([]byte, error) {
	var inByte []byte
	var err error
	if this.EncryptMode == RSA_PUBKEY_ENCRYPT_MODE {
		this.pubkey, err = getPubKey([]byte(this.PubKey))
		if err != nil {
			return nil, err
		}
	}

	if this.EncryptMode == RSA_PRIKEY_ENCRYPT_MODE {
		this.prikey, err = getPriKey([]byte(this.PriKey))
		if err != nil {
			return nil, err
		}
	}

	inByte = []byte(strMesg)

	inByte, err = this.Byte(inByte, this.EncryptMode)
	if err != nil {
		return nil, err
	}
	return inByte, nil
}

func (this *RsaEncrypt) Decrypt(crypted []byte) (decrypted []byte, err error) {
	if this.DecryptMode == RSA_PUBKEY_DECRYPT_MODE {
		this.pubkey, err = getPubKey([]byte(this.PubKey))
		if err != nil {
			return nil, err
		}
	}

	if this.DecryptMode == RSA_PRIKEY_DECRYPT_MODE {
		this.prikey, err = getPriKey([]byte(this.PriKey))
		if err != nil {
			return nil, err
		}
	}

	decrypted, err = base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}

	decrypted, err = this.Byte(decrypted, this.DecryptMode)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (this *RsaEncrypt) Byte(in []byte, mode int) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	err := this.IO(bytes.NewReader(in), out, mode)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(out)
}

func (this *RsaEncrypt) IO(in io.Reader, out io.Writer, mode int) error {
	switch mode {
	case RSA_PUBKEY_ENCRYPT_MODE:
		if key, err := this.getPubKey(); err != nil {
			return err
		} else {
			return pubKeyIO(key, in, out, true)
		}
	case RSA_PUBKEY_DECRYPT_MODE:
		if key, err := this.getPubKey(); err != nil {
			return err
		} else {
			return pubKeyIO(key, in, out, false)
		}
	case RSA_PRIKEY_ENCRYPT_MODE:
		if key, err := this.getPriKey(); err != nil {
			return err
		} else {
			return priKeyIO(key, in, out, true)
		}
	case RSA_PRIKEY_DECRYPT_MODE:
		if key, err := this.getPriKey(); err != nil {
			return err
		} else {
			return priKeyIO(key, in, out, false)
		}
	default:
		return errors.New("mode not found")
	}
}

func (this *RsaEncrypt) getPubKey() (*rsa.PublicKey, error) {

	if this.pubkey == nil {
		return nil, ErrPublicKey
	}
	return this.pubkey, nil

}

func (this *RsaEncrypt) getPriKey() (*rsa.PrivateKey, error) {

	if this.prikey == nil {
		return nil, ErrPrivateKey
	}
	return this.prikey, nil
}

//-----------------------------------------

var (
	ErrDataToLarge     = errors.New("message too long for RSA public key size")
	ErrDataLen         = errors.New("data length error")
	ErrDataBroken      = errors.New("data broken, first byte is not zero")
	ErrKeyPairDismatch = errors.New("data is not encrypted by the private key")
	ErrDecryption      = errors.New("decryption error")
	ErrPublicKey       = errors.New("get public key error")
	ErrPrivateKey      = errors.New("get private key error")
)

/*公钥解密*/
func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, ErrDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, ErrDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, ErrDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, ErrKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

/*私钥加密*/
func priKeyEncrypt(rand io.Reader, priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, ErrDataLen
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], hashed)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

/*公钥加密或解密Reader*/
func pubKeyIO(pub *rsa.PublicKey, in io.Reader, out io.Writer, isEncrytp bool) error {
	k := (pub.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	var err error
	size := 0
	for {
		size, err = in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if isEncrytp {
			b, err = rsa.EncryptPKCS1v15(rand.Reader, pub, b)
		} else {
			b, err = pubKeyDecrypt(pub, b)
		}
		if err != nil {
			return err
		}
		if _, err = out.Write(b); err != nil {
			return err
		}
	}
	return nil
}

/*私钥加密或解密Reader*/
func priKeyIO(pri *rsa.PrivateKey, r io.Reader, w io.Writer, isEncrytp bool) error {
	k := (pri.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	buf := make([]byte, k)
	var err error
	var b []byte
	size := 0
	for {
		size, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if isEncrytp {
			b, err = priKeyEncrypt(rand.Reader, pri, b)
		} else {
			b, err = rsa.DecryptPKCS1v15(rand.Reader, pri, b)
		}

		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

/*公钥加密或解密byte*/
func pubKeyByte(pub *rsa.PublicKey, in []byte, isEncrytp bool) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	if len(in) <= k {
		if isEncrytp {
			return rsa.EncryptPKCS1v15(rand.Reader, pub, in)
		} else {
			return pubKeyDecrypt(pub, in)
		}
	} else {
		iv := make([]byte, k)
		out := bytes.NewBuffer(iv)
		if err := pubKeyIO(pub, bytes.NewReader(in), out, isEncrytp); err != nil {
			return nil, err
		}
		return ioutil.ReadAll(out)
	}
}

/*私钥加密或解密byte*/
func priKeyByte(pri *rsa.PrivateKey, in []byte, isEncrytp bool) ([]byte, error) {
	k := (pri.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	if len(in) <= k {
		if isEncrytp {
			return priKeyEncrypt(rand.Reader, pri, in)
		} else {
			return rsa.DecryptPKCS1v15(rand.Reader, pri, in)
		}
	} else {
		iv := make([]byte, k)
		out := bytes.NewBuffer(iv)
		if err := priKeyIO(pri, bytes.NewReader(in), out, isEncrytp); err != nil {
			return nil, err
		}
		return ioutil.ReadAll(out)
	}
}

/*读取公钥*/
func getPubKey(in []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, ErrPublicKey
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	} else {
		return pub.(*rsa.PublicKey), err
	}

}

/*读取私钥*/
func getPriKey(in []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, ErrPrivateKey
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return pri, nil
	}
	pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	} else {
		return pri2.(*rsa.PrivateKey), nil
	}
}

/*从crypto/rsa复制 */
var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

/*从crypto/rsa复制 */
func encrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

/*从crypto/rsa复制 */
func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	if c.Cmp(priv.N) > 0 {
		err = ErrDecryption
		return
	}
	var ir *big.Int
	if random != nil {
		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N)
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}

	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}
	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

/*从crypto/rsa复制 */
func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

/*从crypto/rsa复制 */
func nonZeroRandomBytes(s []byte, rand io.Reader) (err error) {
	_, err = io.ReadFull(rand, s)
	if err != nil {
		return
	}
	for i := 0; i < len(s); i++ {
		for s[i] == 0 {
			_, err = io.ReadFull(rand, s[i:i+1])
			if err != nil {
				return
			}
			s[i] ^= 0x42
		}
	}
	return
}

/*从crypto/rsa复制 */
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

/*从crypto/rsa复制 */
func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		return
	}
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}
