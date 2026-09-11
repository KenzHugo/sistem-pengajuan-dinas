/* public/js/api.js — Utility API client terpusat */

const API_BASE = '/api';

// Ambil token dari localStorage
function getToken() {
  return localStorage.getItem('sppd_token');
}

// Ambil data user dari localStorage
function getCurrentUser() {
  return JSON.parse(localStorage.getItem('sppd_user') || '{}');
}

// Logout: hapus session dan redirect ke login
function logout() {
  localStorage.removeItem('sppd_token');
  localStorage.removeItem('sppd_user');
  window.location.href = '/index.html';
}

// Guard: redirect ke login jika belum login
function requireAuth() {
  if (!getToken()) {
    window.location.href = '/index.html';
    return false;
  }
  return true;
}

// Wrapper fetch dengan auth header otomatis
async function apiFetch(endpoint, options = {}) {
  const token = getToken();
  const headers = {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
    ...(options.headers || {})
  };

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers
  });

  if (res.status === 401 && endpoint !== '/auth/login') {
    logout();
    return null;
  }

  return res.json();
}

// ── API Methods ──────────────────────────────────────────────────────────────

const api = {
  // Auth
  login: (identifier, password) =>
    apiFetch('/auth/login', { method: 'POST', body: JSON.stringify({ identifier, password }) }),

  // SPPD
  getMyList: (params = {}) => {
    const q = new URLSearchParams(params).toString();
    return apiFetch(`/sppd/mine${q ? '?' + q : ''}`);
  },
  getAllList: (params = {}) => {
    const q = new URLSearchParams(params).toString();
    return apiFetch(`/sppd${q ? '?' + q : ''}`);
  },
  getQueue: (stage, params = {}) => {
    const q = new URLSearchParams({ stage, ...params }).toString();
    return apiFetch(`/sppd/queue?${q}`);
  },
  getSppdByID: (id) => apiFetch(`/sppd/${id}`),
  createSppd: (body) => apiFetch('/sppd', { method: 'POST', body: JSON.stringify(body) }),
  updateSppd: (id, body) => apiFetch(`/sppd/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  submitSppd: (id) => apiFetch(`/sppd/${id}/submit`, { method: 'POST' }),
  deleteSppd: (id) => apiFetch(`/sppd/${id}`, { method: 'DELETE' }),

  // Approval
  managerAction:  (id, body) => apiFetch(`/approval/${id}/manager`, { method: 'POST', body: JSON.stringify(body) }),
  hrgaAction:     (id, body) => apiFetch(`/approval/${id}/hrga`, { method: 'POST', body: JSON.stringify(body) }),
  financeAction:  (id, body) => apiFetch(`/approval/${id}/finance`, { method: 'POST', body: JSON.stringify(body) }),
  sendToDirector: (id) =>       apiFetch(`/approval/${id}/send`, { method: 'POST' }),

  // Tarif
  getAllTarif:  ()        => apiFetch('/tarif'),
  createTarif: (body)    => apiFetch('/tarif', { method: 'POST', body: JSON.stringify(body) }),
  updateTarif: (id, body) => apiFetch(`/tarif/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteTarif: (id)      => apiFetch(`/tarif/${id}`, { method: 'DELETE' }),

  // Users
  getAllUsers:  ()         => apiFetch('/users'),
  createUser:  (body)     => apiFetch('/users', { method: 'POST', body: JSON.stringify(body) }),
  updateUser:  (id, body) => apiFetch(`/users/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteUser:  (id)       => apiFetch(`/users/${id}`, { method: 'DELETE' }),
};

// ── Helpers ──────────────────────────────────────────────────────────────────

function formatRupiah(amount) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
}

function formatDate(dateStr) {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' });
}

function formatDateShort(dateStr) {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
}

const STATUS_CONFIG = {
  'draft':           { label: 'Draft',           color: 'bg-slate-500/20 text-slate-300  border-slate-500/30',  icon: 'ph-file-dashed' },
  'pending_manager': { label: 'Menunggu Manager', color: 'bg-amber-500/20  text-amber-300  border-amber-500/30',  icon: 'ph-clock' },
  'pending_hrga':    { label: 'Menunggu HRGA',    color: 'bg-blue-500/20   text-blue-300   border-blue-500/30',   icon: 'ph-clock' },
  'pending_finance': { label: 'Menunggu Finance', color: 'bg-violet-500/20 text-violet-300 border-violet-500/30', icon: 'ph-clock' },
  'ready_to_send':   { label: 'Siap Dikirim',     color: 'bg-cyan-500/20   text-cyan-300   border-cyan-500/30',   icon: 'ph-paper-plane-right' },
  'approved':        { label: 'Disetujui',         color: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30', icon: 'ph-check-circle' },
  'rejected':        { label: 'Ditolak',           color: 'bg-red-500/20    text-red-300    border-red-500/30',    icon: 'ph-x-circle' },
};

function statusBadge(status) {
  const cfg = STATUS_CONFIG[status] || { label: status, color: 'bg-slate-500/20 text-slate-300 border-slate-500/30', icon: 'ph-question' };
  return `<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border text-xs font-medium ${cfg.color}">
            <i class="ph-bold ${cfg.icon}"></i>${cfg.label}
          </span>`;
}
