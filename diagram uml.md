# Panduan & Kode Diagram UML untuk Draw.io (diagrams.net)

Dokumen ini berisi kode yang siap Anda salin (**copy-paste**) ke dalam **draw.io** (atau [app.diagrams.net](https://app.diagrams.net)).

---

## Cara Memasukkan Kode ke Draw.io

1. Buka [draw.io](https://app.diagrams.net/).
2. Pada menu bar bagian atas, klik tombol **`+`** (atau menu **Arrange / Atur** $\rightarrow$ **Insert / Sisipkan** $\rightarrow$ **Advanced / Lanjutan**).
3. Pilih **Mermaid** (atau **PlantUML** sesuai kode yang dipilih).
4. Hapus teks contoh, lalu **tempel (paste)** kode diagram di bawah ini.
5. Klik tombol **Insert**. Gambar diagram UML akan langsung terbentuk menjadi objek/shape yang bisa Anda geser, ubah warna, dan atur posisinya.

---

## 1. Use Case Diagram

### Opsi A: Format PlantUML (Rekomendasi untuk Use Case karena menghasilkan ikon orang/stickman standar UML)
> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ PlantUML***

```plantuml
@startuml
left to right direction
skinparam packageStyle rectangle
skinparam shadowing false

actor "Publik / Pengunjung" as Publik
actor "Warga" as Warga
actor "Staf Kelurahan" as Staf
actor "Lurah" as Lurah

rectangle "Sistem Informasi SILAYANG" {
  usecase "Melihat Informasi Beranda" as UC1
  usecase "Melacak Permohonan / Aduan\n(Tracking Code)" as UC2
  usecase "Registrasi Akun Warga\n(Wajib Foto KTP & Nama Sesuai)" as UC3
  usecase "Login & Logout" as UC4
  usecase "Perbaiki Data & Verifikasi Ulang\n(Jika Akun Ditolak)" as UC5
  usecase "Mengajukan Surat\n(Pilih Wilayah & Jenis Surat)" as UC6
  usecase "Menyampaikan Aduan & Bukti" as UC7
  usecase "Mencetak / Simpan PDF Surat" as UC8
  usecase "Verifikasi Akun Warga\n(Setujui / Tolak Bersebab)" as UC9
  usecase "Memeriksa Permohonan Surat Masuk" as UC10
  usecase "Memproses Surat\n(No. Resmi, Minta Revisi, Tolak)" as UC11
  usecase "Menindaklanjuti & Menjawab Aduan" as UC12
  usecase "Menyetujui Draft Surat (Approval)" as UC13
  usecase "Melihat Statistik Layanan Kelurahan" as UC14
}

Publik --> UC1
Publik --> UC2
Publik --> UC3
Publik --> UC4

Warga --> UC4
Warga --> UC5
Warga --> UC6
Warga --> UC7
Warga --> UC8
Warga --> UC2

Staf --> UC4
Staf --> UC9
Staf --> UC10
Staf --> UC11
Staf --> UC12
Staf --> UC14

Lurah --> UC4
Lurah --> UC13
Lurah --> UC14
@enduml
```

### Opsi B: Format Mermaid
> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ Mermaid***

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
        UC7["Menyampaikan Aduan Warga & Bukti"]
        UC8["Mencetak / Simpan PDF Surat Resmi"]
        UC9["Verifikasi Akun Warga (Setujui / Tolak Bersebab)"]
        UC10["Memeriksa Berkas Permohonan Surat Masuk"]
        UC11["Memproses Surat (Input No. Resmi, Minta Revisi, Tolak)"]
        UC12["Menindaklanjuti & Menjawab Aduan Warga"]
        UC13["Menyetujui Draft Surat Resmi (Approval Lurah)"]
        UC14["Melihat Statistik & Riwayat Kelurahan"]
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
    Warga --> UC2

    Staf --> UC4
    Staf --> UC9
    Staf --> UC10
    Staf --> UC11
    Staf --> UC12
    Staf --> UC14

    Lurah --> UC4
    Lurah --> UC13
    Lurah --> UC14
```

---

## 2. Activity Diagram (Alur Pengajuan Surat)

> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ Mermaid***

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
    
    ValidateServer -- Berhasil --> SaveData[Simpan Data & Berkas di uploads/ dan data.json]
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

---

## 3. Sequence Diagram (Pengajuan & Penerbitan Surat)

> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ Mermaid***

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
    S-->>B: Tampilkan detail permohonan & nama pemohon
    ST->>B: Input Nomor Surat Resmi & klik "Siapkan draft"
    B->>S: POST /admin/surat/{id} (action: draft, official_number)
    S->>DB: Update status = "draft_ready" & simpan nomor surat

    L->>B: Buka /admin (sebagai Lurah)
    L->>B: Klik "Tinjau draft ↗" (/dokumen/{id})
    B->>S: GET /dokumen/{id}
    S-->>B: Tampilkan dokumen berkop resmi (Watermark DRAFT)
    L->>B: Klik "Setujui" pada panel admin
    B->>S: POST /admin/surat/{id} (action: approve)
    S->>DB: Update status = "selesai"

    W->>B: Buka /dashboard (atau via /tracking)
    B-->>W: Tampilkan status "Selesai" & tombol "Cetak / Simpan PDF →"
    W->>B: Buka /dokumen/{id} & klik "Cetak / Simpan PDF"
    B->>B: Dialog Cetak Browser (Simpan PDF)
```

---

## 4. Class Diagram

> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ Mermaid***

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

## 5. State Machine Diagram (Status Surat & Akun)

> *Di draw.io: Klik **`+` $\rightarrow$ Advanced $\rightarrow$ Mermaid***

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
