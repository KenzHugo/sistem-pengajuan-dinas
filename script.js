// ============================================================
// SIPADIN — SISTEM PENGAJUAN DINAS
// GLOBAL JAVASCRIPT & STATE MANAGEMENT
// v2.0 — HRGA 4-Stage Workflow Edition
// ============================================================

// Tailwind Theme Configuration
if (typeof tailwind !== 'undefined') {
    tailwind.config = {
        theme: {
            extend: {
                fontFamily: {
                    sans: ['"Plus Jakarta Sans"', 'Poppins', 'sans-serif'],
                    poppins: ['Poppins', 'sans-serif'],
                },
                colors: {
                    brand: {
                        50: '#eef2ff',
                        100: '#e0e7ff',
                        200: '#c7d2fe',
                        300: '#a5b4fc',
                        400: '#818cf8',
                        500: '#6366f1',
                        600: '#4f46e5',
                        700: '#4338ca',
                        800: '#3730a3',
                        900: '#312e81',
                        950: '#1e1b4b',
                    },
                    surface: {
                        DEFAULT: '#0b0f19',
                        card: '#1e293b',
                        cardhover: '#24334a',
                        border: '#334155',
                        muted: '#475569',
                    }
                },
                boxShadow: {
                    glow: '0 0 35px -5px rgba(99, 102, 241, 0.45)',
                    glowGreen: '0 0 35px -5px rgba(16, 185, 129, 0.45)',
                    glowRed: '0 0 35px -5px rgba(239, 68, 68, 0.45)',
                    glowPurple: '0 0 35px -5px rgba(168, 85, 247, 0.45)',
                    glowAmber: '0 0 35px -5px rgba(245, 158, 11, 0.45)',
                    card: '0 10px 30px -10px rgba(0, 0, 0, 0.6)',
                },
                backgroundImage: {
                    'hero-gradient': 'radial-gradient(ellipse at top, #1e1b4b 0%, #0f172a 45%, #0b0f19 100%)',
                    'brand-gradient': 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
                    'accent-gradient': 'linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%)',
                    'emerald-gradient': 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                    'rose-gradient': 'linear-gradient(135deg, #f43f5e 0%, #e11d48 100%)',
                    'hrga-gradient': 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
                    'teal-gradient': 'linear-gradient(135deg, #14b8a6 0%, #0d9488 100%)',
                }
            }
        }
    };
}

// ============================================================
// DATE & HELPER UTILITIES
// ============================================================
const today = new Date();
today.setHours(0, 0, 0, 0);

function daysFromNow(n) {
    const d = new Date(today);
    d.setDate(d.getDate() + n);
    return d.toISOString().split('T')[0];
}

function formatDate(dateStr) {
    if (!dateStr) return '-';
    const d = new Date(dateStr + 'T00:00:00');
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
}

function formatDateFull(dateStr) {
    if (!dateStr) return '-';
    const d = new Date(dateStr + 'T00:00:00');
    return d.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
}

function formatRupiah(num) {
    if (isNaN(num)) return 'Rp 0';
    return 'Rp ' + Number(num).toLocaleString('id-ID');
}

function daysUntil(dateStr) {
    if (!dateStr) return 0;
    const target = new Date(dateStr + 'T00:00:00');
    return Math.ceil((target - today) / (1000 * 60 * 60 * 24));
}

function calculateDuration(startStr, endStr) {
    if (!startStr || !endStr) return 1;
    const start = new Date(startStr + 'T00:00:00');
    const end = new Date(endStr + 'T00:00:00');
    const diff = Math.ceil((end - start) / (1000 * 60 * 60 * 24)) + 1;
    return diff > 0 ? diff : 1;
}

function generateNomorSurat(dept) {
    const deptMap = {
        'IT & Infrastructure': 'IT', 'Marketing & Sales': 'MKT', 'Finance & Accounting': 'FIN',
        'Human Resources': 'HRD', 'Production': 'PRD', 'Operations': 'OPS', 'General Affairs': 'GA'
    };
    const kode = deptMap[dept] || 'GEN';
    const now = new Date();
    const yr = now.getFullYear();
    const mo = String(now.getMonth() + 1).padStart(2, '0');
    const list = getPengajuanData();
    const seq = String(list.length + 1).padStart(3, '0');
    return `SPPD/${kode}/${yr}/${mo}/${seq}`;
}

// ============================================================
// STATUS PIPELINE — 4 TAHAP HRGA WORKFLOW
// ============================================================
/**
 * STATUS PIPELINE:
 * 'pending'          → Tahap 1: Karyawan mengajukan, menunggu approval atasan
 * 'approved'         → Tahap 2: Atasan menyetujui, menunggu kontrol HRGA
 * 'rejected'         → Ditolak oleh atasan (terminal)
 * 'verified_hrga'    → Tahap 3: HRGA memverifikasi, siap notifikasi Direksi
 * 'rejected_hrga'    → Ditolak oleh HRGA (terminal)
 * 'notified_direksi' → Tahap 4: HRGA menginformasikan ke Direksi (selesai)
 */
const STATUS_PIPELINE = {
    pending:           { label: 'Menunggu Atasan',           icon: 'ph-clock',              color: 'amber',   step: 1 },
    approved:          { label: 'Pending Kontrol HRGA',      icon: 'ph-check-circle',       color: 'blue',    step: 2 },
    rejected:          { label: 'Ditolak Atasan',            icon: 'ph-x-circle',           color: 'rose',    step: 0 },
    verified_hrga:     { label: 'Terverifikasi HRGA',        icon: 'ph-shield-check',       color: 'teal',    step: 3 },
    rejected_hrga:     { label: 'Ditolak HRGA',              icon: 'ph-warning-octagon',    color: 'orange',  step: 0 },
    notified_direksi:  { label: 'Diinformasikan ke Direksi', icon: 'ph-paper-plane-tilt',   color: 'emerald', step: 4 },
};

function getStatusBadge(status, size) {
    const cfg = STATUS_PIPELINE[status] || STATUS_PIPELINE.pending;
    const sizeClass = size === 'xs' ? 'px-2 py-0.5 text-[10px]' : 'px-3 py-1 text-xs';
    const colorMap = {
        amber:   'bg-amber-500/20 text-amber-400 border-amber-500/30',
        blue:    'bg-blue-500/20 text-blue-400 border-blue-500/30',
        rose:    'bg-rose-500/20 text-rose-400 border-rose-500/30',
        teal:    'bg-teal-500/20 text-teal-400 border-teal-500/30',
        orange:  'bg-orange-500/20 text-orange-400 border-orange-500/30',
        emerald: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30',
    };
    const pulse = (status === 'pending' || status === 'approved') ? 'badge-pulse' : '';
    const colorClass = colorMap[cfg.color] || colorMap.amber;
    return `<span class="inline-flex items-center gap-1.5 ${sizeClass} rounded-full font-semibold border ${colorClass} ${pulse}">
        <i class="ph-fill ${cfg.icon}"></i> ${cfg.label}
    </span>`;
}

