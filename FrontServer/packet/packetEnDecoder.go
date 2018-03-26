// 출처: https://github.com/gonet2/agent

package packet

import (
	"encoding/binary"
	"errors"
)

const (
	packetLimit = 65535
)

type rawPacketDataSt struct {
	pos   int
	data  []byte
	order binary.ByteOrder
}

func newReader(data []byte, isLittleEndian bool) *rawPacketDataSt {
	if isLittleEndian {
		return &rawPacketDataSt{data: data, order: binary.LittleEndian}
	}
	return &rawPacketDataSt{data: data, order: binary.BigEndian}
}

func newWriter(buffer []byte, isLittleEndian bool) *rawPacketDataSt {
	if isLittleEndian {
		return &rawPacketDataSt{data: buffer, order: binary.LittleEndian}
	}
	return &rawPacketDataSt{data: buffer, order: binary.BigEndian}
}

func (p *rawPacketDataSt) Data() []byte {
	return p.data
}

func (p *rawPacketDataSt) Length() int {
	return len(p.data)
}

//=============================================== Readers
func (p *rawPacketDataSt) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

func (p *rawPacketDataSt) ReadS8() (ret int8, err error) {
	_ret, _err := p.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

func (p *rawPacketDataSt) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = p.order.Uint16(buf)
	p.pos += 2
	return
}

func (p *rawPacketDataSt) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *rawPacketDataSt) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = p.order.Uint32(buf)
	p.pos += 4
	return
}

func (p *rawPacketDataSt) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *rawPacketDataSt) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read uint64 failed")
		return
	}

	buf := p.data[p.pos : p.pos+8]
	ret = p.order.Uint64(buf)
	p.pos += 8
	return
}

func (p *rawPacketDataSt) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *rawPacketDataSt) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *rawPacketDataSt) ReadBytes() (ret []byte, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	return
}

func (p *rawPacketDataSt) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

/*
func (p *rawPacketDataSt) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *rawPacketDataSt) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}
*/

//================================================ Writers
func (p *rawPacketDataSt) WriteU16(v uint16) {
	p.order.PutUint16(p.data[p.pos:], v)
	p.pos += 2
}

func (p *rawPacketDataSt) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

func (p *rawPacketDataSt) WriteBytes(v []byte) {
	copy(p.data[p.pos:], v)
}

func (p *rawPacketDataSt) WriteU32(v uint32) {
	p.order.PutUint32(p.data[p.pos:], v)
	p.pos += 4
}

func (p *rawPacketDataSt) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

func (p *rawPacketDataSt) WriteU64(v uint64) {
	p.order.PutUint64(p.data[p.pos:], v)
	p.pos += 4
}

func (p *rawPacketDataSt) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

/*
func (p *rawPacketDataSt) WriteZeros(n int) {
	for i := 0; i < n; i++ {
		p.data = append(p.data, byte(0))
	}
}

func (p *rawPacketDataSt) WriteBool(v bool) {
	if v {
		p.data = append(p.data, byte(1))
	} else {
		p.data = append(p.data, byte(0))
	}
}

func (p *rawPacketDataSt) WriteByte(v byte) {
	p.data = append(p.data, v)
}

func (p *rawPacketDataSt) WriteBytes(v []byte) {
	p.WriteU16(uint16(len(v)))
	p.data = append(p.data, v...)
}

func (p *rawPacketDataSt) WriteRawBytes(v []byte) {
	p.data = append(p.data, v...)
}

func (p *rawPacketDataSt) WriteString(v string) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.data = append(p.data, bytes...)
}

func (p *rawPacketDataSt) WriteS8(v int8) {
	p.WriteByte(byte(v))
}



func (p *rawPacketDataSt) WriteU24(v uint32) {
	p.data = append(p.data, byte(v>>16), byte(v>>8), byte(v))
}

func (p *rawPacketDataSt) WriteFloat32(f float32) {
	v := math.Float32bits(f)
	p.WriteU32(v)
}

func (p *rawPacketDataSt) WriteFloat64(f float64) {
	v := math.Float64bits(f)
	p.WriteU64(v)
}
*/
