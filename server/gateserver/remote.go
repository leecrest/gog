// 对外监听服务

package main

import (
	"net"
	"math/rand"
	"time"
)


// 全局变量
var g_wanserver net.Listener
var g_clients = make(map[uint32]ClientConn)
var g_index uint32 = 0
var g_clientSize uint32 = 0
var g_netID uint32 = 0


func allocClientVfd() (uint32) {
	if g_clientSize >= MAX_VFD {
		return 0
	}
	var id uint32
	for {
		id = (g_index % MAX_VFD) + 1
		_, err := g_clients[id]
		if err {
			g_index++
		} else {
			return id | (g_netID << 30)
		}
	}
	return 0
}

func sendSeed(client *ClientConn) {
	var buff [5] byte
	buff[0] = 4
	//var ptr = (*int32)(&buff[1])
	//*ptr = client.seed
	//Send2Client(client, buff[:])
}

func encode(seed int32, buff []byte, start int32, stop int32) {
	var tmp int32 = 0
	var i int32 = start
	for i <= stop {
		i++
		tmp = (int32)(buff[i])
		buff[i] = (byte)(tmp ^ seed)
		seed += tmp & 0x03
	}
}

func decode(seed int32, buff []byte, start int32, stop int32) {
	var tmp int32 = 0
	var i int32 = start
	for i <= stop {
		tmp = (int32)(buff[i])
		buff[i] = (byte)(tmp ^ seed)
		tmp = (int32)(buff[i])
		seed += tmp & 0x03
	}
}

func onCloseClient(client *ClientConn) {
	logInfo("closeClient: vfd=%d", client.vfd)
	client.conn.Close()
	delete(g_clients, client.vfd)
}


func onListenClient(addr string) {
	logInfo("Listening at %s\n", addr)
	for {
		conn, err := g_wanserver.Accept()
		if err != nil {
			logInfo("accept error: " + err.Error())
			continue
		}
		var client ClientConn
		client.conn = conn
		client.vfd = allocClientVfd()
		client.gameID = 0
		client.seed = rand.Int31()
		client.readBuff = make([]byte, READ_BUFF_SIZE)
		client.readSize = 0
		client.writeBuff = make([]byte, WRITE_BUFF_SIZE)
		client.writeSize = 0
		g_clients[client.vfd] = client
		logInfo("new vfd: %d", client.vfd)
		sendSeed(&client)
		//启动一个新线程
		go onRecvClient(&client)
	}
}

func onRecvClient(client *ClientConn) {
	logInfo("accept vfd: %d", client.vfd)
	for {
		len, err := client.conn.Read(client.readBuff)
		if err != nil || len < 0 {
			onCloseClient(client)
			break
		}
		client.readSize += (int32)(len)
		var buff = client.readBuff
		var size, total int32
		for client.readSize > CLIENT_PACK_MIN {
			// 拆包粘包处理，客户端的数据包定义：2字节长度 + 数据
			size = (int32)(buff[0])
			size = size << 8 + (int32)(buff[1])
			total = size + CLIENT_PACK_HEAD
			if client.readSize < total {
				break
			}
			// 上行数据需要解密
			decode(client.seed, buff, CLIENT_PACK_HEAD, total)

			// 将数据包发送到对应的game进程
			Send2Game(client.gameID, buff[CLIENT_PACK_HEAD:total])

			// 数据前移，已处理的数据需要丢掉
			var i int32 = 0
			for i < size {
				buff[i] = buff[i + total]
			}
			client.readSize -= total
		}
		//str := string(client.readBuff[0:len])

	}
}


func Send2Vfd(vfd uint32, buff []byte) (bool) {
	client, err := g_clients[vfd]
	if !err {
		return false
	}
	return Send2Client(&client, buff)
}

func Send2Client(client *ClientConn, buff []byte) (bool) {
	_, err := client.conn.Write(buff)
	if err != nil {
		return false
	}
	return true
}

func initRemote(host string) (bool) {
	rand.Seed(time.Now().UnixNano())

	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		logInfo("tcp4 addr error: " + err.Error())
		return false
	}
	g_wanserver, err = net.ListenTCP("tcp", addr)
	if err != nil {
		logInfo("listen error: " + err.Error())
		return false
	}
	go onListenClient(host)
	return true
}