function getStatusLabel(status) {
    return (STATUS_PIPELINE[status] || STATUS_PIPELINE.pending).label;
}

function getStatusStep(status) {
    return (STATUS_PIPELINE[status] || STATUS_PIPELINE.pending).step;
}

// ============================================================
// REFERENSI DATA TARIF — dari PerjalananDinas_Tarif.sql
// ============================================================
const TARIF_DATA = {
    'Jawa': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 0.01, uangMakan: 80000,  uangSaku: 60000,  tiket: 0.1, hotel: 500000,  mataUang: 'IDR' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 0.01, uangMakan: 100000, uangSaku: 80000,  tiket: 0.1, hotel: 500000,  mataUang: 'IDR' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 0.01, uangMakan: 120000, uangSaku: 100000, tiket: 0.1, hotel: 650000,  mataUang: 'IDR' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 0.01, uangMakan: 150000, uangSaku: 150000, tiket: 0.1, hotel: 800000,  mataUang: 'IDR' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 0.01, uangMakan: 200000, uangSaku: 150000, tiket: 0.1, hotel: 900000,  mataUang: 'IDR' },
    ],
    'Jakarta / Luar Jawa': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 0.01, uangMakan: 100000, uangSaku: 75000,  tiket: 0.1, hotel: 600000,  mataUang: 'IDR' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 0.01, uangMakan: 150000, uangSaku: 100000, tiket: 0.1, hotel: 600000,  mataUang: 'IDR' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 0.01, uangMakan: 175000, uangSaku: 125000, tiket: 0.1, hotel: 750000,  mataUang: 'IDR' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 0.01, uangMakan: 225000, uangSaku: 150000, tiket: 0.1, hotel: 900000,  mataUang: 'IDR' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 0.1,  uangMakan: 225000, uangSaku: 150000, tiket: 0.1, hotel: 1100000, mataUang: 'IDR' },
    ],
    'Asia': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 20, uangMakan: 20, uangSaku: 12, tiket: 0.2, hotel: 70,  mataUang: 'USD' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 25, uangMakan: 30, uangSaku: 15, tiket: 0.2, hotel: 85,  mataUang: 'USD' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 35, uangMakan: 40, uangSaku: 25, tiket: 0.2, hotel: 105, mataUang: 'USD' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 50, uangMakan: 45, uangSaku: 40, tiket: 0.2, hotel: 145, mataUang: 'USD' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 60, uangMakan: 75, uangSaku: 75, tiket: 0.2, hotel: 0,   mataUang: 'USD' },
    ],
    'Eropa': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 25, uangMakan: 40,  uangSaku: 15,  tiket: 0.2, hotel: 100, mataUang: 'USD' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 25, uangMakan: 55,  uangSaku: 30,  tiket: 0.2, hotel: 115, mataUang: 'USD' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 30, uangMakan: 60,  uangSaku: 30,  tiket: 0.2, hotel: 145, mataUang: 'USD' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 50, uangMakan: 70,  uangSaku: 70,  tiket: 0.2, hotel: 225, mataUang: 'USD' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 90, uangMakan: 120, uangSaku: 110, tiket: 0.2, hotel: 0,   mataUang: 'USD' },
    ],
    'Amerika / Singapura': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 25, uangMakan: 40,  uangSaku: 15,  tiket: 0.2, hotel: 85,  mataUang: 'USD' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 25, uangMakan: 45,  uangSaku: 20,  tiket: 0.2, hotel: 100, mataUang: 'USD' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 40, uangMakan: 50,  uangSaku: 30,  tiket: 0.2, hotel: 115, mataUang: 'USD' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 50, uangMakan: 60,  uangSaku: 40,  tiket: 0.2, hotel: 160, mataUang: 'USD' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 90, uangMakan: 100, uangSaku: 85,  tiket: 0.2, hotel: 0,   mataUang: 'USD' },
    ],
    'Jepang': [
        { noUrut: 1, jabatan: 'Operator - Karu',    transport: 25, uangMakan: 45,  uangSaku: 15,  tiket: 0.2, hotel: 125, mataUang: 'USD' },
        { noUrut: 2, jabatan: 'Assup - Supervisor',  transport: 25, uangMakan: 50,  uangSaku: 20,  tiket: 0.2, hotel: 150, mataUang: 'USD' },
        { noUrut: 3, jabatan: 'Assman - Manager',    transport: 40, uangMakan: 60,  uangSaku: 30,  tiket: 0.2, hotel: 185, mataUang: 'USD' },
        { noUrut: 4, jabatan: 'GM - Wadir',          transport: 50, uangMakan: 85,  uangSaku: 40,  tiket: 0.2, hotel: 275, mataUang: 'USD' },
        { noUrut: 5, jabatan: 'Direktur - Dirut',    transport: 90, uangMakan: 135, uangSaku: 85,  tiket: 0.2, hotel: 0,   mataUang: 'USD' },
    ],
};

// ============================================================
// REFERENSI DATA FASILITAS — dari PerjalananDinas_Fasilitas.sql
// ============================================================
const FASILITAS_DATA = [
    { kriteria: '1. < 200 KM & < 8 Jam',        hakTransport: true, hakUangMakan: false, hakUangSaku: false, hakTiket: false, hakHotel: false },
    { kriteria: '2. < 200 KM & 8 - 12 Jam',     hakTransport: true, hakUangMakan: true,  hakUangSaku: false, hakTiket: false, hakHotel: false },
    { kriteria: '3. > 200 KM atau > 12 Jam',     hakTransport: true, hakUangMakan: true,  hakUangSaku: true,  hakTiket: false, hakHotel: false },
    { kriteria: '4. Menginap',                   hakTransport: true, hakUangMakan: true,  hakUangSaku: true,  hakTiket: true,  hakHotel: true  },
];

function getTarifByAreaAndLevel(area, noUrut) {
    const areaData = TARIF_DATA[area];
    if (!areaData) return null;
    return areaData.find(t => t.noUrut === noUrut) || null;
}

function jabatanToTarifLevel(jabatan) {
    const j = (jabatan || '').toLowerCase();
    if (j.includes('direktur') || j.includes('dirut') || j.includes('president')) return 5;
    if (j.includes('gm') || j.includes('general manager') || j.includes('wakil direktur') || j.includes('wadir')) return 4;
    if (j.includes('manager') || j.includes('manajer') || j.includes('assman')) return 3;
    if (j.includes('supervisor') || j.includes('assup') || j.includes('kepala')) return 2;
    return 1;
}

