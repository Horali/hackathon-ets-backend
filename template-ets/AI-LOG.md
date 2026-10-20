# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

AI yang dipakai: Claude (Claude Code).

## Setup dan B1 (sumber daya)
- Prompt: "bantu saya setup project Go dan implementasi B1"
- Dipakai: `buatSumberDaya`, `daftarSumberDaya`, `detailSumberDaya`, `writeJSON`, dan `writeError` di `main.go`.
- Diubah: versi di `go.mod` saya ubah dari 1.26.2 ke 1.22 supaya bisa di-build dengan Go versi lama. Perintah `go mod edit -go=1.22` sempat gagal di PowerShell, jadi argumennya harus diberi tanda kutip.

## Memecah main.go menjadi beberapa file
- Prompt: "bagaimana cara saya memecah main.go menjadi beberapa file"
- Dipakai: pembagian file per fungsi, yaitu struct di `store.go`, handler peminjaman di `handler_peminjaman.go`, dan validasi di `validasi.go`.
- Diubah: waktu menghapus struct lama dari `main.go`, dua handler sumber daya ikut terhapus sehingga error (Bug 1 di PROSES.md), jadi saya kembalikan dari commit sebelumnya.

## B2 sampai B4 (buat dan lihat peminjaman, validasi)
- Prompt: "bantu saya implementasi B2 sampai B4"
- Dipakai: `buatPeminjaman`, `daftarPeminjaman`, `detailPeminjaman` di `handler_peminjaman.go`, dan `validasiPeminjaman` di `validasi.go`.
- Diubah: status awal disesuaikan dengan persona C menjadi "diajukan". Saat menguji validasi dengan curl, hasilnya sempat salah karena masalah tanda kutip di PowerShell (Bug 2 di PROSES.md).

## B5 (deteksi bentrok)
- Prompt: "bantu saya implementasi B5"
- Dipakai: `menahanSlot`, `terpakai`, dan `adaBentrok` di `bentrok.go`.
- Diubah: dipakai apa adanya. `menahanSlot` hanya menghitung status diajukan dan disetujui sesuai persona C.

## B6 (ubah dan hapus)
- Prompt: "bantu saya implementasi B6"
- Dipakai: `ubahPeminjaman` dan handler DELETE di `handler_peminjaman.go`.
- Diubah: `kecualiID` di pemanggilan `adaBentrok` pada PUT harus diisi `id` peminjaman, bukan 0, supaya tidak bentrok dengan dirinya sendiri (Bug 3 di PROSES.md).

## Fitur khas persona C (setujui, tolak, selesaikan)
- Prompt: "bantu saya implementasi fitur persona C"
- Dipakai: `gantiStatus`, `setujuiPeminjaman`, `selesaikanPeminjaman`, dan `tolakPeminjaman` di `handler_persona.go`.
- Diubah: dipakai apa adanya.

## Change request C (pembatalan)
- Prompt: "bantu saya implementasi change request C"
- Dipakai: `batalkanPeminjaman` di `handler_peminjaman.go`.
- Diubah: fungsi DELETE lama diganti nama dari `hapusPeminjaman` menjadi `batalkanPeminjaman`. Penjelasannya ada di bagian 5 PROSES.md.

## B7 (permintaan bersamaan)
- Prompt: "bantu saya implementasi B7"
- Dipakai: test `TestB7RequestBersamaan` di `main_test.go`.
- Diubah: dipakai apa adanya. Test-nya dijalankan dengan `go test -v .` dan lulus.
