// ============================================================
// SIPADIN — SISTEM PENGAJUAN DINAS
// GLOBAL JAVASCRIPT — 6-STAGE PIPELINE EDITION v4.0 (FINAL FIX)
// Backend: Go/Fiber REST API @ http://localhost:8080/api
// ============================================================

const API_BASE = 'http://localhost:8080/api';

// Global cache for fetched submissions
window.cachedPengajuanList = [];

// ── Tailwind Theme Configuration ────────────────────────────
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
                        50: '#eef2ff', 100: '#e0e7ff', 200: '#c7d2fe',
                        300: '#a5b4fc', 400: '#818cf8', 500: '#6366f1',
                        600: '#4f46e5', 700: '#4338ca', 800: '#3730a3',
                        900: '#312e81', 950: '#1e1b4b',
                    },
                    surface: {
                        DEFAULT: '#0b0f19', card: '#1e293b',
                        cardhover: '#24334a', border: '#334155', muted: '#475569',
                    }
                },
                boxShadow: {
                    glow:        '0 0 35px -5px rgba(99,102,241,0.45)',
                    glowGreen:   '0 0 35px -5px rgba(16,185,129,0.45)',
                    glowRed:     '0 0 35px -5px rgba(239,68,68,0.45)',
                    glowPurple:  '0 0 35px -5px rgba(168,85,247,0.45)',
                    glowAmber:   '0 0 35px -5px rgba(245,158,11,0.45)',
                    card:        '0 10px 30px -10px rgba(0,0,0,0.6)',
                },
                backgroundImage: {
                    'hero-gradient':    'radial-gradient(ellipse at top,#1e1b4b 0%,#0f172a 45%,#0b0f19 100%)',
                    'brand-gradient':   'linear-gradient(135deg,#6366f1 0%,#4f46e5 100%)',
                    'accent-gradient':  'linear-gradient(135deg,#8b5cf6 0%,#6366f1 100%)',
                    'emerald-gradient': 'linear-gradient(135deg,#10b981 0%,#059669 100%)',
                    'rose-gradient':    'linear-gradient(135deg,#f43f5e 0%,#e11d48 100%)',
                    'hrga-gradient':    'linear-gradient(135deg,#f59e0b 0%,#d97706 100%)',
                    'teal-gradient':    'linear-gradient(135deg,#14b8a6 0%,#0d9488 100%)',
                }
            }
        }
    };
}

// ============================================================
// DATE & HELPER UTILITIES
// ============================================================
const _today = new Date();
_today.setHours(0, 0, 0, 0);

function daysFromNow(n) {
    const d = new Date(_today);
    d.setDate(d.getDate() + n);
    return d.toISOString().split('T')[0];
}

function formatDate(dateStr) {
    if (!dateStr) return '-';
    const raw = typeof dateStr === 'string' ? dateStr.split('T')[0] : dateStr;
    const d = new Date(raw + 'T00:00:00');
    if (isNaN(d)) return '-';
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
}

function formatDateFull(dateStr) {
    if (!dateStr) return '-';
    const raw = typeof dateStr === 'string' ? dateStr.split('T')[0] : dateStr;
    const d = new Date(raw + 'T00:00:00');
    if (isNaN(d)) return '-';
    return d.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
}

function formatRupiah(num) {
    if (isNaN(num) || num === null || num === undefined) return 'Rp 0';
    return 'Rp ' + Number(num).toLocaleString('id-ID');
}

function daysUntil(dateStr) {
    if (!dateStr) return 999;
    const raw = typeof dateStr === 'string' ? dateStr.split('T')[0] : dateStr;
    const target = new Date(raw + 'T00:00:00');
    if (isNaN(target)) return 999;
    return Math.ceil((target - _today) / (1000 * 60 * 60 * 24));
}

function calculateDuration(startStr, endStr) {
    if (!startStr || !endStr) return 1;
    const s = new Date((typeof startStr === 'string' ? startStr.split('T')[0] : startStr) + 'T00:00:00');
    const e = new Date((typeof endStr === 'string' ? endStr.split('T')[0] : endStr) + 'T00:00:00');
    if (isNaN(s) || isNaN(e)) return 1;
    const diff = Math.ceil((e - s) / (1000 * 60 * 60 * 24)) + 1;
    return diff > 0 ? diff : 1;
}

// ============================================================
// 6-STAGE STATUS PIPELINE DEFINITION
// ============================================================
const STATUS_PIPELINE = {
    SUBMITTED:          { label: 'Menunggu Atasan',        icon: 'ph-clock',             color: 'amber',   step: 1 },
    APPROVED_ATASAN:    { label: 'Pending Kontrol HRGA',   icon: 'ph-check-circle',      color: 'blue',    step: 2 },
    CONTROLLED_HRD:     { label: 'Terverifikasi HRGA',     icon: 'ph-shield-check',      color: 'teal',    step: 3 },
    WA_SENT_DIREKSI:    { label: 'WA Terkirim ke Direksi', icon: 'ph-paper-plane-tilt',  color: 'violet',  step: 4 },
    CONFIRMED_DIREKSI:  { label: 'Disetujui Direksi',      icon: 'ph-seal-check',        color: 'emerald', step: 5 },
    SURAT_TUGAS_ISSUED: { label: 'Surat Tugas Terbit',     icon: 'ph-file-text',         color: 'indigo',  step: 6 },
    REJECTED:           { label: 'Ditolak',                icon: 'ph-x-circle',          color: 'rose',    step: 0 },
    // Backward-compat aliases
    pending:            { label: 'Menunggu Atasan',        icon: 'ph-clock',             color: 'amber',   step: 1 },
    approved:           { label: 'Pending Kontrol HRGA',   icon: 'ph-check-circle',      color: 'blue',    step: 2 },
    verified_hrga:      { label: 'Terverifikasi HRGA',     icon: 'ph-shield-check',      color: 'teal',    step: 3 },
    notified_direksi:   { label: 'WA Terkirim ke Direksi', icon: 'ph-paper-plane-tilt',  color: 'violet',  step: 4 },
    rejected:           { label: 'Ditolak',                icon: 'ph-x-circle',          color: 'rose',    step: 0 },
};

const STATUS_COLOR_MAP = {
    amber:   'bg-amber-500/20 text-amber-400 border-amber-500/30',
    blue:    'bg-blue-500/20 text-blue-400 border-blue-500/30',
    teal:    'bg-teal-500/20 text-teal-400 border-teal-500/30',
    violet:  'bg-violet-500/20 text-violet-400 border-violet-500/30',
    emerald: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30',
    indigo:  'bg-indigo-500/20 text-indigo-400 border-indigo-500/30',
    rose:    'bg-rose-500/20 text-rose-400 border-rose-500/30',
};

function getStatusConfig(statusFlow) {
    return STATUS_PIPELINE[statusFlow] || STATUS_PIPELINE['SUBMITTED'];
}

function getStatusBadge(statusFlow, size) {
    const cfg = getStatusConfig(statusFlow);
    const sizeClass = size === 'xs' ? 'px-2 py-0.5 text-[10px]' : 'px-3 py-1 text-xs';
    const colorClass = STATUS_COLOR_MAP[cfg.color] || STATUS_COLOR_MAP.amber;
    const pulse = (statusFlow === 'SUBMITTED' || statusFlow === 'APPROVED_ATASAN' || statusFlow === 'pending' || statusFlow === 'approved') ? 'badge-pulse' : '';
    return `<span class="inline-flex items-center gap-1.5 ${sizeClass} rounded-full font-semibold border ${colorClass} ${pulse}">
        <i class="ph-fill ${cfg.icon}"></i> ${cfg.label}
    </span>`;
}

