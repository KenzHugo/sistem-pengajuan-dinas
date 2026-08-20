package com.example.app;

import android.content.SharedPreferences;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;
import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.Toolbar;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import java.lang.reflect.Type;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;

/**
 * ApprovalActivity – dashboard for atasan (manager) to review and approve submissions.
 * Shows H‑3 warning when a pending request is within 3 days of its start date.
 */
public class ApprovalActivity extends AppCompatActivity {
    private RecyclerView rv;
    private SharedPreferences prefs;
    private Gson gson = new Gson();
    private List<FormPengajuanActivity.Pengajuan> submissions;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_approval);

        // Guard: only atasan may access
        SharedPreferences session = getSharedPreferences("session", MODE_PRIVATE);
        if (!"atasan".equals(session.getString("role", ""))) {
            finish();
            return;
        }

        Toolbar toolbar = findViewById(R.id.toolbar);
        toolbar.setNavigationOnClickListener(v -> finish());

        rv = findViewById(R.id.rvPengajuan);
        rv.setLayoutManager(new LinearLayoutManager(this));

        prefs = getSharedPreferences("submissions", MODE_PRIVATE);
        loadSubmissions();
    }

    private void loadSubmissions() {
        String json = prefs.getString("list", null);
        if (json == null) submissions = new ArrayList<>();
        else {
            Type type = new TypeToken<ArrayList<FormPengajuanActivity.Pengajuan>>() {}.getType();
            submissions = gson.fromJson(json, type);
        }
        rv.setAdapter(new PengajuanAdapter());
    }

    /** RecyclerView Adapter for each submission */
    private class PengajuanAdapter extends RecyclerView.Adapter<PengajuanAdapter.VH> {
        class VH extends RecyclerView.ViewHolder {
            TextView tvInfo, tvStatus, tvWarning;
            Button btnApprove;
            VH(View view) {
                super(view);
                tvInfo = view.findViewById(R.id.tvInfo);
                tvStatus = view.findViewById(R.id.tvStatus);
                tvWarning = view.findViewById(R.id.tvWarning);
                btnApprove = view.findViewById(R.id.btnApprove);
            }
        }

        @NonNull
        @Override
        public VH onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
            View v = LayoutInflater.from(parent.getContext())
                    .inflate(R.layout.item_pengajuan, parent, false);
            return new VH(v);
        }

        @Override
        public void onBindViewHolder(@NonNull VH holder, int position) {
            FormPengajuanActivity.Pengajuan p = submissions.get(position);
            holder.tvInfo.setText(
                    "Tujuan: " + p.tujuan + "\n" +
                    "Periode: " + p.start + " → " + p.end + "\n" +
                    "Deskripsi: " + p.deskripsi);
            holder.tvStatus.setText(p.approved ? "✅ Disetujui" : "⏳ Menunggu");
            holder.tvStatus.setVisibility(View.VISIBLE);

            // ----- H‑3 warning (3 days before start) -----
            if (!p.approved) {
                LocalDate today = LocalDate.now();
                LocalDate start = LocalDate.parse(p.start, DateTimeFormatter.ISO_DATE);
                long daysDiff = ChronoUnit.DAYS.between(today, start);
                if (daysDiff >= 0 && daysDiff <= 3) {
                    holder.tvWarning.setVisibility(View.VISIBLE);
                    holder.tvWarning.setText("⚠️ H‑3: Mulai dalam " + daysDiff + " hari.");
                } else {
                    holder.tvWarning.setVisibility(View.GONE);
                }
            } else {
                holder.tvWarning.setVisibility(View.GONE);
            }

            // Approve button (only when pending)
            if (!p.approved) {
                holder.btnApprove.setVisibility(View.VISIBLE);
                holder.btnApprove.setOnClickListener(v -> {
                    p.approved = true;
                    saveChanges();
                    notifyItemChanged(position);
                });
            } else {
                holder.btnApprove.setVisibility(View.GONE);
            }
        }

        @Override
        public int getItemCount() {
            return submissions.size();
        }
    }

    /** Persist changes after an approval */
    private void saveChanges() {
        String json = gson.toJson(submissions);
        prefs.edit().putString("list", json).apply();
    }
}