// ============================================================
// LOCAL STORAGE & DATA MANAGEMENT
// ============================================================
const STORAGE_KEY_PENGAJUAN = 'sipadin_pengajuan_data_v3';
const STORAGE_KEY_USER = 'sipadin_current_user_v2';

// Data awal — mencerminkan semua 4 tahap pipeline
const initialPengajuanList = [
    // Tahap 1: Pending Approval Atasan
    {
        id: 101, nomorSurat: 'SPPD/IT/2026/08/001',
        nama: 'Rizky Darmawan', nip: 'EMP-2022091', dept: 'IT & Infrastructure', jabatan: 'Senior Web Developer',
        asal: 'Surabaya', tujuan: 'Jakarta Pusat', instansi: 'PT Mitra Solusi Nusantara',
        mulai: daysFromNow(2), selesai: daysFromNow(5), durasi: 4,
        transportasi: 'Pesawat Terbang (Garuda Indonesia)', akomodasi: 'Hotel Santika Premiere Slipi',
        keperluan: 'Rapat / Meeting Klien',
        keterangan: 'Presentasi dan implementasi arsitektur sistem ERP terintegrasi kuartal 3 bersama jajaran direksi mitra.',
        biayaTransport: 2400000, biayaPenginapan: 1800000, uangHarian: 1200000, totalBiaya: 5400000,
        lampiran: 'Surat_Undangan_Mitra_Jakarta.pdf', status: 'pending',
        tanggalPengajuan: daysFromNow(-2), catatanAtasan: '', disetujuiOleh: '', tanggalDisetujui: '',
        catatanHRGA: '', diverifikasiOleh: '', tanggalVerifikasiHRGA: '', tanggalNotifikasiDireksi: '', pesanWA: '',
    },
    {
        id: 102, nomorSurat: 'SPPD/MKT/2026/08/002',
        nama: 'Sari Dewi', nip: 'EMP-2023045', dept: 'Marketing & Sales', jabatan: 'Marketing Specialist',
        asal: 'Surabaya', tujuan: 'Denpasar, Bali', instansi: 'Bali Nusa Dua Convention Center',
        mulai: daysFromNow(1), selesai: daysFromNow(4), durasi: 4,
        transportasi: 'Pesawat Terbang (Batik Air)', akomodasi: 'Grand Hyatt Bali',
        keperluan: 'Survey Lapangan',
        keterangan: 'Survey dan koordinasi booth pameran tahunan industri percetakan dan packaging nasional.',
        biayaTransport: 2100000, biayaPenginapan: 2400000, uangHarian: 1400000, totalBiaya: 5900000,
        lampiran: 'Proposal_Survey_Pameran_Bali.pdf', status: 'pending',
        tanggalPengajuan: daysFromNow(-1), catatanAtasan: '', disetujuiOleh: '', tanggalDisetujui: '',
        catatanHRGA: '', diverifikasiOleh: '', tanggalVerifikasiHRGA: '', tanggalNotifikasiDireksi: '', pesanWA: '',
    },
    // Tahap 2: Approved by Atasan — Pending HRGA
    {
        id: 103, nomorSurat: 'SPPD/FIN/2026/08/003',
        nama: 'Andi Kusuma', nip: 'EMP-2021018', dept: 'Finance & Accounting', jabatan: 'Senior Auditor',
        asal: 'Surabaya', tujuan: 'Medan, Sumatera Utara', instansi: 'Kantor Cabang Adiprima Medan',
        mulai: daysFromNow(7), selesai: daysFromNow(11), durasi: 5,
        transportasi: 'Pesawat Terbang (Citilink)', akomodasi: 'JW Marriott Medan',
        keperluan: 'Audit Internal',
        keterangan: 'Pelaksanaan audit kepatuhan keuangan semester 1 dan rekonsiliasi aset logistik regional barat.',
        biayaTransport: 3200000, biayaPenginapan: 3000000, uangHarian: 1750000, totalBiaya: 7950000,
        lampiran: 'Surat_Tugas_Audit_Medan.pdf', status: 'approved',
        tanggalPengajuan: daysFromNow(-3),
        catatanAtasan: 'Disetujui. Pastikan membawa dokumen rekonsiliasi lengkap dari HO.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)', tanggalDisetujui: daysFromNow(-1),
        catatanHRGA: '', diverifikasiOleh: '', tanggalVerifikasiHRGA: '', tanggalNotifikasiDireksi: '', pesanWA: '',
    },
    {
        id: 104, nomorSurat: 'SPPD/HRD/2026/08/004',
        nama: 'Lina Marlina', nip: 'EMP-2022110', dept: 'Human Resources', jabatan: 'HR Development Lead',
        asal: 'Surabaya', tujuan: 'Yogyakarta', instansi: 'Pusat Diklat SDM Ketenagakerjaan',
        mulai: daysFromNow(5), selesai: daysFromNow(8), durasi: 4,
        transportasi: 'Kereta Api (Executive Argo Wilis)', akomodasi: 'Hotel Tentrem Yogyakarta',
        keperluan: 'Pelatihan / Training',
        keterangan: 'Mengikuti sertifikasi BNSP manajemen talenta dan kepemimpinan operasional industri 4.0.',
        biayaTransport: 900000, biayaPenginapan: 1800000, uangHarian: 1200000, totalBiaya: 3900000,
        lampiran: 'Undangan_Diklat_SDM_Jogja.pdf', status: 'approved',
        tanggalPengajuan: daysFromNow(-5),
        catatanAtasan: 'Disetujui penuh. Mohon buat laporan hasil sertifikasi setelah kembali.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)', tanggalDisetujui: daysFromNow(-3),
        catatanHRGA: '', diverifikasiOleh: '', tanggalVerifikasiHRGA: '', tanggalNotifikasiDireksi: '', pesanWA: '',
    },
    // Tahap 3: Verified by HRGA — Ready for Direksi
    {
        id: 105, nomorSurat: 'SPPD/IT/2026/08/005',
        nama: 'Bima Santoso', nip: 'EMP-2023011', dept: 'IT & Infrastructure', jabatan: 'Network Engineer',
        asal: 'Surabaya', tujuan: 'Bandung', instansi: 'Data Center Lintasarta Bandung',
        mulai: daysFromNow(3), selesai: daysFromNow(6), durasi: 4,
        transportasi: 'Kereta Api (Turangga Executive)', akomodasi: 'Aston Tropicana Bandung',
        keperluan: 'Rapat / Meeting Klien',
        keterangan: 'Migrasi core switch server dan konfigurasi VPN private network antar site pabrik.',
        biayaTransport: 1100000, biayaPenginapan: 1600000, uangHarian: 1200000, totalBiaya: 3900000,
        lampiran: 'Work_Order_Network_Bandung.pdf', status: 'verified_hrga',
        tanggalPengajuan: daysFromNow(-8),
        catatanAtasan: 'Pastikan uptime sistem tetap terjaga selama migrasi berlangsung.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)', tanggalDisetujui: daysFromNow(-6),
        catatanHRGA: 'Dokumen lengkap. Estimasi biaya sesuai tarif standar area Jawa. Disetujui untuk dilaporkan ke Direksi.',
        diverifikasiOleh: 'Dewi Rahayu (HRGA Supervisor)', tanggalVerifikasiHRGA: daysFromNow(-4),
        tanggalNotifikasiDireksi: '', pesanWA: '',
    },
    // Tahap 4: Notified to Direksi — Selesai
    {
        id: 106, nomorSurat: 'SPPD/PRD/2026/08/006',
        nama: 'Rizky Darmawan', nip: 'EMP-2022091', dept: 'IT & Infrastructure', jabatan: 'Senior Web Developer',
        asal: 'Surabaya', tujuan: 'Semarang', instansi: 'PT Percetakan Grafika Utama',
        mulai: daysFromNow(-5), selesai: daysFromNow(-2), durasi: 3,
        transportasi: 'Mobil Dinas Operasional', akomodasi: 'Hotel Gumaya Tower Semarang',
        keperluan: 'Pelatihan / Training',
        keterangan: 'Training operasional integrasi barcode inventory mesin cetak digital.',
        biayaTransport: 650000, biayaPenginapan: 1200000, uangHarian: 900000, totalBiaya: 2750000,
        lampiran: 'SOP_Mesin_Cetak_Semarang.pdf', status: 'notified_direksi',
        tanggalPengajuan: daysFromNow(-14), catatanAtasan: 'Pelatihan disetujui.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)', tanggalDisetujui: daysFromNow(-12),
        catatanHRGA: 'Semua persyaratan terpenuhi. Sesuai kebijakan dinas area Jawa.',
        diverifikasiOleh: 'Dewi Rahayu (HRGA Supervisor)', tanggalVerifikasiHRGA: daysFromNow(-10),
        tanggalNotifikasiDireksi: daysFromNow(-9),
        pesanWA: 'Pesan notifikasi Direksi telah dikirim via WhatsApp.',
    },
    // Rejected — Ditolak Atasan
    {
        id: 107, nomorSurat: 'SPPD/MKT/2026/08/007',
        nama: 'Dian Pertiwi', nip: 'EMP-2024003', dept: 'Marketing & Sales', jabatan: 'Brand Associate',
        asal: 'Surabaya', tujuan: 'Makassar', instansi: 'Distributor Nusantara Timur',
        mulai: daysFromNow(10), selesai: daysFromNow(13), durasi: 4,
        transportasi: 'Pesawat Terbang (Lion Air)', akomodasi: 'Hotel Four Points Makassar',
        keperluan: 'Kunjungan Rekanan',
        keterangan: 'Kunjungan rutin ke distributor regional timur untuk evaluasi pencapaian target Q3.',
        biayaTransport: 2800000, biayaPenginapan: 2100000, uangHarian: 1400000, totalBiaya: 6300000,
        lampiran: 'Jadwal_Kunjungan_Distributor_Makassar.pdf', status: 'rejected',
        tanggalPengajuan: daysFromNow(-4),
        catatanAtasan: 'Ditolak. Kunjungan dapat dilakukan secara virtual. Jadwalkan ulang jika benar-benar diperlukan kunjungan fisik.',
        disetujuiOleh: '', tanggalDisetujui: '',
        catatanHRGA: '', diverifikasiOleh: '', tanggalVerifikasiHRGA: '', tanggalNotifikasiDireksi: '', pesanWA: '',
    },
];

