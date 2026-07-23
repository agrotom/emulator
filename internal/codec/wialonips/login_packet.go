package wialonips

import (
	"fmt"
)

type LoginPacket struct {
	IMEI     string
	Password string
}

func (lp *LoginPacket) Encode() ([]byte, error) {
	return fmt.Appendf(nil, "#L#%s;%s\r\n", lp.IMEI, lp.Password), nil
}

func (dp *LoginPacket) ValidateResponse(resp []byte) error {
	if string(resp) != SuccessAL {
		return fmt.Errorf("data packet response error: %s", string(resp))
	}

	return nil
}
