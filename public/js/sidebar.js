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
  const isSpecial = user && user.is_special_access === true;
  const role = isSpecial ? 'Spesial Akses' : (user.role_sistem || '');
  const BASE_PATH = window.location.pathname.startsWith('/perdin') ? '/perdin' : '';

  // Definisi menu: jika user memiliki SpesialAkses, semua menu tampil di sidebar
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
      <!-- Logo Container -->
      <div class="px-5 py-4 border-b border-slate-200 dark:border-white/10 flex items-center justify-center transition-colors duration-300">
        <img 
          src="${BASE_PATH}/logo-light.png" 
          alt="Logo PT Adiprima Suraprinta" 
          class="w-full max-w-[190px] h-auto object-contain block dark:hidden transition-all duration-300" 
        />
        <img 
          src="${BASE_PATH}/logo-dark.png" 
          alt="Logo PT Adiprima Suraprinta" 
          class="w-full max-w-[190px] h-auto object-contain hidden dark:block transition-all duration-300" 
        />
      </div>

      <!-- User Info -->
      <div class="px-4 py-3 border-b border-slate-200 dark:border-white/10 transition-colors duration-300">
        <div class="flex items-center gap-3 p-3 rounded-xl bg-slate-100/80 dark:bg-white/5 transition-colors duration-300">
          <div class="w-8 h-8 rounded-lg ${isSpecial ? 'bg-amber-500/20 text-amber-600 dark:text-amber-400' : 'bg-emerald-500/20 text-emerald-600 dark:text-emerald-300'} flex items-center justify-center flex-shrink-0">
            <i class="ph-bold ${isSpecial ? 'ph-crown' : 'ph-user'} text-sm"></i>
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-slate-800 dark:text-white text-xs font-semibold truncate">${user.nama_lengkap || 'User'}</p>
            <div class="flex items-center gap-1.5 mt-0.5">
              ${isSpecial 
                ? '<span class="inline-flex items-center px-1.5 py-0.2 rounded text-[10px] font-bold bg-amber-500/20 text-amber-700 dark:text-amber-300 border border-amber-500/30">★ Spesial Akses</span>' 
                : `<p class="text-slate-500 dark:text-slate-400 text-xs">${role}</p>`}
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 overflow-y-auto p-3 space-y-1">
        ${menus.filter(m => isSpecial || m.roles.includes(user.role_sistem)).map(m => {
          const isActive = m.page === activePage;
          const isFullAccess = isSpecial || isFullAccessPage(m.page, user.role_sistem);
          const isViewOnly = !isSpecial && m.roles.includes(user.role_sistem) && !isFullAccess && m.page !== 'dashboard';
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
  const user = getCurrentUser();
  if (user && user.is_special_access === true) return true;

  const r = (role || '').trim().toLowerCase();
  const fullAccessMap = {
    'dashboard' : ['karyawan', 'manager', 'hrga', 'finance', 'direksi', 'superadmin'],
    'manager'   : ['manager', 'direksi', 'superadmin'],
    'hrga'      : ['hrga', 'direksi', 'superadmin'],
    'finance'   : ['finance', 'direksi', 'superadmin'],
    'karyawan'  : ['hrga', 'manager', 'direksi', 'superadmin'],
    'tarif'     : ['hrga', 'direksi', 'superadmin'],
    'users'     : ['hrga', 'direksi', 'superadmin'],
    'laporan'   : ['direksi', 'manager', 'hrga', 'finance', 'karyawan', 'superadmin'],
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
    el.className = 'fixed bottom-5 right-5 z-50 flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border text-sm font-medium transition-all duration-300 transform translate-y-2 opacity-0 pointer-events-none';
    document.body.appendChild(el);
  }
  const colors = {
    success: 'bg-emerald-50 border-emerald-200 text-emerald-800 dark:bg-emerald-950/80 dark:border-emerald-800 dark:text-emerald-200',
    error:   'bg-red-50 border-red-200 text-red-800 dark:bg-red-950/80 dark:border-red-800 dark:text-red-200',
    info:    'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-950/80 dark:border-blue-800 dark:text-blue-200',
  };
  const icons = {
    success: 'ph-check-circle text-emerald-600',
    error:   'ph-x-circle text-red-600',
    info:    'ph-info text-blue-600',
  };
  el.className = `fixed bottom-5 right-5 z-50 flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border text-sm font-medium transition-all duration-300 transform ${colors[type] || colors.info}`;
  el.innerHTML = `<i class="ph-bold ${icons[type] || icons.info} text-lg"></i><span>${msg}</span>`;
  setTimeout(() => el.classList.remove('translate-y-2', 'opacity-0', 'pointer-events-none'), 10);
  setTimeout(() => {
    el.classList.add('translate-y-2', 'opacity-0', 'pointer-events-none');
  }, 3500);
}