function getPengajuanData() {
    const raw = localStorage.getItem(STORAGE_KEY_PENGAJUAN);
    if (!raw) {
        savePengajuanData(initialPengajuanList);
        return initialPengajuanList;
    }
    try {
        return JSON.parse(raw);
    } catch (e) {
        savePengajuanData(initialPengajuanList);
        return initialPengajuanList;
    }
}

function savePengajuanData(data) {
    localStorage.setItem(STORAGE_KEY_PENGAJUAN, JSON.stringify(data));
}

// ============================================================
// USER AUTHENTICATION STATE
// ============================================================
const DEMO_ACCOUNTS = {
    karyawan: {
        name: 'Rizky Darmawan', role: 'karyawan', dept: 'IT & Infrastructure',
        jabatan: 'Senior Web Developer', nip: 'EMP-2022091',
        email: 'karyawan@sipadin.com', avatar: 'RD', color: 'brand',
    },
    atasan: {
        name: 'Budi Pratama', role: 'atasan', dept: 'IT & Operations Manager',
        jabatan: 'IT & Ops Manager', nip: 'MGR-2018004',
        email: 'manager@sipadin.com', avatar: 'BP', color: 'purple',
    },
    hrga: {
        name: 'Dewi Rahayu', role: 'hrga', dept: 'Human Resources & General Affairs',
        jabatan: 'HRGA Supervisor', nip: 'HRD-2019007',
        email: 'hrga@sipadin.com', avatar: 'DR', color: 'amber',
    },
};

function getCurrentUser() {
    const raw = localStorage.getItem(STORAGE_KEY_USER);
    if (!raw) {
        setCurrentUser(DEMO_ACCOUNTS.karyawan);
        return DEMO_ACCOUNTS.karyawan;
    }
    try {
        return JSON.parse(raw);
    } catch (e) {
        return DEMO_ACCOUNTS.karyawan;
    }
}

function setCurrentUser(user) {
    localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(user));
}

function logout() {
    localStorage.removeItem(STORAGE_KEY_USER);
    showToast('info', 'Sampai Jumpa', 'Anda telah berhasil keluar dari sistem.');
    setTimeout(() => { window.location.href = 'index.html'; }, 400);
}

/**
 * Inisialisasi data sidebar user secara dinamis via data-attribute
 */
function initSidebarUser() {
    const user = getCurrentUser();
    document.querySelectorAll('[data-user-avatar]').forEach(el => el.textContent = user.avatar || '?');
    document.querySelectorAll('[data-user-name]').forEach(el => el.textContent = user.name || 'Pengguna');
    document.querySelectorAll('[data-user-dept]').forEach(el => el.textContent = user.dept || '-');
    document.querySelectorAll('[data-user-role]').forEach(el => el.textContent = user.jabatan || user.role || '-');
    document.querySelectorAll('[data-user-nip]').forEach(el => el.textContent = user.nip || '-');
}

