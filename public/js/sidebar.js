/* public/js/sidebar.js — Komponen sidebar yang di-render secara dinamis */

function renderSidebar(activePage) {
  const user = getCurrentUser();
  const role = user.role_sistem || '';

  // Definisi menu sesuai PRD: semua menu tampil di sidebar, akses dikontrol oleh role
  const menus = [
    {
      label: 'Dashboard Saya',
      href: '/dashboard.html',
      icon: 'ph-house',
      page: 'dashboard',
      roles: ['Karyawan', 'Manager', 'HRGA', 'Finance', 'Direksi'],
    },
    {
      label: 'Antrian Manager',
      href: '/manager.html',
      icon: 'ph-user-check',
      page: 'manager',
      roles: ['Manager', 'HRGA', 'Finance', 'Direksi'],
    },
    {
      label: 'Verifikasi HRGA',
      href: '/hrga.html',
      icon: 'ph-shield-check',
      page: 'hrga',
      roles: ['HRGA', 'Manager', 'Finance', 'Direksi'],
    },
    {
      label: 'Persetujuan Finance',
      href: '/finance.html',
      icon: 'ph-currency-dollar',
      page: 'finance',
      roles: ['Finance', 'Manager', 'HRGA', 'Direksi'],
    },
    {
      label: 'Laporan & Rekap',
      href: '/laporan.html',
      icon: 'ph-chart-bar',
      page: 'laporan',
      roles: ['Direksi', 'Manager', 'HRGA', 'Finance', 'Karyawan'],
    },
    {
      label: 'Master Tarif',
      href: '/tarif.html',
      icon: 'ph-table',
      page: 'tarif',
      roles: ['HRGA', 'Manager', 'Finance', 'Direksi'],
    },
    {
      label: 'Manajemen User',
      href: '/users.html',
      icon: 'ph-users-three',
      page: 'users',
      roles: ['HRGA'],
    },
  ];

  const sidebarHTML = `
    <aside id="sidebar" class="fixed left-0 top-0 h-full w-64 glass border-r border-white/5 z-40 flex flex-col transition-transform duration-300 lg:translate-x-0 -translate-x-full" id="sidebar">
      <!-- Logo -->
      <div class="p-6 border-b border-white/5">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center flex-shrink-0">
            <i class="ph-bold ph-airplane-takeoff text-indigo-400 text-base"></i>
          </div>
          <div>
            <h2 class="font-bold text-white text-sm">SIPADIN</h2>
            <p class="text-slate-500 text-xs">Perjalanan Dinas</p>
          </div>
        </div>
      </div>

      <!-- User Info -->
      <div class="px-4 py-3 border-b border-white/5">
        <div class="flex items-center gap-3 p-3 rounded-xl bg-white/3">
          <div class="w-8 h-8 rounded-lg bg-indigo-500/30 flex items-center justify-center flex-shrink-0">
            <i class="ph-bold ph-user text-indigo-300 text-sm"></i>
          </div>
          <div class="min-w-0">
            <p class="text-white text-xs font-semibold truncate">${user.nama_lengkap || 'User'}</p>
            <p class="text-slate-500 text-xs">${role}</p>
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
               class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group
                      ${isActive
                        ? 'bg-indigo-600/20 text-indigo-300 border border-indigo-500/20'
                        : 'text-slate-400 hover:bg-white/5 hover:text-slate-200'
                      }">
              <i class="ph-bold ${m.icon} text-base flex-shrink-0 ${isActive ? 'text-indigo-400' : 'group-hover:text-slate-300'}"></i>
              <span class="flex-1 truncate">${m.label}</span>
              ${isViewOnly ? '<span class="text-xs text-slate-600 font-normal">View</span>' : ''}
              ${isActive ? '<i class="ph-bold ph-caret-right text-xs text-indigo-400/50"></i>' : ''}
            </a>
          `;
        }).join('')}
      </nav>

      <!-- Logout -->
      <div class="p-3 border-t border-white/5">
        <button onclick="logout()"
          class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-slate-400 hover:bg-red-500/10 hover:text-red-400 transition-all duration-200 group">
          <i class="ph-bold ph-sign-out text-base group-hover:text-red-400"></i>
          Keluar
        </button>
      </div>
    </aside>

    <!-- Mobile Overlay -->
    <div id="sidebarOverlay" class="fixed inset-0 bg-black/60 z-30 lg:hidden hidden" onclick="toggleSidebar()"></div>

    <!-- Mobile Toggle -->
    <button onclick="toggleSidebar()"
      class="lg:hidden fixed top-4 left-4 z-50 glass border border-white/10 rounded-xl p-2.5 text-slate-300">
      <i class="ph-bold ph-list text-xl" id="sidebarToggleIcon"></i>
    </button>
  `;

  // Inject sidebar ke body
  const container = document.getElementById('sidebarContainer');
  if (container) container.innerHTML = sidebarHTML;
}

// Tentukan apakah role punya Full Access ke halaman tertentu
function isFullAccessPage(page, role) {
  const fullAccessMap = {
    'dashboard' : ['Karyawan', 'Manager', 'HRGA', 'Finance'],
    'manager'   : ['Manager'],
    'hrga'      : ['HRGA'],
    'finance'   : ['Finance'],
    'tarif'     : ['HRGA'],
    'users'     : ['HRGA'],
    'laporan'   : ['Direksi'],
  };
  return (fullAccessMap[page] || []).includes(role);
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
    el.className = 'fixed bottom-6 right-6 z-[9999] px-5 py-3 rounded-xl glass border text-sm font-medium flex items-center gap-2.5 min-w-64 transition-all duration-400';
    el.style.transform = 'translateY(80px)';
    el.style.opacity = '0';
    document.body.appendChild(el);
  }
  const colors = {
    success: { border: 'rgba(34,197,94,0.4)',   icon: 'ph-check-circle',   txt: 'text-emerald-400' },
    error:   { border: 'rgba(239,68,68,0.4)',   icon: 'ph-x-circle',       txt: 'text-red-400' },
    info:    { border: 'rgba(99,102,241,0.4)',  icon: 'ph-info',           txt: 'text-indigo-400' },
    warning: { border: 'rgba(234,179,8,0.4)',   icon: 'ph-warning-circle', txt: 'text-yellow-400' },
  };
  const c = colors[type] || colors.info;
  el.style.borderColor = c.border;
  el.innerHTML = `<i class="ph-bold ${c.icon} text-lg ${c.txt}"></i><span class="text-slate-200">${msg}</span>`;
  setTimeout(() => { el.style.transform = 'translateY(0)'; el.style.opacity = '1'; }, 10);
  setTimeout(() => { el.style.transform = 'translateY(80px)'; el.style.opacity = '0'; }, 3500);
}

// ── Modal Konfirmasi ──────────────────────────────────────────────────────────
function showConfirm(msg, onConfirm) {
  const overlay = document.createElement('div');
  overlay.className = 'fixed inset-0 bg-black/70 z-[9999] flex items-center justify-center p-4';
  overlay.innerHTML = `
    <div class="glass border border-white/10 rounded-2xl p-6 w-full max-w-sm">
      <div class="flex items-start gap-4">
        <div class="w-10 h-10 rounded-xl bg-amber-500/20 flex items-center justify-center flex-shrink-0">
          <i class="ph-bold ph-warning text-amber-400 text-xl"></i>
        </div>
        <div>
          <h3 class="text-white font-semibold mb-1">Konfirmasi</h3>
          <p class="text-slate-400 text-sm">${msg}</p>
        </div>
      </div>
      <div class="flex gap-3 mt-6">
        <button id="cancelBtn" class="flex-1 py-2.5 rounded-xl border border-white/10 text-slate-300 text-sm font-medium hover:bg-white/5 transition">Batal</button>
        <button id="confirmBtn" class="flex-1 py-2.5 rounded-xl bg-red-600 hover:bg-red-500 text-white text-sm font-medium transition">Ya, Lanjutkan</button>
      </div>
    </div>`;
  document.body.appendChild(overlay);
  overlay.querySelector('#cancelBtn').onclick  = () => overlay.remove();
  overlay.querySelector('#confirmBtn').onclick = () => { overlay.remove(); onConfirm(); };
}
