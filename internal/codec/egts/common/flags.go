package common

type Prefix byte

const (
	DefaultPrefix Prefix = iota
)

type Encryption byte

const (
	NoEnc Encryption = iota
	Enc1
	Enc2
	Enc3
)

type Priority byte

const (
	Highest Priority = iota
	High
	Medium
	Low
)

/*type PrefixType string
type EncryptionType string
type PriorityType string
type ServiceOnDeviceType string
type GroupType string
type CompressedType string
type OptionalFlag string

type LongitudeType string
type LatitudeType string
type MovingType string
type BlackBoxType string
type FixType string
type CoordSystem string
type ValidnessType string
type AltitudeSignType string

type SSRAType string

const (
	NorthLongitude LongitudeType = "0"
	SouthLongitude LongitudeType = "1"
)

const (
	EastLatitude LatitudeType = "0"
	WestLatitude LatitudeType = "1"
)

const (
	Stopped MovingType = "0"
	Going   MovingType = "1"
)

const (
	ActualData   BlackBoxType = "0"
	BlackBoxData BlackBoxType = "1"
)

const (
	Fix2D FixType = "0"
	Fix3D FixType = "1"
)

const (
	WGS84System CoordSystem = "0"
	PZ90_02     CoordSystem = "1"
)

const (
	Invalid ValidnessType = "0"
	Valid   ValidnessType = "1"
)

const (
	HigherSea AltitudeSignType = "0"
	LowerSea  AltitudeSignType = "1"
)

const (
	RequestAlgorithm SSRAType = "0"
	SimpleAlgorithm  SSRAType = "1"
)
*/
