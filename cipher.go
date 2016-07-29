// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

type Cipher interface {
	Encrypt(strMsg string) ([]byte, error)
	Decrypt(src []byte) (decrypted []byte, err error)
}
