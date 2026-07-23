package packet

import (
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
)

func Factory(ptype common.PacketType) (common.Binary, error) {
	switch ptype {
	case common.EGTSPTResponse:
		return &PtResponse{}, nil
	case common.EGTSPTAppData:
		return &PtAppData{}, nil
	case common.EGTSPTSignedAppData:
		return &PtSignedAppData{}, nil
	default:
		return nil, fmt.Errorf("error while decoding frame data: unknown packet type %d", ptype)
	}
}