function getStatusLabel(statusFlow) {
    return getStatusConfig(statusFlow).label;
}

function getStatusStep(statusFlow) {
    return getStatusConfig(statusFlow).step;
}

// ============================================================
// DATA NORMALIZATION — Maps DB/API fields to standard UI properties
// ============================================================
function normalizePengajuan(item) {
    if (!item) return null;
    const k = item.karyawan || {};
    const atasan = item.atasan || {};

    const noPengajuan        = item.no_pengajuan      || item.NoPengajuan      || item.id || '';
    const nomorID            = item.nomor_id          || item.NomorID          || (item.user ? item.user.nip : '') || '';
    const nama               = item.nama              || k.nama                || (item.user ? item.user.nama : '') || nomorID || 'Karyawan';
    const kodeDepartemen     = item.kode_departemen   || k.kode_departemen     || item.KodeDepartemen || item.dept || (item.user ? item.user.departemen : '') || '-';
    const kodeJabatan        = item.kode_jabatan      || k.kode_jabatan        || item.KodeJabatan    || item.jabatan || (item.user ? item.user.jabatan : '') || 'Staf';
    const email              = item.email             || k.email               || (item.user ? item.user.email : '') || '';
    const atasanNomorID      = item.atasan_nomor_id   || item.AtasanNomorID    || '';
    const atasanNama         = item.atasan_nama       || atasan.nama           || atasanNomorID || '-';
    const hrdNomorID         = item.hrd_nomor_id      || item.HRDNomorID       || '';
    const areaTujuan         = item.area_tujuan       || item.AreaTujuan       || 'Jawa';
    const kotaTujuan         = item.kota_tujuan       || item.KotaTujuan       || item.tujuan || 'Surabaya';
    const maksudTujuan       = item.maksud_tujuan     || item.MaksudTujuan     || item.keperluan || '-';
    const kriteriaFasilitas  = item.kriteria_fasilitas|| item.KriteriaFasilitas|| item.kriteria || item.transportasi || '1. < 200 KM & < 8 Jam';
    const tanggalBerangkat   = item.tanggal_berangkat || item.TanggalBerangkat || item.tgl_keberangkatan || item.mulai || '';
    const tanggalKembali     = item.tanggal_kembali   || item.TanggalKembali   || item.tgl_kembali || item.selesai || '';
    const jamBerangkat       = item.jam_berangkat     || item.JamBerangkat     || '08:00';
    const jamKembali         = item.jam_kembali       || item.JamKembali       || '17:00';
    const estimasiBiaya      = Number(item.estimasi_biaya !== undefined ? item.estimasi_biaya : (item.EstimasiBiaya !== undefined ? item.EstimasiBiaya : (item.totalBiaya || 0)));
    const statusFlow         = item.status_flow       || item.StatusFlow       || item.status || 'SUBMITTED';
    const catatanHRD         = item.catatan_hrd       || item.CatatanHRD       || item.catatan_hrga || item.catatan_atasan || '';
    const konfirmasiNote     = item.konfirmasi_direksi_note || item.KonfirmasiDireksiNote || '';
    const tglKonfirmasi      = item.tanggal_konfirmasi_direksi || item.TanggalKonfirmasiDireksi || null;
    const nomorSuratTugas    = item.nomor_surat_tugas || item.NomorSuratTugas  || item.nomor_surat || '';
    const createdAt          = item.created_at        || item.CreatedAt        || item.tanggal_pengajuan || '';
    const updatedAt          = item.updated_at        || item.UpdatedAt        || '';

    const durasi = calculateDuration(tanggalBerangkat, tanggalKembali);

    return {
        // Canonical fields
        id:                       noPengajuan,
        noPengajuan,
        nomorID,
        nama,
        kodeDepartemen,
        kodeJabatan,
        email,
        atasanNomorID,
        atasanNama,
        hrdNomorID,
        areaTujuan,
        kotaTujuan,
        maksudTujuan,
        kriteriaFasilitas,
        tanggalBerangkat,
        tanggalKembali,
        jamBerangkat,
        jamKembali,
        estimasiBiaya,
        statusFlow,
        catatanHRD,
        konfirmasiDireksiNote:    konfirmasiNote,
        tanggalKonfirmasiDireksi: tglKonfirmasi,
        nomorSuratTugas,
        createdAt,
        updatedAt,

        // Backward-compatible aliases for dashboard templates
        nip:              nomorID,
        dept:             kodeDepartemen,
        jabatan:          kodeJabatan,
        departemen:       kodeDepartemen,
        tujuan:           kotaTujuan || areaTujuan,
        asal:             'Surabaya (Head Office)',
        instansi:         kotaTujuan,
        keperluan:        maksudTujuan,
        keterangan:       maksudTujuan,
        kriteria:         kriteriaFasilitas,
        transportasi:     kriteriaFasilitas,
        akomodasi:        durasi > 1 ? 'Hotel Sesuai Tarif' : 'Tanpa Menginap',
        mulai:            tanggalBerangkat,
        selesai:          tanggalKembali,
        durasi:           durasi,
        totalBiaya:       estimasiBiaya,
        biayaTransport:   Math.round(estimasiBiaya * 0.4),
        biayaPenginapan:  durasi > 1 ? Math.round(estimasiBiaya * 0.35) : 0,
        uangHarian:       Math.round(estimasiBiaya * 0.25),
        nomorSurat:       nomorSuratTugas || noPengajuan,
        status:           statusFlow,
        catatanAtasan:    catatanHRD,
        catatanHRGA:      catatanHRD,
        disetujuiOleh:    atasanNama || atasanNomorID,
        diverifikasiOleh: hrdNomorID || 'HRGA Supervisor',
        tanggalDisetujui: updatedAt ? String(updatedAt).split('T')[0] : '',
        tanggalVerifikasiHRGA:     updatedAt ? String(updatedAt).split('T')[0] : '',
        tanggalNotifikasiDireksi:  tglKonfirmasi ? String(tglKonfirmasi).split('T')[0] : '',
        tanggalPengajuan:          createdAt ? String(createdAt).split('T')[0] : '',
        pesanWA:          '',
    };
}

// ============================================================
// API CLIENT FUNCTIONS
// ============================================================

/**
 * POST /api/auth/login
 */
async function apiLogin(nama, password) {
    const response = await fetch(`${API_BASE}/auth/login`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ nama, password }),
    });
    const data = await response.json();
    if (!response.ok) throw new Error(data.message || 'Login gagal');
    return data;
}

/** GET /api/master/tarif */
async function apiGetMasterTarif() {
    const res = await fetch(`${API_BASE}/master/tarif`);
    if (!res.ok) throw new Error('Gagal mengambil data tarif');
    return res.json();
}

/** GET /api/master/fasilitas */
async function apiGetMasterFasilitas() {
    const res = await fetch(`${API_BASE}/master/fasilitas`);
    if (!res.ok) throw new Error('Gagal mengambil data fasilitas');
    return res.json();
}

/**
 * GET /api/pengajuan
 */
async function apiGetPengajuan(params = {}) {
    try {
        const query = new URLSearchParams(params).toString();
        const url   = `${API_BASE}/pengajuan${query ? '?' + query : ''}`;
        const res   = await fetch(url);
        if (!res.ok) throw new Error('Gagal mengambil data pengajuan');
        const data = await res.json();
        const normalized = Array.isArray(data) ? data.map(normalizePengajuan) : [];
        window.cachedPengajuanList = normalized;
        return normalized;
    } catch (err) {
        console.warn('API getPengajuan fallback:', err.message);
        return window.cachedPengajuanList || [];
    }
}

