package auth

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/agrotom/emulator/internal/codec/egts/common"
)

const (
	IMEIStringSize  uint16 = 15
	IMSIStringSize  uint16 = 16
	LNGCStringSize  uint16 = 3
	MSISDStringSize uint16 = 15
)

type AuthService struct {
	TerminalIdentifier uint32

	MobileNetworkExists                         bool
	BufferSizeExists                            bool
	NetworkIdentifierExists                     bool
	SSRA                                        bool
	LanguageCodeExists                          bool
	InternationalMobileSubscriberIdentityExists bool
	InternationalMobileEquipmentIdentityExists  bool
	HomeDispatcherIdentifierExists              bool

	HomeDispatcherIdentifier              uint16
	InternationalMobileEquipmentIdentity  string
	InternationalMobileSubscriberIdentity string
	LanguageCode                          string
	NetworkIdentifier                     uint32
	BufferSize                            uint16
	MSISD                                 string
}

func (as *AuthService) Decode(data []byte) error {
	var (
		err   error
		flags byte
	)

	buf := bytes.NewReader(data)

	if value, err := common.ReadUInt32(buf); err != nil {
		return fmt.Errorf("error while decoding terminal identifier of auth service")
	} else {
		as.TerminalIdentifier = value
	}

	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("error while decoding flags of auth service")
	}

	as.MobileNetworkExists = common.Bit(flags, 7)
	as.BufferSizeExists = common.Bit(flags, 6)
	as.NetworkIdentifierExists = common.Bit(flags, 5)
	as.SSRA = common.Bit(flags, 4)
	as.LanguageCodeExists = common.Bit(flags, 3)
	as.InternationalMobileSubscriberIdentityExists = common.Bit(flags, 2)
	as.InternationalMobileEquipmentIdentityExists = common.Bit(flags, 1)
	as.HomeDispatcherIdentifierExists = common.Bit(flags, 0)

	if as.HomeDispatcherIdentifierExists {
		if value, err := common.ReadUInt16(buf); err != nil {
			return fmt.Errorf("error while decoding home dispatcher identifier of auth service")
		} else {
			as.HomeDispatcherIdentifier = value
		}
	}

	if as.InternationalMobileEquipmentIdentityExists {
		if as.InternationalMobileEquipmentIdentity, err = common.ReadStringBinary(buf, IMEIStringSize); err != nil {
			return fmt.Errorf("error while decoding international mobile equipment identity of auth service")
		}
	}

	if as.InternationalMobileSubscriberIdentityExists {
		if as.InternationalMobileSubscriberIdentity, err = common.ReadStringBinary(buf, IMSIStringSize); err != nil {
			return fmt.Errorf("error while decoding international mobile subscriber identity of auth service")
		}
	}

	if as.LanguageCodeExists {
		if as.LanguageCode, err = common.ReadStringBinary(buf, LNGCStringSize); err != nil {
			return fmt.Errorf("error while decoding langugae code of auth service")
		}
	}

	if as.NetworkIdentifierExists {
		if as.NetworkIdentifier, err = common.ReadUInt24(buf); err != nil {
			return fmt.Errorf("error while decoding network identifier of auth service")
		}
	}

	if as.BufferSizeExists {
		if as.BufferSize, err = common.ReadUInt16(buf); err != nil {
			return fmt.Errorf("error while decoding buffer size of auth service")
		}
	}

	if as.MobileNetworkExists {
		if as.MSISD, err = common.ReadStringBinary(buf, MSISDStringSize); err != nil {
			return fmt.Errorf("error while decoding mobile station integrated services digital network number of auth service")
		}
	}

	return nil
}

func (as *AuthService) Encode() ([]byte, error) {
	var (
		err   error
		flags byte
	)

	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, as.TerminalIdentifier); err != nil {
		return nil, fmt.Errorf("error while encoding terminal identifier of auth service")
	}

	common.SetBit(&flags, 7, as.MobileNetworkExists)
	common.SetBit(&flags, 6, as.BufferSizeExists)
	common.SetBit(&flags, 5, as.NetworkIdentifierExists)
	common.SetBit(&flags, 4, as.SSRA)
	common.SetBit(&flags, 3, as.LanguageCodeExists)
	common.SetBit(&flags, 2, as.InternationalMobileSubscriberIdentityExists)
	common.SetBit(&flags, 1, as.InternationalMobileEquipmentIdentityExists)
	common.SetBit(&flags, 0, as.HomeDispatcherIdentifierExists)

	if err = buf.WriteByte(flags); err != nil {
		return nil, fmt.Errorf("error while encoding flags of auth service")
	}

	if as.HomeDispatcherIdentifierExists {
		if err = binary.Write(buf, binary.LittleEndian, as.HomeDispatcherIdentifier); err != nil {
			return nil, fmt.Errorf("error while encoding home dispatcher identifier of auth service")
		}
	}

	if as.InternationalMobileEquipmentIdentityExists {
		if err = common.WriteStringBinary(buf, IMEIStringSize, as.InternationalMobileEquipmentIdentity); err != nil {
			return nil, fmt.Errorf("error while encoding international mobile equipment identity of auth service")
		}
	}

	if as.InternationalMobileSubscriberIdentityExists {
		if err = common.WriteStringBinary(buf, IMSIStringSize, as.InternationalMobileSubscriberIdentity); err != nil {
			return nil, fmt.Errorf("error while encoding international mobile equipment identity of auth service")
		}
	}

	if as.LanguageCodeExists {
		if err = common.WriteStringBinary(buf, LNGCStringSize, as.LanguageCode); err != nil {
			return nil, fmt.Errorf("error while encoding language code of auth service")
		}
	}

	if as.NetworkIdentifierExists {
		if err = common.WriteUInt24(buf, as.NetworkIdentifier); err != nil {
			return nil, fmt.Errorf("error while encoding network identifier of auth service")
		}
	}

	if as.BufferSizeExists {
		if err = binary.Write(buf, binary.LittleEndian, as.BufferSize); err != nil {
			return nil, fmt.Errorf("error while encoding buffer size of auth service")
		}
	}

	if as.MobileNetworkExists {
		if err = common.WriteStringBinary(buf, MSISDStringSize, as.MSISD); err != nil {
			return nil, fmt.Errorf("error while encoding mobile station integrated services digital network number of auth service")
		}
	}

	return buf.Bytes(), nil
}

func (as *AuthService) Length() uint16 {
	if asBytes, err := as.Encode(); err != nil {
		return 0
	} else {
		return uint16(len(asBytes))
	}
}
