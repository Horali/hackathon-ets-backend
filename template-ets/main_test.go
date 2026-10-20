package main

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// B7: 50 request POST bersamaan ke slot yang sama (stok 1) harus menghasilkan tepat satu 201, sisanya 409.
// Diulang 20 kali karena race condition belum tentu muncul di percobaan pertama.
func TestB7RequestBersamaan(t *testing.T) {
	for ulang := 0; ulang < 20; ulang++ {
		h := NewHandler()

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("POST", "/sumber-daya",
			strings.NewReader(`{"nama":"Ruang Hima","jenis":"ruang","stok":1}`)))

		body := `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}`

		var wg sync.WaitGroup
		var mu sync.Mutex
		hitung := map[int]int{}

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, httptest.NewRequest("POST", "/peminjaman", strings.NewReader(body)))
				mu.Lock()
				hitung[rec.Code]++
				mu.Unlock()
			}()
		}
		wg.Wait()

		if hitung[201] != 1 || hitung[409] != 49 {
			t.Fatalf("percobaan %d: harus 1x201 dan 49x409, dapat %v", ulang+1, hitung)
		}
	}
}
