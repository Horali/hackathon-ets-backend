package main

// menahanSlot: persona C, hanya status diajukan dan disetujui yang memakai slot.
// ditolak dan selesai melepaskan slot.
func menahanSlot(status string) bool {
	return status == "diajukan" || status == "disetujui"
}

// terpakai menghitung jumlah unit yang sudah dipakai pada satu slot, di satu sumber daya dan tanggal.
// kecualiID dilewati (dipakai saat PUT supaya tidak bentrok dengan dirinya sendiri; 0 = tidak ada).
// Pemanggil harus memegang s.mu.
func (s *Store) terpakai(sdID int, tanggal string, slot int, kecualiID int) int {
	total := 0
	for _, p := range s.peminjaman {
		if p.ID == kecualiID || !menahanSlot(p.Status) {
			continue
		}
		if p.SumberDayaID == sdID && p.Tanggal == tanggal && p.SlotMulai <= slot && slot <= p.SlotSelesai {
			total += p.Jumlah
		}
	}
	return total
}

// adaBentrok: true kalau di salah satu slot dalam rentang, unit terpakai + jumlah baru melebihi stok.
// Pemanggil harus memegang s.mu.
func (s *Store) adaBentrok(in PeminjamanInput, stok int, kecualiID int) bool {
	for slot := in.SlotMulai; slot <= in.SlotSelesai; slot++ {
		if s.terpakai(in.SumberDayaID, in.Tanggal, slot, kecualiID)+in.Jumlah > stok {
			return true
		}
	}
	return false
}
