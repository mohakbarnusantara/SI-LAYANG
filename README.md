# SILAYANG

SILAYANG adalah portal layanan administrasi warga dan kelurahan berbasis Go. Warga dapat mengajukan surat, menyampaikan aduan, dan melacak permohonan. Staf kelurahan dan lurah mengelola layanan sesuai wilayah akun mereka.

## Mulai cepat

### Persyaratan

- Go **1.22 atau lebih baru**.
- PowerShell di Windows, Terminal di macOS/Linux, atau terminal bawaan VS Code.
- Tidak perlu memasang database atau paket tambahan.

Periksa Go:

```powershell
go version
```

### Jalankan aplikasi

1. Buka folder proyek di VS Code.
2. Buka terminal (**Terminal → New Terminal**). Pastikan terminal berada di folder yang berisi `go.mod`.
3. Jalankan:

   ```powershell
   go run .
   ```

4. Tunggu sampai terminal menampilkan alamat aplikasi, lalu buka **http://localhost:8080** di browser.
5. Untuk menghentikan aplikasi, kembali ke terminal dan tekan **Ctrl+C**.

Jalankan `go run .` dari folder proyek, bukan dengan membuka `main.go` langsung. Aplikasi membuat folder data dan unggahan jika belum tersedia.

### Jika port 8080 sedang digunakan

Pilih port lain sebelum menjalankan aplikasi:

```powershell
$env:SILAYANG_ADDR = ":8081"
go run .
```

Lalu buka **http://localhost:8081**.

### Buat file aplikasi (opsional)

Untuk mengompilasi aplikasi menjadi file `.exe`:

```powershell
go build -o silayang.exe .
.\silayang.exe
```

## Masuk dan mencoba aplikasi

Kata sandi awal akun demo adalah **`silayang123`**.

| Peran | Email demo | Catatan |
| --- | --- | --- |
| Warga | `warga@silayang.local` | Akun demo aktif |
| Staf | `staf@silayang.local` | Akun lama, dibatasi ke Kelurahan Talise |
| Lurah | `lurah@silayang.local` | Akun lama, dibatasi ke Kelurahan Talise |

Setiap satu dari **32 kelurahan** juga mendapat akun staf dan lurah. Format emailnya:

```text
staf.<kecamatan>-<kelurahan>@silayang.local
lurah.<kecamatan>-<kelurahan>@silayang.local
```

Contoh akun Tondo:

```text
staf.mantikulore-tondo@silayang.local
lurah.mantikulore-tondo@silayang.local
```

Semua akun wilayah menggunakan kata sandi awal yang sama, `silayang123`. Akun wilayah dibuat otomatis saat aplikasi dijalankan. Gunakan akun sesuai kelurahan yang ingin diuji.

### Alur singkat

1. **Warga:** masuk dengan akun demo atau daftar akun baru. Saat mendaftar, pilih kecamatan dan kelurahan domisili.
2. **Verifikasi:** akun baru menunggu persetujuan staf dari kelurahan domisili. Warga dapat mengajukan layanan setelah akun aktif.
3. **Surat:** warga memilih wilayah dan jenis surat, mengisi kebutuhan, lalu mengirim permohonan. Staf memeriksa dan menyiapkan draft; lurah menyetujuinya. Setelah status surat **Selesai**, warga dapat memilih **Cetak / Simpan PDF** di riwayat surat, lalu memilih **Simpan sebagai PDF** pada dialog cetak browser.
4. **Aduan:** warga mengirim aduan, lalu memantau tindak lanjutnya melalui dashboard.
5. **Lacak:** gunakan kode tracking yang diberikan setelah permohonan atau aduan dikirim.

Staf hanya dapat melihat pendaftar, foto KTP, dan permohonan surat untuk kelurahan akunnya. Lurah hanya dapat menyetujui draft surat dari kelurahannya. Data warga lama hanya dipetakan otomatis jika alamatnya menyebut nama kelurahan secara jelas.

## Data dan penyimpanan

Aplikasi **tidak memakai database SQL**. Data akun, surat, dan aduan disimpan dalam file JSON:

```text
data/silayang.json
```

Berkas unggahan, seperti foto KTP dan lampiran, disimpan di:

```text
uploads/
```

Data tetap ada setelah aplikasi dihentikan. Jangan menghapus atau membagikan file data maupun unggahan jika berisi informasi pribadi.

## Fitur utama

