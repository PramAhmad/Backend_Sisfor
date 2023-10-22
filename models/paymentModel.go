package models

// import room and mahasiswa
type Payment struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	Total       int64     `json:"total"`
	MahasiswaID int64     `json:"mahasiswa_id"`
	RoomID      int64     `json:"room_id"`
	Addby       int64     `json:"addby"`
	Is_delete   int8      `gorm:"default:0" json:"is_delete"`
	Room        Room      `json:"room"`
	Mahasiswa   Mahasiswa `json:"mahasiswa"`
}
