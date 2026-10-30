# SECURITY (notes for whoever works in this repo)

- **Secrets:** `src/.env` (DB_URI, basic-auth, Elasticsearch/vector settings) lives only on the server and is never committed.
- **Personal data:** WhatsApp sessions/messages in the `storages_data` volume and `logs_data`; never commit them.
- **Exposed:** port 3000 (REST API, basic auth) on the host running the compose (hattin).
- **Before push/deploy:** `git status` shows no `.env`; after a rebuild run `nazhiba_cloud-node scan all` and check brotli-libs >= 1.2.1-r0.

## Temuan & perbaikan

| Tanggal | Temuan | Risiko | Sebelum | Sesudah |
|---|---|---|---|---|
| 2026-10-30 | Image `whatsapp_go` memakai layer Alpine 8 bulan: brotli-libs 1.2.0-r0 (CVE-2026-111497, CRITICAL) | Penyerang bisa memakai bug pustaka dekompresi lama lewat media yang diproses | `fleet-ssh hattin -- sudo -n nazhiba_cloud-node scan all` → 47 HIGH/CRITICAL (44 HIGH, 3 CRITICAL) | `belum diperbaiki` (PR `apk upgrade --no-cache`, lalu redeploy + scan ulang) |
