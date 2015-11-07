package msghandler

type MsgHandlerFunc func(uint16, []byte)

var g_msgHandlers = make(map[uint16]MsgHandlerFunc)

func RegisterHandler() {
	CreateFL("serv.log")
	g_msgHandlers[C2S_FINFO] = c2s_finfo
	g_msgHandlers[S2C_FINFO] = s2c_finfo
	g_msgHandlers[C2S_UPDATE_FILE] = c2s_update_file
	g_msgHandlers[S2C_UPDATE_FILE] = s2c_update_file
}

func GetMsgHandler(msgId uint16) MsgHandlerFunc {
	f, ok := g_msgHandlers[msgId]
	if !ok {
		return nil
	}

	return f
}

func HandleMsg(msgId uint16, msg []byte) {
	f := GetMsgHandler(msgId)
	if f == nil {
		return
	}

	f(msgId, msg)
}
