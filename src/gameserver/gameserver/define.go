package gameserver

type GameCfg struct {
	GameID byte
	GameIDStr string
	ServerID uint32			// 服务器组编号
	WorkSpace string		// 服务器组的工作路径
	LogPath string 			// 日志文件路径，与工作路径的相对路径
	LocalAddr string		// 对内接口
	LogPrint bool
}
