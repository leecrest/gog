package game

import (
	"net"
	"fmt"
	"encoding/binary"
	"os"
)

var g_gates = make(map[byte]GateClient)

func (gate *GateClient) Close() {

}

func (gate *GateClient) Read(cmd byte, gameID byte, vfd uint32, data []byte) {
	switch cmd {
	case CMD_N2G_SYNCGSID:
		if gameID != g_config.id {
			logError("game id error: local=%d, gate alloc=%d", g_config.id, gameID)
			os.Exit(1)
		}
		break
	case CMD_N2G_VFD_ADD:
		logInfo("vfd add:%d", vfd)
		break
	case CMD_N2G_VFD_DEL:
		logInfo("vfd del:%d", vfd)
		break
	case CMD_N2G_VFD_DATA:
		logInfo("vfd data:%d %s", vfd, data)
		break
	case CMD_N2G_VFD_CLOSE:
		logInfo("vfd close:%d", vfd)
		break
	default:
		break
	}
}



func send2Gate() {

}

func Read(gate *GateClient) {
	var buff = gate.readBuff
	var size, total, vfd uint32
	for {
		len, err := gate.conn.Read(buff[gate.readSize:])
		if err != nil || len < 0 {
			if err != nil {
				logError("gate[%d] read error: %s", gate.gateID, err.Error())
			}
			gate.Close()
			break
		}
		gate.readSize += uint32(len)
		for gate.readSize >= CMD_HEAD_SIZE {
			vfd = binary.LittleEndian.Uint32(buff[2:])
			size = binary.LittleEndian.Uint32(buff[6:])
			total = CMD_HEAD_SIZE + size
			if gate.readSize < total {
				break
			}
			if size > 0 {
				gate.Read(buff[0], buff[1], vfd, buff[CMD_HEAD_SIZE : total])
			} else {
				gate.Read(buff[0], buff[1], vfd, nil)
			}
			// 数据前移，已处理的数据需要丢掉
			copy(buff, buff[total:])
			gate.readSize -= total
		}
	}
}

func onConnected(id byte, conn *net.Conn) {
	var gate GateClient
	gate.gateID = id
	gate.conn = *conn
	gate.readBuff = make([]byte, GAME_READ_BUFF_SIZE)
	gate.readSize = 0
	g_gates[id] = gate
	logInfo("connect to gate:%d", id)
	go Read(&gate)
}

func gateInit() bool {
	var addr string
	var port uint32 = g_config.gatePort
	var gate byte
	for gate = 0; gate < g_config.gateNum; gate++ {
		addr = fmt.Sprintf("%s:%d", g_config.gateIP, port)
		port++
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			logError("connect to gate[%d] error: %s", gate, err.Error())
			return false
		}
		onConnected(gate, &conn)
	}
	return true
}