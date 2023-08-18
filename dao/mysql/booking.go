package mysql

import (
	// "bookingBackEnd/model"
	"bookingBackEnd/utils"
	"sync"
)

var onceForBooking sync.Once

type BookingTable struct {
	tableName string
}

var BookingTableInstance *BookingTable

func NewBookingTable() (*BookingTable, error) {
	onceForBooking.Do(func() {
		BookingTableInstance = &BookingTable{
			tableName: utils.ParamsInstance.BookingTableName,
		}
	})
	return BookingTableInstance, nil
}
