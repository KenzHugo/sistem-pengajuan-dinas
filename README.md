# SiPaDin — Workflow Sistem Pengajuan Dinas

Dokumen ini menjelaskan alur kerja aplikasi SiPaDin untuk karyawan, atasan, dan HRGA.

## 1. Tujuan Aplikasi
SiPaDin adalah sistem pengajuan perjalanan dinas yang membantu:
- karyawan mengajukan perjalanan dinas,
- atasan menilai dan menyetujui pengajuan,
- HRGA memverifikasi kebijakan anggaran dan fasilitas,
- direksi menerima notifikasi status pengajuan secara digital.

---

## 2. Flow Login
1. User membuka halaman `index.html`.
2. User memilih role:
   - Karyawan
   - Atasan
   - HRGA
3. Sistem otomatis mengisi email sesuai role demo.
4. User memasukkan password dan klik tombol masuk.
5. Setelah login, browser akan redirect ke dashboard sesuai role:
   - Karyawan → `dashboard-karyawan.html`
   - Atasan → `dashboard-atasan.html`
   - HRGA → `dashboard-hrga.html`

### Role Demo
- Karyawan: `karyawan@sipadin.com`
- Atasan: `manager@sipadin.com`
- HRGA: `hrga@sipadin.com`

Password demo yang dipakai adalah:
- `password123`

---

## 3. Workflow Utama Sistem

### Tahap 1 — Pengajuan oleh Karyawan
1. Karyawan masuk ke dashboard.
2. Klik menu `Pengajuan Dinas`.
3. Isi form pengajuan dengan data seperti:
   - nama dan NIP
   - departemen dan jabatan
   - tujuan dan instansi
   - tanggal keberangkatan dan kembali
   - transportasi dan akomodasi
   - keperluan dan keterangan
   - estimasi biaya
4. Simpan data pengajuan.
5. Status awal pengajuan adalah:
   - `pending`

Artinya: pengajuan sedang menunggu persetujuan atasan.

---

### Tahap 2 — Persetujuan Atasan
1. Atasan login ke dashboard atasan.
2. Sistem menampilkan daftar seluruh pengajuan karyawan yang menunggu review.
3. Atasan dapat:
   - menyetujui pengajuan,
   - menolak pengajuan,
   - melihat detail pengajuan dan estimasi biaya.
4. Jika disetujui:
   - status berubah menjadi `approved`
   - data diteruskan ke HRGA untuk kontrol lanjutan.
5. Jika ditolak:
   - status berubah menjadi `rejected`
   - pengajuan berhenti di tahap ini.

---

### Tahap 3 — Verifikasi HRGA
1. HRGA login ke dashboard HRGA.
2. HRGA mengecek pengajuan yang sudah disetujui atasan.
3. HRGA memeriksa:
   - kesesuaian biaya,
   - fasilitas perjalanan dinas,
   - tarif berdasarkan jabatan dan wilayah tujuan,
   - apakah pengajuan sesuai dengan kebijakan perusahaan.
4. Jika sesuai:
   - status berubah menjadi `verified_hrga`
   - siap dikirim ke Direksi.
5. Jika tidak sesuai:
   - status berubah menjadi `rejected_hrga`
   - pengajuan ditolak oleh HRGA.

---

### Tahap 4 — Notifikasi ke Direksi
1. HRGA mengirimkan notifikasi ke Direksi.
2. Sistem membuat pesan WhatsApp/notification yang berisi:
   - nama karyawan,
   - jabatan,
   - tujuan dinas,
   - estimasi biaya,
   - fasilitas yang diajukan,
   - status pengajuan.
3. Setelah dikirim:
   - status berubah menjadi `notified_direksi`
   - pekerjaan workflow dinas dinyatakan selesai.

---

## 4. Status Workflow Sistem
Aplikasi memakai pipeline status berikut:

- `pending` → Menunggu approval atasan
- `approved` → Disetujui atasan, menunggu kontrol HRGA
- `rejected` → Ditolak atasan
- `verified_hrga` → Terverifikasi HRGA
- `rejected_hrga` → Ditolak HRGA
- `notified_direksi` → Sudah dikirimkan ke Direksi

Flow status dapat digambarkan seperti ini:

```text
Karyawan mengajukan
        ↓
pending
        ↓
Atasan approve/reject
        ↓
approved / rejected
        ↓
HRGA verify/reject
        ↓
verified_hrga / rejected_hrga
        ↓
notified_direksi
```

---

## 5. Struktur Dashboard Berdasarkan Role

### Dashboard Karyawan
Fungsi utama:
- melihat daftar pengajuan sendiri,
- pantau status pengajuan,
- lihat riwayat perjalanan dinas,
- melihat total biaya dan status aktif.

### Dashboard Atasan
Fungsi utama:
- meninjau pengajuan bawahan,
- menyetujui atau menolak pengajuan,
- melihat kebutuhan prioritas berdasarkan tanggal keberangkatan.

### Dashboard HRGA
Fungsi utama:
- menilai kelayakan biaya dan fasilitas,
- mengecek referensi tarif dan kebijakan perusahaan,
- mengirim notifikasi akhir ke Direksi.

---

## 6. Data Referensi yang Dipakai
Aplikasi menggunakan referensi yang terkait dengan SQL yang tersedia di folder proyek:
- `PerjalananDinas_Tarif.sql`
- `PerjalananDinas_Fasilitas.sql`

Data ini digunakan untuk:
- menghitung tarif berdasarkan jabatan dan wilayah,
- memeriksa apakah fasilitas perjalanan dinas sudah sesuai,
- memastikan aturan perusahaan diterapkan saat HRGA memverifikasi biaya.

---

## 7. Alur Kerja dari Sudut Pandang User

### Karyawan
1. Login
2. Buat pengajuan dinas
3. Menunggu persetujuan atasan
4. Lanjutkan jika sudah disetujui
5. Lihat status sampai dinyatakan selesai

### Atasan
1. Login
2. Cek daftar pengajuan
3. Review dokumen dan kebutuhan dinas
4. Setujui atau tolak
5. Lanjutkan ke proses HRGA

### HRGA
1. Login
2. Ambil pengajuan yang sudah disetujui atasan
3. Verifikasi tarif dan fasilitas
4. Lanjutkan ke notifikasi Direksi
5. Tutup workflow jika sudah dikirim

---

## 8. Kesimpulan
SiPaDin memiliki workflow yang jelas dan terstruktur:
- Karyawan mengajukan,
- Atasan memberi approval,
- HRGA melakukan verifikasi kebijakan,
- Direksi menerima notifikasi final.

Dengan model ini, proses perjalanan dinas menjadi lebih terkontrol, terdokumentasi, dan lebih cepat untuk dipantau.

---

## 9. Catatan
File aplikasi utama terdiri dari:
- `index.html` → login
- `dashboard-karyawan.html` → dashboard karyawan
- `dashboard-atasan.html` → dashboard atasan
- `dashboard-hrga.html` → dashboard HRGA
- `pengajuan-dinas.html` → form pengajuan dinas
- `riwayat-pengajuan.html` → riwayat pengajuan
- `script.js` → logika utama aplikasi, status, dan workflow

Semoga dokumentasi ini membantu memahami alur kerja aplikasi secara keseluruhan.
