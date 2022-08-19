// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package confs

var (
	Name    = "Vehicle Management System"
	TCPHost = "0.0.0.0"
	TCPPort = 8999

	MaxConn              = 50000
	MaxPacketSize uint32 = 4096

	MaxWorkerPoolSize uint32 = 10
	MaxWorkerTasksLen uint32 = 1024
	MaxMsgChanLen     uint32 = 1024
)
