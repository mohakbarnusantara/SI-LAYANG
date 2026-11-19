# SILAYANG — Sistem Layanan Kelurahan

**SILAYANG** adalah aplikasi web layanan administrasi kelurahan berbasis **Go** yang dirancang untuk memudahkan warga dalam mengajukan surat keterangan, menyampaikan aduan, serta memantau status permohonan secara daring. Di sisi internal, aplikasi mendukung alur kerja **Staf Kelurahan** dan **Lurah** sesuai wilayah masing-masing (32 kelurahan se-Kota Palu).

---

## Anggota Kelompok 11

| No | Nama | NIM |
| :---: | --- | --- |
| 1 | Amisha | `F5512510003` |
| 2 | Muh. Aqsha Syahrir | `F5512510033` |
| 3 | Moh. Akbar Nusantara | `F5512520037` |

---

## Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Teknologi & Struktur Proyek](#teknologi--struktur-proyek)
- [Cara Menjalankan Aplikasi](#cara-menjalankan-aplikasi)
- [Akun Demo & Pengujian](#akun-demo--pengujian)
- [Alur Kerja Sistem](#alur-kerja-sistem)
- [Data, Penyimpanan & Pengaturan](#data-penyimpanan--pengaturan)
- [Menjalankan Unit Test](#menjalankan-unit-test)
- [Dokumentasi Diagram UML](#dokumentasi-diagram-uml)
- [Catatan Keamanan](#catatan-keamanan)

---

## Fitur Utama

- **Registrasi & Verifikasi Akun Warga**: Pendaftaran warga dilengkapi unggah foto KTP, validasi format input, serta alur penolakan bersebab dan verifikasi ulang jika data ditolak staf.
- **30 Jenis Layanan Surat**: Mencakup pengajuan untuk diri sendiri maupun orang lain/anggota keluarga, termasuk form terstruktur (seperti *Surat Keterangan Penghasilan* dengan data orang tua, data anak, dan rincian penghasilan).
- **Alur Revisi Surat Dua Arah**:
  - Staf dapat meminta revisi beserta catatan perbaikan.
  - Warga mendapat tombol **Perbaiki surat**, dengan **penanda warna merah otomatis** pada kolom yang perlu diperbaiki (berubah hijau setelah diedit).
  - Panel admin menampilkan status jelas (**Belum direvisi warga** / **Sudah direvisi warga**) serta rincian lengkap data pemohon dan surat terbaru di urutan paling atas.
- **Persetujuan Lurah & Cetak PDF**: Staf menyiapkan draf dengan nomor surat resmi, lalu lurah meninjau draf A4 berkop dinamis kelurahan dan menyetujui dokumen hingga siap dicetak/disimpan PDF oleh warga.
- **Aduan & Pelacakan Publik**: Warga dapat mengirim aduan beserta bukti, sementara status surat dan aduan dapat dipantau lewat dashboard maupun halaman `/tracking` menggunakan kode pelacakan (`REQ-...` / `ADU-...`).

---

## Teknologi & Struktur Proyek

- **Backend**: Go (`net/http`, `html/template`, penyimpanan file JSON dengan `sync.RWMutex`)
- **Frontend**: HTML Templates, CSS responsif (desktop & mobile), Vanilla JavaScript

```text
si-layang/
├── main.go                 # Server HTTP, routing, logika bisnis, model & store JSON
├── main_test.go            # Unit test & integration test alur aplikasi
├── go.mod                  # Definisi modul Go
├── Dockerfile              # Konfigurasi container Docker
├── docker-compose.yml      # Orkestrasi container
├── static/
│   ├── app.css             # Stylesheet utama & tampilan responsif
│   └── app.js              # Interaktivitas form, validasi & highlight revisi
├── templates/
│   ├── layout.html         # Kerangka layout utama
│   ├── home.html           # Halaman beranda
│   ├── auth.html           # Halaman masuk & daftar akun
│   ├── dashboard.html      # Dashboard warga / ringkasan petugas
│   ├── surat_form.html     # Form pengajuan 30 layanan surat
│   ├── surat_edit.html     # Form perbaikan (revisi) surat warga
│   ├── complaint_form.html # Form pengajuan aduan warga
│   ├── tracking.html       # Halaman lacak status permohonan/aduan
│   ├── admin.html          # Panel administrasi staf & lurah
│   └── document.html       # Pratinjau & cetak dokumen surat resmi A4
├── data/                   # Penyimpanan data JSON (dibuat otomatis)
└── uploads/                # Berkas unggahan KTP & lampiran (dibuat otomatis)
```

---

## Cara Menjalankan Aplikasi

### Persyaratan

- **Go 1.22** atau lebih baru (tidak memerlukan database eksternal maupun paket pihak ketiga).
- PowerShell di Windows, Terminal di macOS/Linux, atau terminal bawaan VS Code.

Periksa versi Go:

```powershell
go version
```

### Menjalankan Secara Lokal

1. Buka folder proyek dan pastikan terminal berada di direktori yang berisi `go.mod`.
2. Jalankan perintah berikut:

   ```powershell
   go run .
   ```

3. Buka **http://localhost:8080** di browser.
4. Tekan **Ctrl+C** di terminal untuk menghentikan server.

### Mengganti Port (Jika Port 8080 Terpakai)

```powershell
$env:SILAYANG_ADDR = ":8081"
go run .
```

Lalu akses **http://localhost:8081**.

### Build ke Binary (Opsional)

```powershell
go build -o silayang.exe .
.\silayang.exe
```

---

## Akun Demo & Pengujian

Semua akun bawaan menggunakan kata sandi awal: **`silayang123`**.

| Peran | Email | Keterangan |
| --- | --- | --- |
| Warga | `warga@silayang.local` | Akun demo warga aktif *(dibuat saat database awal masih kosong)* |
| Staf | `staf@silayang.local` | Akun demo staf (Kelurahan Talise) |
| Lurah | `lurah@silayang.local` | Akun demo lurah (Kelurahan Talise) |

Setiap kelurahan dari **32 kelurahan** di Kota Palu juga memiliki akun staf dan lurah yang dibuat otomatis dengan format:

```text
staf.<kecamatan>-<kelurahan>@silayang.local
lurah.<kecamatan>-<kelurahan>@silayang.local
```

Contoh untuk Kelurahan **Talise** (Kec. Mantikulore) dan **Tatura Utara** (Kec. Palu Selatan):

```text
staf.mantikulore-talise@silayang.local
lurah.mantikulore-talise@silayang.local

staf.palu-selatan-tatura-utara@silayang.local
lurah.palu-selatan-tatura-utara@silayang.local
```

---

## Alur Kerja Sistem

1. **Pendaftaran & Verifikasi Warga**: Warga mendaftar di `/register` sambil memilih kecamatan/kelurahan domisili dan mengunggah foto KTP. Staf kelurahan terkait meninjau pendaftar di `/admin` lalu menyetujui atau menolak (dengan alasan agar warga bisa memperbaiki data).
2. **Pengajuan & Pemeriksaan Surat**: Warga aktif mengajukan surat di `/surat/baru`. Staf memeriksa data lengkap pemohon di `/admin`, lalu memilih **Siapkan draft** (mengisi nomor surat resmi), **Minta revisi** (warga memperbaiki kolom yang ditandai merah), atau **Tolak**.
3. **Persetujuan Lurah & Cetak PDF**: Lurah meninjau draf surat di `/admin` → **Tinjau draft ↗**, lalu klik **Setujui**. Setelah status menjadi **Selesai**, warga dapat membuka dokumen dan menyimpannya sebagai PDF lewat dialog cetak browser.
4. **Aduan & Pelacakan**: Warga mengirim aduan di `/aduan/baru` dan petugas menindaklanjutinya di `/admin`. Status permohonan surat dan aduan dapat dicek kapan saja melalui `/tracking`.

---

## Data, Penyimpanan & Pengaturan

Aplikasi menyimpan seluruh data akun, permohonan surat, dan aduan dalam format JSON tanpa database SQL:

- Data aplikasi: `data/silayang.json`
- Berkas unggahan (foto KTP, pengantar RT, bukti aduan): `uploads/`

Variabel lingkungan opsional yang dapat diatur sebelum `go run .`:

| Variabel | Nilai Bawaan | Fungsi |
| --- | --- | --- |
| `SILAYANG_ADDR` | `:8080` | Alamat dan port web server |
| `SILAYANG_DATA` | `data/silayang.json` | Lokasi file data JSON |
| `SILAYANG_UPLOADS` | `uploads` | Direktori penyimpanan berkas unggahan |
| `SILAYANG_SECRET` | Kunci development bawaan | Kunci penandatanganan sesi cookie |
| `SILAYANG_COOKIE_SECURE` | `false` | Set `true` jika dijalankan di balik HTTPS |

---

## Menjalankan Unit Test

Dari direktori utama proyek, jalankan:

```powershell
go test -count=1 ./...
```

---

## Dokumentasi Diagram UML

Bagian ini memuat dokumentasi rancangan sistem **SILAYANG** menggunakan **Unified Modeling Language (UML)**.

### 1. Use Case Diagram

```mermaid
flowchart LR
    subgraph Aktor ["Aktor Pengguna"]
        Publik["Publik / Pengunjung"]
        Warga["Warga"]
        Staf["Staf Kelurahan"]
        Lurah["Lurah"]
    end

    subgraph System ["Sistem Informasi SILAYANG"]
        UC1["Melihat Informasi & Beranda"]
        UC2["Melacak Status Permohonan/Aduan (Tracking Code)"]
        UC3["Registrasi Akun Warga (Unggah Foto KTP)"]
        UC4["Login & Logout"]
        UC5["Perbaiki Data & Verifikasi Ulang Akun (Jika Ditolak)"]
        UC6["Mengajukan Permohonan Surat (30 Layanan)"]
        UC7["Memperbaiki Surat yang Diminta Revisi"]
        UC8["Menyampaikan Aduan Warga & Bukti"]
        UC9["Mencetak / Simpan PDF Surat Resmi"]
        UC10["Verifikasi Akun Warga (Setujui / Tolak Bersebab)"]
        UC11["Memeriksa Berkas Permohonan Surat Masuk"]
        UC12["Memproses Surat (Input No. Resmi, Minta Revisi, Tolak)"]
        UC13["Menindaklanjuti & Menjawab Aduan Warga"]
        UC14["Menyetujui Draft Surat Resmi (Approval Lurah)"]
        UC15["Melihat Statistik & Riwayat Kelurahan"]
    end

    Publik --> UC1
    Publik --> UC2
    Publik --> UC3
    Publik --> UC4

    Warga --> UC4
    Warga --> UC5
    Warga --> UC6
    Warga --> UC7
    Warga --> UC8
    Warga --> UC9
    Warga --> UC2

    Staf --> UC4
    Staf --> UC10
    Staf --> UC11
    Staf --> UC12
    Staf --> UC13
    Staf --> UC15

    Lurah --> UC4
    Lurah --> UC14
    Lurah --> UC15
```

### 2. Class Diagram

```mermaid
classDiagram
    class User {
        +string ID
        +string NIK
        +string Name
        +string Email
        +string Phone
        +string BirthPlace
        +string BirthDate
        +string Gender
        +string Address
        +string Occupation
        +string KTPPath
        +string HomeDistrict
        +string HomeNeighborhood
        +string Role
        +string AccountStatus
        +string RejectionReason
        +string OfficeDistrict
        +string OfficeNeighborhood
        +string PasswordHash
        +time.Time CreatedAt
    }

    class SuratRequest {
        +string ID
        +string TrackingCode
        +string UserID
        +string ServiceCode
        +string ServiceName
        +string District
        +string Neighborhood
        +bool IsSelf
        +map Applicant
        +map Specific
        +string Purpose
        +string RTLetterNumber
        +string[] Attachments
        +string OfficialNumber
        +string Status
        +string AdminNote
        +int RevisionCount
        +time.Time* RevisedAt
        +time.Time CreatedAt
        +time.Time UpdatedAt
    }

    class Complaint {
        +string ID
        +string TrackingCode
        +string UserID
        +string Category
        +string District
        +string Neighborhood
        +string Description
        +string EvidencePath
        +string Status
        +string AdminReply
        +time.Time CreatedAt
        +time.Time UpdatedAt
    }

    class Service {
        +string Code
        +string Category
        +string Name
        +string[] Fields
        +string ExtraDoc
    }

    class LetterOffice {
        +string Neighborhood
        +string District
        +string Name
        +string Mayor
        +string Address
    }

    class Store {
        -sync.RWMutex mu
        -string path
        -database data
        +openStore(path) Store
        +snapshot() database
        +update(fn) error
    }

    class App {
        -Store store
        -Template templates
        -byte[] secret
        -string uploadDir
        -bool secureCookie
        +register(w, r)
        +resubmitVerification(w, r)
        +createSurat(w, r)
        +updateSurat(w, r)
        +createComplaint(w, r)
        +reviewResident(w, r)
        +reviewSurat(w, r)
        +reviewComplaint(w, r)
        +document(w, r)
    }

    Store --> User : stores
    Store --> SuratRequest : stores
    Store --> Complaint : stores
    App --> Store : manages
    SuratRequest ..> Service : references code
    SuratRequest --> User : belongs to (UserID)
    Complaint --> User : submitted by (UserID)
    SuratRequest ..> LetterOffice : office verification
```

### 3. Sequence Diagram

#### A. Registrasi, Penolakan, dan Verifikasi Ulang Akun Warga

```mermaid
sequenceDiagram
    autonumber
    actor W as Warga
    participant B as Web Browser / UI
    participant S as Server SILAYANG
    participant DB as File Storage (JSON & Uploads)
    actor ST as Staf Kelurahan

    W->>B: Isi form registrasi & upload Foto KTP
    B->>B: Validasi karakter nama, NIK 16 digit, format foto
    B->>S: POST /register (Multipart Form)
    S->>S: Validasi regex nama, NIK unik, validasi file KTP
    S->>DB: Simpan berkas KTP & buat Akun (Status: pending)
    S-->>B: Redirect ke /dashboard?notice=pending

    ST->>B: Buka Panel Admin /admin
    B->>S: GET /admin
    S->>DB: Ambil data pendaftar sewilayah
    S-->>B: Tampilkan antrean warga pending
    ST->>B: Periksa Foto KTP & Klik Tolak (isi alasan)
    B->>S: POST /admin/warga/{id} (action: reject, reason)
    S->>DB: Update status akun = "ditolak", simpan alasan

    W->>B: Buka /dashboard
    B->>S: GET /dashboard
    S-->>B: Tampilkan panel merah penolakan & form "Perbaiki Data"
    W->>B: Perbaiki data/NIK & klik "Kirim perbaikan & verifikasi ulang"
    B->>S: POST /akun/verifikasi-ulang
    S->>S: Validasi ulang data
    S->>DB: Update status akun = "pending", hapus alasan penolakan
    S-->>B: Redirect ke /dashboard?notice=reverification-submitted
```

#### B. Pengajuan Surat, Revisi, sampai Disetujui dan Dicetak

```mermaid
sequenceDiagram
    autonumber
    actor W as Warga (Aktif)
    participant B as Web Browser / UI
    participant S as Server SILAYANG
    participant DB as File Storage (JSON & Uploads)
    actor ST as Staf Kelurahan
    actor L as Lurah

    W->>B: Akses form /surat/baru
    W->>B: Pilih kelurahan, jenis surat, isi data subjek & upload berkas RT
    B->>S: POST /surat
    S->>DB: Simpan permohonan (Status: pending_verifikasi, generate Tracking Code)
    S-->>B: Redirect ke /dashboard

    ST->>B: Buka /admin (Inbox Surat)
    B->>S: GET /admin
    S-->>B: Tampilkan surat terbaru di atas & data lengkap pemohon
    alt Perlu Revisi
        ST->>B: Isi catatan revisi & klik "Minta revisi"
        B->>S: POST /admin/surat/{id} (action: revision, note)
        S->>DB: Update status = "perlu_revisi"
        W->>B: Klik "Perbaiki surat →" (/surat/{id}/edit)
        B-->>W: Tandai merah kolom sesuai catatan staf
        W->>B: Perbaiki data & kirim ulang
        B->>S: POST /surat/{id}/edit
        S->>DB: Update status = "pending_verifikasi", catat RevisionCount & RevisedAt
    end

    ST->>B: Input Nomor Surat Resmi & klik "Siapkan draft"
    B->>S: POST /admin/surat/{id} (action: draft, official_number)
    S->>DB: Update status = "draft_ready" & nomor surat

    L->>B: Buka /admin (sebagai Lurah) & klik "Tinjau draft ↗"
    B->>S: GET /dokumen/{id}
    S-->>B: Tampilkan dokumen berkop resmi (Watermark DRAFT)
    L->>B: Klik "Setujui" pada panel admin
    B->>S: POST /admin/surat/{id} (action: approve)
    S->>DB: Update status = "selesai"

    W->>B: Buka /dashboard atau /tracking
    B-->>W: Tampilkan status "Selesai" & tombol "Cetak / Simpan PDF →"
    W->>B: Buka /dokumen/{id} & klik "Cetak / Simpan PDF"
```

### 4. State Machine Diagram

#### A. Status Akun Warga (`AccountStatus`)

```mermaid
stateDiagram-v2
    [*] --> Pending : Registrasi Warga Baru (KTP Diunggah)
    Pending --> Aktif : Disetujui oleh Staf Kelurahan
    Pending --> Ditolak : Ditolak oleh Staf Kelurahan (Ada Alasan Penolakan)
    Ditolak --> Pending : Verifikasi Ulang (Warga Memperbaiki Data / KTP)
    Aktif --> [*] : Akun Aktif Dapat Mengajukan Layanan
```

#### B. Status Permohonan Surat (`SuratRequest.Status`)

```mermaid
stateDiagram-v2
    [*] --> pending_verifikasi : Warga Mengajukan Surat (/surat)
    pending_verifikasi --> perlu_revisi : Staf Meminta Perbaikan Berkas
    perlu_revisi --> pending_verifikasi : Warga Memperbaiki Berkas
    pending_verifikasi --> ditolak : Ditolak oleh Staf Kelurahan
    pending_verifikasi --> draft_ready : Staf Mengisi No. Surat Resmi & Siapkan Draft
    draft_ready --> ditolak : Ditolak oleh Lurah
    draft_ready --> selesai : Disetujui oleh Lurah
    selesai --> [*] : Surat Siap Dicetak / Disimpan PDF
```

### 5. Activity Diagram — Alur Pengajuan Layanan Surat

```mermaid
flowchart TD
    Start([Mulai]) --> LoginCheck{Warga Login & Akun Aktif?}
    LoginCheck -- Tidak --> BlockAccess[Arahkan ke Dashboard / Login]
    LoginCheck -- Ya --> OpenForm[Buka Halaman Form Surat /surat/baru]

    OpenForm --> FillData[Pilih Wilayah, Jenis Surat, Isi Data & Unggah Pengantar RT]
    FillData --> SubmitForm[Klik Kirim Permohonan]

    SubmitForm --> ValidateClient{Validasi Form di Browser}
    ValidateClient -- Gagal --> ShowPopup[Tampilkan Peringatan pada Kolom Kosong/Salah]
    ShowPopup --> FillData

    ValidateClient -- Berhasil --> PostServer[Kirim Data ke Server POST /surat]
    PostServer --> ValidateServer{Validasi Server & Penyimpanan}
    ValidateServer -- Gagal --> ReturnError[Tampilkan Galat di Halaman]
    ReturnError --> FillData

    ValidateServer -- Berhasil --> SaveData[Simpan Data & File di uploads/ dan data/silayang.json]
    SaveData --> StatusPending[Status: pending_verifikasi]

    StatusPending --> StaffReview[Staf Kelurahan Meninjau Berkas di /admin]
    StaffReview --> StaffDecision{Keputusan Staf}

    StaffDecision -- Minta Revisi --> StatusRevisi[Status: perlu_revisi]
    StatusRevisi --> ResidentEdit[Warga Perbaiki Kolom Bertanda Merah]
    ResidentEdit --> StatusPending
    StaffDecision -- Tolak --> StatusReject[Status: ditolak]
    StaffDecision -- Lengkap --> InputNomor[Input Nomor Surat Resmi & Generate Draft]

    InputNomor --> StatusDraft[Status: draft_ready]
    StatusDraft --> LurahReview[Lurah Meninjau Dokumen Draft]

    LurahReview --> LurahDecision{Persetujuan Lurah}
    LurahDecision -- Tolak --> StatusReject
    LurahDecision -- Setuju --> StatusSelesai[Status: selesai]

    StatusSelesai --> PrintReady[Surat Resmi Terbit, Warga Cetak / Simpan PDF]
    StatusReject --> End([Selesai])
    PrintReady --> End
```

---

## Catatan Keamanan

Aplikasi ini ditujukan untuk simulasi, pembelajaran, dan demonstrasi tugas perkuliahan. QR code, NIP, dan tanda tangan pada dokumen yang disetujui bersifat **simulasi/dummy** (bukan TTE resmi). Untuk penggunaan produksi nyata diperlukan autentikasi yang lebih ketat, HTTPS, basis data produksi, serta integrasi TTE resmi.
