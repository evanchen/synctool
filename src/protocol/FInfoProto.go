//This file is automatically created by protocol generator.
//Any manual changes are not suggested.

package protocol

//===========================protocol ClientFInfo===========================
type ClientFInfo struct {
	TarPath      string
	SrcPath      string
	ModTime      uint64
	TotalFileNum uint32
}

func CreateClientFInfo() *ClientFInfo {
	obj := &ClientFInfo{}
	return obj
}

func (this *ClientFInfo) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_string(this.TarPath)...)
	buf = append(buf, Encode_string(this.SrcPath)...)
	buf = append(buf, Encode_uint64(this.ModTime)...)
	buf = append(buf, Encode_uint32(this.TotalFileNum)...)
	return buf
}

func (this *ClientFInfo) Unmarshal(Data []byte) []byte {
	this.TarPath, Data = Decode_string(Data)
	this.SrcPath, Data = Decode_string(Data)
	this.ModTime, Data = Decode_uint64(Data)
	this.TotalFileNum, Data = Decode_uint32(Data)
	return Data
}

//===========================protocol ClientFInfoList===========================
type ClientFInfoList struct {
	Flist []ClientFInfo
}

func CreateClientFInfoList() *ClientFInfoList {
	obj := &ClientFInfoList{}
	return obj
}

func (this *ClientFInfoList) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_array_ClientFInfo(this.Flist)...)
	return buf
}

func (this *ClientFInfoList) Unmarshal(Data []byte) []byte {
	this.Flist, Data = Decode_array_ClientFInfo(Data)
	return Data
}

func Encode_array_ClientFInfo(Flist []ClientFInfo) []byte {
	buf := make([]byte, 0, 16)
	size := uint16(len(Flist))
	buf = append(buf, Encode_uint16(size)...)
	for _, obj := range Flist {
		buf = append(buf, obj.Marshal()...)
	}
	return buf
}

func Decode_array_ClientFInfo(Data []byte) ([]ClientFInfo, []byte) {
	var size uint16
	size, Data = Decode_uint16(Data)
	Flist := make([]ClientFInfo, 0, size)
	var obj *ClientFInfo
	for i := uint16(0); i < size; i++ {
		obj = &ClientFInfo{}
		Data = obj.Unmarshal(Data)
		Flist = append(Flist, *obj)
	}
	return Flist, Data
}

//===========================protocol ClientDataBuff===========================
type ClientDataBuff struct {
	MsgId uint16
	Buff  []uint8
	Bpos  uint32
	Total uint32
	Path  string
}

func CreateClientDataBuff() *ClientDataBuff {
	obj := &ClientDataBuff{}
	return obj
}

func (this *ClientDataBuff) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_uint16(this.MsgId)...)
	buf = append(buf, Encode_array_uint8(this.Buff)...)
	buf = append(buf, Encode_uint32(this.Bpos)...)
	buf = append(buf, Encode_uint32(this.Total)...)
	buf = append(buf, Encode_string(this.Path)...)
	return buf
}

func (this *ClientDataBuff) Unmarshal(Data []byte) []byte {
	this.MsgId, Data = Decode_uint16(Data)
	this.Buff, Data = Decode_array_uint8(Data)
	this.Bpos, Data = Decode_uint32(Data)
	this.Total, Data = Decode_uint32(Data)
	this.Path, Data = Decode_string(Data)
	return Data
}

//===========================protocol ServNeedClientData===========================
type ServNeedClientData struct {
	Bpos uint32
	Path string
}

func CreateServNeedClientData() *ServNeedClientData {
	obj := &ServNeedClientData{}
	return obj
}

func (this *ServNeedClientData) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_uint32(this.Bpos)...)
	buf = append(buf, Encode_string(this.Path)...)
	return buf
}

func (this *ServNeedClientData) Unmarshal(Data []byte) []byte {
	this.Bpos, Data = Decode_uint32(Data)
	this.Path, Data = Decode_string(Data)
	return Data
}

//===========================protocol FileMD5Info===========================
type FileMD5Info struct {
	Path string
	Sum  string
}

func CreateFileMD5Info() *FileMD5Info {
	obj := &FileMD5Info{}
	return obj
}

func (this *FileMD5Info) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_string(this.Path)...)
	buf = append(buf, Encode_string(this.Sum)...)
	return buf
}

func (this *FileMD5Info) Unmarshal(Data []byte) []byte {
	this.Path, Data = Decode_string(Data)
	this.Sum, Data = Decode_string(Data)
	return Data
}