// ============================================================
// TOAST NOTIFICATION SYSTEM
// ============================================================
const toastConfigs = {
    success: { icon: 'ph-check-circle', color: 'text-emerald-400', bg: 'bg-surface-card border-emerald-500/40 text-emerald-100' },
    error:   { icon: 'ph-x-circle',     color: 'text-rose-400',    bg: 'bg-surface-card border-rose-500/40 text-rose-100'    },
    warning: { icon: 'ph-warning',      color: 'text-amber-400',   bg: 'bg-surface-card border-amber-500/40 text-amber-100'  },
    info:    { icon: 'ph-info',         color: 'text-brand-400',   bg: 'bg-surface-card border-brand-500/40 text-brand-100'  },
};

function showToast(type, title, message) {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        container.className = 'fixed top-5 right-5 z-[99999] space-y-3 w-80 max-w-[calc(100vw-2.5rem)] pointer-events-none';
        document.body.appendChild(container);
    }

    const cfg = toastConfigs[type] || toastConfigs.info;
    const toast = document.createElement('div');
    toast.className = `pointer-events-auto toast-enter flex items-start gap-3 p-4 rounded-2xl border ${cfg.bg} shadow-2xl backdrop-blur-xl transition-all duration-300`;
    toast.innerHTML = `
        <div class="p-1 rounded-xl bg-white/5 flex-shrink-0">
            <i class="ph ${cfg.icon} ${cfg.color} text-2xl"></i>
        </div>
        <div class="flex-1 min-w-0 pt-0.5">
            <h4 class="text-sm font-bold text-white tracking-tight">${title}</h4>
            <p class="text-xs text-slate-300 mt-1 leading-relaxed">${message}</p>
        </div>
        <button onclick="this.closest('div[class*=toast]').remove()" class="text-slate-400 hover:text-white transition-colors p-1 flex-shrink-0">
            <i class="ph ph-x text-sm"></i>
        </button>
    `;
    container.appendChild(toast);
    setTimeout(() => {
        toast.classList.remove('toast-enter');
        toast.classList.add('toast-exit');
        setTimeout(() => toast.remove(), 350);
    }, 4500);
}

// ============================================================
// MOBILE SIDEBAR CONTROLS
// ============================================================
function openMobileSidebar() {
    const sidebar = document.getElementById('mobile-sidebar') || document.getElementById('atasan-mobile-sidebar') || document.getElementById('hrga-mobile-sidebar');
    const overlay = document.getElementById('sidebar-overlay') || document.getElementById('atasan-sidebar-overlay') || document.getElementById('hrga-sidebar-overlay');
    if (sidebar) sidebar.classList.add('open');
    if (overlay) overlay.classList.remove('hidden');
}

function closeMobileSidebar() {
    const sidebar = document.getElementById('mobile-sidebar') || document.getElementById('atasan-mobile-sidebar') || document.getElementById('hrga-mobile-sidebar');
    const overlay = document.getElementById('sidebar-overlay') || document.getElementById('atasan-sidebar-overlay') || document.getElementById('hrga-sidebar-overlay');
    if (sidebar) sidebar.classList.remove('open');
    if (overlay) overlay.classList.add('hidden');
}

// ============================================================
// WORKFLOW ACTIONS — TAHAP 2: ATASAN
// ============================================================
function approveByAtasan(id, catatan) {
    const user = getCurrentUser();
    const list = getPengajuanData();
    const idx = list.findIndex(p => p.id === Number(id));
    if (idx === -1) { showToast('error', 'Error', 'Data tidak ditemukan.'); return false; }
    list[idx].status = 'approved';
    list[idx].catatanAtasan = catatan || 'Disetujui.';
    list[idx].disetujuiOleh = `${user.name} (${user.jabatan || user.dept})`;
    list[idx].tanggalDisetujui = new Date().toISOString().split('T')[0];
    savePengajuanData(list);
    showToast('success', 'Pengajuan Disetujui', `SPPD ${list[idx].nomorSurat} telah disetujui & diteruskan ke HRGA.`);
    return true;
}

function rejectByAtasan(id, catatan) {
    const user = getCurrentUser();
    const list = getPengajuanData();
    const idx = list.findIndex(p => p.id === Number(id));
    if (idx === -1) { showToast('error', 'Error', 'Data tidak ditemukan.'); return false; }
    list[idx].status = 'rejected';
    list[idx].catatanAtasan = catatan || 'Tidak disetujui.';
    list[idx].tanggalDisetujui = new Date().toISOString().split('T')[0];
    savePengajuanData(list);
    showToast('warning', 'Pengajuan Ditolak', `SPPD ${list[idx].nomorSurat} telah ditolak.`);
    return true;
}

// ============================================================
// WORKFLOW ACTIONS — TAHAP 3: HRGA
// ============================================================
function verifyByHRGA(id, catatan) {
    const user = getCurrentUser();
    const list = getPengajuanData();
    const idx = list.findIndex(p => p.id === Number(id));
    if (idx === -1) { showToast('error', 'Error', 'Data tidak ditemukan.'); return false; }
    if (list[idx].status !== 'approved') { showToast('warning', 'Status Tidak Valid', 'Hanya pengajuan berstatus Approved yang dapat diverifikasi HRGA.'); return false; }
    list[idx].status = 'verified_hrga';
    list[idx].catatanHRGA = catatan || 'Dokumen dan anggaran telah diverifikasi sesuai kebijakan.';
    list[idx].diverifikasiOleh = `${user.name} (${user.jabatan || user.dept})`;
    list[idx].tanggalVerifikasiHRGA = new Date().toISOString().split('T')[0];
    savePengajuanData(list);
    showToast('success', 'Verifikasi HRGA Berhasil', `SPPD ${list[idx].nomorSurat} siap diinformasikan ke Direksi.`);
    return true;
}

function rejectByHRGA(id, catatan) {
    const user = getCurrentUser();
    const list = getPengajuanData();
    const idx = list.findIndex(p => p.id === Number(id));
    if (idx === -1) { showToast('error', 'Error', 'Data tidak ditemukan.'); return false; }
    list[idx].status = 'rejected_hrga';
    list[idx].catatanHRGA = catatan || 'Tidak memenuhi kebijakan perjalanan dinas.';
    list[idx].diverifikasiOleh = `${user.name} (${user.jabatan || user.dept})`;
    list[idx].tanggalVerifikasiHRGA = new Date().toISOString().split('T')[0];
    savePengajuanData(list);
    showToast('warning', 'Pengajuan Ditolak HRGA', `SPPD ${list[idx].nomorSurat} ditolak oleh HRGA.`);
    return true;
}

