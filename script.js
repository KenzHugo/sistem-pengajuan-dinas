// ============================================================
// SIPADIN — SISTEM PENGAJUAN DINAS
// GLOBAL JAVASCRIPT & STATE MANAGEMENT
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
                    card: '0 10px 30px -10px rgba(0, 0, 0, 0.6)',
                },
                backgroundImage: {
                    'hero-gradient': 'radial-gradient(ellipse at top, #1e1b4b 0%, #0f172a 45%, #0b0f19 100%)',
                    'brand-gradient': 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
                    'accent-gradient': 'linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%)',
                    'emerald-gradient': 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                    'rose-gradient': 'linear-gradient(135deg, #f43f5e 0%, #e11d48 100%)',
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

// ============================================================
// LOCAL STORAGE & DATA MANAGEMENT
// ============================================================
const STORAGE_KEY_PENGAJUAN = 'sipadin_pengajuan_data_v2';
const STORAGE_KEY_USER = 'sipadin_current_user_v2';

const initialPengajuanList = [
    {
        id: 101,
        nomorSurat: 'SPPD/IT/2026/08/001',
        nama: 'Rizky Darmawan',
        nip: 'EMP-2022091',
        dept: 'IT & Infrastructure',
        jabatan: 'Senior Web Developer',
        asal: 'Surabaya',
        tujuan: 'Jakarta Pusat',
        instansi: 'PT Mitra Solusi Nusantara',
        mulai: daysFromNow(2),
        selesai: daysFromNow(5),
        durasi: 4,
        transportasi: 'Pesawat Terbang (Garuda Indonesia)',
        akomodasi: 'Hotel Santika Premiere Slipi',
        keperluan: 'Rapat / Meeting Klien',
        keterangan: 'Presentasi dan implementasi arsitektur sistem ERP terintegrasi kuartal 3 bersama jajaran direksi mitra.',
        biayaTransport: 2400000,
        biayaPenginapan: 1800000,
        uangHarian: 1200000,
        totalBiaya: 5400000,
        lampiran: 'Surat_Undangan_Mitra_Jakarta.pdf',
        status: 'pending',
        tanggalPengajuan: daysFromNow(-2),
        catatanAtasan: '',
        disetujuiOleh: '',
        tanggalDisetujui: ''
    },
    {
        id: 102,
        nomorSurat: 'SPPD/MKT/2026/08/002',
        nama: 'Sari Dewi',
        nip: 'EMP-2023045',
        dept: 'Marketing & Sales',
        jabatan: 'Marketing Specialist',
        asal: 'Surabaya',
        tujuan: 'Denpasar, Bali',
        instansi: 'Bali Nusa Dua Convention Center',
        mulai: daysFromNow(1),
        selesai: daysFromNow(4),
        durasi: 4,
        transportasi: 'Pesawat Terbang (Batik Air)',
        akomodasi: 'Grand Hyatt Bali',
        keperluan: 'Survey Lapangan',
        keterangan: 'Survey dan koordinasi booth pameran tahunan industri percetakan dan packaging nasional.',
        biayaTransport: 2100000,
        biayaPenginapan: 2400000,
        uangHarian: 1400000,
        totalBiaya: 5900000,
        lampiran: 'Proposal_Survey_Pameran_Bali.pdf',
        status: 'pending',
        tanggalPengajuan: daysFromNow(-1),
        catatanAtasan: '',
        disetujuiOleh: '',
        tanggalDisetujui: ''
    },
    {
        id: 103,
        nomorSurat: 'SPPD/FIN/2026/08/003',
        nama: 'Andi Kusuma',
        nip: 'EMP-2021018',
        dept: 'Finance & Accounting',
        jabatan: 'Senior Auditor',
        asal: 'Surabaya',
        tujuan: 'Medan, Sumatera Utara',
        instansi: 'Kantor Cabang Adiprima Medan',
        mulai: daysFromNow(7),
        selesai: daysFromNow(11),
        durasi: 5,
        transportasi: 'Pesawat Terbang (Citilink)',
        akomodasi: 'JW Marriott Medan',
        keperluan: 'Audit Internal',
        keterangan: 'Pelaksanaan audit kepatuhan keuangan semester 1 dan rekonsiliasi aset logistik regional barat.',
        biayaTransport: 3200000,
        biayaPenginapan: 3000000,
        uangHarian: 1750000,
        totalBiaya: 7950000,
        lampiran: 'Surat_Tugas_Audit_Medan.pdf',
        status: 'pending',
        tanggalPengajuan: daysFromNow(-3),
        catatanAtasan: '',
        disetujuiOleh: '',
        tanggalDisetujui: ''
    },
    {
        id: 104,
        nomorSurat: 'SPPD/HRD/2026/08/004',
        nama: 'Lina Marlina',
        nip: 'EMP-2022110',
        dept: 'Human Resources',
        jabatan: 'HR Development Lead',
        asal: 'Surabaya',
        tujuan: 'Yogyakarta',
        instansi: 'Pusat Diklat SDM Ketenagakerjaan',
        mulai: daysFromNow(-6),
        selesai: daysFromNow(-3),
        durasi: 4,
        transportasi: 'Kereta Api (Executive Argo Wilis)',
        akomodasi: 'Hotel Tentrem Yogyakarta',
        keperluan: 'Pelatihan / Training',
        keterangan: 'Mengikuti sertifikasi BNSP manajemen talenta dan kepemimpinan operasional industri 4.0.',
        biayaTransport: 900000,
        biayaPenginapan: 1800000,
        uangHarian: 1200000,
        totalBiaya: 3900000,
        lampiran: 'Undangan_Diklat_SDM_Jogja.pdf',
        status: 'approved',
        tanggalPengajuan: daysFromNow(-10),
        catatanAtasan: 'Disetujui penuh. Mohon buat laporan hasil sertifikasi setelah kembali.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)',
        tanggalDisetujui: daysFromNow(-8)
    },
    {
        id: 105,
        nomorSurat: 'SPPD/IT/2026/08/005',
        nama: 'Bima Santoso',
        nip: 'EMP-2023011',
        dept: 'IT & Infrastructure',
        jabatan: 'Network Engineer',
        asal: 'Surabaya',
        tujuan: 'Bandung',
        instansi: 'Data Center Lintasarta Bandung',
        mulai: daysFromNow(-12),
        selesai: daysFromNow(-9),
        durasi: 4,
        transportasi: 'Kereta Api (Turangga Executive)',
        akomodasi: 'Aston Tropicana Bandung',
        keperluan: 'Rapat / Meeting Klien',
        keterangan: 'Migrasi core switch server dan konfigurasi VPN private network antar site pabrik.',
        biayaTransport: 1100000,
        biayaPenginapan: 1600000,
        uangHarian: 1200000,
        totalBiaya: 3900000,
        lampiran: 'Work_Order_Network_Bandung.pdf',
        status: 'approved',
        tanggalPengajuan: daysFromNow(-16),
        catatanAtasan: 'Pastikan uptime sistem tetap terjaga selama migrasi berlangsung.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)',
        tanggalDisetujui: daysFromNow(-14)
    },
    {
        id: 106,
        nomorSurat: 'SPPD/PRD/2026/08/006',
        nama: 'Rizky Darmawan',
        nip: 'EMP-2022091',
        dept: 'IT & Infrastructure',
        jabatan: 'Senior Web Developer',
        asal: 'Surabaya',
        tujuan: 'Semarang',
        instansi: 'PT Percetakan Grafika Utama',
        mulai: daysFromNow(-18),
        selesai: daysFromNow(-16),
        durasi: 3,
        transportasi: 'Mobil Dinas Operasional',
        akomodasi: 'Hotel Gumaya Tower Semarang',
        keperluan: 'Pelatihan / Training',
        keterangan: 'Training operasional integrasi barcode inventory mesin cetak digital.',
        biayaTransport: 650000,
        biayaPenginapan: 1200000,
        uangHarian: 900000,
        totalBiaya: 2750000,
        lampiran: 'SOP_Mesin_Cetak_Semarang.pdf',
        status: 'approved',
        tanggalPengajuan: daysFromNow(-22),
        catatanAtasan: 'Pelatihan disetujui.',
        disetujuiOleh: 'Budi Pratama (IT & Ops Manager)',
        tanggalDisetujui: daysFromNow(-20)
    }
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

// User Authentication State
function getCurrentUser() {
    const raw = localStorage.getItem(STORAGE_KEY_USER);
    if (!raw) {
        // default to karyawan
        const defaultUser = {
            name: 'Rizky Darmawan',
            role: 'karyawan',
            dept: 'IT & Infrastructure',
            nip: 'EMP-2022091',
            email: 'karyawan@sipadin.com',
            avatar: 'RD'
        };
        setCurrentUser(defaultUser);
        return defaultUser;
    }
    try {
        return JSON.parse(raw);
    } catch (e) {
        return {
            name: 'Rizky Darmawan',
            role: 'karyawan',
            dept: 'IT & Infrastructure',
            nip: 'EMP-2022091',
            email: 'karyawan@sipadin.com',
            avatar: 'RD'
        };
    }
}

function setCurrentUser(user) {
    localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(user));
}

function logout() {
    localStorage.removeItem(STORAGE_KEY_USER);
    showToast('info', 'Sampai Jumpa', 'Anda telah berhasil keluar dari sistem.');
    setTimeout(() => {
        window.location.href = 'index.html';
    }, 400);
}

// ============================================================
// TOAST NOTIFICATION SYSTEM
// ============================================================
const toastConfigs = {
    success: { icon: 'ph-check-circle', color: 'text-emerald-400', bg: 'bg-surface-card border-emerald-500/40 text-emerald-100' },
    error: { icon: 'ph-x-circle', color: 'text-rose-400', bg: 'bg-surface-card border-rose-500/40 text-rose-100' },
    warning: { icon: 'ph-warning', color: 'text-amber-400', bg: 'bg-surface-card border-amber-500/40 text-amber-100' },
    info: { icon: 'ph-info', color: 'text-brand-400', bg: 'bg-surface-card border-brand-500/40 text-brand-100' },
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
        <button onclick="this.closest('.toast-enter').remove()" class="text-slate-400 hover:text-white transition-colors p-1 flex-shrink-0">
            <i class="ph ph-x text-sm"></i>
        </button>
    `;

    container.appendChild(toast);

    setTimeout(() => {
        toast.classList.remove('toast-enter');
        toast.classList.add('toast-exit');
        setTimeout(() => toast.remove(), 350);
    }, 4200);
}

// ============================================================
// MOBILE SIDEBAR CONTROLS
// ============================================================
function openMobileSidebar() {
    const sidebar = document.getElementById('mobile-sidebar') || document.getElementById('atasan-mobile-sidebar');
    const overlay = document.getElementById('sidebar-overlay') || document.getElementById('atasan-sidebar-overlay');
    if (sidebar) sidebar.classList.add('open');
    if (overlay) overlay.classList.remove('hidden');
}

function closeMobileSidebar() {
    const sidebar = document.getElementById('mobile-sidebar') || document.getElementById('atasan-mobile-sidebar');
    const overlay = document.getElementById('sidebar-overlay') || document.getElementById('atasan-sidebar-overlay');
    if (sidebar) sidebar.classList.remove('open');
    if (overlay) overlay.classList.add('hidden');
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
        modal.className = 'fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md transition-opacity';
        document.body.appendChild(modal);
    }

    const diff = daysUntil(item.mulai);
    const isUrgent = item.status === 'pending' && diff <= 3 && diff >= 0;

    let statusHtml = '';
    if (item.status === 'approved') {
        statusHtml = `<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/20 text-emerald-400 border border-emerald-500/30"><i class="ph-fill ph-check-circle"></i> Disetujui (Approved)</span>`;
    } else if (item.status === 'rejected') {
        statusHtml = `<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-rose-500/20 text-rose-400 border border-rose-500/30"><i class="ph-fill ph-x-circle"></i> Ditolak (Rejected)</span>`;
    } else {
        statusHtml = `<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-amber-500/20 text-amber-400 border border-amber-500/30 badge-pulse"><i class="ph-fill ph-clock"></i> Menunggu Persetujuan</span>`;
    }

    modal.innerHTML = `
        <div class="w-full max-w-2xl bg-surface-card border border-surface-border rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-200">
            <!-- Modal Header -->
            <div class="p-6 border-b border-surface-border flex items-center justify-between bg-surface/60">
                <div class="flex items-center gap-3">
                    <div class="w-11 h-11 rounded-2xl bg-brand-gradient flex items-center justify-center shadow-glow">
                        <i class="ph ph-file-text text-xl text-white"></i>
                    </div>
                    <div>
                        <div class="flex items-center gap-2 flex-wrap">
                            <h3 class="text-lg font-bold text-white">${item.nomorSurat || 'SPPD-' + item.id}</h3>
                            ${statusHtml}
                        </div>
                        <p class="text-xs text-slate-400 mt-0.5">Diajukan pada: ${formatDateFull(item.tanggalPengajuan || item.mulai)}</p>
                    </div>
                </div>
                <button onclick="closeDetailModal()" class="w-9 h-9 rounded-xl bg-surface-border/50 text-slate-400 hover:text-white hover:bg-surface-border flex items-center justify-center transition-colors">
                    <i class="ph ph-x text-lg"></i>
                </button>
            </div>

            <!-- Modal Content (Scrollable) -->
            <div class="p-6 overflow-y-auto space-y-6 flex-1 text-sm">
                <!-- Info Karyawan -->
                <div class="bg-surface/50 border border-surface-border rounded-2xl p-4 flex items-center justify-between flex-wrap gap-3">
                    <div class="flex items-center gap-3.5">
                        <div class="w-12 h-12 rounded-2xl bg-brand-500/20 border border-brand-500/30 flex items-center justify-center text-brand-300 font-bold text-base">
                            ${item.nama.split(' ').map(n=>n[0]).join('').slice(0,2)}
                        </div>
                        <div>
                            <p class="font-bold text-white text-base">${item.nama}</p>
                            <p class="text-xs text-slate-400 font-mono">${item.nip || 'EMP-2022091'} · ${item.dept}</p>
                        </div>
                    </div>
                    <div class="text-right">
                        <span class="text-xs text-slate-400 block">Jabatan</span>
                        <span class="text-xs font-semibold text-slate-200">${item.jabatan || 'Staf Profesional'}</span>
                    </div>
                </div>

                <!-- Rute & Jadwal -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                        <div class="flex items-center gap-2 text-brand-400 font-semibold text-xs uppercase tracking-wider">
                            <i class="ph ph-map-pin-line text-base"></i> Rute & Instansi Tujuan
                        </div>
                        <div>
                            <p class="text-xs text-slate-400">Rute Perjalanan</p>
                            <p class="text-base font-bold text-white flex items-center gap-2 mt-0.5">
                                <span>${item.asal || 'Surabaya'}</span>
                                <i class="ph ph-arrow-right text-brand-400 text-sm"></i>
                                <span class="text-brand-300">${item.tujuan}</span>
                            </p>
                        </div>
                        <div>
                            <p class="text-xs text-slate-400">Instansi / Klien Tujuan</p>
                            <p class="text-sm font-semibold text-slate-200 mt-0.5">${item.instansi || 'Kantor Cabang / Rekanan'}</p>
                        </div>
                    </div>

                    <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                        <div class="flex items-center gap-2 text-brand-400 font-semibold text-xs uppercase tracking-wider">
                            <i class="ph ph-calendar-blank text-base"></i> Jadwal & Durasi
                        </div>
                        <div class="flex items-center justify-between">
                            <div>
                                <p class="text-xs text-slate-400">Tanggal Mulai - Selesai</p>
                                <p class="text-sm font-bold text-white mt-0.5">${formatDate(item.mulai)} — ${formatDate(item.selesai)}</p>
                            </div>
                            <div class="text-right">
                                <p class="text-xs text-slate-400">Durasi</p>
                                <p class="text-sm font-bold text-brand-400 mt-0.5">${item.durasi || calculateDuration(item.mulai, item.selesai)} Hari</p>
                            </div>
                        </div>
                        <div>
                            <p class="text-xs text-slate-400">Keperluan</p>
                            <p class="text-sm font-semibold text-white mt-0.5">${item.keperluan}</p>
                        </div>
                    </div>
                </div>

                <!-- Transportasi & Akomodasi -->
                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p class="text-xs text-slate-400 flex items-center gap-1.5"><i class="ph ph-airplane text-brand-400"></i> Moda Transportasi</p>
                        <p class="text-sm font-semibold text-slate-200 mt-1">${item.transportasi || 'Pesawat Terbang / Kereta Api'}</p>
                    </div>
                    <div>
                        <p class="text-xs text-slate-400 flex items-center gap-1.5"><i class="ph ph-buildings text-brand-400"></i> Akomodasi Penginapan</p>
                        <p class="text-sm font-semibold text-slate-200 mt-1">${item.akomodasi || 'Hotel Bintang Rekanan Dinas'}</p>
                    </div>
                </div>

                <!-- Keterangan & Rincian Anggaran -->
                <div class="space-y-2">
                    <p class="text-xs font-semibold text-slate-300 uppercase tracking-wider">Keterangan & Rencana Kegiatan</p>
                    <div class="p-4 rounded-2xl bg-surface/50 border border-surface-border text-slate-300 text-xs leading-relaxed">
                        ${item.keterangan || 'Tidak ada keterangan tambahan.'}
                    </div>
                </div>

                <!-- Estimasi Anggaran Biaya -->
                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-3">
                    <p class="text-xs font-semibold text-slate-300 uppercase tracking-wider flex items-center justify-between">
                        <span>Rincian Estimasi Biaya</span>
                        <span class="text-emerald-400 font-mono font-bold">${formatRupiah(item.totalBiaya || 0)}</span>
                    </p>
                    <div class="space-y-1.5 text-xs text-slate-400 border-t border-surface-border/50 pt-2.5">
                        <div class="flex justify-between">
                            <span>Biaya Tiket & Transportasi</span>
                            <span class="text-slate-200 font-mono">${formatRupiah(item.biayaTransport || 0)}</span>
                        </div>
                        <div class="flex justify-between">
                            <span>Biaya Akomodasi / Penginapan</span>
                            <span class="text-slate-200 font-mono">${formatRupiah(item.biayaPenginapan || 0)}</span>
                        </div>
                        <div class="flex justify-between">
                            <span>Uang Harian / Konsumsi</span>
                            <span class="text-slate-200 font-mono">${formatRupiah(item.uangHarian || 0)}</span>
                        </div>
                    </div>
                </div>

                ${item.catatanAtasan ? `
                    <div class="p-4 rounded-2xl ${item.status === 'rejected' ? 'bg-rose-500/10 border border-rose-500/30' : 'bg-emerald-500/10 border border-emerald-500/30'}">
                        <p class="text-xs font-bold ${item.status === 'rejected' ? 'text-rose-300' : 'text-emerald-300'} mb-1">
                            <i class="ph ${item.status === 'rejected' ? 'ph-warning-octagon' : 'ph-chat-circle-dots'} mr-1"></i> Catatan Atasan / Manager:
                        </p>
                        <p class="text-xs text-slate-300 leading-relaxed">${item.catatanAtasan}</p>
                        ${item.disetujuiOleh ? `<p class="text-[11px] text-slate-400 mt-2 italic">— Disetujui oleh: ${item.disetujuiOleh} (${formatDate(item.tanggalDisetujui)})</p>` : ''}
                    </div>
                ` : ''}
            </div>

            <!-- Modal Footer -->
            <div class="p-5 border-t border-surface-border bg-surface/80 flex items-center justify-between gap-3">
                <button onclick="closeDetailModal()" class="px-5 py-2.5 rounded-xl border border-surface-border text-slate-300 hover:text-white hover:bg-white/5 text-sm font-semibold transition-all">
                    Tutup
                </button>
                <div class="flex gap-2.5">
                    ${item.status === 'approved' ? `
                        <button onclick="printSPPDSlip(${item.id})" class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-brand-gradient hover:opacity-90 active:scale-95 text-white text-sm font-semibold transition-all shadow-glow">
                            <i class="ph ph-printer text-base"></i> Cetak SPPD Resmi
                        </button>
                    ` : ''}
                </div>
            </div>
        </div>
    `;
    modal.classList.remove('hidden');
}

