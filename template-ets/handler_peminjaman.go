package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// cariSumberDaya: panggil hanya saat sudah memegang s.mu
func (s *Store) cariSumberDaya(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

// cariIndexPeminjaman: panggil hanya saat sudah memegang s.mu. -1 kalau tidak ada.
func (s *Store) cariIndexPeminjaman(id int) int {
	for i, p := range s.peminjaman {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// B2: POST /peminjaman
func (s *Store) buatPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in PeminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "body JSON tidak valid")
		return
	}
	if pesan := validasiPeminjaman(&in); pesan != "" {
		writeError(w, 400, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sd, ada := s.cariSumberDaya(in.SumberDayaID)
	if !ada {
		writeError(w, 400, "sumber_daya_id tidak ditemukan") // 400, bukan 404: ini isi body
		return
	}
	if in.Jumlah > sd.Stok {
		writeError(w, 400, "jumlah melebihi stok sumber daya")
		return
	}
	// B5: cek bentrok dan simpan ada di dalam lock yang sama
	if s.adaBentrok(in, sd.Stok, 0) {
		writeError(w, 409, "slot sudah penuh, kapasitas terlampaui")
		return
	}

	p := Peminjaman{
		ID:           s.nextPJID,
		SumberDayaID: in.SumberDayaID,
		Peminjam:     in.Peminjam,
		Tanggal:      in.Tanggal,
		SlotMulai:    in.SlotMulai,
		SlotSelesai:  in.SlotSelesai,
		Jumlah:       in.Jumlah,
		Status:       "diajukan", // persona C: harus disetujui admin dulu
		Alasan:       "",
	}
	s.nextPJID++
	s.peminjaman = append(s.peminjaman, p)
	writeJSON(w, 201, p)
}

// B3: GET /peminjaman dengan filter opsional sumber_daya_id, tanggal, status
func (s *Store) daftarPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filterSD := 0
	if v := q.Get("sumber_daya_id"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, 400, "sumber_daya_id harus angka")
			return
		}
		filterSD = n
	}
	filterTanggal := q.Get("tanggal")
	filterStatus := q.Get("status")

	s.mu.Lock()
	hasil := []Peminjaman{} // bukan nil, supaya JSON-nya [] bukan null
	for _, p := range s.peminjaman {
		if filterSD != 0 && p.SumberDayaID != filterSD {
			continue
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && p.Status != filterStatus {
			continue
		}
		hasil = append(hasil, p)
	}
	s.mu.Unlock()
	writeJSON(w, 200, hasil)
}

// B3: GET /peminjaman/{id}
func (s *Store) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.cariIndexPeminjaman(id)
	if i < 0 {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, 200, s.peminjaman[i])
}

// B6: PUT /peminjaman/{id} mengganti enam field, id dan status tidak berubah.
// Urutan: id bukan angka 400 -> tidak ada 404 -> body tidak valid 400 -> bentrok 409.
func (s *Store) ubahPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	// cek 404 dulu sebelum membaca body
	s.mu.Lock()
	ada := s.cariIndexPeminjaman(id) >= 0
	s.mu.Unlock()
	if !ada {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}

	var in PeminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "body JSON tidak valid")
		return
	}
	if pesan := validasiPeminjaman(&in); pesan != "" {
		writeError(w, 400, pesan)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	i := s.cariIndexPeminjaman(id) // cari lagi, siapa tahu dihapus selama lock dilepas
	if i < 0 {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}
	sd, ada := s.cariSumberDaya(in.SumberDayaID)
	if !ada {
		writeError(w, 400, "sumber_daya_id tidak ditemukan")
		return
	}
	if in.Jumlah > sd.Stok {
		writeError(w, 400, "jumlah melebihi stok sumber daya")
		return
	}
	// id sendiri dikecualikan supaya tidak bentrok dengan dirinya sendiri
	if s.adaBentrok(in, sd.Stok, id) {
		writeError(w, 409, "slot sudah penuh, kapasitas terlampaui")
		return // data lama tidak diubah
	}

	p := &s.peminjaman[i]
	p.SumberDayaID = in.SumberDayaID
	p.Peminjam = in.Peminjam
	p.Tanggal = in.Tanggal
	p.SlotMulai = in.SlotMulai
	p.SlotSelesai = in.SlotSelesai
	p.Jumlah = in.Jumlah
	writeJSON(w, 200, *p)
}

// B6: DELETE /peminjaman/{id} menghapus data, 204 tanpa body.
func (s *Store) hapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.cariIndexPeminjaman(id)
	if i < 0 {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}
	s.peminjaman = append(s.peminjaman[:i], s.peminjaman[i+1:]...)
	w.WriteHeader(204)
}