// ============================================================
// WORKFLOW ACTIONS — TAHAP 4: NOTIFIKASI DIREKSI
// ============================================================
function generateWAMessage(item) {
    const isMenginap = item.statusMenginap || ((item.akomodasi && !item.akomodasi.toLowerCase().includes('tidak')) || (item.durasi > 1) ? 'Ya' : 'Tidak');
    const jarak = item.jarak || (item.durasi * 150);
    const lamaJam = item.lamaJam || ((item.durasi || 1) * 24);
    
    // Determinasi Kriteria Fasilitas
    let kriteria = item.kriteria || '1. < 200 KM & < 8 Jam';
    if (!item.kriteria) {
        if (isMenginap === 'Ya') {
            kriteria = '4. Menginap';
        } else if (jarak > 200 || lamaJam > 12) {
            kriteria = '3. > 200 KM atau > 12 Jam';
        } else if (jarak < 200 && lamaJam >= 8 && lamaJam <= 12) {
            kriteria = '2. < 200 KM & 8 - 12 Jam';
        }
    }
    
    // Fasilitas aktif berdasarkan kriteria & komponen
    const fasilitasList = [];
    if (item.biayaTransport > 0 || item.transportasi) fasilitasList.push('Transport');
    if (item.uangHarian > 0) fasilitasList.push('Makan', 'Saku');
    if (isMenginap === 'Ya' || item.biayaPenginapan > 0) fasilitasList.push('Tiket', 'Hotel');
    
    const fasilitasStr = fasilitasList.length ? fasilitasList.join(', ') : 'Transport, Makan, Saku';

    return `📲 *PERMOHONAN APPROVAL PERJALANAN DINAS*
--------------------------------------------------
*Karyawan:* ${item.nama} (${item.jabatan || 'Staf'} - ${item.dept})
*Tujuan:* ${item.tujuan} (${jarak} KM / ${lamaJam} Jam)
*Kriteria:* ${kriteria}
*Status Menginap:* ${isMenginap}
*Rincian Fasilitas:* ${fasilitasStr}
*Total Estimasi CA:* *${formatRupiah(item.totalBiaya)}*

Mohon konfirmasi balasan: *OK / Disetujui* atau *Ditolak*.
--------------------------------------------------`;
}

function notifyDireksi(id) {
    const list = getPengajuanData();
    const idx = list.findIndex(p => p.id === Number(id));
    if (idx === -1) { showToast('error', 'Error', 'Data tidak ditemukan.'); return null; }
    if (list[idx].status !== 'verified_hrga') { showToast('warning', 'Status Tidak Valid', 'Hanya pengajuan berstatus Terverifikasi HRGA yang dapat dikirim ke Direksi.'); return null; }

    const pesan = generateWAMessage(list[idx]);
    list[idx].status = 'notified_direksi';
    list[idx].tanggalNotifikasiDireksi = new Date().toISOString().split('T')[0];
    list[idx].pesanWA = pesan;
    savePengajuanData(list);
    showToast('success', 'Notifikasi Dikirim', `Pesan WhatsApp untuk SPPD ${list[idx].nomorSurat} telah disiapkan.`);
    return { waUrl: `https://wa.me/?text=${encodeURIComponent(pesan)}`, pesan, item: list[idx] };
}

