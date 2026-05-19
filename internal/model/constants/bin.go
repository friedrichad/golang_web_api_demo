package constants

// Bin Status
const (
	BinStatusAvailable = 0
	BinStatusOccupied  = 1
	BinStatusReserved  = 2
)

func IsValidBinStatus(status int) bool {
	switch status {
	case BinStatusAvailable, BinStatusOccupied, BinStatusReserved:
		return true
	default:
		return false
	}
}
