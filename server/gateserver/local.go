package main

import (
	"net"
	"fmt"
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

func onCloseGame(game *GameConn) {

}

func CloseGame(game * GameConn) {

}

func onRecvGame(game *GameConn) {
	fmt.Println("[INFO]onRecvGame:" + game.conn.RemoteAddr().String())
	for {
		len, err := game.conn.Read(game.readBuff)
		if err != nil || len < 0 {
			onCloseGame(game)
			break
		}
		game.readSize += (int32)(len)
		var buff = game.readBuff
		var size, total int32
		for game.readSize > CLIENT_PACK_MIN {
			// 拆包粘包处理，客户端的数据包定义：2字节长度 + 数据
			size = (int32)(buff[0])
			size = size << 8 + (int32)(buff[1])
			total = size + CLIENT_PACK_HEAD
			if game.readSize < total {
				break
			}
			// 上行数据需要解密
			//decode(game.seed, buff, CLIENT_PACK_HEAD, total)

			// 将数据包发送到对应的game进程
			//Send2Game(game.gameID, buff[CLIENT_PACK_HEAD:total])

			// 数据前移，已处理的数据需要丢掉
			var i int32 = 0
			for i < size {
				buff[i] = buff[i + total]
			}
			game.readSize -= total
		}
		//str := string(client.readBuff[0:len])

	}
}


func onListenGame(addr string) {
	fmt.Printf("Listening at %s ...\n", addr)
	for {
		conn, err := g_lanserver.Accept()
		if err != nil {
			fmt.Println("[ERROR]Accept: " + err.Error())
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
		//启动一个新线程
		go onRecvGame(&game)
	}
}

func initLocal(host string) (bool) {
	var addr, err = net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("[ERROR]ResolveTCPAddr:" + err.Error())
		return false
	}
	g_lanserver, err = net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("[ERROR]ListenTCP:" + err.Error())
		return false
	}
	onListenGame(host)
	return true
}






func Send2Game(id byte, buff []byte) (bool) {

	//_, err := conn.Write(buff)
	//if err != nil {
	//	return false
	//}
	return true
}