// ============================================================
// MODAL DETAIL & SPPD PRINT GENERATOR
// ============================================================
function openDetailModal(id) {
    const list = getPengajuanData();
    const item = list.find(p => p.id === Number(id));
    if (!item) return;

    let modal = document.getElementById('modal-detail-pengajuan');
    if (!modal) {
        modal = document.createElement('div');
        modal.id = 'modal-detail-pengajuan';
        modal.className = 'fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md';
        document.body.appendChild(modal);
    }

    const statusBadge = getStatusBadge(item.status);

    // Timeline helper
    function tlStep(done, rejected, iconDone, iconPending, label, sub, note) {
        const c = rejected ? 'rose' : (done ? 'emerald' : 'slate');
        const ic = rejected ? 'ph-x-circle' : (done ? iconDone : iconPending);
        return `<div class="flex gap-3">
            <div class="flex flex-col items-center">
                <div class="w-8 h-8 rounded-full flex items-center justify-center bg-${c}-500/20 border border-${c}-500/40 text-${c}-400 flex-shrink-0">
                    <i class="ph-fill ${ic} text-sm"></i>
                </div>
            </div>
            <div class="pb-4 flex-1">
                <p class="text-xs font-bold text-white">${label}</p>
                ${sub ? `<p class="text-[11px] text-slate-400 mt-0.5">${sub}</p>` : `<p class="text-[11px] text-slate-500 mt-0.5">Menunggu...</p>`}
                ${note ? `<p class="text-[11px] text-slate-300 mt-1 italic">"${note}"</p>` : ''}
            </div>
        </div>`;
    }

    const s = item.status;
    const timelineHtml = [
        tlStep(true, false, 'ph-check-circle', 'ph-clock', 'Pengajuan Karyawan',
            `${formatDate(item.tanggalPengajuan)} · ${item.nama}`, ''),
        tlStep(['approved','rejected','verified_hrga','rejected_hrga','notified_direksi'].includes(s), s==='rejected',
            'ph-check-circle', 'ph-clock', 'Approval Atasan / Manager',
            item.disetujuiOleh ? `${formatDate(item.tanggalDisetujui)} · ${item.disetujuiOleh}` : '',
            item.catatanAtasan),
        tlStep(['verified_hrga','rejected_hrga','notified_direksi'].includes(s), s==='rejected_hrga',
            'ph-shield-check', 'ph-clock', 'Kontrol HRGA/HRD',
            item.diverifikasiOleh ? `${formatDate(item.tanggalVerifikasiHRGA)} · ${item.diverifikasiOleh}` : '',
            item.catatanHRGA),
        tlStep(s==='notified_direksi', false, 'ph-paper-plane-tilt', 'ph-clock', 'Notifikasi Direksi',
            s==='notified_direksi' ? `${formatDate(item.tanggalNotifikasiDireksi)} · Pesan WhatsApp terkirim` : '', ''),
    ].join('');

    modal.innerHTML = `
        <div class="w-full max-w-2xl bg-surface-card border border-surface-border rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[92vh]">
            <div class="p-6 border-b border-surface-border flex items-center justify-between bg-surface/60 flex-shrink-0">
                <div class="flex items-center gap-3">
                    <div class="w-11 h-11 rounded-2xl bg-brand-gradient flex items-center justify-center shadow-glow">
                        <i class="ph ph-file-text text-xl text-white"></i>
                    </div>
                    <div>
                        <div class="flex items-center gap-2 flex-wrap">
                            <h3 class="text-lg font-bold text-white">${item.nomorSurat || 'SPPD-' + item.id}</h3>
                            ${statusBadge}
                        </div>
                        <p class="text-xs text-slate-400 mt-0.5">Diajukan: ${formatDateFull(item.tanggalPengajuan || item.mulai)}</p>
                    </div>
                </div>
                <button onclick="closeDetailModal()" class="w-9 h-9 rounded-xl bg-surface-border/50 text-slate-400 hover:text-white hover:bg-surface-border flex items-center justify-center transition-colors flex-shrink-0">
                    <i class="ph ph-x text-lg"></i>
                </button>
            </div>
            <div class="p-6 overflow-y-auto space-y-5 flex-1 text-sm">
                <div class="bg-surface/50 border border-surface-border rounded-2xl p-4 flex items-center justify-between flex-wrap gap-3">
                    <div class="flex items-center gap-3.5">
                        <div class="w-12 h-12 rounded-2xl bg-brand-500/20 border border-brand-500/30 flex items-center justify-center text-brand-300 font-bold text-base">
                            ${item.nama.split(' ').map(n=>n[0]).join('').slice(0,2)}
                        </div>
                        <div>
                            <p class="font-bold text-white text-base">${item.nama}</p>
                            <p class="text-xs text-slate-400 font-mono">${item.nip || '-'} · ${item.dept}</p>
                        </div>
                    </div>
                    <div class="text-right">
                        <span class="text-xs text-slate-400 block">Jabatan</span>
                        <span class="text-xs font-semibold text-slate-200">${item.jabatan || '-'}</span>
                    </div>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                        <p class="flex items-center gap-2 text-brand-400 font-semibold text-xs uppercase tracking-wider"><i class="ph ph-map-pin-line"></i> Rute & Tujuan</p>
                        <p class="text-base font-bold text-white flex items-center gap-2">${item.asal || 'Surabaya'} <i class="ph ph-arrow-right text-brand-400 text-sm"></i> <span class="text-brand-300">${item.tujuan}</span></p>
                        <div><p class="text-xs text-slate-400">Instansi</p><p class="text-sm font-semibold text-slate-200 mt-0.5">${item.instansi || '-'}</p></div>
                    </div>
                    <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                        <p class="flex items-center gap-2 text-brand-400 font-semibold text-xs uppercase tracking-wider"><i class="ph ph-calendar-blank"></i> Jadwal</p>
                        <div class="flex justify-between">
                            <div><p class="text-xs text-slate-400">Tanggal</p><p class="text-sm font-bold text-white mt-0.5">${formatDate(item.mulai)} — ${formatDate(item.selesai)}</p></div>
                            <div class="text-right"><p class="text-xs text-slate-400">Durasi</p><p class="text-sm font-bold text-brand-400 mt-0.5">${item.durasi || calculateDuration(item.mulai, item.selesai)} Hari</p></div>
                        </div>
                        <div><p class="text-xs text-slate-400">Keperluan</p><p class="text-sm font-semibold text-white mt-0.5">${item.keperluan}</p></div>
                    </div>
                </div>
                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div><p class="text-xs text-slate-400 flex items-center gap-1.5"><i class="ph ph-airplane text-brand-400"></i> Transportasi</p><p class="text-sm font-semibold text-slate-200 mt-1">${item.transportasi || '-'}</p></div>
                    <div><p class="text-xs text-slate-400 flex items-center gap-1.5"><i class="ph ph-buildings text-brand-400"></i> Akomodasi</p><p class="text-sm font-semibold text-slate-200 mt-1">${item.akomodasi || '-'}</p></div>
                </div>
                <div><p class="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Keterangan</p>
                    <div class="p-4 rounded-2xl bg-surface/50 border border-surface-border text-slate-300 text-xs leading-relaxed">${item.keterangan || '-'}</div>
                </div>
                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                    <p class="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center justify-between">
                        <span>Estimasi Biaya</span><span class="text-emerald-400 font-mono font-bold">${formatRupiah(item.totalBiaya || 0)}</span>
                    </p>
                    <div class="space-y-1.5 text-xs text-slate-400 border-t border-surface-border/50 pt-2.5">
                        <div class="flex justify-between"><span>Transportasi</span><span class="text-slate-200 font-mono">${formatRupiah(item.biayaTransport || 0)}</span></div>
                        <div class="flex justify-between"><span>Penginapan</span><span class="text-slate-200 font-mono">${formatRupiah(item.biayaPenginapan || 0)}</span></div>
                        <div class="flex justify-between"><span>Uang Harian</span><span class="text-slate-200 font-mono">${formatRupiah(item.uangHarian || 0)}</span></div>
                    </div>
                </div>
                <div><p class="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2 flex items-center gap-2"><i class="ph ph-git-branch text-brand-400"></i> Riwayat Approval</p>
                    <div class="p-4 rounded-2xl bg-surface/50 border border-surface-border">${timelineHtml}</div>
                </div>
                ${item.pesanWA ? `
                    <div class="p-4 rounded-2xl bg-teal-500/10 border border-teal-500/30">
                        <p class="text-xs font-bold text-teal-300 mb-2 flex items-center gap-1.5"><i class="ph ph-chat-circle-text"></i> Pesan Notifikasi Direksi:</p>
                        <pre class="text-[11px] text-slate-300 whitespace-pre-wrap leading-relaxed font-sans">${item.pesanWA}</pre>
                    </div>` : ''}
            </div>
            <div class="p-5 border-t border-surface-border bg-surface/80 flex items-center justify-between gap-3 flex-shrink-0">
                <button onclick="closeDetailModal()" class="px-5 py-2.5 rounded-xl border border-surface-border text-slate-300 hover:text-white hover:bg-white/5 text-sm font-semibold transition-all">Tutup</button>
                <div class="flex gap-2.5">
                    ${['approved','verified_hrga','notified_direksi'].includes(item.status) ? `
                        <button onclick="printSPPDSlip(${item.id})" class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-brand-gradient hover:opacity-90 active:scale-95 text-white text-sm font-semibold transition-all shadow-glow">
                            <i class="ph ph-printer text-base"></i> Cetak SPPD
                        </button>` : ''}
                </div>
            </div>
        </div>
    `;
    modal.classList.remove('hidden');
}

