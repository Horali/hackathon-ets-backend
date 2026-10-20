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

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