- Pendaftaran dan verifikasi akun warga.
- 30 jenis layanan surat dan pratinjau draf A4.
- Pembuatan, pemeriksaan, dan persetujuan surat sesuai wilayah akun.
- Aduan warga dan pelacakan status.
- Kop surat dinamis dengan data kantor kelurahan.

QR dan NIP pada dokumen yang disetujui adalah **dummy untuk tampilan saja**. QR tidak dapat dipakai untuk verifikasi dan tanda tangan bukan TTE resmi. Kode tracking tetap digunakan untuk melacak status permohonan.

## Pengaturan opsional

Atur variabel berikut sebelum menjalankan `go run .`:

| Variabel | Nilai bawaan | Fungsi |
| --- | --- | --- |
| `SILAYANG_ADDR` | `:8080` | Alamat dan port web server |
| `SILAYANG_DATA` | `data/silayang.json` | Lokasi file data JSON |
| `SILAYANG_UPLOADS` | `uploads` | Lokasi penyimpanan unggahan |
| `SILAYANG_SECRET` | Kunci development bawaan | Kunci sesi login |
| `SILAYANG_COOKIE_SECURE` | `false` | Gunakan `true` jika aplikasi diakses melalui HTTPS |

Contoh mengatur lokasi data di PowerShell:

```powershell
$env:SILAYANG_DATA = "data\uji.json"
go run .
```

## Menjalankan tes

Dari folder proyek, jalankan:

```powershell
go test ./...
```

## Catatan keamanan

Aplikasi ini ditujukan untuk demonstrasi dan pengembangan lokal. Kata sandi bersama, kunci sesi bawaan, HTTP tanpa HTTPS, dan file JSON lokal **tidak cocok untuk deployment produksi**. Layanan nyata memerlukan autentikasi dan pengelolaan kata sandi yang aman, HTTPS, penyimpanan data produksi, backup, serta TTE resmi.




# Dokumentasi Diagram UML — SILAYANG (Sistem Layanan Kelurahan)

Dokumen ini memuat kumpulan diagram **Unified Modeling Language (UML)** untuk sistem **SILAYANG**, yang dirancang sesuai dengan implementasi arsitektur sistem informasi berbasis Go, arsitektur multi-wilayah (32 kelurahan se-Kota Palu), serta alur kerja layanan administrasi publik.

---

## 1. Use Case Diagram

Diagram ini mendefinisikan fungsionalitas sistem dari sudut pandang para aktor: **Publik / Pengunjung**, **Warga**, **Staf Kelurahan**, dan **Lurah**.

```mermaid
flowchart LR
    %% Aktor
    subgraph Aktor ["Aktor Pengguna"]
        Publik["Publik / Pengunjung"]
        Warga["Warga"]
        Staf["Staf Kelurahan"]
        Lurah["Lurah"]
    end

    %% Batas Sistem
    subgraph System ["Sistem Informasi SILAYANG"]
        UC1["Melihat Informasi & Beranda"]
        UC2["Melacak Status Permohonan/Aduan (Tracking Code)"]
        UC3["Registrasi Akun Warga (Unggah Foto KTP)"]
        UC4["Login & Logout"]
        UC5["Perbaiki Data & Verifikasi Ulang Akun (Jika Ditolak)"]
        UC6["Mengajukan Permohonan Surat (30 Layanan)"]
        UC7["Menyampaikan Aduan Warga & Bukti"]
        UC8["Mencetak / Simpan PDF Surat Resmi"]
        UC9["Verifikasi Akun Warga (Setujui / Tolak Bersebab)"]
        UC10["Memeriksa Berkas Permohonan Surat Masuk"]
        UC11["Memproses Surat (Input No. Resmi, Minta Revisi, Tolak)"]
        UC12["Menindaklanjuti & Menjawab Aduan Warga"]
        UC13["Menyetujui Draft Surat Resmi (Approval Lurah)"]
        UC14["Melihat Statistik & Riwayat Kelurahan"]
    end

    %% Relasi Aktor Publik
    Publik --> UC1
    Publik --> UC2
    Publik --> UC3
    Publik --> UC4

    %% Relasi Warga
    Warga --> UC4
    Warga --> UC5
    Warga --> UC6
    Warga --> UC7
    Warga --> UC8
    Warga --> UC2

    %% Relasi Staf Kelurahan
    Staf --> UC4
    Staf --> UC9
    Staf --> UC10
    Staf --> UC11
    Staf --> UC12
    Staf --> UC14

    %% Relasi Lurah
    Lurah --> UC4
    Lurah --> UC13
    Lurah --> UC14
```

