package jt808

import "fmt"

type MsgHead struct {
	MsgID    uint16
	Property uint16
	PhoneNo  string
	FlowID   uint16

	Crypto     bool
	BodyLength uint16
	Mux        bool

	PhoneNoBCD []byte
}

type LocMsg struct {
	Warn   uint32
	Status uint32

	Lat       uint32
	Lng       uint32
	Elevation uint16
	Speed     uint16
	Direction uint16

	UploadTime string
}

type AuthMsg struct {
	AuthCode string
}

type LoginMsg struct {
	ProvinceID uint16
	CityID     uint16

	ManufacturerID string
	TerminalType   string
	TerminalID     string

	VehicleColor uint16
	VehicleNo    string
}

type TerminalCommonRespBody struct {
	FlowID     uint16
	ResponseID uint16
	Result     byte
}

const (
	MsgHeadMuxMask        = 0x2000
	MsgHeadCryptoMask     = 0x0400
	MsgHeadBodyLengthMask = 0x03ff
)

var (
	ErrNotSupport = fmt.Errorf("Not Support")
)

func DecodeMsgHead(bd *BufDecoder) (*MsgHead, error) {

	head := new(MsgHead)

	var err error

	if head.MsgID, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if head.Property, err = bd.Uint16(); err != nil {
		return nil, err
	}

	head.PhoneNoBCD = make([]byte, 6, 6)
	if _, err = bd.Read(head.PhoneNoBCD); err != nil {
		return nil, err
	}
	head.PhoneNo = bcd2str(head.PhoneNoBCD)

	if head.FlowID, err = bd.Uint16(); err != nil {
		return nil, err
	}

	head.BodyLength = head.Property & MsgHeadBodyLengthMask
	head.Crypto = (head.Property & MsgHeadCryptoMask) != 0
	head.Mux = (head.Property & MsgHeadMuxMask) != 0

	if head.Mux {
		return nil, ErrNotSupport
	}

	return head, nil
}

func DecodeLocMsg(head *MsgHead, bd *BufDecoder) (*LocMsg, error) {

	msg := new(LocMsg)
	var err error

	if msg.Warn, err = bd.Uint32(); err != nil {
		return nil, err
	}

	if msg.Status, err = bd.Uint32(); err != nil {
		return nil, err
	}

	if msg.Lat, err = bd.Uint32(); err != nil {
		return nil, err
	}

	if msg.Lng, err = bd.Uint32(); err != nil {
		return nil, err
	}

	if msg.Elevation, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if msg.Speed, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if msg.Direction, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if msg.UploadTime, err = bd.BCD(6); err != nil {
		return nil, err
	}

	return msg, nil
}

func DecodeLoginMsg(head *MsgHead, bd *BufDecoder) (*LoginMsg, error) {
	msg := new(LoginMsg)
	var err error

	if msg.ProvinceID, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if msg.CityID, err = bd.Uint16(); err != nil {
		return nil, err
	}

	if msg.ManufacturerID, err = bd.BCD(5); err != nil {
		return nil, err
	}

	if msg.TerminalType, err = bd.BCD(20); err != nil {
		return nil, err
	}

	if msg.TerminalID, err = bd.BCD(7); err != nil {
		return nil, err
	}

	return msg, nil
}

func DecodeAuthMsg(head *MsgHead, bd *BufDecoder) (*AuthMsg, error) {
	return &AuthMsg{AuthCode: "7777777"}, nil
}

func DecodeTerminalCommonRespBody(head *MsgHead, bd *BufDecoder) (*TerminalCommonRespBody, error) {
	msg := new(TerminalCommonRespBody)
	var err error
	if msg.FlowID, err = bd.Uint16(); err != nil {
		return nil, err
	}
	if msg.ResponseID, err = bd.Uint16(); err != nil {
		return nil, err
	}
	if msg.Result, err = bd.ReadByte(); err != nil {
		return nil, err
	}
	return nil, nil
}

const (
	// 终端通用应答
	MsgIdTerminalCommonResp uint16 = 0x0001
	// 终端心跳
	MsgIdTerminalHeartBeat uint16 = 0x0002
	// 终端注册
	MsgIdTerminalLogin uint16 = 0x0100
	// 终端注销
	MsgIdTerminalLogOut uint16 = 0x0003
	// 终端鉴权
	MsgIdTerminalAuthentication uint16 = 0x0102
	// 位置信息汇报
	MsgIdTerminalLocationInfoUpload uint16 = 0x0200
	// 胎压数据透传
	MsgIdTerminalTransmissionTyrePressure uint16 = 0x0600
	// 查询终端参数应答
	MsgIdTerminalParamQueryResp uint16 = 0x0104
)