/**
 * GET /api/pengajuan/detail?id=... (with list lookup fallback)
 */
async function apiGetPengajuanByID(id) {
    if (!id) return null;
    const cleanId = String(id).trim();
    try {
        const res = await fetch(`${API_BASE}/pengajuan/detail?id=${encodeURIComponent(cleanId)}`);
        if (res.ok) {
            const data = await res.json();
            return normalizePengajuan(data);
        }
    } catch (err) {
        console.warn('apiGetPengajuanByID detail query failed, fallback to list lookup:', err);
    }
    
    // Fallback: search in cached list or fetch fresh list
    let list = window.cachedPengajuanList;
    if (!list || !list.length) {
        list = await apiGetPengajuan();
    }
    const found = list.find(p => p.noPengajuan === cleanId || p.id === cleanId || p.nomorSurat === cleanId || p.nomorSuratTugas === cleanId);
    if (found) return normalizePengajuan(found);
    throw new Error('Pengajuan tidak ditemukan');
}

/**
 * POST /api/pengajuan
 */
async function apiCreatePengajuan(payload) {
    const user = getCurrentUser();
    if (!payload.nomor_id && user) {
        payload.nomor_id = user.nip || user.nomor_id || '';
    }

    const res  = await fetch(`${API_BASE}/pengajuan`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(payload),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal membuat pengajuan');
    return normalizePengajuan(data.data || data);
}

/**
 * PUT /api/pengajuan/action/atasan-approve?id=...
 */
async function apiAtasanApprove(id, atasanNomorID, catatan) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/atasan-approve?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ atasan_nomor_id: atasanNomorID, catatan: catatan || 'Disetujui' }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal menyetujui pengajuan');
    return normalizePengajuan(data.data || data);
}

/**
 * PUT /api/pengajuan/action/atasan-reject?id=...
 */
async function apiAtasanReject(id, catatan) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/atasan-reject?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ catatan_hrd: catatan, catatan: catatan }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal menolak pengajuan');
    return normalizePengajuan(data.data || data);
}

/**
 * PUT /api/pengajuan/action/hrd-control?id=...
 */
async function apiHRDControl(id, payload) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/hrd-control?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(payload),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal mengontrol pengajuan');
    return normalizePengajuan(data.data || data);
}

/**
 * PUT /api/pengajuan/action/hrd-reject?id=...
 */
async function apiHRDReject(id, note) {
    const user = getCurrentUser();
    const res  = await fetch(`${API_BASE}/pengajuan/action/hrd-reject?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({
            catatan_hrd:  note || 'Tidak memenuhi kebijakan HRGA.',
            hrd_nomor_id: user?.nip || user?.nomor_id || '',
        }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal menolak pengajuan HRGA');
    return normalizePengajuan(data.data || data);
}

/**
 * GET /api/pengajuan/action/wa-text?id=...
 */
async function apiGetWAText(id) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/wa-text?id=${encodeURIComponent(id)}`);
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal menyiapkan teks WA');
    return data;
}

/**
 * PUT /api/pengajuan/action/direksi-confirm?id=...
 */
async function apiDireksiConfirm(id, note) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/direksi-confirm?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ konfirmasi_note: note || 'OK / Disetujui' }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal merekam konfirmasi Direksi');
    return normalizePengajuan(data.data || data);
}

/**
 * PUT /api/pengajuan/action/direksi-reject?id=...
 */
