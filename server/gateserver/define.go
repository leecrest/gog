package main

import "net"

// 客户端链接定义
type ClientConn struct {
	vfd uint32				// 链接编号
	gameID byte				// 所属服务器编号
	seed int32				// 种子
	conn net.Conn			// 链接对象
	readBuff []byte			// 读缓冲区
	readSize int32			// 读缓冲区大小
	writeBuff []byte		// 写缓冲区
	writeSize int32			// 写缓冲区大小
}


// 游戏服链接对象
type GameConn struct {
	id byte
	conn net.Conn			// 链接对象
	readBuff []byte			// 读缓冲区
	readSize int32			// 读缓冲区大小
	writeBuff []byte		// 写缓冲区
	writeSize int32			// 写缓冲区大小
}


// 常量
const MAX_VFD uint32 = 102400
const READ_BUFF_SIZE int32 = 1024		// 读取来自客户端的请求的缓冲区大小
const WRITE_BUFF_SIZE int32 = 10240		// 发送到客户端的数据的缓冲区大小
const CLIENT_PACK_HEAD int32 = 2		// 包头，2字节表现包体的长度
const CLIENT_PACK_MIN int32 = 3			// 包体最小长度


// 进程之间通信的数据结构
type GameCmd struct {
	cmd byte		// 指令类型
	id uint32		// 操作id
	value byte		// 附加数值
	size uint32		// 附加数据长度
}

const GAME_CMD_SIZE int32 = 10	// sizeof(GameCmd)


// 从vfd中获取netid
func VFD2NID(vfd uint32) (byte) {
	return (byte)(vfd >> 30)
}

// 查询vfd在net中的序号
func VFD2IDX(vfd uint32) (uint32) {
	var tmp uint32 = 0xFFFF
	return (tmp >> 2) & vfd
}