function closeDetailModal() {
    const modal = document.getElementById('modal-detail-pengajuan');
    if (modal) modal.remove();
}

// ============================================================
// SPPD DOCUMENT PRINT VIEW GENERATOR
// ============================================================
function printSPPDSlip(id) {
    const list = getPengajuanData();
    const item = list.find(p => p.id === Number(id));
    if (!item) return;

    let printArea = document.getElementById('sppd-print-area');
    if (!printArea) {
        printArea = document.createElement('div');
        printArea.id = 'sppd-print-area';
        document.body.appendChild(printArea);
    }

    const mgrName = item.disetujuiOleh ? item.disetujuiOleh.split('(')[0].trim() : 'Budi Pratama';
    const hrgaName = item.diverifikasiOleh ? item.diverifikasiOleh.split('(')[0].trim() : 'Dewi Rahayu';

    printArea.innerHTML = `
        <div style="font-family: 'Times New Roman', serif; max-width: 800px; margin: 0 auto; color: #111; padding: 20px; line-height: 1.6;">
            <div style="display:flex;align-items:center;border-bottom:3px double #000;padding-bottom:12px;margin-bottom:20px;">
                <div style="flex:1;text-align:center;">
                    <h2 style="margin:0;font-size:18pt;text-transform:uppercase;font-weight:bold;letter-spacing:1px;">PT ADIPRIMA SURAPRINTA</h2>
                    <p style="margin:3px 0;font-size:10pt;">Kawasan Industri Driyorejo, Gresik, Jawa Timur — Indonesia</p>
                    <p style="margin:0;font-size:9pt;color:#444;">Telp: (031) 7507888 | Email: corporate@adiprima.co.id | Website: www.adiprima.co.id</p>
                </div>
            </div>
            <div style="text-align:center;margin-bottom:25px;">
                <h3 style="margin:0;font-size:14pt;text-decoration:underline;font-weight:bold;text-transform:uppercase;">SURAT PERINTAH PERJALANAN DINAS (SPPD)</h3>
                <p style="margin:4px 0 0 0;font-size:11pt;">Nomor: <strong>${item.nomorSurat || 'SPPD/2026/08/' + item.id}</strong></p>
            </div>
            <table style="width:100%;border-collapse:collapse;margin-bottom:25px;font-size:11pt;" border="1" cellpadding="8">
                <tr><td style="width:5%;text-align:center;font-weight:bold;">1</td><td style="width:35%;">Pejabat Berwenang yang Memberi Perintah</td><td><strong>${item.disetujuiOleh || '-'}</strong></td></tr>
                <tr><td style="text-align:center;font-weight:bold;">2</td><td>Nama / NIP Pegawai yang Ditugaskan</td><td><strong>${item.nama}</strong> / NIP: ${item.nip || '-'}</td></tr>
                <tr><td style="text-align:center;font-weight:bold;">3</td><td>Pangkat / Jabatan / Departemen</td><td>${item.jabatan || '-'} / ${item.dept}</td></tr>
                <tr><td style="text-align:center;font-weight:bold;">4</td><td>Maksud Perjalanan Dinas</td><td>${item.keperluan} — <em>${item.keterangan}</em></td></tr>
                <tr><td style="text-align:center;font-weight:bold;">5</td><td>Alat Angkut / Transportasi</td><td>${item.transportasi || '-'}</td></tr>
                <tr><td style="text-align:center;font-weight:bold;">6</td><td>a. Tempat Berangkat<br>b. Tempat Tujuan / Instansi</td><td>a. ${item.asal || 'Surabaya (Head Office)'}<br>b. ${item.tujuan} (${item.instansi || '-'})</td></tr>
                <tr><td style="text-align:center;font-weight:bold;">7</td><td>a. Lamanya Perjalanan Dinas<br>b. Tanggal Berangkat<br>c. Tanggal Harus Kembali</td><td>a. ${item.durasi || calculateDuration(item.mulai, item.selesai)} Hari<br>b. ${formatDateFull(item.mulai)}<br>c. ${formatDateFull(item.selesai)}</td></tr>
                <tr><td style="text-align:center;font-weight:bold;">8</td><td>Rincian Estimasi Anggaran</td><td>- Transportasi: ${formatRupiah(item.biayaTransport || 0)}<br>- Penginapan: ${formatRupiah(item.biayaPenginapan || 0)}<br>- Uang Harian: ${formatRupiah(item.uangHarian || 0)}<br><strong>Total: ${formatRupiah(item.totalBiaya || 0)}</strong></td></tr>
                <tr><td style="text-align:center;font-weight:bold;">9</td><td>Diverifikasi HRGA</td><td>${item.diverifikasiOleh || '-'} (${formatDate(item.tanggalVerifikasiHRGA) || '-'})<br><em>${item.catatanHRGA || '-'}</em></td></tr>
            </table>
            <div style="display:flex;justify-content:space-between;margin-top:40px;font-size:11pt;">
                <div style="text-align:center;width:30%;"><p style="margin-bottom:60px;">Pegawai yang Melaksanakan Tugas,</p><p style="margin:0;font-weight:bold;text-decoration:underline;">${item.nama}</p><p style="margin:0;font-size:9pt;">NIP: ${item.nip || '-'}</p></div>
                <div style="text-align:center;width:30%;"><p style="margin:0;">Dikeluarkan di: Surabaya</p><p style="margin:0 0 60px 0;">Pada tanggal: ${formatDate(item.tanggalDisetujui || new Date().toISOString().split('T')[0])}</p><p style="margin:0;font-weight:bold;text-decoration:underline;">${mgrName}</p><p style="margin:0;font-size:9pt;">Manager / Pejabat Berwenang</p></div>
                <div style="text-align:center;width:30%;"><p style="margin:0;">Mengetahui,</p><p style="margin:0 0 60px 0;">HRD/HRGA PT Adiprima Suraprinta</p><p style="margin:0;font-weight:bold;text-decoration:underline;">${hrgaName}</p><p style="margin:0;font-size:9pt;">HRGA Supervisor</p></div>
            </div>
        </div>
    `;
    window.print();
}

// ============================================================
// AUTO INIT ON LOAD
// ============================================================
document.addEventListener('DOMContentLoaded', () => {
    // Highlight active sidebar link
    const currentPath = window.location.pathname.split('/').pop() || 'index.html';
    document.querySelectorAll('.sidebar-link').forEach(link => {
        const href = link.getAttribute('href');
        if (href && (href === currentPath || (currentPath === '' && href === 'index.html'))) {
            link.classList.add('active');
        } else {
            link.classList.remove('active');
        }
    });
    // Inisialisasi data user dinamis di sidebar
    initSidebarUser();
});
