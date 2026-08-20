// script.js – core logic for the web version
// -------------------------------------------------
// Simple in‑memory user store (hard‑coded for demo)
const USERS = {
    employee: { password: "emp123", role: "karyawan" },
    manager:  { password: "mgr123", role: "atasan" }
};

// ------------------------------------------------------------------
// ROUTING – decide which initializer runs based on file name

document.addEventListener("DOMContentLoaded", () => {
    const page = location.pathname.split("/").pop();
    if (page === "" || page === "index.html" ) initLogin();
    else if (page === "form.html") initForm();
    else if (page === "manager.html") initManager();
});

// ------------------------------------------------------------------
// 1. LOGIN PAGE
function initLogin() {
    const form = document.getElementById("loginForm");
    if (!form) return; // not a login page
    form.addEventListener("submit", e => {
        e.preventDefault();
        const user = document.getElementById("username").value.trim();
        const pass = document.getElementById("password").value;
        const rec = USERS[user];
        const msg = document.getElementById("msg");
        if (rec && rec.password === pass) {
            localStorage.setItem("session", JSON.stringify({ user, role: rec.role }));
            // redirect based on role
            if (rec.role === "karyawan") location.href = "form.html";
            else location.href = "manager.html";
        } else {
            msg.textContent = "Username atau password salah.";
        }
    });
}

// ------------------------------------------------------------------
// 2. FORM PAGE – Karyawan submits a request
function initForm() {
    const sess = JSON.parse(localStorage.getItem("session") || "{}");
    if (sess.role !== "karyawan") { location.href = "index.html"; return; }
    const form = document.getElementById("pengajuanForm");
    const msg = document.getElementById("msg");
    form.addEventListener("submit", e => {
        e.preventDefault();
        const submission = {
            id: Date.now(),
            tujuan: document.getElementById("tujuan").value.trim(),
            start: document.getElementById("start").value,
            end: document.getElementById("end").value,
            deskripsi: document.getElementById("deskripsi").value.trim(),
            approved: false,
            submitter: sess.user,
            submittedAt: new Date().toISOString()
        };
        const list = JSON.parse(localStorage.getItem("submissions") || "[]");
        list.push(submission);
        localStorage.setItem("submissions", JSON.stringify(list));
        msg.textContent = "Pengajuan berhasil disimpan. Redirecting...";
        setTimeout(() => location.href = "index.html", 1500);
    });
}

// ------------------------------------------------------------------
// 3. MANAGER DASHBOARD – review and approve
function initManager() {
    const sess = JSON.parse(localStorage.getItem("session") || "{}");
    if (sess.role !== "atasan") { location.href = "index.html"; return; }
    document.getElementById("logoutBtn").addEventListener("click", () => {
        localStorage.removeItem("session");
        location.href = "index.html";
    });

    const listEl = document.getElementById("list");
    const submissions = JSON.parse(localStorage.getItem("submissions") || "[]");
    if (!submissions.length) { listEl.innerHTML = "<li>Belum ada pengajuan.</li>"; return; }

    submissions.forEach(item => {
        const li = document.createElement("li");
        li.className = "pengajuan-item";
        const status = item.approved ? "✅ Disetujui" : "⏳ Menunggu";
        li.innerHTML = `
            <div><strong>Tujuan:</strong> ${item.tujuan}</div>
            <div><strong>Periode:</strong> ${item.start} → ${item.end}</div>
            <div><strong>Deskripsi:</strong> ${item.deskripsi}</div>
            <div class="status">${status}</div>
        `;
        // ---- H‑3 warning (3 days before start) ----
        const today = new Date();
        const startDate = new Date(item.start);
        const diffDays = Math.ceil((startDate - today) / (1000*60*60*24));
        if (!item.approved && diffDays >= 0 && diffDays <= 3) {
            const warn = document.createElement("div");
            warn.className = "warning";
            warn.textContent = `⚠️ H‑3: Mulai dalam ${diffDays} hari.`;
            li.appendChild(warn);
        }
        // Approve button if pending
        if (!item.approved) {
            const btn = document.createElement("button");
            btn.className = "approve-btn";
            btn.textContent = "Setujui";
            btn.addEventListener("click", () => {
                item.approved = true;
                localStorage.setItem("submissions", JSON.stringify(submissions));
                initManager(); // refresh UI
            });
            li.appendChild(btn);
        }
        listEl.appendChild(li);
    });
}

// End of script.js
