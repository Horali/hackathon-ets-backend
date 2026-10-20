# PROSES

## 1. Persona dan masalahnya
Persona saya adalah C, admin himpunan. Masalahnya, ruang himpunan tidak boleh langsung terpakai saat ada yang meminjam, tetapi harus disetujui admin terlebih dahulu. Karena itu, setiap peminjaman baru masuk dengan status "diajukan", lalu admin bisa menyetujui atau menolaknya (penolakan wajib disertai alasan), dan setelah selesai dipakai statusnya diubah menjadi "selesai" agar slotnya kembali tersedia.

## 2. Rancangan struktur data
Rencana awal saya, semua data disimpan di satu struct `Store` yang berisi:
- slice `SumberDaya` (id, nama, jenis, stok)
- slice `Peminjaman` (field sesuai kontrak, ditambah `status` dan `alasan` untuk persona C)
- counter id untuk masing-masing jenis data
- satu `sync.Mutex` untuk mengunci data saat diakses bersamaan

Untuk sementara saya pilih slice karena daftar harus urut berdasarkan id, dan append dengan id yang terus naik sudah otomatis urut. Kemungkinan masih berubah, misalnya kalau nanti ternyata lebih praktis pakai map untuk mencari data berdasarkan id.

Update: setelah B3 selesai, saya tetap pakai slice. Mencari data cukup dengan loop (`cariSumberDaya`, `cariIndexPeminjaman`), dan struct-nya saya pindah ke `store.go` supaya `main.go` tidak terlalu panjang.

## 3. Satu keputusan desain
Saat membuat peminjaman, saya mengunci mutex sekali di awal (`s.mu.Lock()` lalu `defer s.mu.Unlock()`) dan baru dilepas setelah data selesai disimpan. Jadi pengecekan sumber daya, stok, dan nanti bentrok, sampai append ke slice terjadi dalam satu kunci yang sama.

Pilihan lain yang sempat saya pikirkan adalah mengunci hanya di bagian append saja. Tapi kalau begitu, dua request yang datang bersamaan bisa sama-sama lolos pengecekan dulu sebelum salah satunya menyimpan, sehingga dua peminjaman bisa masuk ke slot yang sama. Validasi body (trim, format tanggal, slot) tetap saya lakukan sebelum lock karena tidak menyentuh data bersama, jadi request yang body-nya salah tidak perlu menunggu kunci.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: setelah memindahkan struct dari `main.go` ke `store.go`, program tidak bisa di-build. `go vet` menampilkan error `s.daftarSumberDaya undefined (type *Store has no field or method daftarSumberDaya)`.
- Cara menemukan: menjalankan `go vet ./...` setelah memecah file, lalu membuka baris yang disebut di pesan errornya (pendaftaran route di `NewHandler`).
- Penyebab dan solusi: saat menghapus struct lama dari `main.go`, fungsi `daftarSumberDaya` dan `detailSumberDaya` ikut terhapus tanpa sengaja, padahal route-nya masih terdaftar. Saya kembalikan kedua fungsi itu dari commit checkpoint 1, dan sejak itu selalu menjalankan `go vet` setiap habis memindahkan kode.

### Bug 2
- Gejala: saat menguji validasi B4 dengan curl, request dengan nama kosong (`"   "`) memang dapat 400, tetapi pesannya "body JSON tidak valid", bukan "nama wajib". Di terminal juga muncul error `curl: (3) URL rejected`.
- Cara menemukan: saya cek isi pesan error-nya, bukan hanya kode statusnya. Ternyata semua request POST dari PowerShell dapat pesan yang sama, jadi masalahnya bukan di kode validasi.
- Penyebab dan solusi: PowerShell 5.1 membuang tanda kutip `\"` di dalam argumen `-d`, sehingga yang sampai ke server adalah JSON yang rusak. Solusinya menambahkan `--%` setelah `curl.exe` supaya PowerShell tidak mengubah argumennya. Dari sini saya jadi selalu cek pesan error-nya juga, karena kode status yang benar belum tentu muncul karena alasan yang benar.

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
