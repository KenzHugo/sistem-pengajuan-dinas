/**
 * public/js/theme.js — Manajemen Global Light / Dark Mode SiPaDin
 * PT Adiprima Suraprinta
 */

const THEME_STORAGE_KEY = 'sipadin_theme';

// Inisialisasi tema segera saat script dimuat (mencegah kedipan / FOUC)
(function initThemeImmediately() {
  try {
    const savedTheme = localStorage.getItem(THEME_STORAGE_KEY);
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  } catch (e) {
    console.warn('Error reading theme from localStorage:', e);
  }
})();

/**
 * Menerapkan tema (dark / light)
 * @param {'dark' | 'light'} theme
 */
function applyTheme(theme) {
  const root = document.documentElement;
  const isDark = theme === 'dark';
  
  if (isDark) {
    root.classList.add('dark');
  } else {
    root.classList.remove('dark');
  }
  
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  } catch (e) {
    console.warn('Error saving theme to localStorage:', e);
  }

  updateToggleButtons(isDark);
}

/**
 * Toggle antara tema dark dan light
 */
function toggleTheme() {
  const isCurrentlyDark = document.documentElement.classList.contains('dark');
  applyTheme(isCurrentlyDark ? 'light' : 'dark');
}

/**
 * Memperbarui tampilan ikon dan tooltip pada tombol toggle di header
 * @param {boolean} isDark
 */
function updateToggleButtons(isDark) {
  const toggleButtons = document.querySelectorAll('#themeToggleBtn, .theme-toggle-btn');
  toggleButtons.forEach(btn => {
    btn.setAttribute('title', isDark ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap');
    btn.setAttribute('aria-label', isDark ? 'Beralih ke Mode Terang' : 'Beralih ke Mode Gelap');
    
    const sunIcon = btn.querySelector('.theme-icon-sun');
    const moonIcon = btn.querySelector('.theme-icon-moon');
    
    if (sunIcon && moonIcon) {
      if (isDark) {
        sunIcon.classList.remove('hidden');
        moonIcon.classList.add('hidden');
      } else {
        sunIcon.classList.add('hidden');
        moonIcon.classList.remove('dark:block', 'hidden');
        moonIcon.classList.remove('hidden');
        sunIcon.classList.add('hidden');
      }
    }
  });
}

// Sinkronisasi antar-tab browser secara real-time
window.addEventListener('storage', (e) => {
  if (e.key === THEME_STORAGE_KEY) {
    applyTheme(e.newValue === 'dark' ? 'dark' : 'light');
  }
});

// Inisialisasi event listener setelah DOM selesai dimuat
document.addEventListener('DOMContentLoaded', () => {
  const isDark = document.documentElement.classList.contains('dark');
  updateToggleButtons(isDark);

  // Sinkronisasi jika OS theme preference berubah
  if (window.matchMedia) {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (!localStorage.getItem(THEME_STORAGE_KEY)) {
        applyTheme(e.matches ? 'dark' : 'light');
      }
    });
  }
});

// Event delegation agar toggle tombol selalu responsif
document.addEventListener('click', (e) => {
  const btn = e.target.closest('#themeToggleBtn, .theme-toggle-btn');
  if (btn) {
    e.preventDefault();
    toggleTheme();
  }
});

