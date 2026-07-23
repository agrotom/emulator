package auth

import (
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
)

func Factory(rtype common.EgtsService) (common.Binary, error) {
	switch rtype {
	case common.EgtsSrTermIdentity:
		return &AuthService{}, nil
	default:
		return nil, fmt.Errorf("unknown type of auth service: %d", rtype)
	}
}
