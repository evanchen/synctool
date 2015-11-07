//This file is automatically created by protocol generator.
//Any manual changes are not suggested.

package protocol

//===========================protocol FInfo===========================
type FInfo struct {
	Path    string
	Modtime uint64
}

func CreateFInfo() *FInfo {
	obj := &FInfo{}
	return obj
}

func (this *FInfo) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_string(this.Path)...)
	buf = append(buf, Encode_uint64(this.Modtime)...)
	return buf
}

func (this *FInfo) Unmarshal(Data []byte) []byte {
	this.Path, Data = Decode_string(Data)
	this.Modtime, Data = Decode_uint64(Data)
	return Data
}

//===========================protocol FInfoList===========================
type FInfoList struct {
	FinfoList []FInfo
}

func CreateFInfoList() *FInfoList {
	obj := &FInfoList{}
	return obj
}

func (this *FInfoList) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, Encode_array_FInfo(this.FinfoList)...)
	return buf
}

func (this *FInfoList) Unmarshal(Data []byte) []byte {
	this.FinfoList, Data = Decode_array_FInfo(Data)
	return Data
}

func Encode_array_FInfo(FinfoList []FInfo) []byte {
	buf := make([]byte, 0, 16)
	size := uint16(len(FinfoList))
	buf = append(buf, Encode_uint16(size)...)
	for _, obj := range FinfoList {
		buf = append(buf, obj.Marshal()...)
	}
	return buf
}

func Decode_array_FInfo(Data []byte) ([]FInfo, []byte) {
	var size uint16
	size, Data = Decode_uint16(Data)
	FinfoList := make([]FInfo, 0, size)
	var obj *FInfo
	for i := uint16(0); i < size; i++ {
		obj = &FInfo{}
		Data = obj.Unmarshal(Data)
		FinfoList = append(FinfoList, *obj)
	}
	return FinfoList, Data
}
