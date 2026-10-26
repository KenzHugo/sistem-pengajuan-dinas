/* public/js/sidebar.js — Komponen sidebar & styling tombol responsif global */

// ── Global Button & Interactive Responsiveness Enhancement ────────────────────
(function injectGlobalButtonStyles() {
  if (document.getElementById('sipadin-responsive-btn-styles')) return;
  const style = document.createElement('style');
  style.id = 'sipadin-responsive-btn-styles';
  style.textContent = `
    /* Universal Button & Action Responsiveness */
    button, [role="button"], a[class*="btn-"], .btn-approve, .btn-reject, .btn-send, .btn-action {
      touch-action: manipulation;
      -webkit-tap-highlight-color: transparent;
      user-select: none;
      -webkit-user-select: none;
    }
    button:not(:disabled), [role="button"]:not(.disabled-btn), a[class*="btn-"]:not(.disabled-btn) {
      cursor: pointer;
    }
    button:active:not(:disabled), [role="button"]:active:not(.disabled-btn), a[class*="btn-"]:active:not(.disabled-btn) {
      transform: scale(0.95) !important;
      transition: transform 0.1s ease-out, box-shadow 0.1s ease-out !important;
    }
    button:disabled, .disabled-btn {
      opacity: 0.45 !important;
      cursor: not-allowed !important;
      pointer-events: none !important;
      transform: none !important;
    }
    /* Touch target sizing on mobile */
    @media (max-width: 640px) {
      button, .btn-approve, .btn-reject, .btn-send {
        min-height: 38px;
      }
      .table-action-btn {
        min-height: 34px;
        min-width: 34px;
      }
    }
  `;
  document.head.appendChild(style);
})();