async function apiDireksiReject(id, note) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/direksi-reject?id=${encodeURIComponent(id)}`, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ konfirmasi_note: note || 'Ditolak oleh Direksi' }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal merekam penolakan Direksi');
    return normalizePengajuan(data.data || data);
}

/**
 * POST /api/pengajuan/action/surat-tugas?id=...
 */
async function apiIssueSuratTugas(id) {
    const res  = await fetch(`${API_BASE}/pengajuan/action/surat-tugas?id=${encodeURIComponent(id)}`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Gagal menerbitkan Surat Tugas');
    return data;
}

// ============================================================
// SESSION MANAGEMENT (Multi-key Sync: currentUser, user, sipadin_user_v3)
// ============================================================
const STORAGE_KEYS = ['currentUser', 'user', 'sipadin_user_v3'];

function getCurrentUser() {
    for (const key of STORAGE_KEYS) {
        try {
            const raw = localStorage.getItem(key);
            if (raw) {
                const u = JSON.parse(raw);
                if (u && (u.nama || u.name || u.nomor_id || u.nip)) {
                    if (!u.name && u.nama)             u.name = u.nama;
                    if (!u.nama && u.name)             u.nama = u.name;
                    if (!u.dept && u.kode_departemen)  u.dept = u.kode_departemen;
                    if (!u.departemen && u.kode_departemen) u.departemen = u.kode_departemen;
                    if (!u.nip && u.nomor_id)          u.nip = u.nomor_id;
                    if (!u.nomor_id && u.nip)          u.nomor_id = u.nip;
                    if (!u.avatar) u.avatar = _buildAvatarFromName(u.nama || u.name || 'U');
                    return u;
                }
            }
        } catch {}
    }
    return null;
}

function setCurrentUser(user) {
    if (!user) return;
    if (!user.name && user.nama)             user.name = user.nama;
    if (!user.nama && user.name)             user.nama = user.name;
    if (!user.dept && user.kode_departemen)  user.dept = user.kode_departemen;
    if (!user.departemen && user.kode_departemen) user.departemen = user.kode_departemen;
    if (!user.nip && user.nomor_id)          user.nip = user.nomor_id;
    if (!user.nomor_id && user.nip)          user.nomor_id = user.nip;
    if (!user.avatar) user.avatar = _buildAvatarFromName(user.nama || user.name || 'U');

    const json = JSON.stringify(user);
    STORAGE_KEYS.forEach(k => localStorage.setItem(k, json));
}

function _buildAvatarFromName(name) {
    return (name || 'U').split(' ').filter(Boolean).map(n => n[0]).join('').slice(0, 2).toUpperCase();
}

function logout() {
    STORAGE_KEYS.forEach(k => localStorage.removeItem(k));
    showToast('info', 'Sampai Jumpa', 'Anda telah berhasil keluar dari sistem.');
    setTimeout(() => { window.location.href = 'index.html'; }, 400);
}

function redirectToDashboard(role) {
    const r = (role || '').toUpperCase();
    if (r === 'ATASAN') {
        window.location.href = 'dashboard-atasan.html';
    } else if (r === 'HRGA' || r === 'HRD') {
        window.location.href = 'dashboard-hrga.html';
    } else if (r === 'DIREKSI' || r === 'DIR' || r === 'BOD') {
        window.location.href = 'dashboard-hrga.html';
    } else {
        window.location.href = 'dashboard-karyawan.html';
    }
}

// ============================================================
// 1. PROFILE SYNCHRONIZATION FOR ALL ROLES (syncUserProfile)
// ============================================================
function syncUserProfile() {
    const user = getCurrentUser();
    const currentPath = window.location.pathname.split('/').pop().toLowerCase() || 'index.html';
    const isIndex = currentPath === 'index.html' || currentPath === '';

    // If not on login page and no session, redirect to index.html
    if (!user && !isIndex) {
        window.location.href = 'index.html';
        return;
    }

    if (!user) return;

    const initial = user.avatar || _buildAvatarFromName(user.nama || user.name || 'U');
    const fullName = user.nama || user.name || 'Pengguna';
    const dept = user.kode_departemen || user.dept || user.departemen || '-';
    const roleText = user.role || 'Karyawan';
    const nip = user.nomor_id || user.nip || '-';

    // Update uniform standard IDs
    const elInitial = document.getElementById('userAvatarInitial');
    if (elInitial) elInitial.textContent = initial;

    const elName = document.getElementById('userProfileName');
    if (elName) elName.textContent = fullName;

    const elDept = document.getElementById('userProfileDept');
    if (elDept) elDept.textContent = dept;

    const elRole = document.getElementById('userProfileRole');
    if (elRole) {
        elRole.textContent = roleText;
    }

    // Update greeting banner name if present
    const elGreeting = document.getElementById('userGreetingName');
    if (elGreeting) elGreeting.textContent = fullName;

    // Backward-compat data attribute selectors
    document.querySelectorAll('[data-user-avatar]').forEach(el => el.textContent = initial);
    document.querySelectorAll('[data-user-name]').forEach(el => el.textContent = fullName);
    document.querySelectorAll('[data-user-dept]').forEach(el => el.textContent = dept);
    document.querySelectorAll('[data-user-role]').forEach(el => el.textContent = roleText);
    document.querySelectorAll('[data-user-nip]').forEach(el => el.textContent = nip);

    // Auto-fill & lock on pengajuan-dinas.html
    if (document.getElementById('main-pengajuan-form') || document.getElementById('inp-pemohon-nama')) {
        const inpNama = document.getElementById('inp-pemohon-nama');
        if (inpNama) {
            inpNama.value = fullName;
            inpNama.readOnly = true;
        }

        const inpDept = document.getElementById('inp-pemohon-dept');
        if (inpDept) {
            inpDept.value = dept;
            inpDept.readOnly = true;
        }

        const inpId = document.getElementById('inp-pemohon-id');
        if (inpId) {
            inpId.value = nip;
            inpId.readOnly = true;
        }

        const hiddenId = document.getElementById('nomor_id');
        if (hiddenId) {
            hiddenId.value = nip;
        }
    }
}
window.syncUserProfile = syncUserProfile;

// ============================================================
// TOAST NOTIFICATION SYSTEM
// ============================================================
const _toastConfigs = {
    success: { icon: 'ph-check-circle', color: 'text-emerald-400', bg: 'bg-surface-card border-emerald-500/40' },
    error:   { icon: 'ph-x-circle',     color: 'text-rose-400',    bg: 'bg-surface-card border-rose-500/40'    },
    warning: { icon: 'ph-warning',      color: 'text-amber-400',   bg: 'bg-surface-card border-amber-500/40'   },
    info:    { icon: 'ph-info',         color: 'text-brand-400',   bg: 'bg-surface-card border-brand-500/40'   },
};

function showToast(type, title, message) {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        container.className = 'fixed top-5 right-5 z-[99999] space-y-3 w-80 max-w-[calc(100vw-2.5rem)] pointer-events-none';
        document.body.appendChild(container);
    }
    const cfg   = _toastConfigs[type] || _toastConfigs.info;
    const toast = document.createElement('div');
    toast.className = `pointer-events-auto toast-enter flex items-start gap-3 p-4 rounded-2xl border ${cfg.bg} shadow-2xl backdrop-blur-xl transition-all duration-300`;
    toast.innerHTML = `
        <div class="p-1 rounded-xl bg-white/5 flex-shrink-0"><i class="ph ${cfg.icon} ${cfg.color} text-2xl"></i></div>
        <div class="flex-1 min-w-0 pt-0.5">
            <h4 class="text-sm font-bold text-white tracking-tight">${title}</h4>
            <p class="text-xs text-slate-300 mt-1 leading-relaxed">${message}</p>
        </div>
        <button onclick="this.closest('div[class*=toast]').remove()" class="text-slate-400 hover:text-white transition-colors p-1 flex-shrink-0">
            <i class="ph ph-x text-sm"></i>
        </button>`;
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
    ['mobile-sidebar','atasan-mobile-sidebar','hrga-mobile-sidebar'].forEach(id => {
        const el = document.getElementById(id); if (el) el.classList.add('open');
    });
    ['sidebar-overlay','atasan-sidebar-overlay','hrga-sidebar-overlay'].forEach(id => {
        const el = document.getElementById(id); if (el) el.classList.remove('hidden');
    });
}

function closeMobileSidebar() {
    ['mobile-sidebar','atasan-mobile-sidebar','hrga-mobile-sidebar'].forEach(id => {
        const el = document.getElementById(id); if (el) el.classList.remove('open');
    });
    ['sidebar-overlay','atasan-sidebar-overlay','hrga-sidebar-overlay'].forEach(id => {
        const el = document.getElementById(id); if (el) el.classList.add('hidden');
    });
}

// ============================================================
// WORKFLOW ACTIONS — HIGH-LEVEL WRAPPERS
// ============================================================

/** Stage 2: Atasan approves */
async function approveByAtasan(noPengajuan, catatan) {
    const user = getCurrentUser();
    try {
        const updated = await apiAtasanApprove(noPengajuan, user?.nip || user?.nomor_id || '', catatan);
        showToast('success', 'Pengajuan Disetujui', `${updated.noPengajuan} → diteruskan ke HRGA.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Menyetujui', err.message);
        return false;
    }
}

