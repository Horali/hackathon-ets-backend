package main

import (
	"strings"
	"time"
	"unicode/utf8"
)

// PeminjamanInput: enam field body POST/PUT /peminjaman
type PeminjamanInput struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
}

// validasiPeminjaman men-trim teks lalu memeriksa field. Mengembalikan pesan error, atau "" kalau valid.
// Cek sumber_daya_id ada dan jumlah <= stok dilakukan di handler karena butuh data store.
func validasiPeminjaman(in *PeminjamanInput) string {
	in.Peminjam = strings.TrimSpace(in.Peminjam)
	in.Tanggal = strings.TrimSpace(in.Tanggal)

	if in.Peminjam == "" || utf8.RuneCountInString(in.Peminjam) > 50 {
		return "peminjam wajib, maksimal 50 karakter"
	}
	if !tanggalValid(in.Tanggal) {
		return "tanggal harus format YYYY-MM-DD dan tanggalnya harus ada"
	}
	if in.SlotMulai < 1 || in.SlotMulai > 8 || in.SlotSelesai < 1 || in.SlotSelesai > 8 {
		return "slot harus antara 1 dan 8"
	}
	if in.SlotMulai > in.SlotSelesai {
		return "slot_mulai tidak boleh lebih besar dari slot_selesai"
	}
	if in.Jumlah < 1 {
		return "jumlah minimal 1"
	}
	return ""
}

// tanggalValid: time.Parse otomatis menolak tanggal yang tidak ada, misalnya 2026-02-30
func tanggalValid(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}
