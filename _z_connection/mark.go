// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package _z_connection

//bufPre := make([]byte, 1)
//_, errTag := io.ReadFull(c.Conn, bufPre)
//if errTag != nil {
//	go z.SetOErrToZ(fmt.Sprintf("reader Prefix err:%v", errTag.Error()))
//	break
//}
//
//if bytes.Equal(bufPre, []byte{0x7e}) == false {
//	if bytes.Equal(bufPre, []byte{0x00}) == false {
//		go z.SetOErrToZ(fmt.Sprintf("Tag != 0x7e, sys break, Tag:%v", hex.EncodeToString(bufPre)))
//	}
//	break
//}
//
//var buf bytes.Buffer
//c.ioReader(&buf)
//res := buf.Bytes()
//buf.Reset()
//
//if bytes.Contains(res, []byte{0x7d, 0x02}) {
//	res = bytes.Replace(res, []byte{0x7d, 0x02}, []byte{0x7e}, -1)
//}
//
//if bytes.Contains(res, []byte{0x7d, 0x01}) {
//	res = bytes.Replace(res, []byte{0x7d, 0x01}, []byte{0x7d}, -1)
//}
//
//c.SendReqToTaskQueue(c.MsgType(hex.EncodeToString(res[:2])), res, uint32(len(res)))

//msgHeaderLen := 12
//bufSrcHeader := make([]byte, msgHeaderLen)
//_, errSrcHeader := io.ReadFull(c.Conn, bufSrcHeader)
//if errSrcHeader != nil {
//	go z.SetOErrToZ(fmt.Sprintf("reader Msg Source errSrcHeader err:%v", errSrcHeader.Error()))
//	break
//}
//
//msgDatumxLen, _ := ztools.BytesToIntU(bufSrcHeader[2:4])
//bufSrcDatumx := make([]byte, msgDatumxLen)
//_, errSrcDatumx := io.ReadFull(c.Conn, bufSrcDatumx)
//if errSrcDatumx != nil {
//	go z.SetOErrToZ(fmt.Sprintf("reader Msg Source errSrcDatumx err:%v", errSrcDatumx.Error()))
//	break
//}
//
//bufSuf := make([]byte, 1)
//_, errSuf := io.ReadFull(c.Conn, bufSuf)
//if errSuf != nil {
//	go z.SetOErrToZ(fmt.Sprintf("reader Suffix err:%v", errSuf.Error()))
//	break
//}

//bufSource := make([]byte, 1024)
//numSource, errSource := c.Conn.Read(bufSource)
//if errSource != nil {
//	return
//}
//
//c.SendMsgToLogQueue(
//	hex.EncodeToString(bufSource[1:3]),
//	hex.EncodeToString(bufSource))
//
//arrSource := bytes.Split(bufSource[:numSource], []byte{0x7e})
//if len(arrSource) < 3 {
//	break
//}
//
//startNum := 0
//if len(arrSource[startNum]) != 0 {
//	arrSource = arrSource[1:]
//}
//
//endedNum := len(arrSource) - 1
//if len(arrSource[endedNum]) != 0 {
//	arrSource = arrSource[:endedNum-1]
//}
//
//if len(arrSource) <= 0 {
//	break
//}
//
//res := make([][]byte, 0)
//
//for _, src := range arrSource {
//	if len(src) != 0 {
//		res = append(res, src)
//	}
//}
//
//if len(res) < 1 {
//	break
//}
//
//for _, v := range res {
//	if bytes.Contains(v, []byte{0x7d, 0x02}) {
//		v = bytes.Replace(v, []byte{0x7d, 0x02}, []byte{0x7e}, -1)
//	}
//
//	if bytes.Contains(v, []byte{0x7d, 0x01}) {
//		v = bytes.Replace(v, []byte{0x7d, 0x01}, []byte{0x7d}, -1)
//	}
//
//	c.SendReqToTaskQueue(c.MsgType(hex.EncodeToString(v[:2])), v, uint32(len(v)))
//}
