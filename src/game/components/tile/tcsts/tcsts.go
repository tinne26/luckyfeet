package tcsts // tileconstants

const (
	LayerBack = iota
	LayerBackDecor
	LayerMain
	LayerMainDecor
	LayerFront
	LayerFrontDecor
	LayerSpecial
	LayerCountSentinel
)

const (
	MainGround = iota + 1
	MainGroundRaiser
	MainGroundSide
	MainGroundCorner
	MainGroundMark
	MainGroundMarkCorner
	MainSinglePlatform
	MainGrassSide
	MainGrassSideFull
	MainGrassCorner
	MainGrassCornerFull

	BackGround
	BackGroundSide
	BackGroundCorner
	BackGroundMark
	BackGroundMarkCorner

	FrontGround
	FrontGroundRaiser
	FrontGroundSide
	FrontGroundCorner
	FrontGroundMark
	FrontGroundMarkCorner
	FrontSinglePlatform
	FrontGrassSide
	FrontGrassSideFull
	FrontGrassCorner
	FrontGrassCornerFull
	
	StartPoint // not placeable directly as a tile, displayed separately
	RaceGoal

	CarrotOrange
	CarrotYellow
	CarrotPurple
	CarrotMissing

	MainOrangePlatSingle
	MainOrangePlatSingleFill
	MainOrangePlatLeft
	MainOrangePlatLeftFill
	MainOrangePlatRight
	MainOrangePlatRightFill
	MainYellowPlatSingle
	MainYellowPlatSingleFill
	MainYellowPlatLeft
	MainYellowPlatLeftFill
	MainYellowPlatRight
	MainYellowPlatRightFill
	MainPurplePlatSingle
	MainPurplePlatSingleFill
	MainPurplePlatLeft
	MainPurplePlatLeftFill
	MainPurplePlatRight
	MainPurplePlatRightFill

	TransferUp
	TransferUpA
	TransferUpB
	TransferUpC
	TransferLeft
	TransferLeftA
	TransferLeftB
	TransferLeftC
	TransferRight
	TransferRightA
	TransferRightB
	TransferRightC
	TransferDown
	TransferDownA
	TransferDownB
	TransferDownC

	TileTypeMax
	TileNone // out of range, for hacky purposes
)

const (
	GeometryNone = iota + 1
	Geometry20x20
	GeometryBL20x19
	GeometryTR19x19
	GeometryBL20x9
	GeometryBR19x9
	GeometryMT18x17 // single platforms
	GeometryBR19x20 // hackily used for raisers too
	GeometryBR18x16 // carrot left plats
	GeometryBL17x16 // carrot right plats
	GeometryBL1_17x16 // carrot single plats
	GeometryMM4x4 // special target for carrots, goals and transfers

	GeometryMaxSentinel
)

var GeometryTable [TileTypeMax]uint8
func init() {
	// assign no geometry by default
	for i, _ := range GeometryTable {
		GeometryTable[i] = GeometryNone
	}

	GeometryTable[MainGround] = Geometry20x20
	GeometryTable[MainGroundSide] = GeometryBL20x19
	GeometryTable[MainGroundCorner] = GeometryTR19x19
	GeometryTable[MainSinglePlatform] = GeometryMT18x17
	GeometryTable[MainGrassSide] = GeometryBL20x9
	GeometryTable[MainGrassSideFull] = GeometryBL20x9
	GeometryTable[MainGrassCorner] = GeometryBR19x9
	GeometryTable[MainGrassCornerFull] = GeometryBR19x9

	GeometryTable[FrontGround] = Geometry20x20
	GeometryTable[FrontGroundSide] = GeometryBL20x19
	GeometryTable[FrontGroundCorner] = GeometryTR19x19
	GeometryTable[FrontSinglePlatform] = GeometryMT18x17
	GeometryTable[FrontGrassSide] = GeometryBL20x9
	GeometryTable[FrontGrassSideFull] = GeometryBL20x9
	GeometryTable[FrontGrassCorner] = GeometryBR19x9
	GeometryTable[FrontGrassCornerFull] = GeometryBR19x9

	GeometryTable[BackGround] = Geometry20x20
	GeometryTable[BackGroundSide] = GeometryBL20x19
	GeometryTable[BackGroundCorner] = GeometryTR19x19

	GeometryTable[RaceGoal] = GeometryMM4x4
	GeometryTable[CarrotOrange] = GeometryMM4x4
	GeometryTable[CarrotYellow] = GeometryMM4x4
	GeometryTable[CarrotPurple] = GeometryMM4x4
	for i := TransferUp; i < TileTypeMax; i++ {
		GeometryTable[i] = GeometryMM4x4 // all transfers
	}

	GeometryTable[MainOrangePlatSingle] = GeometryBL1_17x16
	GeometryTable[MainYellowPlatSingle] = GeometryBL1_17x16
	GeometryTable[MainPurplePlatSingle] = GeometryBL1_17x16
	GeometryTable[MainOrangePlatLeft] = GeometryBR18x16
	GeometryTable[MainYellowPlatLeft] = GeometryBR18x16
	GeometryTable[MainPurplePlatLeft] = GeometryBR18x16
	GeometryTable[MainOrangePlatRight] = GeometryBL17x16
	GeometryTable[MainYellowPlatRight] = GeometryBL17x16
	GeometryTable[MainPurplePlatRight] = GeometryBL17x16
}