function renderSidebar(activePage) {
  const user = getCurrentUser();
  const role = user.role_sistem || '';
  const BASE_PATH = window.location.pathname.startsWith('/perdin') ? '/perdin' : '';

  // Definisi menu sesuai PRD: semua menu tampil di sidebar, akses dikontrol oleh role
  const menus = [
    {
      label: 'Dashboard Saya',
      href: `${BASE_PATH}/dashboard.html`,
      icon: 'ph-house',
      page: 'dashboard',
      roles: ['Karyawan', 'Manager', 'HRGA', 'Finance', 'Direksi'],
    },
    {
      label: 'Verifikasi Manager',
      href: `${BASE_PATH}/manager.html`,
      icon: 'ph-user-check',
      page: 'manager',
      roles: ['Manager', 'HRGA', 'Finance', 'Direksi'],
    },
    {
      label: 'Verifikasi HRGA',
      href: `${BASE_PATH}/hrga.html`,
      icon: 'ph-shield-check',
      page: 'hrga',
      roles: ['HRGA', 'Manager', 'Finance', 'Direksi'],
    },
    {
      label: 'Verifikasi Finance',
      href: `${BASE_PATH}/finance.html`,
      icon: 'ph-currency-dollar',
      page: 'finance',
      roles: ['Finance', 'Manager', 'HRGA', 'Direksi'],
    },
    {
      label: 'Laporan & Rekap',
      href: `${BASE_PATH}/laporan.html`,
      icon: 'ph-chart-bar',
      page: 'laporan',
      roles: ['Direksi', 'Manager', 'HRGA', 'Finance', 'Karyawan'],
    },
    {
      label: 'Direktori Karyawan',
      href: `${BASE_PATH}/karyawan.html`,
      icon: 'ph-users',
      page: 'karyawan',
      roles: ['Manager', 'HRGA', 'Finance', 'Direksi'],
    },
    {
      label: 'Master Tarif',
      href: `${BASE_PATH}/tarif.html`,
      icon: 'ph-table',
      page: 'tarif',
      roles: ['HRGA', 'Manager', 'Finance', 'Direksi'],
    },
    {
      label: 'Manajemen User',
      href: `${BASE_PATH}/users.html`,
      icon: 'ph-users-three',
      page: 'users',
      roles: ['HRGA'],
    },
  ];

  const sidebarHTML = `
    <aside id="sidebar" class="fixed left-0 top-0 h-full w-64 bg-white dark:bg-[#151c2c] border-r border-slate-200 dark:border-white/5 z-40 flex flex-col transition-all duration-300 lg:translate-x-0 -translate-x-full shadow-lg lg:shadow-none">
      <!-- Logo Container (No white background wrapper, responsive and clear in both light & dark modes) -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-white/10 flex items-center justify-center transition-colors duration-300">
        <!-- Logo Light Mode (Original corporate brand colors) -->
        <img 
          src="${BASE_PATH}/logo-light.png" 
          alt="Logo PT Adiprima Suraprinta" 
          class="w-full max-w-[190px] h-auto object-contain block dark:hidden transition-all duration-300" 
        />
        <!-- Logo Dark Mode (Crisp white text with brand icon) -->
        <img 
          src="${BASE_PATH}/logo-dark.png" 
          alt="Logo PT Adiprima Suraprinta" 
          class="w-full max-w-[190px] h-auto object-contain hidden dark:block transition-all duration-300" 
        />
      </div>

      <!-- User Info -->
      <div class="px-4 py-3 border-b border-slate-200 dark:border-white/10 transition-colors duration-300">
        <div class="flex items-center gap-3 p-3 rounded-xl bg-slate-100/80 dark:bg-white/5 transition-colors duration-300">
          <div class="w-8 h-8 rounded-lg bg-emerald-500/20 text-emerald-600 dark:bg-emerald-500/20 dark:text-emerald-300 flex items-center justify-center flex-shrink-0">
            <i class="ph-bold ph-user text-sm"></i>
          </div>
          <div class="min-w-0">
            <p class="text-slate-800 dark:text-white text-xs font-semibold truncate">${user.nama_lengkap || 'User'}</p>
            <p class="text-slate-500 dark:text-slate-400 text-xs">${role}</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 overflow-y-auto p-3 space-y-1">
        ${menus.filter(m => m.roles.includes(role)).map(m => {
          const isActive = m.page === activePage;
          const isFullAccess = isFullAccessPage(m.page, role);
          const isViewOnly = m.roles.includes(role) && !isFullAccess && m.page !== 'dashboard';
          return `
            <a href="${m.href}"
               class="flex items-center gap-3 px-3.5 py-2.5 rounded-r-xl text-sm font-medium transition-all duration-300 group
                      ${isActive
                        ? 'bg-blue-50 text-blue-600 border-l-4 border-emerald-500 font-semibold shadow-sm dark:bg-emerald-500/10 dark:text-emerald-300 dark:border-l-4 dark:border-emerald-400'
                        : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900 border-l-4 border-transparent dark:text-slate-400 dark:hover:bg-white/5 dark:hover:text-white'
                      }">
              <i class="ph-bold ${m.icon} text-base flex-shrink-0 ${isActive ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-300'}"></i>
              <span class="flex-1 truncate">${m.label}</span>
              ${isViewOnly ? '<span class="text-xs text-slate-400 dark:text-slate-500 font-normal">View</span>' : ''}
              ${isActive ? '<i class="ph-bold ph-caret-right text-xs text-emerald-600 dark:text-emerald-400"></i>' : ''}
            </a>
          `;
        }).join('')}
      </nav>

      <!-- Logout -->
      <div class="p-3 border-t border-slate-200 dark:border-white/5 transition-colors duration-300">
        <button onclick="logout()"
          class="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-r-xl text-sm font-medium text-slate-600 dark:text-slate-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-500/10 dark:hover:text-red-400 cursor-pointer active:scale-95 transition-all duration-300 group border-l-4 border-transparent">
          <i class="ph-bold ph-sign-out text-base text-slate-400 group-hover:text-red-600 dark:group-hover:text-red-400"></i>
          Keluar
        </button>
      </div>
    </aside>

    <!-- Mobile Overlay -->
    <div id="sidebarOverlay" class="fixed inset-0 bg-black/60 z-30 lg:hidden hidden" onclick="toggleSidebar()"></div>

    <!-- Mobile Toggle -->
    <button onclick="toggleSidebar()"
      class="lg:hidden fixed top-4 left-4 z-50 bg-white/90 dark:bg-slate-800/90 backdrop-blur-md border border-slate-200 dark:border-white/10 rounded-xl p-2.5 text-slate-700 dark:text-slate-300 cursor-pointer active:scale-95 transition shadow-lg">
      <i class="ph-bold ph-list text-xl" id="sidebarToggleIcon"></i>
    </button>
  `;

  // Inject sidebar ke body
  const container = document.getElementById('sidebarContainer');
  if (container) container.innerHTML = sidebarHTML;
}

// Tentukan apakah role punya Full Access ke halaman tertentu
function isFullAccessPage(page, role) {
  const r = (role || '').trim().toLowerCase();
  const fullAccessMap = {
    'dashboard' : ['karyawan', 'manager', 'hrga', 'finance', 'direksi'],
    'manager'   : ['manager', 'direksi'],
    'hrga'      : ['hrga', 'direksi'],
    'finance'   : ['finance', 'direksi'],
    'karyawan'  : ['hrga', 'manager', 'direksi'],
    'tarif'     : ['hrga', 'direksi'],
    'users'     : ['hrga', 'direksi'],
    'laporan'   : ['direksi', 'manager', 'hrga', 'finance', 'karyawan'],
  };
  return (fullAccessMap[page] || []).includes(r);
}

function toggleSidebar() {
  const sidebar = document.getElementById('sidebar');
  const overlay = document.getElementById('sidebarOverlay');
  const icon    = document.getElementById('sidebarToggleIcon');
  if (sidebar.classList.contains('-translate-x-full')) {
    sidebar.classList.remove('-translate-x-full');
    overlay.classList.remove('hidden');
    icon.className = 'ph-bold ph-x text-xl';
  } else {
    sidebar.classList.add('-translate-x-full');
    overlay.classList.add('hidden');
    icon.className = 'ph-bold ph-list text-xl';
  }
}

