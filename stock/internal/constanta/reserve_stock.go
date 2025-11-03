package constanta

type (
	ReservedStockStatus string
)

const (
	ReservedStockStatusReserved  ReservedStockStatus = "reserved"
	ReservedStockStatusReleased  ReservedStockStatus = "released"
	ConfirmedStockStatusReleased ReservedStockStatus = "confirmed"
)