---

## 2. Class Diagram

Class diagram ini memetakan struktur entitas data, model relasi, tipe penyimpanan file JSON (`Store`), dan controller aplikasi (`App`).

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

---

## 3. Sequence Diagram

### A. Alur Registrasi, Penolakan, dan Verifikasi Ulang Akun Warga
Diagram ini menggambarkan siklus validasi KTP, penolakan oleh staf, dan mekanisme perbaikan akun.

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
    ST->>B: Periksa Foto KTP & Klik Tolak (isi alasan: "NIK tidak sesuai")
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

---

### B. Alur Pengajuan Surat sampai Disetujui dan Dicetak
Diagram ini menggambarkan proses bisnis pengajuan surat, pemeriksaan oleh staf, persetujuan lurah, dan pencetakan dokumen PDF.

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
    B->>B: Validasi pop-up kolom wajib
    B->>S: POST /surat
    S->>DB: Simpan permohonan (Status: pending_verifikasi, generate Tracking Code)
    S-->>B: Redirect ke /dashboard (Tampil riwayat permohonan)

    ST->>B: Buka /admin (Tab Inbox Surat)
    B->>S: GET /admin
    S-->>B: Tampilkan detail permohonan & pemohon
    ST->>B: Input Nomor Surat Resmi & klik "Siapkan draft"
    B->>S: POST /admin/surat/{id} (action: draft, official_number)
    S->>DB: Update status = "draft_ready" & nomor surat

    L->>B: Buka /admin (sebagai Lurah)
    L->>B: Klik "Tinjau draft ↗" (/dokumen/{id})
    B->>S: GET /dokumen/{id}
    S-->>B: Tampilkan format dokumen berkop resmi (Watermark DRAFT)
    L->>B: Klik "Setujui" pada panel admin
    B->>S: POST /admin/surat/{id} (action: approve)
    S->>DB: Update status = "selesai"

    W->>B: Buka /dashboard (atau via /tracking)
    B-->>W: Tampilkan status "Selesai" & tombol "Cetak / Simpan PDF →"
    W->>B: Buka /dokumen/{id} & klik "Cetak / Simpan PDF"
    B->>B: Browser Print Dialog (Save as PDF)
```

---

## 4. State Machine Diagram

### A. Status Akun Warga (`AccountStatus`)
```mermaid
stateDiagram-v2
    [*] --> Pending : Registrasi Warga Baru (KTP Diunggah)
    Pending --> Aktif : Disetujui oleh Staf Kelurahan
    Pending --> Ditolak : Ditolak oleh Staf Kelurahan (Ada Alasan Penolakan)
    Ditolak --> Pending : Verifikasi Ulang (Warga Memperbaiki Data / KTP)
    Aktif --> [*] : Akun Aktif Dapat Mengajukan Layanan
```

### B. Status Permohonan Surat (`SuratRequest.Status`)
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

---

## 5. Activity Diagram
 
### Alur Proses Pengajuan Layanan Surat
```mermaid
flowchart TD
    Start([Mulai]) --> LoginCheck{Warga Login & Akun Aktif?}
    LoginCheck -- Tidak --> BlockAccess[Arahkan ke Dashboard / Login]
    LoginCheck -- Ya --> OpenForm[Buka Halaman Form Surat /surat/baru]
    
    OpenForm --> FillData[Pilih Wilayah, Jenis Surat, Isi Data & Unggah Pengantar RT]
    FillData --> SubmitForm[Klik Kirim Permohonan]
    
    SubmitForm --> ValidateClient{Validasi Form di Browser}
    ValidateClient -- Gagal --> ShowPopup[Munculkan Pop-up Peringatan pada Kolom Kosong/Salah]
    ShowPopup --> FillData
    
    ValidateClient -- Berhasil --> PostServer[Kirim Data ke Server POST /surat]
    PostServer --> ValidateServer{Validasi Server & Penyimpanan}
    ValidateServer -- Gagal --> ReturnError[Tampilkan Galat di Halaman]
    ReturnError --> FillData
    
    ValidateServer -- Berhasil --> SaveData[Simpan Data & File di uploads/ dan data.json]
    SaveData --> StatusPending[Status: pending_verifikasi]
    
    StatusPending --> StaffReview[Staf Kelurahan Meninjau Berkas di /admin]
    StaffReview --> StaffDecision{Keputusan Staf}
    
    StaffDecision -- Minta Revisi --> StatusRevisi[Status: perlu_revisi]
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
