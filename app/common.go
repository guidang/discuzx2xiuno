package app

import (
	"encoding/binary"
	"net"
)

const (
	PreInsertErr = "预插入 %s 失败: %s"
	InsertErr = "插入 %s 失败: %s"
	SelectErr = "查询 %s 失败: %s"
	InsertSuccess = "插入 %s, 成功 %d 条"
	ClearErrMsg = "清理 %s 失败: %s"
)

func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}