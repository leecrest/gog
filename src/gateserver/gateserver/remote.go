// 对外监听服务

package gateserver

import (
	"net"
	"math/rand"
	"time"
	"encoding/binary"
	"fmt"
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

func (client *ClientConn) SendSeed() {
	var buff = make([]byte, 5)
	buff[0] = 4
	//binary.LittleEndian.PutUint32(buff[1:], client.seed)
	client.Send2Self(buff)
}

func (client *ClientConn) Encode(buff []byte, start int, stop int) {
	var tmp uint32 = 0
	var seed = client.seed
	for i := start; i <= stop; i++ {
		tmp = (uint32)(buff[i])
		buff[i] = (byte)(tmp ^ seed)
		seed += tmp & 0x03
	}
}

func (client *ClientConn) Decode(buff []byte, start int, stop int) {
	var tmp uint32 = 0
	var seed = client.seed
	for i := start; i < stop; i++ {
		tmp = (uint32)(buff[i])
		buff[i] = (byte)(tmp ^ seed)
		tmp = (uint32)(buff[i])
		seed += tmp & 0x03
	}
}

func (client *ClientConn) Close() {
	logInfo("close vfd: %d", client.vfd)
	client.conn.Close()
	delete(g_clients, client.vfd)
	client.Send2Game(CMD_N2G_VFD_CLOSE, client.gameID, nil, 0)
}

func (client *ClientConn) Send2Self(buff []byte) (bool) {
	_, err := client.conn.Write(buff)
	if err != nil {
		return false
	}
	return true
}

func (client *ClientConn) Send2Game(cmd byte, gameID byte, data []byte, size int) bool {
	var game, err = g_games[gameID]
	if !err {
		return false
	}
	var buff = game.writeBuff
	buff[0] = cmd
	buff[1] = gameID
	binary.LittleEndian.PutUint32(buff[2:], client.vfd)
	if data == nil {
		size = 0
	}
	binary.LittleEndian.PutUint32(buff[6:], uint32(size))
	if size > 0 {
		copy(buff[CMD_HEAD_SIZE:], data)
	}
	game.conn.Write(buff[:CMD_HEAD_SIZE + size])
	return true
}

func (client *ClientConn) ChangeGame(gameID byte) {
	client.Send2Game(CMD_N2G_VFD_DEL, client.gameID, nil, 0)
	client.gameID = gameID
	client.Send2Game(CMD_N2G_VFD_ADD, gameID, nil, 0)
}






func send2Vfd(vfd uint32, data []byte) (bool) {
	client, err := g_clients[vfd]
	if !err {
		return false
	}
	return client.Send2Self(data)
}


func onRecvClient(client *ClientConn) {
	logInfo("new vfd: %d", client.vfd)
	client.SendSeed()
	// 通知游戏进程，新链接进入
	client.Send2Game(CMD_N2G_VFD_ADD,0, nil, 0)

	var buff = client.readBuff
	var size, total int
	for {
		len, err := client.conn.Read(buff[client.readSize:])
		if err != nil || len < 0 {
			client.Close()
			break
		}
		client.readSize += len
		for client.readSize > CLIENT_PACK_MIN {
			// 拆包粘包处理，客户端的数据包定义：2字节长度 + 数据
			fmt.Printf("recv client:[%d] %s", client.vfd, string(buff))
			size = (int)(binary.LittleEndian.Uint16(buff))
			total = size + CLIENT_PACK_HEAD
			if client.readSize < total {
				break
			}
			// 上行数据需要解密
			client.Decode(buff, CLIENT_PACK_HEAD, total)
			// 将数据包发送到对应的game进程
			fmt.Printf("recv client:[%d] %s", client.vfd, string(buff[CLIENT_PACK_HEAD:total]))
			client.Send2Game(CMD_N2G_VFD_DATA, 0, buff[CLIENT_PACK_HEAD:total], size)

			// 数据前移，已处理的数据需要丢掉
			for i := 0; i < size; i++ {
				buff[i] = buff[i + total]
			}
			client.readSize -= total
		}

	}
}

func onListenClient(addr string) {
	logInfo("listening at %s", addr)
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
		client.seed = rand.Uint32()
		client.readBuff = make([]byte, READ_BUFF_SIZE)
		client.readSize = 0
		client.writeBuff = make([]byte, WRITE_BUFF_SIZE)
		client.writeSize = 0
		g_clients[client.vfd] = client
		go onRecvClient(&client)
	}
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