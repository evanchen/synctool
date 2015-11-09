package msghandler

import (
	"errors"
	"fmt"
	"io"
	"net"
	"protocol"
)

type MsgHandlerFunc func(uint16, []byte, net.Conn)

var g_msgHandlers = make(map[uint16]MsgHandlerFunc)

type Marshaller interface {
	Marshal() []byte
}

func RegisterHandler() {
	g_msgHandlers[C2S_FINFO] = c2s_finfo
	g_msgHandlers[S2C_FINFO] = s2c_finfo
	g_msgHandlers[C2S_UPDATE_FILE] = c2s_update_file
	g_msgHandlers[S2C_UPDATE_FILE] = s2c_update_file
	g_msgHandlers[S2C_DONE] = s2c_done
}

func GetMsgHandler(msgId uint16) MsgHandlerFunc {
	f, ok := g_msgHandlers[msgId]
	if !ok {
		fmt.Println("wtf: ", msgId, g_msgHandlers)
		return nil
	}

	return f
}

func HandleMsg(msgId uint16, msg []byte, conn net.Conn) {
	f := GetMsgHandler(msgId)
	if f == nil {
		return
	}

	f(msgId, msg, conn)
}

func DoRecv(conn net.Conn) (uint16, []byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return 0, nil, err
	}

	msgLen, header1 := protocol.Decode_uint16(header)
	msgId, _ := protocol.Decode_uint16(header1)
	var content []byte
	if !(msgLen >= 0 && msgLen < 65535) {
		return 0, nil, errors.New("len error!")
	}

	content = make([]byte, msgLen)
	_, err = io.ReadFull(conn, content)
	if err != nil {
		return 0, nil, err
	}

	return msgId, content, nil
}

func Marshal(msgId uint16, itf Marshaller) []byte {
	buff := itf.Marshal()
	csz := len(buff) //content len
	tsz := 4 + csz   //+ header len
	total := make([]byte, tsz)

	//len first
	tmp := protocol.Encode_uint16(uint16(csz))
	copy(total, tmp)

	//msgId
	tmp = protocol.Encode_uint16(msgId)
	copy(total[2:], tmp)
	fmt.Println(msgId, csz, tmp, total[:4])

	//content
	copy(total[4:], buff)
	return total
}
