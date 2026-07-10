# Backup & Restore Database

Panduan backup database sebelum melakukan migrate/update, dan cara restore-nya.

## Cara 1 (Rekomendasi): `mysqldump` — snapshot database utuh

Cocok untuk skenario "backup dulu sebelum migrate, lalu restore persis seperti semula". Jalankan dari root project (folder yang berisi `docker-compose.yml`).

### Backup

```bash
docker compose exec mysql sh -c 'exec mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' > backup_$(date +%Y%m%d_%H%M%S).sql
```

```bash
docker compose exec mysql sh -c 'exec mysqldump -uroot -p"changeme_root" "md_note"' > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Restore

```bash
docker compose exec -T mysql sh -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' < backup_20260710_xxxxxx.sql
```

```bash
docker compose exec -T mysql sh -c 'exec mysql -uroot -p"changeme_root" "md_note"' < backup_20260710_xxxxxx.sql
```

Ganti `backup_20260710_xxxxxx.sql` dengan nama file backup yang sesuai.

Perintah ini pakai `docker compose exec` (bukan `docker exec` dengan nama container manual), dan memanfaatkan env var `$MYSQL_ROOT_PASSWORD` / `$MYSQL_DATABASE` yang sudah ada di dalam container `mysql` (diset dari `docker-compose.yml`) — jadi tidak perlu mengetik password secara manual.

**Catatan migrate:** migration (`golang-migrate`, folder `backend/migrations/`) berjalan otomatis setiap kali container `backend` start. Dump `mysqldump` di atas valid sebagai rollback point selama migration tidak mengubah struktur tabel secara tidak kompatibel. Kalau ada migration yang mengubah/menghapus kolom, restore dump lama ke skema baru bisa gagal/perlu penyesuaian manual.

## Cara 2: `mysqldump` langsung ke port yang di-export (tanpa `docker compose exec`)

Karena port MySQL di-export ke host (`DB_HOST_PORT` di `.env`, default `3309` → `3306` di container), database bisa diakses langsung dari luar container — baik dari mesin host itu sendiri maupun dari komputer lain di jaringan yang sama, tanpa perlu masuk lewat `docker compose exec`. Syaratnya, mesin yang menjalankan perintah ini sudah terinstall MySQL client (`mysqldump`/`mysql`, contoh install: `sudo apt install mysql-client`).

Jalankan dari root project (folder yang berisi `.env`), supaya variabel `DB_HOST_PORT`, `DB_NAME`, `DB_ROOT_PASSWORD` otomatis terbaca dari `.env` tanpa perlu ditulis manual:

### Backup

```bash
set -a && source .env && set +a
mysqldump -h 127.0.0.1 -P "$DB_HOST_PORT" -u root -p"$DB_ROOT_PASSWORD" "$DB_NAME" > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Restore

```bash
set -a && source .env && set +a
mysql -h 127.0.0.1 -P "$DB_HOST_PORT" -u root -p"$DB_ROOT_PASSWORD" "$DB_NAME" < backup_20260710_xxxxxx.sql
```

Kalau backup diambil dari komputer lain (bukan home server-nya langsung), ganti `127.0.0.1` dengan alamat IP home server (misal `192.168.0.143`, sesuaikan dengan `HOST` di `.env`) — pastikan port tersebut memang bisa diakses dari jaringan tersebut (firewall/router mengizinkan).

## Cara 3 (Alternatif): fitur backup/restore bawaan aplikasi (JSON)

Ada di halaman **Settings → Backup & Restore** di aplikasi (hanya bisa diakses oleh akun dengan role admin).

- **Export** → tombol download, hasilnya file `md-note-backup-<tanggal>.json` (berisi users termasuk password hash, tags, notes, note_tags, cors_origins).
- **Import** → upload file JSON yang sama.

Via API langsung (butuh token admin dari `POST /api/auth/login`):

```bash
# Export
curl -H "Authorization: Bearer $TOKEN" http://localhost:8888/api/backup -o md-note-backup.json

# Restore
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
     --data @md-note-backup.json http://localhost:8888/api/restore
```

**Catatan penting:** restore di fitur ini pakai `INSERT IGNORE` per baris — kalau baris dengan `id` yang sama masih ada di database, isinya **tidak ditimpa** (data lama dipertahankan, bukan dikembalikan ke isi backup). Jadi ini cocok untuk menggabungkan backup lama ke database baru/kosong, **bukan** untuk mengembalikan database persis ke kondisi sebelum migrate.

## Kesimpulan

Gunakan **Cara 1 atau Cara 2 (`mysqldump`)** sebagai pengaman utama sebelum migrate — keduanya sama-sama menghasilkan snapshot database yang utuh, tinggal pilih mana yang lebih praktis (lewat `docker compose exec` di server, atau langsung dari mesin manapun yang bisa akses port MySQL yang di-export). Cara 3 (JSON) boleh dipakai sebagai backup tambahan/portable (misal untuk pindah server), tapi jangan diandalkan sebagai rollback point.
