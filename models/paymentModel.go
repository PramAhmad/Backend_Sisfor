package models

// import room and mahasiswa
type Payment struct {
	ID          int64 `gorm:"primary_key"`
	Total       int64
	MahasiswaID int64
	RoomID      int64
	Addby       int64
	Is_delete   int8 `gorm:"default:0" json:"is_delete"`
	Room        Room
	Mahasiswa   Mahasiswa
}
