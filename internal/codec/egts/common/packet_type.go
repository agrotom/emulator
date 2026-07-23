package common

type PacketType byte

const (
	EGTSPTResponse PacketType = iota
	EGTSPTAppData
	EGTSPTSignedAppData
)