/** Stage 2: Atasan rejects */
async function rejectByAtasan(noPengajuan, catatan) {
    try {
        const updated = await apiAtasanReject(noPengajuan, catatan || 'Tidak disetujui.');
        showToast('warning', 'Pengajuan Ditolak', `${updated.noPengajuan} telah ditolak.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Menolak', err.message);
        return false;
    }
}

/** Stage 3: HRGA controls & verifies */
async function verifyByHRGA(noPengajuan, catatan, estimasiBiaya) {
    const user = getCurrentUser();
    try {
        const updated = await apiHRDControl(noPengajuan, {
            catatan_hrd:    catatan       || 'Dokumen dan anggaran telah diverifikasi sesuai kebijakan.',
            estimasi_biaya: estimasiBiaya || 0,
            hrd_nomor_id:   user?.nip     || user?.nomor_id || '',
        });
        showToast('success', 'Verifikasi HRGA Berhasil', `${updated.noPengajuan} siap dikirim ke Direksi.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Verifikasi HRGA', err.message);
        return false;
    }
}
const controlByHRD = verifyByHRGA;

/** Stage 3 (reject): HRGA rejects */
async function rejectByHRGA(noPengajuan, catatan) {
    try {
        const updated = await apiHRDReject(noPengajuan, catatan || 'Tidak memenuhi kebijakan HRGA.');
        showToast('warning', 'Pengajuan Ditolak HRGA', `${updated.noPengajuan} ditolak oleh HRGA.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Menolak HRGA', err.message);
        return false;
    }
}

/** Stage 4: HRGA sends WhatsApp to Direksi */
async function sendWADireksi(noPengajuan) {
    try {
        const result = await apiGetWAText(noPengajuan);
        const waUrl = result.whatsapp_url || result.wa_url || `https://api.whatsapp.com/send?text=${encodeURIComponent(result.text || '')}`;
        window.open(waUrl, '_blank');
        showToast('success', 'WhatsApp Dibuka', `Pesan untuk ${noPengajuan} siap dikirim. Status: WA_SENT_DIREKSI.`);
        return result;
    } catch (err) {
        showToast('error', 'Gagal Membuka WhatsApp', err.message);
        return null;
    }
}
const notifyDireksi = sendWADireksi;

/** Stage 5: HRGA records Direksi confirmation */
async function confirmDireksi(noPengajuan, note) {
    try {
        const updated = await apiDireksiConfirm(noPengajuan, note || 'OK / Disetujui');
        showToast('success', 'Konfirmasi Direksi Tersimpan', `${updated.noPengajuan} → CONFIRMED_DIREKSI.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Rekam Konfirmasi', err.message);
        return false;
    }
}

/** Stage 5 (reject): Direksi rejects */
async function rejectByDireksi(noPengajuan, note) {
    try {
        const updated = await apiDireksiReject(noPengajuan, note || 'Ditolak oleh Direksi');
        showToast('warning', 'Pengajuan Ditolak Direksi', `${updated.noPengajuan} ditolak oleh Direksi.`);
        return true;
    } catch (err) {
        showToast('error', 'Gagal Menolak Direksi', err.message);
        return false;
    }
}

/** Stage 6: HRGA issues Surat Tugas */
async function issueSuratTugas(noPengajuan) {
    try {
        const result = await apiIssueSuratTugas(noPengajuan);
        showToast('success', 'Surat Tugas Diterbitkan', `Nomor: ${result.surat_tugas?.nomor_surat_tugas || '-'}`);
        return result;
    } catch (err) {
        showToast('error', 'Gagal Terbitkan Surat Tugas', err.message);
        return null;
    }
}

// Global Aliases for event delegation & inline callbacks
window.approveByAtasan   = approveByAtasan;
window.rejectByAtasan    = rejectByAtasan;
window.verifyByHRGA      = verifyByHRGA;
window.rejectByHRGA      = rejectByHRGA;
window.sendWADireksi     = sendWADireksi;
window.confirmDireksi    = confirmDireksi;
window.rejectByDireksi   = rejectByDireksi;
window.issueSuratTugas   = issueSuratTugas;

// ============================================================
// 2. POPUP / MODAL DETAIL DINAMIS (window.showDetail)
// ============================================================
async function showDetail(noPengajuan) {
    if (!noPengajuan) {
        showToast('warning', 'Peringatan', 'Nomor pengajuan tidak ditemukan.');
        return;
    }

    let item;
    try {
        if (typeof noPengajuan === 'object' && (noPengajuan.noPengajuan || noPengajuan.id)) {
            item = normalizePengajuan(noPengajuan);
        } else {
            item = await apiGetPengajuanByID(noPengajuan);
        }
    } catch (err) {
        showToast('error', 'Error', 'Pengajuan tidak ditemukan: ' + err.message);
        return;
    }

    // Modal container (dynamically created if not present)
    let modal = document.getElementById('modalDetailPengajuan');
    if (!modal) {
        modal = document.getElementById('modalDetail');
    }
    if (!modal) {
        modal = document.createElement('div');
        modal.id = 'modalDetailPengajuan';
        modal.className = 'fixed inset-0 z-[9999] flex items-center justify-center p-4 bg-black/80 backdrop-blur-md transition-all';
        document.body.appendChild(modal);
    }

    const dur = calculateDuration(item.tanggalBerangkat, item.tanggalKembali);
    const hDiff = daysUntil(item.tanggalBerangkat);
    const hBadge = hDiff >= 0 && hDiff <= 3 && item.statusFlow === 'SUBMITTED'
        ? `<span class="px-2 py-0.5 rounded-md bg-rose-500/20 text-rose-300 border border-rose-500/30 text-[10px] font-bold badge-pulse">Mendesak H-${hDiff}</span>`
        : '';

    const isRejected = item.statusFlow === 'REJECTED' || item.status === 'rejected';
    const isDireksiRejected = isRejected && (Boolean(item.konfirmasiDireksiNote) || Boolean(item.tanggalKonfirmasiDireksi));
    const isHRGARejected = isRejected && !isDireksiRejected && (Boolean(item.hrdNomorID) || (item.catatanHRD && item.catatanHRD.toLowerCase().includes('hrga')));
    const isAtasanRejected = isRejected && !isDireksiRejected && !isHRGARejected;

    const stages = [
        { label: '1. Pengajuan Karyawan',      done: true, note: '', sub: `${formatDate(item.createdAt || item.tanggalBerangkat)} · ${item.nama}` },
        { 
            label: '2. Approval Atasan',         
            done: ['APPROVED_ATASAN','CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) || isHRGARejected || isDireksiRejected,
            rej: isAtasanRejected, 
            note: isAtasanRejected ? item.catatanHRD : '', 
            sub: isAtasanRejected ? (item.atasanNama ? `Ditolak oleh ${item.atasanNama}` : 'Ditolak oleh Atasan') : (item.atasanNama || item.atasanNomorID ? `Disetujui oleh: ${item.atasanNama || item.atasanNomorID}` : '') 
        },
        { 
            label: '3. Kontrol & Verifikasi HRGA',
            done: ['CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) || isDireksiRejected,
            rej: isHRGARejected,
            note: isHRGARejected ? item.catatanHRD : (['CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) ? item.catatanHRD : ''), 
            sub: isHRGARejected ? (item.hrdNomorID ? `Ditolak oleh HRGA (${item.hrdNomorID})` : 'Ditolak oleh HRGA') : (item.hrdNomorID ? `Verifikator: ${item.hrdNomorID}` : '') 
        },
        { 
            label: '4. WhatsApp ke Direksi',     
            done: ['WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) || isDireksiRejected,
            note: '', 
            sub: item.statusFlow === 'WA_SENT_DIREKSI' || ['CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) || isDireksiRejected ? 'Pesan resmi telah dikirim' : '' 
        },
        { 
            label: '5. Konfirmasi Direksi',      
            done: ['CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow),
            rej: isDireksiRejected,
            note: item.konfirmasiDireksiNote, 
            sub: isDireksiRejected ? 'Ditolak oleh Direksi' : (item.tanggalKonfirmasiDireksi ? `Disetujui: ${formatDate(item.tanggalKonfirmasiDireksi)}` : '') 
        },
        { 
            label: '6. Surat Tugas Terbit',     
            done: item.statusFlow === 'SURAT_TUGAS_ISSUED',
            note: item.nomorSuratTugas, 
            sub: item.nomorSuratTugas ? `No: ${item.nomorSuratTugas}` : '' 
        },
    ];

    const timelineHtml = stages.map(s => {
        const c  = s.rej ? 'rose' : s.done ? 'emerald' : 'slate';
        const ic = s.rej ? 'ph-x-circle' : s.done ? 'ph-check-circle' : 'ph-clock';
        return `
            <div class="flex gap-3">
                <div class="w-7 h-7 rounded-full flex items-center justify-center bg-${c}-500/20 border border-${c}-500/40 text-${c}-400 flex-shrink-0 mt-0.5">
                    <i class="ph-fill ${ic} text-xs"></i>
                </div>
                <div class="pb-3 flex-1">
                    <p class="text-xs font-bold text-white">${s.label}</p>
                    ${s.sub ? `<p class="text-[11px] text-slate-400 mt-0.5">${s.sub}</p>` : `<p class="text-[11px] text-slate-500 mt-0.5">Menunggu proses...</p>`}
                    ${s.note ? `<p class="text-[11px] text-amber-300/90 mt-0.5 italic">"${s.note}"</p>` : ''}
                </div>
            </div>`;
    }).join('');

    modal.innerHTML = `
    <div class="w-full max-w-2xl bg-surface-card border border-surface-border rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[92vh]" onclick="event.stopPropagation()">
        <!-- Header -->
        <div class="p-5 border-b border-surface-border flex items-center justify-between bg-surface/80 flex-shrink-0">
            <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-2xl bg-brand-gradient flex items-center justify-center shadow-glow flex-shrink-0">
                    <i class="ph-bold ph-file-text text-lg text-white"></i>
                </div>
                <div>
                    <div class="flex items-center gap-2 flex-wrap">
                        <h3 class="text-base font-bold text-white font-mono">${item.noPengajuan}</h3>
                        ${getStatusBadge(item.statusFlow, 'xs')}
                        ${hBadge}
                    </div>
                    <p class="text-[11px] text-slate-400 mt-0.5">Diajukan pada: ${formatDateFull(item.createdAt || item.tanggalBerangkat)}</p>
                </div>
            </div>
            <button onclick="closeDetailModal()" class="w-8 h-8 rounded-xl bg-surface-border/50 text-slate-400 hover:text-white hover:bg-surface-border flex items-center justify-center transition-colors flex-shrink-0">
                <i class="ph-bold ph-x text-base"></i>
            </button>
        </div>

        <!-- Body Scrollable -->
        <div class="p-6 overflow-y-auto space-y-5 flex-1 text-sm">
            <!-- Profil Pegawai -->
            <div class="bg-surface/50 border border-surface-border rounded-2xl p-4 flex items-center justify-between flex-wrap gap-3">
                <div class="flex items-center gap-3.5">
                    <div class="w-11 h-11 rounded-2xl bg-brand-500/20 border border-brand-500/30 flex items-center justify-center text-brand-300 font-bold text-sm">
                        ${_buildAvatarFromName(item.nama)}
                    </div>
                    <div>
                        <p class="font-bold text-white text-sm">${item.nama}</p>
                        <p class="text-xs text-slate-400 font-mono">${item.nomorID} · ${item.kodeDepartemen}</p>
                    </div>
                </div>
                <div class="text-right">
                    <span class="text-[10px] text-slate-400 block uppercase font-bold tracking-wider">Jabatan</span>
                    <span class="text-xs font-semibold text-slate-200">${item.kodeJabatan || 'Staf'}</span>
                </div>
            </div>

            <!-- Rincian Rute & Jadwal -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-2.5">
                    <p class="flex items-center gap-1.5 text-brand-400 font-bold text-xs uppercase tracking-wider">
                        <i class="ph-bold ph-map-pin"></i> Rute &amp; Tujuan
                    </p>
                    <p class="text-sm font-bold text-white">${item.kotaTujuan} <span class="text-xs text-slate-400 font-normal">(${item.areaTujuan || 'Jawa'})</span></p>
                    <div>
                        <span class="text-[11px] text-slate-400 block font-semibold">Keperluan / Maksud Tugas</span>
                        <p class="text-xs text-slate-200 mt-0.5 leading-relaxed">${item.maksudTujuan || '-'}</p>
                    </div>
                    <div>
                        <span class="text-[11px] text-slate-400 block font-semibold">Hak Kriteria Fasilitas</span>
                        <p class="text-xs text-brand-300 font-medium mt-0.5">${item.kriteriaFasilitas || '1. < 200 KM & < 8 Jam'}</p>
                    </div>
                </div>

                <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-2.5">
                    <p class="flex items-center gap-1.5 text-brand-400 font-bold text-xs uppercase tracking-wider">
                        <i class="ph-bold ph-calendar-blank"></i> Jadwal Perjalanan
                    </p>
                    <div>
                        <span class="text-[11px] text-slate-400 block font-semibold">Periode Tanggal</span>
                        <p class="text-xs font-bold text-white mt-0.5">${formatDate(item.tanggalBerangkat)} s.d. ${formatDate(item.tanggalKembali)}</p>
                    </div>
                    <div class="grid grid-cols-3 gap-2 pt-1 border-t border-surface-border/50 text-xs">
                        <div>
                            <span class="text-[10px] text-slate-400 block">Berangkat</span>
                            <span class="font-semibold text-slate-200">${item.jamBerangkat || '08:00'}</span>
                        </div>
                        <div>
                            <span class="text-[10px] text-slate-400 block">Kembali</span>
                            <span class="font-semibold text-slate-200">${item.jamKembali || '17:00'}</span>
                        </div>
                        <div>
                            <span class="text-[10px] text-slate-400 block">Durasi</span>
                            <span class="font-bold text-brand-300 font-mono">${dur} Hari</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Rincian Biaya -->
            <div class="bg-surface/40 border border-surface-border rounded-2xl p-4 space-y-2">
                <div class="flex items-center justify-between">
                    <span class="text-xs font-bold text-slate-300 uppercase tracking-wider">Total Estimasi Anggaran (CA)</span>
                    <span class="text-emerald-400 font-mono font-extrabold text-base">${formatRupiah(item.estimasiBiaya)}</span>
                </div>
                <div class="grid grid-cols-3 gap-2 pt-2 border-t border-surface-border/50 text-[11px] text-slate-400">
                    <div>Transport: <strong class="text-slate-200 font-mono">${formatRupiah(item.biayaTransport)}</strong></div>
                    <div>Penginapan: <strong class="text-slate-200 font-mono">${formatRupiah(item.biayaPenginapan)}</strong></div>
                    <div>Uang Harian: <strong class="text-slate-200 font-mono">${formatRupiah(item.uangHarian)}</strong></div>
                </div>
            </div>

            <!-- Dokumen Lampiran Info jika ada -->
            ${(() => {
                let lampiranObj = null;
                try {
                    const storedLampiran = localStorage.getItem('sipadin_lampiran_' + item.noPengajuan);
                    if (storedLampiran) lampiranObj = JSON.parse(storedLampiran);
                } catch (e) {}

                let lampiranName = lampiranObj ? lampiranObj.name : '';
                if (!lampiranName && item.maksudTujuan && item.maksudTujuan.includes('[Lampiran:')) {
                    const match = item.maksudTujuan.match(/\[Lampiran:\s*([^\]]+)\]/);
                    if (match) lampiranName = match[1];
                }

                if (!lampiranName) return '';
                return `
                <div class="bg-surface/40 border border-teal-500/20 rounded-2xl p-3.5 flex items-center justify-between gap-3">
                    <div class="flex items-center gap-2.5 overflow-hidden">
                        <div class="w-8 h-8 rounded-xl bg-teal-500/20 text-teal-300 flex items-center justify-center flex-shrink-0">
                            <i class="ph-bold ph-paperclip text-base"></i>
                        </div>
                        <div class="truncate">
                            <p class="text-[10px] text-teal-300/80 font-bold uppercase tracking-wider">Dokumen Pendukung / Lampiran</p>
                            <p class="text-xs font-bold text-white truncate mt-0.5">${lampiranName} ${lampiranObj && lampiranObj.size ? `<span class="text-[10px] text-slate-400 font-normal">(${lampiranObj.size})</span>` : ''}</p>
                        </div>
                    </div>
                    ${lampiranObj && lampiranObj.base64 ? `
                    <a href="${lampiranObj.base64}" download="${lampiranObj.name}" target="_blank" class="px-3 py-1.5 rounded-xl bg-teal-600/20 hover:bg-teal-600/30 text-teal-300 hover:text-white text-xs font-semibold flex items-center gap-1.5 transition-all flex-shrink-0 border border-teal-500/30">
                        <i class="ph-bold ph-download-simple"></i> Unduh File
                    </a>` : `
                    <span class="px-2.5 py-1 rounded-lg bg-white/5 text-slate-400 text-[11px] font-medium">Terlampir</span>
                    `}
                </div>`;
            })()}

            <!-- Surat Tugas Resmi Info jika sudah terbit -->
            ${item.nomorSuratTugas ? `
            <div class="bg-indigo-500/10 border border-indigo-500/30 rounded-2xl p-3.5 flex items-center justify-between gap-3">
                <div class="flex items-center gap-2.5">
                    <i class="ph-fill ph-file-text text-indigo-400 text-2xl flex-shrink-0"></i>
                    <div>
                        <p class="text-[11px] text-indigo-300 font-bold uppercase tracking-wider">Surat Tugas Resmi</p>
                        <p class="text-xs font-mono font-bold text-white mt-0.5">${item.nomorSuratTugas}</p>
                    </div>
                </div>
                <button onclick="printSuratTugas('${item.noPengajuan}')" class="px-3 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold transition-all flex items-center gap-1 shadow-sm">
                    <i class="ph-bold ph-printer"></i> Cetak
                </button>
            </div>` : ''}

            <!-- Riwayat 6 Tahap Workflow -->
            <div>
                <p class="text-xs font-bold text-slate-300 uppercase tracking-wider mb-2.5 flex items-center gap-1.5">
                    <i class="ph-bold ph-git-branch text-brand-400"></i> Riwayat Persetujuan &amp; Pipeline 6 Tahap
                </p>
                <div class="p-4 rounded-2xl bg-surface/50 border border-surface-border">
                    ${timelineHtml}
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div class="p-4 border-t border-surface-border bg-surface/90 flex items-center justify-between gap-3 flex-shrink-0">
            <button onclick="closeDetailModal()" class="px-4 py-2 rounded-xl border border-surface-border text-slate-300 hover:text-white hover:bg-white/5 text-xs font-semibold transition-all">
                Tutup (Esc)
            </button>
            <div class="flex items-center gap-2">
                ${['APPROVED_ATASAN','CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED'].includes(item.statusFlow) ? `
                    <button onclick="printSPPDSlip('${item.noPengajuan}')" class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-brand-gradient hover:opacity-90 text-white text-xs font-bold transition-all shadow-glow">
                        <i class="ph-bold ph-printer"></i> Cetak Dokumen SPPD
                    </button>
                ` : ''}
            </div>
        </div>
    </div>`;

    modal.classList.remove('hidden');
    modal.onclick = (e) => {
        if (e.target === modal) closeDetailModal();
    };
}

function closeDetailModal() {
    ['modalDetailPengajuan', 'modalDetail', 'modal-detail-pengajuan'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.remove();
    });
}
window.showDetail = showDetail;
window.openDetailModal = showDetail;
window.closeDetailModal = closeDetailModal;
window.closeDetail = closeDetailModal;

// Close modal on Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        closeDetailModal();
    }
});

// ============================================================
// 3. TABLE EVENT DELEGATION SYSTEM
// ============================================================
function initTableEventDelegation() {
    document.addEventListener('click', async (e) => {
        const btn = e.target.closest('button[data-action], [data-action]');
        if (!btn) return;

        const action = btn.getAttribute('data-action');
        const id     = btn.getAttribute('data-id') || btn.getAttribute('data-no');
        if (!action) return;

        e.preventDefault();
        e.stopPropagation();

        switch (action) {
            case 'detail':
                showDetail(id);
                break;
            case 'approve':
                if (typeof window.openReviewModal === 'function') {
                    window.openReviewModal(id);
                } else {
                    await approveByAtasan(id);
                    if (typeof renderAtasanTable === 'function') renderAtasanTable();
                }
                break;
            case 'reject':
                if (typeof window.openReviewModal === 'function') {
                    window.openReviewModal(id);
                } else if (typeof window.openRejectModal === 'function') {
                    window.openRejectModal(id, 'atasan');
                } else {
                    const catatan = prompt('Masukkan alasan penolakan:');
                    if (catatan) {
                        await rejectByAtasan(id, catatan);
                        if (typeof renderAtasanTable === 'function') renderAtasanTable();
                    }
                }
                break;
            case 'hrd-control':
                if (typeof window.openVerifikasiModal === 'function') {
                    window.openVerifikasiModal(id);
                } else {
                    await verifyByHRGA(id);
                    if (typeof renderHRGATable === 'function') renderHRGATable();
                }
                break;
            case 'hrd-reject':
                if (typeof window.openRejectModal === 'function') {
                    window.openRejectModal(id, 'hrga');
                } else {
                    const catatan = prompt('Masukkan alasan penolakan HRGA:');
                    if (catatan) {
                        await rejectByHRGA(id, catatan);
                        if (typeof renderHRGATable === 'function') renderHRGATable();
                    }
                }
                break;
            case 'wa-direksi':
                if (typeof window.openWAModal === 'function') {
                    window.openWAModal(id);
                } else {
                    await sendWADireksi(id);
                    if (typeof renderHRGATable === 'function') renderHRGATable();
                }
                break;
            case 'confirm-direksi':
                if (typeof window.openDireksiModal === 'function') {
                    window.openDireksiModal(id, 'approve');
                } else {
                    await confirmDireksi(id);
                    if (typeof renderHRGATable === 'function') renderHRGATable();
                }
                break;
            case 'direksi-reject':
                if (typeof window.openDireksiModal === 'function') {
                    window.openDireksiModal(id, 'reject');
                } else {
                    const catatan = prompt('Masukkan alasan penolakan Direksi:');
                    if (catatan) {
                        await rejectByDireksi(id, catatan);
                        if (typeof renderHRGATable === 'function') renderHRGATable();
                    }
                }
                break;
            case 'surat-tugas':
                if (typeof window.handleIssueSuratTugas === 'function') {
                    window.handleIssueSuratTugas(id);
                } else {
                    await issueSuratTugas(id);
                    if (typeof renderHRGATable === 'function') renderHRGATable();
                }
                break;
            case 'print-st':
                printSuratTugas(id);
                break;
            case 'print-sppd':
                printSPPDSlip(id);
                break;
        }
    });
}

// ============================================================
// 4. STATUS TAB FILTER & STATS CALCULATION
// ============================================================
function filterDataByTab(statusKey, list) {
    const data = list || window.cachedPengajuanList || [];
    const key = String(statusKey || '').toLowerCase();

    if (key === 'semua' || key === 'all' || key === '') {
        return data;
    }
    if (key === 'menunggu' || key === 'pending') {
        return data.filter(p => {
            const s = p.statusFlow || p.status;
            return s === 'SUBMITTED' || s === 'APPROVED_ATASAN' || s === 'pending' || s === 'approved';
        });
    }
    if (key === 'disetujui' || key === 'approved') {
        return data.filter(p => {
            const s = p.statusFlow || p.status;
            return ['CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED','verified_hrga','notified_direksi'].includes(s);
        });
    }
    if (key === 'ditolak' || key === 'rejected') {
        return data.filter(p => {
            const s = p.statusFlow || p.status;
            return s === 'REJECTED' || s === 'rejected';
        });
    }
    return data.filter(p => (p.statusFlow || p.status) === statusKey);
}
window.filterDataByTab = filterDataByTab;

function calculateStatCards(list) {
    const data = list || window.cachedPengajuanList || [];
    let total = data.length;
    let pending = 0;
    let approved = 0;
    let rejected = 0;
    let urgent = 0;
    let totalBudget = 0;

    data.forEach(p => {
        totalBudget += (p.estimasiBiaya || p.totalBiaya || 0);
        const s = p.statusFlow || p.status;
        if (s === 'SUBMITTED' || s === 'pending') {
            pending++;
            const diff = daysUntil(p.tanggalBerangkat || p.mulai);
            if (diff >= 0 && diff <= 3) urgent++;
        } else if (s === 'APPROVED_ATASAN' || s === 'approved') {
            pending++;
        } else if (s === 'REJECTED' || s === 'rejected') {
            rejected++;
        } else if (['CONTROLLED_HRD','WA_SENT_DIREKSI','CONFIRMED_DIREKSI','SURAT_TUGAS_ISSUED','verified_hrga','notified_direksi'].includes(s)) {
            approved++;
        }
    });

    return { total, pending, approved, rejected, urgent, totalBudget };
}
window.calculateStatCards = calculateStatCards;

// ============================================================
// SURAT TUGAS & SPPD PRINT VIEW GENERATOR
// ============================================================
async function printSuratTugas(noPengajuan) {
    let item;
    try {
        item = await apiGetPengajuanByID(noPengajuan);
    } catch (err) {
        showToast('error', 'Error', 'Data tidak ditemukan: ' + err.message);
        return;
    }

    // Clean up any existing print container
    const old = document.getElementById('sppd-print-area');
    if (old) old.remove();

    const dur = calculateDuration(item.tanggalBerangkat, item.tanggalKembali);
    const printArea = document.createElement('div');
    printArea.id = 'sppd-print-area';

    printArea.innerHTML = `
    <div style="font-family:'Times New Roman',serif;max-width:800px;margin:0 auto;color:#111;padding:20px;line-height:1.6;background:#fff;">
        <div style="display:flex;align-items:center;border-bottom:3px double #000;padding-bottom:12px;margin-bottom:20px;">
            <div style="flex:1;text-align:center;">
                <h2 style="margin:0;font-size:18pt;text-transform:uppercase;font-weight:bold;letter-spacing:1px;color:#000;">PT ADIPRIMA SURAPRINTA</h2>
                <p style="margin:3px 0;font-size:10pt;color:#333;">Kawasan Industri Driyorejo, Gresik, Jawa Timur — Indonesia</p>
                <p style="margin:0;font-size:9pt;color:#555;">Telp: (031) 7507888 | Email: corporate@adiprima.co.id | Website: www.adiprima.co.id</p>
            </div>
        </div>
        <div style="text-align:center;margin-bottom:25px;">
            <h3 style="margin:0;font-size:14pt;text-decoration:underline;font-weight:bold;text-transform:uppercase;color:#000;">SURAT TUGAS PERJALANAN DINAS</h3>
            <p style="margin:4px 0 0 0;font-size:11pt;color:#222;">Nomor: <strong>${item.nomorSuratTugas || item.noPengajuan}</strong></p>
        </div>
        <table style="width:100%;border-collapse:collapse;margin-bottom:25px;font-size:11pt;color:#000;" border="1" cellpadding="8">
            <tr><td style="width:5%;text-align:center;font-weight:bold;">1</td><td style="width:35%;">Nama / Nomor ID Pegawai</td><td><strong>${item.nama}</strong> / ${item.nomorID}</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">2</td><td>Departemen / Jabatan</td><td>${item.kodeDepartemen} / ${item.kodeJabatan}</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">3</td><td>Tujuan Perjalanan Dinas</td><td><strong>${item.kotaTujuan}</strong> (${item.areaTujuan || 'Jawa'})</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">4</td><td>Maksud / Keperluan</td><td>${item.maksudTujuan}</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">5</td><td>Kriteria Fasilitas</td><td>${item.kriteriaFasilitas}</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">6</td>
                <td>a. Tanggal Berangkat<br>b. Tanggal Kembali<br>c. Lama Perjalanan</td>
                <td>a. ${formatDateFull(item.tanggalBerangkat)}<br>b. ${formatDateFull(item.tanggalKembali)}<br>c. ${dur} Hari</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">7</td><td>Estimasi Biaya</td><td><strong>${formatRupiah(item.estimasiBiaya)}</strong></td></tr>
            <tr><td style="text-align:center;font-weight:bold;">8</td><td>Catatan HRD/HRGA</td><td>${item.catatanHRD || 'Sesuai kebijakan perusahaan.'}</td></tr>
            <tr><td style="text-align:center;font-weight:bold;">9</td><td>Persetujuan Direksi</td><td>${item.konfirmasiDireksiNote || 'Disetujui'}</td></tr>
        </table>
        <div style="display:flex;justify-content:space-between;margin-top:40px;font-size:11pt;color:#000;">
            <div style="text-align:center;width:30%;"><p style="margin-bottom:60px;">Pegawai yang Melaksanakan Tugas,</p><p style="margin:0;font-weight:bold;text-decoration:underline;">${item.nama}</p><p style="margin:0;font-size:9pt;">NIP: ${item.nomorID}</p></div>
            <div style="text-align:center;width:30%;"><p style="margin:0 0 60px 0;">Dikeluarkan di: Gresik<br>Pada tanggal: ${formatDateFull(new Date().toISOString())}</p><p style="margin:0;font-weight:bold;text-decoration:underline;">${item.atasanNama || 'Atasan Langsung'}</p><p style="margin:0;font-size:9pt;">Atasan Langsung</p></div>
            <div style="text-align:center;width:30%;"><p style="margin:0 0 60px 0;">Mengetahui,<br>HRD/HRGA PT Adiprima Suraprinta</p><p style="margin:0;font-weight:bold;text-decoration:underline;">Dewi Rahayu</p><p style="margin:0;font-size:9pt;">HRGA Supervisor</p></div>
        </div>
    </div>`;

    document.body.appendChild(printArea);

    const cleanup = () => {
        const el = document.getElementById('sppd-print-area');
        if (el) el.remove();
        window.removeEventListener('afterprint', cleanup);
    };
    window.addEventListener('afterprint', cleanup, { once: true });

    setTimeout(() => {
        window.print();
        setTimeout(cleanup, 1000);
    }, 150);
}
const printSPPDSlip = printSuratTugas;
window.printSuratTugas = printSuratTugas;
window.printSPPDSlip = printSuratTugas;

// ============================================================
// PENGAJUAN FORM — Populate Dropdowns from Master API
// ============================================================
async function populateMasterDropdowns() {
    try {
        const [tarif, fasilitas] = await Promise.all([apiGetMasterTarif(), apiGetMasterFasilitas()]);

        const areaSelect = document.getElementById('area-tujuan');
        if (areaSelect) {
            const areas = [...new Set(tarif.map(t => t.area))];
            areas.forEach(area => {
                const opt = document.createElement('option');
                opt.value = area;
                opt.textContent = area;
                areaSelect.appendChild(opt);
            });
        }

        const kriteriaSelect = document.getElementById('kriteria-fasilitas');
        if (kriteriaSelect) {
            fasilitas.forEach(f => {
                const opt = document.createElement('option');
                opt.value = f.kriteria;
                opt.textContent = f.kriteria;
                kriteriaSelect.appendChild(opt);
            });
        }
    } catch (err) {
        console.warn('Gagal memuat data master:', err.message);
    }
}

// ============================================================
// AUTO INIT ON LOAD (Profile Sync & Event Delegation)
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

    // Run Profile Synchronization
    syncUserProfile();

    // Initialize Global Table Event Delegation
    initTableEventDelegation();

    if (document.getElementById('area-tujuan') || document.getElementById('kriteria-fasilitas')) {
        populateMasterDropdowns();
    }
});
