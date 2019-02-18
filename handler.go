package jt808

import (
	"errors"
	"io"
)

type Handler interface {
	Raw(raw []byte)
	TerminalCommonResp(head *MsgHead, body *TerminalCommonRespBody) error
	Heartbeat(w io.Writer, head *MsgHead) error
	Login(w io.Writer, head *MsgHead, body *LoginMsg) error
	Logout(w io.Writer, head *MsgHead) error
	Auth(w io.Writer, head *MsgHead, body *AuthMsg) error
	UploadLocation(w io.Writer, head *MsgHead, body *LocMsg) error

	Err(w io.Writer, err error) error

	OtherMsg(w io.Writer, head *MsgHead, bd *BufDecoder) error
}

type (
	DefaultHandler struct{}
)

func (h *DefaultHandler) Raw(raw []byte) {
}

func (h *DefaultHandler) TerminalCommonResp(head *MsgHead, body *TerminalCommonRespBody) error {
	return nil
}

func (h *DefaultHandler) Heartbeat(w io.Writer, head *MsgHead) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func (h *DefaultHandler) Login(w io.Writer, head *MsgHead, body *LoginMsg) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func (h *DefaultHandler) Logout(w io.Writer, head *MsgHead) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func (h *DefaultHandler) Auth(w io.Writer, head *MsgHead, body *AuthMsg) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func (handler *DefaultHandler) UploadLocation(w io.Writer, head *MsgHead, body *LocMsg) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func (h *DefaultHandler) Err(w io.Writer, err error) error {
	if err == ErrContinue {
		return nil
	}
	return err
}

func (h *DefaultHandler) OtherMsg(w io.Writer, head *MsgHead, bd *BufDecoder) error {
	respHead, respBody := NewCommonMsgResp(head.FlowID, head.MsgID, ResultSuccess, head.PhoneNoBCD)
	return JT808SendMsg(w, respHead, respBody)
}

func NewDefaultHandler() *DefaultHandler {
	return new(DefaultHandler)
}

var (
	ErrContinue error = errors.New("err continue")
)

func Processor(w io.Writer, raw []byte, h Handler) error {
	var data []byte
	var err error

	if data, err = Restore(raw); err != nil {
		return nil
	}

	h.Raw(raw)

	bd := NewBufDecoder(data)

	var head *MsgHead
	if head, err = DecodeMsgHead(bd); err != nil {
		return nil
	}

	switch head.MsgID {
	case MsgIdTerminalLocationInfoUpload:
		err = doUploadLocation(w, h, head, bd)
	case MsgIdTerminalHeartBeat:
		err = h.Err(w, h.Heartbeat(w, head))
	case MsgIdTerminalAuthentication:
		err = doAuth(w, h, head, bd)
	case MsgIdTerminalLogin:
	case MsgIdTerminalLogOut:
	case MsgIdTerminalCommonResp:
	default:
		err = h.Err(w, h.OtherMsg(w, head, bd))
	}

	return err
}

func doUploadLocation(w io.Writer, h Handler, head *MsgHead, bd *BufDecoder) error {
	var body *LocMsg
	var err error

	if body, err = DecodeLocMsg(head, bd); err != nil {
		return ErrContinue
	}

	return h.Err(w, h.UploadLocation(w, head, body))
}

func doAuth(w io.Writer, h Handler, head *MsgHead, bd *BufDecoder) error {
	var body *AuthMsg
	var err error

	if body, err = DecodeAuthMsg(head, bd); err != nil {
		return ErrContinue
	}

	return h.Err(w, h.Auth(w, head, body))
}

func NewCommonMsgResp(flowID, respID uint16, result byte, phoneNo []byte) (*MsgRespHead, *CommonRespMsg) {
	body := NewCommonRespMsg(flowID, respID, result)
	head := NewMsgRespHead(MsgCmdCommonResp, body.Len(), TimeFlowID(), phoneNo)
	return head, body
}

const (
	ResultSuccess     byte = 0x00
	ResultFailure     byte = 0x01
	ResultError       byte = 0x02
	ResultUnsupported byte = 0x03
	ResultWarnning    byte = 0x04
)
