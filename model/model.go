package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Nama     string
	Username string
	Email    string
	Password string
	Role     string
	Login Login `gorm:"foreignKey:SessionId;constraint:OnDelete:CASCADE"`
}

type Order struct {
	gorm.Model
	OrderId int
	Tanggal_masuk  string
	Tanggal_keluar string
	Total_harga    int
	OrderDetail []OrderDetail `gorm:"foreignKey:OrderId;references:ID"`
}

type OrderDetail struct {
	gorm.Model
	OrderDetailId int
	OrderId       int
	Jenis_layanan string
	Berat         int
	Harga_per_kg  int
	Subtotal      int
}

type Customer struct {
	gorm.Model
	CustomerId int
	Nama   string
	Username string
	Password string
	No_hp  string
	Alamat string
	Order []Order `gorm:"foreignKey:OrderId"`
}

type Payment struct {
	gorm.Model
	PaymentId int
	Tanggal_bayar     string
	Metode            string
	Jumlah_bayar      int
	Status_pembayaran string
	OrderId Order `gorm:"foreignKey:OrderId"`
}

type Report struct {
	gorm.Model
	ReportId int
	Periode          string
	Total_pendapatan int
}

type Login struct{
	HashedPassword string
	SessionToken string
	CSRFToken string
	SessionId uint
}