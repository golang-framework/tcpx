// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ztools

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesByJoin(res *[]byte, d ...[]byte) {
	if len(d) < 1 {
		return
	}

	var buf bytes.Buffer

	for _, b := range d {
		buf.Write(b)
	}

	*res = buf.Bytes()
	buf.Reset()
}

func BytesToIntU(d []byte) (int, error) {
	if len(d) == 3 {
		d = append([]byte{0}, d...)
	}

	buf := bytes.NewBuffer(d)
	switch len(d) {
	case 1:
		var tmp uint8
		err := binary.Read(buf, binary.BigEndian, &tmp)
		return int(tmp), err

	case 2:
		var tmp uint16
		err := binary.Read(buf, binary.BigEndian, &tmp)
		return int(tmp), err

	case 4:
		var tmp uint32
		err := binary.Read(buf, binary.BigEndian, &tmp)
		return int(tmp), err

	default:
		return 0, fmt.Errorf("%s", "BytesToIntU bytes length is invalid.")
	}
}
