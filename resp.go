package jt808

import (
	"bytes"
	"encoding/binary"
)

type MsgRespHead struct {
	MsgID    uint16
	Property uint16
	PhoneNo  []byte
	FlowID   uint16
}

func NewMsgRespHead(msgID, property, flowID uint16, phoneNo []byte) *MsgRespHead {
	return &MsgRespHead{
		MsgID:    msgID,
		Property: property,
		PhoneNo:  phoneNo,
		FlowID:   flowID,
	}
}

func (msg *MsgRespHead) MarshalBinary() ([]byte, error) {
	bs := make([]byte, 12, 12)
	binary.BigEndian.PutUint16(bs[:2], msg.MsgID)
	binary.BigEndian.PutUint16(bs[2:4], msg.Property)
	bs[4] = msg.PhoneNo[0]
	bs[5] = msg.PhoneNo[1]
	bs[6] = msg.PhoneNo[2]
	bs[7] = msg.PhoneNo[3]
	bs[8] = msg.PhoneNo[4]
	bs[9] = msg.PhoneNo[5]
	binary.BigEndian.PutUint16(bs[10:12], msg.FlowID)
	return bs, nil
}

type CommonRespMsg struct {
	FlowID     uint16
	ResponseID uint16
	Result     byte
}

func NewCommonRespMsg(f, r uint16, result byte) *CommonRespMsg {
	return &CommonRespMsg{
		FlowID:     f,
		ResponseID: r,
		Result:     result,
	}
}

func (msg *CommonRespMsg) MarshalBinary() ([]byte, error) {
	bs := make([]byte, 5, 5)
	binary.BigEndian.PutUint16(bs[:2], msg.FlowID)
	binary.BigEndian.PutUint16(bs[2:4], msg.ResponseID)
	bs[4] = msg.Result
	return bs, nil
}

func (msg *CommonRespMsg) Len() uint16 {
	return 5
}

type LoginRespMsg struct {
	FlowID   uint16
	Result   byte
	AuthCode string
}

func NewLoginRespMsg(fid uint16, result byte, code string) *LoginRespMsg {
	return &LoginRespMsg{FlowID: fid, Result: result, AuthCode: code}
}

func (msg *LoginRespMsg) MarshalBinary() ([]byte, error) {
	bs := make([]byte, 3, 3)
	binary.BigEndian.PutUint16(bs[:2], msg.FlowID)
	bs[2] = msg.Result
	if msg.Result != 0 {
		return bs, nil
	}

	var buf bytes.Buffer
	buf.Write(bs)
	buf.WriteString(msg.AuthCode)
	return buf.Bytes(), nil
}

func (msg *LoginRespMsg) Len() uint16 {
	return uint16(len(msg.AuthCode) + 3)
}

var EmptyRespMsg *emptyRespMsg = new(emptyRespMsg)

type emptyRespMsg struct{}

func (msg *emptyRespMsg) MarshalBinary() ([]byte, error) {
	return []byte{}, nil
}

func (msg *emptyRespMsg) Len() uint16 {
	return 0
}

const (
	// 平台通用应答
	MsgCmdCommonResp uint16 = 0x8001
	// 终端注册应答
	MsgCmdTerminalRegisterResp uint16 = 0x8100
	// 设置终端参数
	MsgCmdTerminalParamSettings uint16 = 0X8103
	// 查询终端参数
	MsgCmdTerminalParamQuery uint16 = 0x8104
)
