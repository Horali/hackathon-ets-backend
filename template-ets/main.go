package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, pesan string) {
	writeJSON(w, code, map[string]string{"error": pesan})
}

func (s *Store) buatSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in SumberDaya
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, 400, "body JSON tidak valid")
		return
	}
	in.Nama = strings.TrimSpace(in.Nama)
	in.Jenis = strings.TrimSpace(in.Jenis)
	if in.Nama == "" || utf8.RuneCountInString(in.Nama) > 50 {
		writeError(w, 400, "nama wajib, maksimal 50 karakter")
		return
	}
	if in.Jenis != "ruang" && in.Jenis != "alat" {
		writeError(w, 400, "jenis harus ruang atau alat")
		return
	}
	if in.Stok < 1 {
		writeError(w, 400, "stok minimal 1")
		return
	}
	s.mu.Lock()
	in.ID = s.nextSDID
	s.nextSDID++
	s.sumberDaya = append(s.sumberDaya, in)
	s.mu.Unlock()
	writeJSON(w, 201, in)
}

func (s *Store) daftarSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	hasil := append([]SumberDaya{}, s.sumberDaya...) // [] bukan null
	s.mu.Unlock()
	writeJSON(w, 200, hasil)
}

func (s *Store) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "id harus angka")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			writeJSON(w, 200, sd)
			return
		}
	}
	writeError(w, 404, "sumber daya tidak ditemukan")
}

func NewHandler() http.Handler {
	s := &Store{nextSDID: 1, nextPJID: 1}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", s.buatSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.daftarSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.detailSumberDaya)
	mux.HandleFunc("POST /peminjaman", s.buatPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.daftarPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.detailPeminjaman)
	mux.HandleFunc("PUT /peminjaman/{id}", s.ubahPeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", s.batalkanPeminjaman)
	mux.HandleFunc("POST /peminjaman/{id}/setujui", s.setujuiPeminjaman)
	mux.HandleFunc("POST /peminjaman/{id}/tolak", s.tolakPeminjaman)
	mux.HandleFunc("POST /peminjaman/{id}/selesaikan", s.selesaikanPeminjaman)
	return mux
}

func main() {
	http.ListenAndServe(":8080", NewHandler())
}
