package gate

import (
	"net"
	"encoding/binary"
)

// 全局变量
var g_gameserver net.Listener
var g_games = make([]*GameConn, MAX_GAME)


func allocGameID() (byte) {
	var idx byte = 0
	for idx = 0; idx < MAX_GAME; idx++ {
		if g_games[idx] == nil {
			return idx
		}
	}
	return 0
}

func (game *GameConn) Close() {
	logInfo("close game: %d", game.id)
	g_games[game.id] = nil
}

func (game *GameConn) Send2Self(pack *GameCmd) {
	//var buff = game.writeBuff
	var buff = make([]byte, pack.size+CMD_HEAD_SIZE)
	buff[0] = pack.cmd
	buff[1] = pack.gameID
	binary.LittleEndian.PutUint32(buff[2:], pack.vfd)
	binary.LittleEndian.PutUint32(buff[6:], pack.size)
	if pack.size > 0 {
		copy(buff[CMD_HEAD_SIZE:], pack.data)
	}
	_, err := game.conn.Write(buff)
	if err != nil {
		return
	}
}

func (game *GameConn) Read(cmd byte, gameID byte, vfd uint32, data []byte) {
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



func onReadGame(game *GameConn) {
	var buff = game.readBuff
	var size, total, vfd uint32
	for {
		len, err := game.conn.Read(buff[game.readSize:])
		if err != nil || len < 0 {
			game.Close()
			break
		}
		game.readSize += uint32(len)
		for game.readSize >= CMD_HEAD_SIZE {
			vfd = binary.LittleEndian.Uint32(buff[2:])
			size = binary.LittleEndian.Uint32(buff[6:])
			total = CMD_HEAD_SIZE + size
			if game.readSize < total {
				break
			}
			if size > 0 {
				game.Read(buff[0], buff[1], vfd, buff[CMD_HEAD_SIZE : total])
			} else {
				game.Read(buff[0], buff[1], vfd, nil)
			}
			// 数据前移，已处理的数据需要丢掉
			copy(buff, buff[total:])
			game.readSize -= total
		}
	}
}

func onAcceptGame(conn *net.Conn) {
	var game GameConn
	game.conn = *conn
	game.id = allocGameID()
	game.readBuff = make([]byte, GAME_READ_BUFF_SIZE)
	game.readSize = 0
	game.writeBuff = make([]byte, GAME_WRITE_BUFF_SIZE)
	g_games[game.id] = &game
	logInfo("new game: %d", game.id)

	var pack GameCmd
	pack.cmd = CMD_N2G_SYNCGSID
	pack.vfd = 0
	pack.gameID = game.id
	pack.data = nil
	pack.size = 0
	game.Send2Self(&pack)

	onReadGame(&game)
}

func onListenGame(addr string) {
	logInfo("listening at %s", addr)
	for {
		conn, err := g_gameserver.Accept()
		if err != nil {
			logError("accept error: %s", err.Error())
			continue
		}
		go onAcceptGame(&conn)
	}
}

func initGame(host string, sync bool) bool {
	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		logError("tcp4 addr error: " + err.Error())
		return false
	}
	g_gameserver, err = net.ListenTCP("tcp", addr)
	if err != nil {
		logError("tcp listen error: " + err.Error())
		return false
	}
	if sync {
		onListenGame(host)
	} else {
		go onListenGame(host)
	}
	return true
}
