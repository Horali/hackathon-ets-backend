package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Persona C: alur persetujuan diajukan -> disetujui/ditolak -> selesai.

// gantiStatus dipakai setujui dan selesaikan (tanpa body).
// Urutan: id bukan angka 400 -> tidak ada 404 -> status tidak sesuai 409.
func (s *Store) gantiStatus(w http.ResponseWriter, r *http.Request, statusSyarat, statusBaru string) {
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
	if s.peminjaman[i].Status != statusSyarat {
		writeError(w, 409, "status saat ini harus "+statusSyarat)
		return // data tidak berubah
	}
	s.peminjaman[i].Status = statusBaru
	writeJSON(w, 200, s.peminjaman[i])
}

// POST /peminjaman/{id}/setujui: diajukan -> disetujui
func (s *Store) setujuiPeminjaman(w http.ResponseWriter, r *http.Request) {
	s.gantiStatus(w, r, "diajukan", "disetujui")
}

// POST /peminjaman/{id}/selesaikan: disetujui -> selesai
func (s *Store) selesaikanPeminjaman(w http.ResponseWriter, r *http.Request) {
	s.gantiStatus(w, r, "disetujui", "selesai")
}

// POST /peminjaman/{id}/tolak dengan body {"alasan": "..."}: diajukan -> ditolak
// Urutan: id bukan angka 400 -> tidak ada 404 -> body tidak valid 400 -> status tidak sesuai 409.
func (s *Store) tolakPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}

	// body dibaca dulu, tapi error-nya baru dilaporkan setelah cek 404
	var in struct {
		Alasan string `json:"alasan"`
	}
	pesanBody := ""
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pesanBody = "body JSON tidak valid"
	} else {
		in.Alasan = strings.TrimSpace(in.Alasan)
		if in.Alasan == "" || utf8.RuneCountInString(in.Alasan) > 100 {
			pesanBody = "alasan wajib, maksimal 100 karakter"
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.cariIndexPeminjaman(id)
	if i < 0 {
		writeError(w, 404, "peminjaman tidak ditemukan")
		return
	}
	if pesanBody != "" {
		writeError(w, 400, pesanBody)
		return
	}
	if s.peminjaman[i].Status != "diajukan" {
		writeError(w, 409, "status saat ini harus diajukan")
		return
	}
	s.peminjaman[i].Status = "ditolak"
	s.peminjaman[i].Alasan = in.Alasan
	writeJSON(w, 200, s.peminjaman[i])
}
