package config

import (
	"fmt"
	"sipadin-backend/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		hashedPassword = []byte("$2a$10$7/kMsz8V7zZpU3o7XpWq4eYfW2QyV5K3X7Gg9K2p8P9z3P1x4p6l6")
	}

	if userCount == 0 {
		users := []models.User{
			{
				ID:         1,
				Nama:       "Rizky Darmawan",
				NIP:        "EMP-2022091",
				Email:      "karyawan@sipadin.com",
				Password:   string(hashedPassword),
				Role:       "karyawan",
				Departemen: "IT & Infrastructure",
				Jabatan:    "Senior Web Developer",
				CreatedAt:  time.Now(),
			},
			{
				ID:         2,
				Nama:       "Budi Pratama",
				NIP:        "MGR-2018004",
				Email:      "manager@sipadin.com",
				Password:   string(hashedPassword),
				Role:       "atasan",
				Departemen: "IT & Operations",
				Jabatan:    "IT & Ops Manager",
				CreatedAt:  time.Now(),
			},
			{
				ID:         3,
				Nama:       "Dewi Rahayu",
				NIP:        "HRD-2019007",
				Email:      "hrga@sipadin.com",
				Password:   string(hashedPassword),
				Role:       "hrga",
				Departemen: "Human Resources & General Affairs",
				Jabatan:    "HRGA Supervisor",
				CreatedAt:  time.Now(),
			},
		}

		for _, u := range users {
			DB.Create(&u)
		}
		fmt.Println("Seed default users dengan Bcrypt berhasil!")
	} else {
		// Pastikan semua user yang sudah ada di database juga terenkripsi Bcrypt
		var existingUsers []models.User
		DB.Find(&existingUsers)
		for _, u := range existingUsers {
			if !strings.HasPrefix(u.Password, "$2a$") && !strings.HasPrefix(u.Password, "$2b$") && !strings.HasPrefix(u.Password, "$2y$") {
				plain := u.Password
				if plain == "" {
					plain = "password123"
				}
				if hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost); err == nil {
					DB.Model(&models.User{}).Where("id = ?", u.ID).Update("password", string(hash))
				}
			}
		}
	}

	var pengajuanCount int64
	DB.Model(&models.PengajuanDinas{}).Count(&pengajuanCount)
	if pengajuanCount == 0 {
		today := time.Now()
		d := func(offset int) string {
			return today.AddDate(0, 0, offset).Format("2006-01-02")
		}

		initialPengajuan := []models.PengajuanDinas{
			{
				ID:               101,
				NomorSurat:       "SPPD/IT/2026/08/001",
				UserID:           1,
				Nama:             "Rizky Darmawan",
				NIP:              "EMP-2022091",
				Dept:             "IT & Infrastructure",
				Jabatan:          "Senior Web Developer",
				Asal:             "Surabaya",
				Tujuan:           "Jakarta Pusat",
				Instansi:         "PT Mitra Solusi Nusantara",
				TglKeberangkatan: d(2),
				TglKembali:       d(5),
				JamBerangkat:     "08:00",
				Jarak:            780,
				LamaJam:          96,
				Durasi:           4,
				StatusMenginap:   "Ya",
				Kriteria:         "4. Menginap",
				Transportasi:     "Pesawat Terbang (Garuda Indonesia)",
				Akomodasi:        "Hotel Santika Premiere Slipi",
				Keperluan:        "Rapat / Meeting Klien",
				Keterangan:       "Presentasi dan implementasi arsitektur sistem ERP terintegrasi kuartal 3 bersama jajaran direksi mitra.",
				BiayaTransport:   2400000,
				BiayaPenginapan:  1800000,
				UangHarian:       1200000,
				EstimasiBiaya:    5400000,
				Lampiran:         "Surat_Undangan_Mitra_Jakarta.pdf",
				Status:           "pending",
				TanggalPengajuan: d(-2),
				CreatedAt:        today.AddDate(0, 0, -2),
				UpdatedAt:        today.AddDate(0, 0, -2),
			},
			{
				ID:               102,
				NomorSurat:       "SPPD/MKT/2026/08/002",
				UserID:           1,
				Nama:             "Sari Dewi",
				NIP:              "EMP-2023045",
				Dept:             "Marketing & Sales",
				Jabatan:          "Marketing Specialist",
				Asal:             "Surabaya",
				Tujuan:           "Denpasar, Bali",
				Instansi:         "Bali Nusa Dua Convention Center",
				TglKeberangkatan: d(1),
				TglKembali:       d(4),
				JamBerangkat:     "09:00",
				Jarak:            420,
				LamaJam:          96,
				Durasi:           4,
				StatusMenginap:   "Ya",
				Kriteria:         "4. Menginap",
				Transportasi:     "Pesawat Terbang (Batik Air)",
				Akomodasi:        "Grand Hyatt Bali",
				Keperluan:        "Survey Lapangan / Venue",
				Keterangan:       "Survey dan koordinasi booth pameran tahunan industri percetakan dan packaging nasional.",
				BiayaTransport:   2100000,
				BiayaPenginapan:  2400000,
				UangHarian:       1400000,
				EstimasiBiaya:    5900000,
				Lampiran:         "Proposal_Survey_Pameran_Bali.pdf",
				Status:           "pending",
				TanggalPengajuan: d(-1),
				CreatedAt:        today.AddDate(0, 0, -1),
				UpdatedAt:        today.AddDate(0, 0, -1),
			},
			{
				ID:               103,
				NomorSurat:       "SPPD/FIN/2026/08/003",
				UserID:           1,
				Nama:             "Andi Kusuma",
				NIP:              "EMP-2021018",
				Dept:             "Finance & Accounting",
				Jabatan:          "Senior Auditor",
				Asal:             "Surabaya",
				Tujuan:           "Medan, Sumatera Utara",
				Instansi:         "Kantor Cabang Adiprima Medan",
				TglKeberangkatan: d(7),
				TglKembali:       d(11),
				JamBerangkat:     "07:30",
				Jarak:            1800,
				LamaJam:          120,
				Durasi:           5,
				StatusMenginap:   "Ya",
				Kriteria:         "4. Menginap",
				Transportasi:     "Pesawat Terbang (Citilink)",
				Akomodasi:        "JW Marriott Medan",
				Keperluan:        "Audit Internal Cabang",
				Keterangan:       "Pelaksanaan audit kepatuhan keuangan semester 1 dan rekonsiliasi aset logistik regional barat.",
				BiayaTransport:   3200000,
				BiayaPenginapan:  3000000,
				UangHarian:       1750000,
				EstimasiBiaya:    7950000,
				Lampiran:         "Surat_Tugas_Audit_Medan.pdf",
				Status:           "approved",
				TanggalPengajuan: d(-3),
				CatatanAtasan:    "Disetujui. Pastikan membawa dokumen rekonsiliasi lengkap dari HO.",
				DisetujuiOleh:    "Budi Pratama (IT & Ops Manager)",
				TanggalDisetujui: d(-1),
				CreatedAt:        today.AddDate(0, 0, -3),
				UpdatedAt:        today.AddDate(0, 0, -1),
			},
			{
				ID:               104,
				NomorSurat:       "SPPD/HRD/2026/08/004",
				UserID:           1,
				Nama:             "Lina Marlina",
				NIP:              "EMP-2022110",
				Dept:             "Human Resources",
				Jabatan:          "HR Development Lead",
				Asal:             "Surabaya",
				Tujuan:           "Yogyakarta",
				Instansi:         "Pusat Diklat SDM Ketenagakerjaan",
				TglKeberangkatan: d(5),
				TglKembali:       d(8),
				JamBerangkat:     "06:30",
				Jarak:            330,
				LamaJam:          96,
				Durasi:           4,
				StatusMenginap:   "Ya",
				Kriteria:         "4. Menginap",
				Transportasi:     "Kereta Api (Executive Argo Wilis)",
				Akomodasi:        "Hotel Tentrem Yogyakarta",
				Keperluan:        "Pelatihan / Training SDM",
				Keterangan:       "Mengikuti sertifikasi BNSP manajemen talenta dan kepemimpinan operasional industri 4.0.",
				BiayaTransport:   900000,
				BiayaPenginapan:  1800000,
				UangHarian:       1200000,
				EstimasiBiaya:    3900000,
				Lampiran:         "Undangan_Diklat_SDM_Jogja.pdf",
				Status:           "approved",
				TanggalPengajuan: d(-5),
				CatatanAtasan:    "Disetujui penuh. Mohon buat laporan hasil sertifikasi setelah kembali.",
				DisetujuiOleh:    "Budi Pratama (IT & Ops Manager)",
				TanggalDisetujui: d(-3),
				CreatedAt:        today.AddDate(0, 0, -5),
				UpdatedAt:        today.AddDate(0, 0, -3),
			},
			{
				ID:                    105,
				NomorSurat:            "SPPD/IT/2026/08/005",
				UserID:                1,
				Nama:                  "Bima Santoso",
				NIP:                   "EMP-2023011",
				Dept:                  "IT & Infrastructure",
				Jabatan:               "Network Engineer",
				Asal:                  "Surabaya",
				Tujuan:                "Bandung",
				Instansi:              "Data Center Lintasarta Bandung",
				TglKeberangkatan:      d(3),
				TglKembali:            d(6),
				JamBerangkat:          "19:00",
				Jarak:                 700,
				LamaJam:               96,
				Durasi:                4,
				StatusMenginap:        "Ya",
				Kriteria:              "4. Menginap",
				Transportasi:          "Kereta Api (Turangga Executive)",
				Akomodasi:             "Aston Tropicana Bandung",
				Keperluan:             "Maintenance & Proyek",
				Keterangan:            "Migrasi core switch server dan konfigurasi VPN private network antar site pabrik.",
				BiayaTransport:        1100000,
				BiayaPenginapan:       1600000,
				UangHarian:            1200000,
				EstimasiBiaya:         3900000,
				Lampiran:              "Work_Order_Network_Bandung.pdf",
				Status:                "verified_hrga",
				TanggalPengajuan:      d(-8),
				CatatanAtasan:         "Pastikan uptime sistem tetap terjaga selama migrasi berlangsung.",
				DisetujuiOleh:         "Budi Pratama (IT & Ops Manager)",
				TanggalDisetujui:      d(-6),
				CatatanHRGA:           "Dokumen lengkap. Estimasi biaya sesuai tarif standar area Jawa. Disetujui untuk dilaporkan ke Direksi.",
				DiverifikasiOleh:      "Dewi Rahayu (HRGA Supervisor)",
				TanggalVerifikasiHRGA: d(-4),
				CreatedAt:             today.AddDate(0, 0, -8),
				UpdatedAt:             today.AddDate(0, 0, -4),
			},
			{
				ID:                       106,
				NomorSurat:               "SPPD/PRD/2026/08/006",
				UserID:                   1,
				Nama:                     "Rizky Darmawan",
				NIP:                      "EMP-2022091",
				Dept:                     "IT & Infrastructure",
				Jabatan:                  "Senior Web Developer",
				Asal:                     "Surabaya",
				Tujuan:                   "Semarang",
				Instansi:                 "PT Percetakan Grafika Utama",
				TglKeberangkatan:         d(-5),
				TglKembali:               d(-2),
				JamBerangkat:             "08:00",
				Jarak:                    350,
				LamaJam:                  72,
				Durasi:                   3,
				StatusMenginap:           "Ya",
				Kriteria:                 "4. Menginap",
				Transportasi:             "Mobil Dinas Operasional",
				Akomodasi:                "Hotel Gumaya Tower Semarang",
				Keperluan:                "Pelatihan / Training SDM",
				Keterangan:               "Training operasional integrasi barcode inventory mesin cetak digital.",
				BiayaTransport:           650000,
				BiayaPenginapan:          1200000,
				UangHarian:               900000,
				EstimasiBiaya:            2750000,
				Lampiran:                 "SOP_Mesin_Cetak_Semarang.pdf",
				Status:                   "notified_direksi",
				TanggalPengajuan:         d(-14),
				CatatanAtasan:            "Pelatihan disetujui.",
				DisetujuiOleh:            "Budi Pratama (IT & Ops Manager)",
				TanggalDisetujui:         d(-12),
				CatatanHRGA:              "Semua persyaratan terpenuhi. Sesuai kebijakan dinas area Jawa.",
				DiverifikasiOleh:         "Dewi Rahayu (HRGA Supervisor)",
				TanggalVerifikasiHRGA:    d(-10),
				TanggalNotifikasiDireksi: d(-9),
				PesanWA:                  "Pesan notifikasi Direksi telah dikirim via WhatsApp.",
				CreatedAt:                today.AddDate(0, 0, -14),
				UpdatedAt:                today.AddDate(0, 0, -9),
			},
			{
				ID:               107,
				NomorSurat:       "SPPD/MKT/2026/08/007",
				UserID:           1,
				Nama:             "Dian Pertiwi",
				NIP:              "EMP-2024003",
				Dept:             "Marketing & Sales",
				Jabatan:          "Brand Associate",
				Asal:             "Surabaya",
				Tujuan:           "Makassar",
				Instansi:         "Distributor Nusantara Timur",
				TglKeberangkatan: d(10),
				TglKembali:       d(13),
				JamBerangkat:     "08:00",
				Jarak:            1100,
				LamaJam:          96,
				Durasi:           4,
				StatusMenginap:   "Ya",
				Kriteria:         "4. Menginap",
				Transportasi:     "Pesawat Terbang (Lion Air)",
				Akomodasi:        "Hotel Four Points Makassar",
				Keperluan:        "Survey Lapangan / Venue",
				Keterangan:       "Kunjungan rutin ke distributor regional timur untuk evaluasi pencapaian target Q3.",
				BiayaTransport:   2800000,
				BiayaPenginapan:  2100000,
				UangHarian:       1400000,
				EstimasiBiaya:    6300000,
				Lampiran:         "Jadwal_Kunjungan_Distributor_Makassar.pdf",
				Status:           "rejected",
				TanggalPengajuan: d(-4),
				CatatanAtasan:    "Ditolak. Kunjungan dapat dilakukan secara virtual. Jadwalkan ulang jika benar-benar diperlukan kunjungan fisik.",
				CreatedAt:        today.AddDate(0, 0, -4),
				UpdatedAt:        today.AddDate(0, 0, -2),
			},
		}

		for _, p := range initialPengajuan {
			DB.Create(&p)
		}
		fmt.Println("Seed default pengajuan berhasil!")
	}
}
