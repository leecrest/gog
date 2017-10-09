// 对外监听服务

package gate

import (
	"net"
	"math/rand"
	"time"
	"encoding/binary"
	"io"
)


// 全局变量
var g_clientServer net.Listener
var g_clients = make(map[uint32]ClientConn)
var g_clientVfdIdx uint32 = 0
var g_clientSize uint32 = 0
var g_netID uint32 = 0
var g_clientVfds = make([]bool, MAX_VFD)


// vfd的组成：高2bit是网关编号，低30bit是在本网关的编号
func allocClientVfd() (uint32) {
	if g_clientSize >= MAX_VFD {
		return 0
	}
	g_clientVfdIdx++
	if g_clientVfdIdx >= MAX_VFD {
		g_clientVfdIdx = 1
	}
	for {
		if !g_clientVfds[g_clientVfdIdx] {
			return g_clientVfdIdx | (g_netID << 30)
		} else {
			g_clientVfdIdx++
			if g_clientVfdIdx >= MAX_VFD {
				g_clientVfdIdx = 1
			}
		}
	}
	return 0
}

func (client *ClientConn) SendSeed() {
	var buff = make([]byte, 5)
	buff[0] = 4
	binary.LittleEndian.PutUint32(buff[1:], client.seed)
	client.Send2Self(buff)
}

func (client *ClientConn) Encode(buff []byte, start uint32, stop uint32) {
	var tmp uint32 = 0
	var seed = client.seed
	for i := start; i <= stop; i++ {
		tmp = (uint32)(buff[i])
		buff[i] = (byte)(tmp ^ seed)
		seed += tmp & 0x03
	}
}

func (client *ClientConn) Decode(buff []byte, start uint32, stop uint32) {
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
	var vfd = client.vfd
	logInfo("close vfd: %d", vfd)
	client.conn.Close()
	delete(g_clients, vfd)
	g_clientVfds[VFD2IDX(vfd)] = false
	client.Send2Game(CMD_N2G_VFD_CLOSE, client.gameID, nil)
}

func (client *ClientConn) Send2Self(buff []byte) (bool) {
	_, err := client.conn.Write(buff)
	if err != nil {
		return false
	}
	return true
}


func (client *ClientConn) Send2Game(cmd byte, from byte, data []byte) bool {
	var game = g_games[client.gameID]
	if game == nil {
		return false
	}
	var pack GameCmd
	pack.cmd = cmd
	pack.gameID = from
	pack.vfd = client.vfd
	pack.size = uint32(len(data))
	pack.data = data
	game.Send2Self(&pack)
	return true
}


func (client *ClientConn) ChangeGame(gameID byte) {
	if client.gameID == gameID {
		return
	}
	client.Send2Game(CMD_N2G_VFD_DEL, client.gameID, nil)
	client.gameID = gameID
	client.Send2Game(CMD_N2G_VFD_ADD, gameID, nil)
}






func send2Vfd(vfd uint32, data []byte) (bool) {
	client, err := g_clients[vfd]
	if !err {
		return false
	}
	return client.Send2Self(data)
}


func onReadClient(client *ClientConn) {
	client.SendSeed()
	// 通知游戏进程，新链接进入
	client.Send2Game(CMD_N2G_VFD_ADD,0, nil)

	var buff = client.readBuff
	var size, total uint32
	for {
		len, err := client.conn.Read(buff[client.readSize:])
		if err != nil || len < 0 {
			if err != io.EOF {
				logError("read error: %s", err.Error())
			}
			client.Close()
			break
		}
		client.readSize += uint32(len)
		for client.readSize > CLIENT_PACK_MIN {
			// 拆包粘包处理，客户端的数据包定义：2字节长度 + 数据
			size = (uint32)(binary.LittleEndian.Uint16(buff))
			total = size + CLIENT_PACK_HEAD
			if client.readSize < total {
				break
			}
			// 上行数据需要解密
			client.Decode(buff, CLIENT_PACK_HEAD, total)
			// 将数据包发送到对应的game进程
			client.Send2Game(CMD_N2G_VFD_DATA, 0, buff[CLIENT_PACK_HEAD:total])

			// 数据前移，已处理的数据需要丢掉
			copy(buff[:], buff[total:])
			client.readSize -= total
		}

	}
}

func onAcceptClient(conn net.Conn) {
	var vfd = allocClientVfd()
	if vfd <= 0 {
		logError("allocClientVfd error, vfd is full!")
		conn.Close()
		return
	}
	var client ClientConn
	client.conn = conn
	client.vfd = vfd
	client.gameID = 0
	client.seed = 0//rand.Uint32()
	client.readBuff = make([]byte, CLIENT_READ_BUFF_SIZE)
	client.readSize = 0
	g_clients[vfd] = client
	g_clientVfds[VFD2IDX(vfd)] = true
	logInfo("new vfd: %d", vfd)
	onReadClient(&client)
}

func onListenClient(addr string) {
	logInfo("listening at %s", addr)
	for {
		conn, err := g_clientServer.Accept()
		if err != nil {
			logError("accept error: " + err.Error())
			continue
		}
		go onAcceptClient(conn)
	}
}

func initClient(host string, sync bool) bool {
	rand.Seed(time.Now().UnixNano())

	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		logError("tcp4 addr error: " + err.Error())
		return false
	}
	g_clientServer, err = net.ListenTCP("tcp", addr)
	if err != nil {
		logError("tcp listen error: " + err.Error())
		return false
	}
	if sync {
		onListenClient(host)
	} else {
		go onListenClient(host)
	}
	return true
}