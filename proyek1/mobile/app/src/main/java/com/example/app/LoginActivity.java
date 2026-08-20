package com.example.app;

import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import androidx.appcompat.app.AppCompatActivity;

/**
 * LoginActivity – handles login for both employee (karyawan) and manager (atasan).
 * Hard‑coded credentials for demo purpose:
 *   employee / emp123 → role "karyawan"
 *   manager  / mgr123 → role "atasan"
 */
public class LoginActivity extends AppCompatActivity {
    private EditText etUsername, etPassword;
    private TextView tvError;
    private SharedPreferences prefs;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);

        etUsername = findViewById(R.id.etUsername);
        etPassword = findViewById(R.id.etPassword);
        tvError    = findViewById(R.id.tvError);
        prefs = getSharedPreferences("session", MODE_PRIVATE);

        // If user already logged in, jump to proper screen
        String role = prefs.getString("role", null);
        if (role != null) {
            navigateToRole(role);
            return;
        }

        Button btnLogin = findViewById(R.id.btnLogin);
        btnLogin.setOnClickListener(v -> attemptLogin());
    }

    private void attemptLogin() {
        String user = etUsername.getText().toString().trim();
        String pass = etPassword.getText().toString();
        if (user.isEmpty() || pass.isEmpty()) {
            tvError.setText("Harap isi semua kolom.");
            return;
        }
        // Simple credential check
        if ("employee".equals(user) && "emp123".equals(pass)) {
            saveSession(user, "karyawan");
            navigateToRole("karyawan");
        } else if ("manager".equals(user) && "mgr123".equals(pass)) {
            saveSession(user, "atasan");
            navigateToRole("atasan");
        } else {
            tvError.setText("Username atau password salah.");
        }
    }

    private void saveSession(String user, String role) {
        prefs.edit()
                .putString("user", user)
                .putString("role", role)
                .apply();
    }

    private void navigateToRole(String role) {
        Intent intent;
        if ("karyawan".equals(role)) {
            intent = new Intent(this, FormPengajuanActivity.class);
        } else {
            intent = new Intent(this, ApprovalActivity.class);
        }
        startActivity(intent);
        finish(); // prevent back to login
    }
}
