package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string    `json:"nama" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Phone       string    `json:"no_telp" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	DateOfBirth *time.Time `json:"tanggal_lahir"`
	Gender      string    `json:"jenis_kelamin"`
	About       string    `json:"tentang"`
	Job         string    `json:"pekerjaan"`
	ProvinceID  string    `json:"id_provinsi"`
	CityID      string    `json:"id_kota"`
	IsAdmin     bool      `json:"is_admin" gorm:"default:false"`
	
	Store       *Store        `json:"toko,omitempty"`
	Addresses   []Address     `json:"alamat,omitempty"`
	Transactions []Transaction `json:"trx,omitempty"`
}

type Store struct {
	gorm.Model
	UserID    uint   `json:"id_user" gorm:"not null;unique"`
	Name      string `json:"nama_toko" gorm:"not null"`
	PhotoURL  string `json:"url_foto"`
	
	User     User      `json:"user,omitempty"`
	Products []Product `json:"produk,omitempty"`
}

type Address struct {
	gorm.Model
	UserID        uint   `json:"id_user" gorm:"not null"`
	Title         string `json:"judul_alamat" gorm:"not null"`
	RecipientName string `json:"nama_penerima" gorm:"not null"`
	Phone         string `json:"no_telp" gorm:"not null"`
	Detail        string `json:"detail_alamat" gorm:"not null"`
}

type Category struct {
	gorm.Model
	Name     string    `json:"nama_category" gorm:"not null"`
	Products []Product `json:"produk,omitempty"`
}

type Product struct {
	gorm.Model
	StoreID       uint   `json:"id_toko" gorm:"not null"`
	CategoryID    uint   `json:"id_category" gorm:"not null"`
	Name          string `json:"nama_produk" gorm:"not null"`
	Slug          string `json:"slug" gorm:"unique;not null"`
	ResellerPrice string `json:"harga_reseller" gorm:"not null"`
	ConsumerPrice string `json:"harga_konsumen" gorm:"not null"`
	Stock         int    `json:"stok" gorm:"not null"`
	Description   string `json:"deskripsi" gorm:"type:text"`
	
	Store      Store       `json:"toko,omitempty"`
	Category   Category    `json:"category,omitempty"`
	Photos     []ProductPhoto `json:"foto_produk,omitempty"`
}

type ProductPhoto struct {
	gorm.Model
	ProductID uint   `json:"id_produk" gorm:"not null"`
	URL       string `json:"url" gorm:"not null"`
	
	Product Product `json:"produk,omitempty"`
}

type Transaction struct {
	gorm.Model
	UserID      uint   `json:"id_user" gorm:"not null"`
	AddressID   uint   `json:"alamat_kirim" gorm:"not null"`
	TotalPrice  int    `json:"harga_total" gorm:"not null"`
	Invoice     string `json:"kode_invoice" gorm:"unique;not null"`
	PaymentMethod string `json:"metode_bayar" gorm:"not null"`
	
	User    User               `json:"user,omitempty"`
	Details []TransactionDetail `json:"detail_trx,omitempty"`
	Logs    []ProductLog       `json:"log_produk,omitempty"`
}

type TransactionDetail struct {
	gorm.Model
	TransactionID uint `json:"id_trx" gorm:"not null"`
	ProductID     uint `json:"id_produk" gorm:"not null"`
	StoreID       uint `json:"id_toko" gorm:"not null"`
	Quantity      int  `json:"kuantitas" gorm:"not null"`
	TotalPrice    int  `json:"harga_total" gorm:"not null"`
	
	Transaction Transaction `json:"trx,omitempty"`
	Product     Product     `json:"produk,omitempty"`
	Store       Store       `json:"toko,omitempty"`
}

type ProductLog struct {
	gorm.Model
	TransactionID uint   `json:"id_trx" gorm:"not null"`
	ProductName   string `json:"nama_produk" gorm:"not null"`
	ProductPrice  int    `json:"harga_produk" gorm:"not null"`
	Quantity      int    `json:"kuantitas" gorm:"not null"`
	StockBefore   int    `json:"stok_sebelum" gorm:"not null"`
	StockAfter    int    `json:"stok_sesudah" gorm:"not null"`
	
	Transaction Transaction `json:"trx,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"nama" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"no_telp" validate:"required"`
	Password string `json:"kata_sandi" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"kata_sandi" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"auth"`
	User  *User  `json:"user"`
} 