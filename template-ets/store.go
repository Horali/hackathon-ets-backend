package main

import "sync"

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
	Alasan       string `json:"alasan"` // persona C: "" sampai ditolak
}

// Store: semua data di memori, dijaga satu mutex
type Store struct {
	mu         sync.Mutex
	sumberDaya []SumberDaya
	peminjaman []Peminjaman
	nextSDID   int
	nextPJID   int
}