function closeDetailModal() {
    const modal = document.getElementById('modal-detail-pengajuan');
    if (modal) modal.classList.add('hidden');
}

// SPPD Document Print View Generator
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

    printArea.innerHTML = `
        <div style="font-family: 'Times New Roman', serif; max-width: 800px; margin: 0 auto; color: #111; padding: 20px; line-height: 1.6;">
            <!-- Header Surat -->
            <div style="display: flex; align-items: center; border-bottom: 3px double #000; padding-bottom: 12px; margin-bottom: 20px;">
                <div style="flex: 1; text-align: center;">
                    <h2 style="margin: 0; font-size: 18pt; text-transform: uppercase; font-weight: bold; letter-spacing: 1px;">PT ADIPRIMA SURAPRINTA</h2>
                    <p style="margin: 3px 0; font-size: 10pt;">Kawasan Industri Driyorejo, Gresik, Jawa Timur — Indonesia</p>
                    <p style="margin: 0; font-size: 9pt; color: #444;">Telp: (031) 7507888 | Email: corporate@adiprima.co.id | Website: www.adiprima.co.id</p>
                </div>
            </div>

            <!-- Judul Surat -->
            <div style="text-align: center; margin-bottom: 25px;">
                <h3 style="margin: 0; font-size: 14pt; text-decoration: underline; font-weight: bold; text-transform: uppercase;">SURAT PERINTAH PERJALANAN DINAS (SPPD)</h3>
                <p style="margin: 4px 0 0 0; font-size: 11pt;">Nomor: <strong>${item.nomorSurat || 'SPPD/2026/08/' + item.id}</strong></p>
            </div>

            <!-- Tabel Detail Perjalanan -->
            <table style="width: 100%; border-collapse: collapse; margin-bottom: 25px; font-size: 11pt;" border="1" cellpadding="8">
                <tr>
                    <td style="width: 5%; text-align: center; font-weight: bold;">1</td>
                    <td style="width: 35%;">Pejabat Berwenang yang Memberi Perintah</td>
                    <td style="width: 60%;"><strong>${item.disetujuiOleh || 'Budi Pratama (IT & Ops Manager)'}</strong></td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">2</td>
                    <td>Nama / NIP Pegawai yang Ditugaskan</td>
                    <td><strong>${item.nama}</strong> / NIP: ${item.nip || 'EMP-2022091'}</td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">3</td>
                    <td>Pangkat / Jabatan / Departemen</td>
                    <td>${item.jabatan || 'Senior Staff'} / ${item.dept}</td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">4</td>
                    <td>Maksud Perjalanan Dinas</td>
                    <td>${item.keperluan} — <em>${item.keterangan}</em></td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">5</td>
                    <td>Alat Angkut / Transportasi</td>
                    <td>${item.transportasi || 'Pesawat Terbang / Kereta Api'}</td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">6</td>
                    <td>a. Tempat Berangkat<br>b. Tempat Tujuan / Instansi</td>
                    <td>a. ${item.asal || 'Surabaya (Head Office)'}<br>b. ${item.tujuan} (${item.instansi || 'Rekanan Resmi'})</td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">7</td>
                    <td>a. Lamanya Perjalanan Dinas<br>b. Tanggal Berangkat<br>c. Tanggal Harus Kembali</td>
                    <td>a. ${item.durasi || calculateDuration(item.mulai, item.selesai)} Hari<br>b. ${formatDateFull(item.mulai)}<br>c. ${formatDateFull(item.selesai)}</td>
                </tr>
                <tr>
                    <td style="text-align: center; font-weight: bold;">8</td>
                    <td>Total Estimasi Anggaran Disetujui</td>
                    <td><strong>${formatRupiah(item.totalBiaya || 0)}</strong></td>
                </tr>
            </table>

            <!-- Tanda Tangan Pengesahan -->
            <div style="display: flex; justify-content: space-between; margin-top: 40px; font-size: 11pt;">
                <div style="text-align: center; width: 40%;">
                    <p style="margin-bottom: 60px;">Pegawai yang Melaksanakan Tugas,</p>
                    <p style="margin: 0; font-weight: bold; text-decoration: underline;">${item.nama}</p>
                    <p style="margin: 0; font-size: 9pt;">NIP: ${item.nip || 'EMP-2022091'}</p>
                </div>
                <div style="text-align: center; width: 40%;">
                    <p style="margin: 0;">Dikeluarkan di: Surabaya</p>
                    <p style="margin: 0 0 60px 0;">Pada tanggal: ${formatDate(item.tanggalDisetujui || item.tanggalPengajuan || item.mulai)}</p>
                    <p style="margin: 0; font-weight: bold; text-decoration: underline;">${item.disetujuiOleh ? item.disetujuiOleh.split(' (')[0] : 'Budi Pratama'}</p>
                    <p style="margin: 0; font-size: 9pt;">Manager Departemen Operasional</p>
                </div>
            </div>
        </div>
    `;

    window.print();
}

// Auto init on load
document.addEventListener('DOMContentLoaded', () => {
    // Highlight active link automatically based on current pathname
    const currentPath = window.location.pathname.split('/').pop() || 'index.html';
    document.querySelectorAll('.sidebar-link').forEach(link => {
        const href = link.getAttribute('href');
        if (href && (href === currentPath || (currentPath === '' && href === 'index.html'))) {
            link.classList.add('active');
        } else {
            link.classList.remove('active');
        }
    });
});
