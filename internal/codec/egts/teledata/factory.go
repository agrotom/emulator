package teledata

import (
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
)

func Factory(rtype common.EgtsService) (common.Binary, error) {
	switch rtype {
	case common.EgtsSrPosData:
		return &SrPosData{}, nil
	default:
		return nil, fmt.Errorf("unknown type of teledata service: %d", rtype)
	}
}
