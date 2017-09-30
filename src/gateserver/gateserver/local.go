package gateserver

import (
	"net"
	"encoding/binary"
)

// 全局变量
var g_lanserver net.Listener
var g_games = make(map[byte]GameConn)
var g_gameID byte = 0


func allocGameID() (byte) {
	var idx byte = 0
	for idx < g_gameID {
		_, err := g_games[idx]
		if !err {
			return idx
		}
		idx++
	}
	g_gameID++
	return g_gameID
}

func (game *GameConn) Close() {
	logInfo("close game: %d", game.id)
	delete(g_games, game.id)
}

func (game *GameConn) Send2Client() {

}

func (game *GameConn) Recv(cmd byte, gameID byte, vfd uint32, data []byte) {
	switch cmd {
	case CMD_G2N_VFD_SEND:
		send2Vfd(vfd, data)
		break
	case CMD_G2N_VFD_SENDS:
		var buf = data[vfd * 4:]
		var i, pos, idx uint32 = 0, 0, 0
		for i = 0; i < vfd; i++ {
			idx = binary.LittleEndian.Uint32(data[pos:pos+4])
			pos += 4
			send2Vfd(idx, buf)
		}
		break
	case CMD_G2N_VFD_GSID:
		var client, err = g_clients[vfd]
		if !err {
			return
		}
		client.ChangeGame(gameID)
		break
	default:
		break
	}
}



func onRecvGame(game *GameConn) {
	logInfo("recv game: " + game.conn.RemoteAddr().String())
	var buff = game.readBuff
	var size, total int
	var vfd uint32
	for {
		len, err := game.conn.Read(buff[game.readSize:])
		if err != nil || len < 0 {
			game.Close()
			break
		}
		game.readSize += len
		for game.readSize >= CMD_HEAD_SIZE {
			vfd = binary.LittleEndian.Uint32(buff[2:])
			size = int(binary.LittleEndian.Uint32(buff[6:]))
			total = CMD_HEAD_SIZE + size
			if game.readSize < total {
				break
			}
			game.Recv(buff[0], buff[1], vfd, buff[CMD_HEAD_SIZE:total])

			// 数据前移，已处理的数据需要丢掉
			copy(buff, buff[total:])
			game.readSize -= total
		}
		//str := string(client.readBuff[0:len])

	}
}

func onListenGame(addr string) {
	logInfo("listening at %s", addr)
	for {
		conn, err := g_lanserver.Accept()
		if err != nil {
			logInfo("accept error: %s", err.Error())
			continue
		}
		var game GameConn
		game.conn = conn
		game.id = allocGameID()
		game.readBuff = make([]byte, READ_BUFF_SIZE)
		game.readSize = READ_BUFF_SIZE
		game.writeBuff = make([]byte, WRITE_BUFF_SIZE)
		game.writeSize = WRITE_BUFF_SIZE
		g_games[game.id] = game
		go onRecvGame(&game)
	}
}

func initLocal(host string) (bool) {
	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		logInfo("tcp4 addr error: " + err.Error())
		return false
	}
	g_lanserver, err = net.ListenTCP("tcp", addr)
	if err != nil {
		logInfo("tcp listen error: " + err.Error())
		return false
	}
	onListenGame(host)
	return true
}