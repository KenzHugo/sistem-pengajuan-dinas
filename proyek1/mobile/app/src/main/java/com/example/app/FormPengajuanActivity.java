package com.example.app;

import android.content.SharedPreferences;
import android.os.Bundle;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import androidx.appcompat.app.AppCompatActivity;
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import java.lang.reflect.Type;
import java.util.ArrayList;
import java.util.List;

/**
 * FormPengajuanActivity – UI for employee to submit a travel request.
 * Uses SharedPreferences (JSON) as a lightweight storage.
 */
public class FormPengajuanActivity extends AppCompatActivity {
    private EditText etTujuan, etStart, etEnd, etDeskripsi;
    private TextView tvResult;
    private SharedPreferences prefs;
    private Gson gson = new Gson();

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_form);

        // Ensure only karyawan can access this screen
        SharedPreferences session = getSharedPreferences("session", MODE_PRIVATE);
        if (!"karyawan".equals(session.getString("role", ""))) {
            finish(); // illegal access
            return;
        }

        etTujuan = findViewById(R.id.etTujuan);
        etStart  = findViewById(R.id.etStart);
        etEnd    = findViewById(R.id.etEnd);
        etDeskripsi = findViewById(R.id.etDeskripsi);
        tvResult = findViewById(R.id.tvResult);
        prefs = getSharedPreferences("submissions", MODE_PRIVATE);

        Button submitBtn = findViewById(R.id.btnSubmit);
        submitBtn.setOnClickListener(v -> submitPengajuan());
    }

    private void submitPengajuan() {
        // ---------- Basic validation ----------
        if (etTujuan.getText().toString().trim().isEmpty()
                || etStart.getText().toString().isEmpty()
                || etEnd.getText().toString().isEmpty()
                || etDeskripsi.getText().toString().trim().isEmpty()) {
            tvResult.setText("Semua kolom harus diisi.");
            return;
        }

        // ---------- Build Pengajuan object ----------
        String user = getSharedPreferences("session", MODE_PRIVATE).getString("user", "unknown");
        Pengajuan p = new Pengajuan(
                System.currentTimeMillis(),
                etTujuan.getText().toString().trim(),
                etStart.getText().toString(),
                etEnd.getText().toString(),
                etDeskripsi.getText().toString().trim(),
                false, // belum disetujui
                user
        );

        // ---------- Store in SharedPreferences ----------
        List<Pengajuan> list = loadList();
        list.add(p);
        saveList(list);
        tvResult.setText("Pengajuan berhasil disimpan.");
        clearForm();
    }

    private List<Pengajuan> loadList() {
        String json = prefs.getString("list", null);
        if (json == null) return new ArrayList<>();
        Type type = new TypeToken<ArrayList<Pengajuan>>() {}.getType();
        return gson.fromJson(json, type);
    }

    private void saveList(List<Pengajuan> list) {
        String json = gson.toJson(list);
        prefs.edit().putString("list", json).apply();
    }

    private void clearForm() {
        etTujuan.setText("");
        etStart.setText("");
        etEnd.setText("");
        etDeskripsi.setText("");
    }

    /** Simple POJO representing one submission */
    static class Pengajuan {
        long id;
        String tujuan;
        String start; // YYYY‑MM‑DD
        String end;
        String deskripsi;
        boolean approved;
        String submitter;

        Pengajuan(long id, String tujuan, String start, String end,
                  String deskripsi, boolean approved, String submitter) {
            this.id = id;
            this.tujuan = tujuan;
            this.start = start;
            this.end = end;
            this.deskripsi = deskripsi;
            this.approved = approved;
            this.submitter = submitter;
        }
    }
}
