# SECURITY (notes for whoever works in this repo)

- **Secrets:** `src/.env` (DB_URI, basic-auth, Elasticsearch/vector settings) lives only on the server and is never committed.
- **Personal data:** WhatsApp sessions/messages in the `storages_data` volume and `logs_data`; never commit them.
- **Exposed:** port 3000 (REST API, basic auth) on the host running the compose.
- **Before push/deploy:** `git status` shows no `.env`; after a rebuild run `nazhiba_cloud-node scan all` and check brotli-libs >= 1.2.1-r0.

## Temuan & perbaikan

| Tanggal | Temuan | Risiko | Sebelum | Sesudah |
|---|---|---|---|---|
| 2026-10-30 | Image `whatsapp_go` memakai layer Alpine 8 bulan: brotli-libs 1.2.0-r0 (CVE-2026-111497, CRITICAL) | Penyerang bisa memakai bug pustaka dekompresi lama lewat media yang diproses | `nazhiba_cloud-node scan all` → 47 HIGH/CRITICAL (44 HIGH, 3 CRITICAL) | `belum diperbaiki` (PR `apk upgrade --no-cache`, lalu redeploy + scan ulang) |
| 2026-10-30 | Binary Go `app/whatsapp` (image `whatsapp_go`): 45 HIGH di stdlib Go 1.25.7, x/net 0.50, x/crypto 0.48, x/text 0.34, x/image 0.36, klauspost/compress 1.18.4, buger/jsonparser 1.1.1 | Penyerang bisa membuat server macet (DoS) lewat HTTP/TLS/XML/gambar/JSON buatan, atau menyusup lewat sertifikat x509 palsu (MITM) | `nazhiba_cloud-node scan go-whatsapp-web-multidevice-source-whatsapp_go` → `app/whatsapp (gobinary) 45` | `nazhiba_cloud-node scan go-whatsapp-web-multidevice-source-whatsapp_go` (30 Okt, sesudah rilis PR #5, merge 29035af) → `alpine 0`, `app/whatsapp (gobinary) 0` (builder golang:1.26.10 + modul dinaikkan). Catatan: kontainer crash-loop pada koneksi DB basic-auth (bukan dari diff ini), ditangani kartu insiden terpisah |