// ── Toast global ─────────────────────────────────────────────────────────────
function showToast(msg, type = 'info') {
  let el = document.getElementById('globalToast');
  if (!el) {
    el = document.createElement('div');
    el.id = 'globalToast';
    el.className = 'fixed bottom-6 right-6 z-[9999] px-5 py-3 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-800 dark:text-white shadow-xl text-sm font-medium flex items-center gap-2.5 min-w-64 transition-all duration-400';
    el.style.transform = 'translateY(80px)';
    el.style.opacity = '0';
    document.body.appendChild(el);
  }
  const colors = {
    success: { border: 'rgba(34,197,94,0.4)',   icon: 'ph-check-circle',   txt: 'text-emerald-500' },
    error:   { border: 'rgba(239,68,68,0.4)',   icon: 'ph-x-circle',       txt: 'text-red-500' },
    info:    { border: 'rgba(99,102,241,0.4)',  icon: 'ph-info',           txt: 'text-indigo-500' },
    warning: { border: 'rgba(234,179,8,0.4)',   icon: 'ph-warning-circle', txt: 'text-amber-500' },
  };
  const c = colors[type] || colors.info;
  el.style.borderColor = c.border;
  el.innerHTML = `<i class="ph-bold ${c.icon} text-lg ${c.txt}"></i><span class="text-slate-800 dark:text-white font-medium">${msg}</span>`;
  setTimeout(() => { el.style.transform = 'translateY(0)'; el.style.opacity = '1'; }, 10);
  setTimeout(() => { el.style.transform = 'translateY(80px)'; el.style.opacity = '0'; }, 3500);
}

// ── Modal Konfirmasi ──────────────────────────────────────────────────────────
function showConfirm(msg, onConfirm) {
  const overlay = document.createElement('div');
  overlay.className = 'fixed inset-0 bg-black/70 z-[9999] flex items-center justify-center p-4';
  overlay.innerHTML = `
    <div class="bg-white dark:bg-[#1e293b] border border-slate-200 dark:border-white/10 rounded-2xl p-6 w-full max-w-sm shadow-2xl transition-colors duration-300">
      <div class="flex items-start gap-4">
        <div class="w-10 h-10 rounded-xl bg-amber-500/20 flex items-center justify-center flex-shrink-0">
          <i class="ph-bold ph-warning text-amber-500 text-xl"></i>
        </div>
        <div>
          <h3 class="text-slate-800 dark:text-white font-semibold mb-1">Konfirmasi</h3>
          <p class="text-slate-600 dark:text-slate-400 text-sm">${msg}</p>
        </div>
      </div>
      <div class="flex gap-3 mt-6">
        <button id="cancelBtn" class="flex-1 py-2.5 rounded-xl border border-slate-200 dark:border-white/10 text-slate-700 dark:text-slate-300 text-sm font-medium hover:bg-slate-100 dark:hover:bg-white/5 transition">Batal</button>
        <button id="confirmBtn" class="flex-1 py-2.5 rounded-xl bg-red-600 hover:bg-red-500 text-white text-sm font-medium transition shadow-lg shadow-red-600/20">Ya, Lanjutkan</button>
      </div>
    </div>`;
  document.body.appendChild(overlay);
  overlay.querySelector('#cancelBtn').onclick  = () => overlay.remove();
  overlay.querySelector('#confirmBtn').onclick = () => { overlay.remove(); onConfirm(); };
}

// ── Theme Manager (Light / Dark Mode) ─────────────────────────────────────────
(function setupTheme() {
  const THEME_KEY = 'sipadin_theme';
  const saved = localStorage.getItem(THEME_KEY);
  const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
  if (saved === 'dark' || (!saved && prefersDark)) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
})();

function toggleGlobalTheme() {
  const isDark = document.documentElement.classList.toggle('dark');
  localStorage.setItem('sipadin_theme', isDark ? 'dark' : 'light');
  syncThemeIcons();
}

function syncThemeIcons() {
  const isDark = document.documentElement.classList.contains('dark');
  const btns = document.querySelectorAll('#themeToggleBtn, .theme-toggle-btn');
  btns.forEach(btn => {
    btn.setAttribute('title', isDark ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap');
    const sun = btn.querySelector('.theme-icon-sun');
    const moon = btn.querySelector('.theme-icon-moon');
    if (sun && moon) {
      if (isDark) {
        sun.classList.remove('hidden');
        moon.classList.add('hidden');
      } else {
        sun.classList.add('hidden');
        moon.classList.remove('hidden');
      }
    }
  });
}

// Sinkronisasi antar-tab browser secara real-time
window.addEventListener('storage', (e) => {
  if (e.key === 'sipadin_theme') {
    const isDark = e.newValue === 'dark';
    if (isDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    syncThemeIcons();
  }
});

// Jalankan sinkronisasi ikon saat script dimuat & saat DOM siap
syncThemeIcons();
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', syncThemeIcons);
}

// Event delegation agar selalu aktif meski elemen di-render ulang
document.addEventListener('click', (e) => {
  const btn = e.target.closest('#themeToggleBtn, .theme-toggle-btn');
  if (btn) {
    e.preventDefault();
    toggleGlobalTheme();
  }
});

