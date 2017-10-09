/*
使用tcp实现进程间的通信rpc
 */

package gamerpc

import (
	"net"
	"encoding/binary"
)

// 进程之间通信的数据结构
type RpcCmd struct {
	cmd byte
	nameSize uint32
	nameData string
	argSize uint32
	argData string
}

// 链接客户端
type RpcClient struct {
	host string
	gameID byte
	conn net.Conn
	readBuff []byte			// 读缓冲区
	readSize uint32			// 读缓冲区大小
	writeBuff []byte		// 写缓冲区
	writeSize int			// 写缓冲区大小
}

const CMD_HEAD_SIZE uint32 = 9			// RpcCmd包头长度
const READ_BUFF_SIZE uint32 = 1024		// 读取来自客户端的请求的缓冲区大小
const WRITE_BUFF_SIZE uint32 = 10240	// 发送到客户端的数据的缓冲区大小

// 本地监听的server socket
var g_rpcServer net.Listener
// 链接其他server的socket
var g_rpcClients = make(map[byte]RpcClient)
// 链接到本server的socket
var g_rpcConns = make(map[byte]RpcClient)
// 主机字符串到gsid的映射
var g_host2id = make(map[string]byte)








func (client *RpcClient) Close() {

}

func (client *RpcClient) Cmd(cmd byte, name string, args string) {

}


func (client *RpcClient) onRead() {
	var buff = client.readBuff
	var size, pos uint32
	var name, args string
	for {
		len, err := client.conn.Read(buff[client.readSize:])
		if err != nil || len < 0 {
			client.Close()
			break
		}
		client.readSize += uint32(len)
		for client.readSize >= CMD_HEAD_SIZE {
			pos = 1
			size = binary.LittleEndian.Uint32(buff[1:])
			pos += 4
			name = string(buff[pos:pos+size])
			pos += size
			size = binary.LittleEndian.Uint32(buff[pos:])
			pos += 4
			if client.readSize < pos + size {
				break
			}
			args = string(buff[pos:pos+size])
			pos += size
			client.Cmd(buff[0], name, args)
			// 数据前移，已处理的数据需要丢掉
			copy(buff, buff[pos:])
			client.readSize -= pos
		}
	}
}

func onAccept(host string, conn *net.Conn) {
	var client RpcClient
	client.host = host
	client.gameID = 0
	client.conn = *conn
	client.readBuff = make([]byte, 0, READ_BUFF_SIZE)
	client.readSize = 0
	client.writeBuff = make([]byte, 0, WRITE_BUFF_SIZE)
	g_rpcConns[client.gameID] = client
	client.onRead()
}

func onListen(host string) {
	for {
		conn, err := g_rpcServer.Accept()
		if err != nil {
			continue
		}
		go onAccept(host, &conn)
	}
}

// 初始化本地server
func InitRpcServer(host string, sync bool) error {
	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		return err
	}
	g_rpcServer, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	if sync {
		onListen(host)
	} else {
		go onListen(host)
	}
	return err
}

// 发送请求到指定服务器进程，无返回值
func Send2Host(gameID byte, name string, args string) {

}

// 发送请求到指定服务器进程，有返回值
func Send2HostReturn(gameID byte, name string, args string) {

}