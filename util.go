package jt808

import "time"

func TimeFlowID() uint16 {
	h, m, s := time.Now().Clock()
	return uint16(h*3600 + m*60 + s)
}

const (
	FlagAcc      = 0x0001
	FlagLoc      = 0x0002
	FlagSouthLat = 0x0004
	FlagWestLng  = 0x0008
)
