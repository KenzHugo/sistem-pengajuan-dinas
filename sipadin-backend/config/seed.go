package config

// SeedData is retained as a no-op placeholder.
// All master data (Karyawan_copy1, PerjalananDinas_Tarif, PerjalananDinas_Fasilitas)
// is pre-populated via the SQL dump files provided in the project root:
//   - Karyawan_copy1.sql
//   - PerjalananDinas_Tarif.sql
//   - PerjalananDinas_Fasilitas.sql
//
// The transactional table (PerjalananDinas_Pengajuan) is created automatically
// by AutoMigrate in database.go and is populated through normal application usage.
func SeedData() {}
