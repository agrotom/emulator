package common

type EgtsServiceType byte

const (
	EgtsAuth     EgtsServiceType = 0
	EgtsTeledata EgtsServiceType = 1
)

type EgtsService byte

const (
	EgtsSrRecordResponse     EgtsService = 0
	EgtsSrTermIdentity       EgtsService = 1
	EgtsSrPosData            EgtsService = 16
	EgtsSrExtPosData         EgtsService = 17
	EgtsSrAdSensorsData      EgtsService = 18
	EgtsSrCountersData       EgtsService = 19
	EgtsSrStateData          EgtsService = 20
	EgtsSrLoopinData         EgtsService = 22
	EgtsSrAbsDigSensData     EgtsService = 23
	EgtsSrAbsAnsensData      EgtsService = 24
	EgtsSrAbsCntrData        EgtsService = 25
	EgtsSrAbsLoopinData      EgtsService = 26
	EgtsSrLiquidLevelSensor  EgtsService = 27
	EgtsSrPassangersCounters EgtsService = 28
)
