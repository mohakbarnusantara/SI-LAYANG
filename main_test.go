package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSubjectValueFallsBackToSpecific(t *testing.T) {
	req := SuratRequest{
		Applicant: map[string]string{"nama": "Rasya"},
		Specific:  map[string]string{"Status perkawinan": "Belum pernah menikah", "Pekerjaan": "Sales"},
	}
	if got := subjectValue(req, "status_perkawinan"); got != "Belum pernah menikah" {
		t.Fatalf("status_perkawinan = %q", got)
	}
	applicant := map[string]string{"pekerjaan": "Guru"}
	mergeSpecificIntoApplicant(applicant, req.Specific)
	if applicant["status_perkawinan"] != "Belum pernah menikah" || applicant["pekerjaan"] != "Guru" {
		t.Fatalf("merge salah: %#v", applicant)
	}
	for _, f := range serviceFields("SK-BELUM-MENIKAH") {
		if isApplicantField(f) {
			t.Fatalf("field duplikat %q masih tampil di data khusus layanan", f)
		}
	}
}

func TestPasswordHash(t *testing.T) {
	hash := hashPassword("contoh-kata-sandi")
	if !verifyPassword(hash, "contoh-kata-sandi") {
		t.Fatal("kata sandi yang benar tidak diterima")
	}
	if verifyPassword(hash, "kata-sandi-lain") {
		t.Fatal("kata sandi yang salah diterima")
	}
}

func TestAllServicesHaveUniqueCodes(t *testing.T) {
	if len(services) != 30 {
		t.Fatalf("jumlah jenis surat = %d, ingin 30", len(services))
	}
	seen := make(map[string]bool, len(services))
	for _, service := range services {
		if service.Code == "" || service.Name == "" || len(service.Fields) == 0 {
			t.Errorf("spesifikasi layanan tidak lengkap: %#v", service)
		}
		if seen[service.Code] {
			t.Errorf("kode layanan duplikat: %s", service.Code)
		}
		seen[service.Code] = true
	}
	if len(complaintCategories) != 12 {
		t.Fatalf("jumlah kategori aduan = %d, ingin 12", len(complaintCategories))
	}
}

func TestStorePersistsUpdates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "data.json")
	store, err := openStore(path)
	if err != nil {
		t.Fatal(err)
	}
	user := User{ID: "warga-1", Email: "warga@example.test", Role: "warga", CreatedAt: time.Now()}
	if err := store.update(func(db *database) error {
		db.Users[user.ID] = user
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	reloaded, err := openStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := reloaded.snapshot().Users[user.ID]; got.Email != user.Email {
		t.Fatalf("pengguna tidak tersimpan: %#v", got)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded database
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("berkas data tidak valid JSON: %v", err)
	}
}

func TestSessionSignatureAndExpiration(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	user := User{ID: "citizen-1", Email: "warga@example.test", Role: "warga", AccountStatus: "aktif"}
	if err := store.update(func(db *database) error {
		db.Users[user.ID] = user
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{store: store, secret: []byte("test-secret")}
	response := httptest.NewRecorder()
	if err := app.setSession(response, user.ID); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("GET", "/dashboard", nil)
	request.AddCookie(response.Result().Cookies()[0])
	if got := app.currentUser(request); got == nil || got.ID != user.ID {
		t.Fatal("sesi valid tidak ditemukan")
	}
	request.Header.Set("Cookie", sessionName+"=tampered.signature")
	if got := app.currentUser(request); got != nil {
		t.Fatal("sesi dengan tanda tangan tidak valid diterima")
	}
}

func TestSeedCreatesScopedStaffAndMayorAccounts(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := seedDemoUsers(store); err != nil {
		t.Fatal(err)
	}
	users := store.snapshot().Users
	if len(users) != 3+2*len(letterOffices) {
		t.Fatalf("jumlah akun = %d, ingin %d", len(users), 3+2*len(letterOffices))
	}
	for _, office := range letterOffices {
		for _, role := range []string{"staf", "lurah"} {
			account, ok := users[officeAccountID(role, office)]
			if !ok {
				t.Fatalf("akun %s untuk %s tidak dibuat", role, office.Neighborhood)
			}
			if account.Role != role || account.OfficeDistrict != office.District ||
				account.OfficeNeighborhood != office.Neighborhood || account.AccountStatus != "aktif" ||
				account.Email != officeAccountEmail(role, office) || account.PasswordHash == "" {
				t.Errorf("akun %s tidak memiliki scope yang benar: %#v", role, account)
			}
		}
	}
	if !verifyPassword(users[officeAccountID("staf", letterOffices[0])].PasswordHash, "silayang123") {
		t.Fatal("password demo akun kelurahan tidak sesuai")
	}
	for _, id := range []string{"demo-staf", "demo-lurah"} {
		if user := users[id]; user.OfficeDistrict != "Mantikulore" || user.OfficeNeighborhood != "Talise" {
			t.Errorf("akun lama %s tidak dibatasi ke Talise: %#v", id, user)
		}
	}
}

func TestSeedMapsLegacyLetterOnlyWhenAddressNamesKelurahan(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	requests := []SuratRequest{
		{ID: "known", Applicant: map[string]string{"alamat": "Jl. Mawar, Kel. Talise, Palu"}},
		{ID: "unknown", Applicant: map[string]string{"alamat": "Jalan Mawar, Kota Palu"}},
	}
	if err := store.update(func(db *database) error {
		db.Users["resident"] = User{ID: "resident", Role: "warga"}
		db.Users["known-resident"] = User{
			ID: "known-resident", Role: "warga", Address: "Jl. Mawar, Kel. Talise, Palu",
		}
		db.Users["unknown-resident"] = User{
			ID: "unknown-resident", Role: "warga", Address: "Jalan Mawar, Kota Palu",
		}
		for _, request := range requests {
			db.Requests[request.ID] = request
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := seedDemoUsers(store); err != nil {
		t.Fatal(err)
	}
	snapshot := store.snapshot()
	if got := snapshot.Requests["known"]; got.District != "Mantikulore" || got.Neighborhood != "Talise" {
		t.Errorf("surat lama yang punya kelurahan eksplisit tidak dipetakan: %#v", got)
	}
	if got := snapshot.Requests["unknown"]; got.District != "" || got.Neighborhood != "" {
		t.Errorf("surat tanpa kelurahan eksplisit dipetakan secara spekulatif: %#v", got)
	}
	if got := snapshot.Users["known-resident"]; got.HomeDistrict != "Mantikulore" || got.HomeNeighborhood != "Talise" {
		t.Errorf("domisili warga lama yang menyebut kelurahan tidak dipetakan: %#v", got)
	}
	if got := snapshot.Users["unknown-resident"]; got.HomeDistrict != "" || got.HomeNeighborhood != "" {
		t.Errorf("domisili tanpa kelurahan eksplisit dipetakan secara spekulatif: %#v", got)
	}
}

func registerRequest(t *testing.T, fields map[string]string, withKTP bool) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatal(err)
		}
	}
	if withKTP {
		part, err := writer.CreateFormFile("ktp", "ktp.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte("\x89PNG\r\n\x1a\n0000IHDR")); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/register", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func validRegisterFields() map[string]string {
	return map[string]string{
		"name": "Warga Tatura", "nik": "1234567890123456", "email": "tatura@example.test",
		"phone": "081234567890", "birth_date": "1990-01-01", "gender": "L",
		"address": "Jl. Anoa", "home_district": "Palu Selatan",
		"home_neighborhood": "Tatura Utara", "password": "rahasia123",
	}
}

func TestRegisterStoresSelectedHomeOffice(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{store: store, secret: []byte("test-secret"), uploadDir: t.TempDir()}
	fields := validRegisterFields()
	request := registerRequest(t, fields, true)
	response := httptest.NewRecorder()

	app.register(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("pendaftaran berstatus %d, ingin 303: %s", response.Code, response.Body.String())
	}
	var resident User
	for _, user := range store.snapshot().Users {
		if user.Email == fields["email"] {
			resident = user
			break
		}
	}
	if resident.ID == "" || resident.HomeDistrict != "Palu Selatan" || resident.HomeNeighborhood != "Tatura Utara" {
		t.Fatalf("domisili pilihan tidak tersimpan: %#v", resident)
	}
	if resident.KTPPath == "" {
		t.Fatalf("foto KTP tidak tersimpan: %#v", resident)
	}
}

func TestRegisterRejectsInvalidNameAndMissingKTP(t *testing.T) {
	cases := map[string]struct {
		override map[string]string
		withKTP  bool
		message  string
	}{
		"nama berisi angka":     {map[string]string{"name": "Muh Akbar Nusantara1233"}, true, "tidak boleh mengandung angka"},
		"nama berisi simbol":    {map[string]string{"name": "Warga@Tatura!"}, true, "hanya boleh berisi huruf"},
		"tanda beruntun":        {map[string]string{"name": "Warga--Tatura"}, true, "hanya boleh berisi huruf"},
		"nomor hp berisi huruf": {map[string]string{"phone": "0812abcd5678"}, true, "Nomor HP / WhatsApp"},
		"email tidak valid":     {map[string]string{"email": "bukan-email"}, true, "Format email tidak valid"},
		"tempat lahir angka":    {map[string]string{"birth_place": "Palu123"}, true, "Tempat lahir hanya boleh berisi huruf"},
		"pekerjaan simbol":      {map[string]string{"occupation": "PNS@2026"}, true, "Pekerjaan hanya boleh berisi huruf"},
		"alamat simbol aneh":    {map[string]string{"address": "Jl. Mawar <script>"}, true, "Alamat domisili mengandung karakter"},
		"foto KTP tidak diisi":  {nil, false, "Foto KTP wajib diunggah"},
	}
	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
			if err != nil {
				t.Fatal(err)
			}
			app := &App{
				store: store, secret: []byte("test-secret"), uploadDir: t.TempDir(),
				templates: template.Must(template.New("auth").Parse(`{{define "auth"}}{{.Error}}{{end}}`)),
			}
			fields := validRegisterFields()
			for k, v := range tc.override {
				fields[k] = v
			}
			before := len(store.snapshot().Users)
			response := httptest.NewRecorder()

			app.register(response, registerRequest(t, fields, tc.withKTP))

			if response.Code == http.StatusSeeOther {
				t.Fatalf("pendaftaran seharusnya ditolak")
			}
			if !strings.Contains(response.Body.String(), tc.message) {
				t.Fatalf("pesan galat %q tidak memuat %q", response.Body.String(), tc.message)
			}
			if after := len(store.snapshot().Users); after != before {
				t.Fatalf("akun tetap tersimpan meski data tidak valid (%d -> %d)", before, after)
			}
		})
	}
}

func TestValidPersonName(t *testing.T) {
	valid := []string{
		"Budi", "Siti Aminah", "ANDI baso", "H. Ahmad", "Muh.Rizal", "Ahmad S.",
		"Ma'ruf Amin", "O'Neil", "Siti-Aminah", "Abd. Rahman Al-Habsyi",
		"Budi Santoso, S.Kom", "Budi Santoso, S.Kom.", "Dr. Ir. Budi, S.Kom., M.T.",
		"Andi Baso, S.Pd, M.Pd", "Rina, A.Md.Kom", "Prof. Dr. H. Ahmad Al-Habsyi, S.H., M.H.",
		"Budi Santoso , S.Kom",
	}
	invalid := []string{
		"", " ", "Budi1", "Andi_", "Budi@", ".", "-Budi", "'Budi", ",Budi",
		"Budi'", "Budi-", "Budi,", "Budi--Santoso", "Budi..", "Budi.-Santoso", "Budi - Santoso", "Budi ' Santoso",
		"Budi,, S.Kom", "Budi, , S.Kom", "Budi,.S.Kom", "Budi, S1", "Budi (S.Kom)",
	}
	for _, value := range valid {
		if !validPersonName(normalizePersonName(value)) {
			t.Errorf("%q seharusnya valid", value)
		}
	}
	for _, value := range invalid {
		if validPersonName(normalizePersonName(value)) {
			t.Errorf("%q seharusnya tidak valid", value)
		}
	}
	if got := normalizePersonName("  Ma\u2019ruf   Amin "); got != "Ma'ruf Amin" {
		t.Errorf("normalisasi nama menghasilkan %q", got)
	}
	if got := normalizePersonName("Budi Santoso  , S.Kom"); got != "Budi Santoso, S.Kom" {
		t.Errorf("normalisasi spasi sebelum koma menghasilkan %q", got)
	}
}

func TestDashboardPendingNoticeUsesLatestAccountStatus(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	users := []User{
		{ID: "pending-resident", Role: "warga", AccountStatus: "pending"},
		{ID: "verified-resident", Role: "warga", AccountStatus: "aktif"},
	}
	if err := store.update(func(db *database) error {
		for _, user := range users {
			db.Users[user.ID] = user
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{
		store: store, secret: []byte("test-secret"),
		templates: template.Must(template.New("dashboard").Parse(
			`{{define "dashboard"}}{{.Notice}}{{end}}`,
		)),
	}
	for _, testCase := range []struct {
		userID  string
		want    string
		notWant string
	}{
		{
			userID: "pending-resident",
			want:   "Akun Anda menunggu verifikasi staf kelurahan",
		},
		{
			userID:  "verified-resident",
			want:    "Akun Anda sudah diverifikasi",
			notWant: "menunggu verifikasi",
		},
	} {
		cookieResponse := httptest.NewRecorder()
		if err := app.setSession(cookieResponse, testCase.userID); err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest("GET", "/dashboard?notice=pending", nil)
		request.AddCookie(cookieResponse.Result().Cookies()[0])
		response := httptest.NewRecorder()
		app.dashboard(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), testCase.want) {
			t.Errorf("pesan untuk %s tidak sesuai (%d): %q", testCase.userID, response.Code, response.Body.String())
		}
		if testCase.notWant != "" && strings.Contains(response.Body.String(), testCase.notWant) {
			t.Errorf("pesan lama masih terlihat untuk %s: %q", testCase.userID, response.Body.String())
		}
	}
}

func TestDashboardShowsPrintLinkOnlyForCompletedLetters(t *testing.T) {
	templates, err := template.New("dashboard").Funcs(template.FuncMap{
		"statusClass": func(string) string { return "success" },
		"statusName":  func(status string) string { return status },
		"date":        func(value time.Time) string { return value.Format("02 Jan 2006, 15:04") },
	}).ParseFiles(filepath.Join("templates", "dashboard.html"))
	if err != nil {
		t.Fatal(err)
	}
	data := PageData{
		User: &User{Role: "warga", AccountStatus: "aktif"},
		Requests: []SuratRequest{
			{ID: "completed-letter", ServiceName: "Surat Selesai", TrackingCode: "SK-001", Status: "selesai"},
			{ID: "pending-letter", ServiceName: "Surat Menunggu", TrackingCode: "SK-002", Status: "pending_verifikasi"},
		},
	}
	var output bytes.Buffer
	if err := templates.ExecuteTemplate(&output, "dashboard-content", data); err != nil {
		t.Fatal(err)
	}
	rendered := output.String()
	if !strings.Contains(rendered, `href="/dokumen/completed-letter"`) ||
		!strings.Contains(rendered, "Cetak / Simpan PDF") {
		t.Fatalf("surat selesai tidak memiliki tautan cetak: %s", rendered)
	}
	if strings.Contains(rendered, `href="/dokumen/pending-letter"`) {
		t.Fatal("tautan cetak tampil sebelum surat selesai")
	}
}

func TestStaffResidentVerificationAndKTPAreOfficeScoped(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	staff := User{
		ID: "staff-tondo", Role: "staf", AccountStatus: "aktif",
		OfficeDistrict: "Mantikulore", OfficeNeighborhood: "Tondo",
	}
	ownResident := User{
		ID: "resident-tondo", Name: "Warga Tondo", Role: "warga", AccountStatus: "pending",
		HomeDistrict: "Mantikulore", HomeNeighborhood: "Tondo", KTPPath: "/uploads/tondo.png",
	}
	foreignResident := User{
		ID: "resident-tatura", Name: "Warga Tatura Utara", Role: "warga", AccountStatus: "pending",
		HomeDistrict: "Palu Selatan", HomeNeighborhood: "Tatura Utara", KTPPath: "/uploads/tatura.png",
	}
	if err := store.update(func(db *database) error {
		db.Users[staff.ID] = staff
		db.Users[ownResident.ID] = ownResident
		db.Users[foreignResident.ID] = foreignResident
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	uploadDir := t.TempDir()
	for _, name := range []string{"tondo.png", "tatura.png"} {
		if err := os.WriteFile(filepath.Join(uploadDir, name), []byte(name), 0600); err != nil {
			t.Fatal(err)
		}
	}
	app := &App{
		store: store, secret: []byte("test-secret"), uploadDir: uploadDir,
		templates: template.Must(template.New("dashboard").Parse(
			`{{define "admin"}}{{range .Residents}}{{.ID}}|{{end}}{{end}}{{define "dashboard"}}{{.Stats.Residents}}{{end}}`,
		)),
	}
	cookieResponse := httptest.NewRecorder()
	if err := app.setSession(cookieResponse, staff.ID); err != nil {
		t.Fatal(err)
	}
	cookie := cookieResponse.Result().Cookies()[0]

	pageRequest := httptest.NewRequest("GET", "/admin", nil)
	pageRequest.AddCookie(cookie)
	pageResponse := httptest.NewRecorder()
	app.adminPage(pageResponse, pageRequest)
	if pageResponse.Code != http.StatusOK || strings.Contains(pageResponse.Body.String(), foreignResident.ID) ||
		!strings.Contains(pageResponse.Body.String(), ownResident.ID) {
		t.Fatalf("daftar verifikasi tidak terbatasi ke Tondo: %d %q", pageResponse.Code, pageResponse.Body.String())
	}
	dashboardRequest := httptest.NewRequest("GET", "/dashboard", nil)
	dashboardRequest.AddCookie(cookie)
	dashboardResponse := httptest.NewRecorder()
	app.dashboard(dashboardResponse, dashboardRequest)
	if dashboardResponse.Code != http.StatusOK || dashboardResponse.Body.String() != "1" {
		t.Fatalf("jumlah warga menunggu tidak terbatasi ke Tondo: %d %q", dashboardResponse.Code, dashboardResponse.Body.String())
	}

	reviewRequest := httptest.NewRequest("POST", "/admin/warga/"+foreignResident.ID,
		strings.NewReader("action=approve"))
	reviewRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reviewRequest.AddCookie(cookie)
	reviewRequest.SetPathValue("id", foreignResident.ID)
	reviewResponse := httptest.NewRecorder()
	app.reviewResident(reviewResponse, reviewRequest)
	if reviewResponse.Code != http.StatusBadRequest {
		t.Fatalf("verifikasi warga kelurahan lain berstatus %d, ingin ditolak", reviewResponse.Code)
	}
	if got := store.snapshot().Users[foreignResident.ID].AccountStatus; got != "pending" {
		t.Fatalf("status warga kelurahan lain berubah menjadi %q", got)
	}

	ownReviewRequest := httptest.NewRequest("POST", "/admin/warga/"+ownResident.ID,
		strings.NewReader("action=approve"))
	ownReviewRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ownReviewRequest.AddCookie(cookie)
	ownReviewRequest.SetPathValue("id", ownResident.ID)
	ownReviewResponse := httptest.NewRecorder()
	app.reviewResident(ownReviewResponse, ownReviewRequest)
	if ownReviewResponse.Code != http.StatusSeeOther {
		t.Fatalf("verifikasi warga Tondo berstatus %d, ingin 303", ownReviewResponse.Code)
	}

	foreignKTPRequest := httptest.NewRequest("GET", "/uploads/tatura.png", nil)
	foreignKTPRequest.AddCookie(cookie)
	foreignKTPRequest.SetPathValue("name", "tatura.png")
	foreignKTPResponse := httptest.NewRecorder()
	app.upload(foreignKTPResponse, foreignKTPRequest)
	if foreignKTPResponse.Code != http.StatusForbidden {
		t.Fatalf("akses KTP kelurahan lain berstatus %d, ingin 403", foreignKTPResponse.Code)
	}

	ownKTPRequest := httptest.NewRequest("GET", "/uploads/tondo.png", nil)
	ownKTPRequest.AddCookie(cookie)
	ownKTPRequest.SetPathValue("name", "tondo.png")
	ownKTPResponse := httptest.NewRecorder()
	app.upload(ownKTPResponse, ownKTPRequest)
	if ownKTPResponse.Code != http.StatusOK || ownKTPResponse.Body.String() != "tondo.png" {
		t.Fatalf("akses KTP Tondo berstatus %d: %q", ownKTPResponse.Code, ownKTPResponse.Body.String())
	}
}

func TestOfficeStaffCannotReviewViewOrDownloadAnotherOfficeLetter(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	staff := User{
		ID: "staff-tondo", Role: "staf", AccountStatus: "aktif",
		OfficeDistrict: "Mantikulore", OfficeNeighborhood: "Tondo",
	}
	foreignRequest := SuratRequest{
		ID: "letter-talise", UserID: "resident-1", District: "Mantikulore",
		Neighborhood: "Talise", Status: "selesai", Attachments: []string{"/uploads/private.pdf"},
	}
	ownRequest := SuratRequest{
		ID: "letter-tondo", UserID: "resident-2", District: "Mantikulore",
		Neighborhood: "Tondo", Status: "pending_verifikasi",
	}
	if err := store.update(func(db *database) error {
		db.Users[staff.ID] = staff
		db.Requests[foreignRequest.ID] = foreignRequest
		db.Requests[ownRequest.ID] = ownRequest
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{store: store, secret: []byte("test-secret"), uploadDir: t.TempDir()}
	cookieResponse := httptest.NewRecorder()
	if err := app.setSession(cookieResponse, staff.ID); err != nil {
		t.Fatal(err)
	}
	cookie := cookieResponse.Result().Cookies()[0]

	reviewRequest := httptest.NewRequest("POST", "/admin/surat/letter-talise",
		strings.NewReader("action=draft&official_number=TEST-001"))
	reviewRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reviewRequest.AddCookie(cookie)
	reviewRequest.SetPathValue("id", foreignRequest.ID)
	reviewResponse := httptest.NewRecorder()
	app.reviewSurat(reviewResponse, reviewRequest)
	if reviewResponse.Code != http.StatusBadRequest {
		t.Fatalf("aksi surat lintas kelurahan berstatus %d, ingin ditolak", reviewResponse.Code)
	}
	if got := store.snapshot().Requests[foreignRequest.ID].Status; got != "selesai" {
		t.Fatalf("status surat kelurahan lain berubah menjadi %q", got)
	}

	ownReviewRequest := httptest.NewRequest("POST", "/admin/surat/letter-tondo",
		strings.NewReader("action=draft&official_number=TONDO-001"))
	ownReviewRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ownReviewRequest.AddCookie(cookie)
	ownReviewRequest.SetPathValue("id", ownRequest.ID)
	ownReviewResponse := httptest.NewRecorder()
	app.reviewSurat(ownReviewResponse, ownReviewRequest)
	if ownReviewResponse.Code != http.StatusSeeOther {
		t.Fatalf("aksi surat kelurahan sendiri berstatus %d, ingin 303: %s", ownReviewResponse.Code, ownReviewResponse.Body.String())
	}
	if got := store.snapshot().Requests[ownRequest.ID].Status; got != "draft_ready" {
		t.Fatalf("status surat kelurahan sendiri = %q, ingin draft_ready", got)
	}

	documentRequest := httptest.NewRequest("GET", "/dokumen/letter-talise", nil)
	documentRequest.AddCookie(cookie)
	documentRequest.SetPathValue("id", foreignRequest.ID)
	documentResponse := httptest.NewRecorder()
	app.document(documentResponse, documentRequest)
	if documentResponse.Code != http.StatusNotFound {
		t.Fatalf("dokumen kelurahan lain berstatus %d, ingin 404", documentResponse.Code)
	}

	uploadRequest := httptest.NewRequest("GET", "/uploads/private.pdf", nil)
	uploadRequest.AddCookie(cookie)
	uploadRequest.SetPathValue("name", "private.pdf")
	uploadResponse := httptest.NewRecorder()
	app.upload(uploadResponse, uploadRequest)
	if uploadResponse.Code != http.StatusForbidden {
		t.Fatalf("lampiran kelurahan lain berstatus %d, ingin 403", uploadResponse.Code)
	}
}

func TestApprovedDocumentShowsDummyQRAndNIPWithoutVerificationCode(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	office, ok := letterOffice("Palu Selatan", "Tatura Utara")
	if !ok {
		t.Fatal("kantor Tatura Utara tidak ditemukan")
	}
	resident := User{ID: "resident-1", Role: "warga", AccountStatus: "aktif"}
	request := SuratRequest{
		ID: "approved-letter", TrackingCode: "SK-TEST-001", UserID: "resident-1",
		ServiceName: "Surat Keterangan", District: office.District, Neighborhood: office.Neighborhood,
		Applicant: map[string]string{
			"nama": "Warga Contoh", "nik": "1234567890123456", "tanggal_lahir": "Palu, 01 Januari 2000",
			"jenis_kelamin": "Laki-laki", "alamat": "Jl. Anoa",
		},
		OfficialNumber: "001/2026", Status: "selesai", UpdatedAt: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC),
	}
	if err := store.update(func(db *database) error {
		db.Users[resident.ID] = resident
		db.Requests[request.ID] = request
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	templates, err := template.New("document").Funcs(template.FuncMap{
		"date": func(value time.Time) string { return value.Format("02 Jan 2006, 15:04") },
	}).ParseFiles(filepath.Join("templates", "document.html"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{store: store, secret: []byte("test-secret"), templates: templates}
	cookieResponse := httptest.NewRecorder()
	if err := app.setSession(cookieResponse, resident.ID); err != nil {
		t.Fatal(err)
	}
	httpRequest := httptest.NewRequest("GET", "/dokumen/"+request.ID, nil)
	httpRequest.AddCookie(cookieResponse.Result().Cookies()[0])
	httpRequest.SetPathValue("id", request.ID)
	response := httptest.NewRecorder()

	app.document(response, httpRequest)

	output := response.Body.String()
	if response.Code != http.StatusOK {
		t.Fatalf("dokumen berstatus %d: %s", response.Code, output)
	}
	for _, expected := range []string{
		"QR SIMULASI", "NIP. 000000000000000000 (DUMMY)", office.Mayor,
		"tidak dapat digunakan untuk verifikasi",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("dokumen tidak menampilkan %q", expected)
		}
	}
	for _, removed := range []string{"Kode verifikasi:", "TTE DIGITAL (SIMULASI)", "qr_hash"} {
		if strings.Contains(output, removed) {
			t.Errorf("dokumen masih menampilkan %q", removed)
		}
	}
}

func TestLetterOfficeDirectoryHasUniqueAreas(t *testing.T) {
	if len(letterOffices) != 32 {
		t.Fatalf("jumlah kantor kelurahan = %d, ingin 32 dari spreadsheet", len(letterOffices))
	}
	seen := make(map[string]bool, len(letterOffices))
	for _, office := range letterOffices {
		key := office.District + "/" + office.Neighborhood
		if office.District == "" || office.Neighborhood == "" || office.Name == "" || office.Mayor == "" || office.Address == "" {
			t.Errorf("data kantor tidak lengkap: %#v", office)
		}
		if seen[key] {
			t.Errorf("pasangan kecamatan/kelurahan duplikat: %s", key)
		}
		seen[key] = true
	}
	tondo, ok := letterOffice("Mantikulore", "Tondo")
	if !ok || tondo.Address != "Jl. Soekarno-Hatta No. 1, Kota Palu, Sulawesi Tengah 94118" ||
		officeLetterEmail(tondo) != "tondo@palukota.go.id" {
		t.Errorf("alamat atau email Tondo tidak mengikuti spesifikasi kop: %#v", tondo)
	}
}

func TestSuratFormDoesNotExposeDraftPreviewToCitizens(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("templates", "surat_form.html"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, forbidden := range []string{"preview-draft", "draft-preview", "Pratinjau draf", "draft-paper"} {
		if strings.Contains(text, forbidden) {
			t.Errorf("formulir surat masih menampilkan pratinjau draf (%q)", forbidden)
		}
	}
}

func TestRejectedResidentCanResubmitVerification(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		ID: "rejected-1", NIK: "7271010101900099", Name: "Warga Lama",
		Email: "rejected@example.test", Phone: "081234567890", BirthPlace: "Palu",
		BirthDate: "1990-01-01", Gender: "L", Occupation: "Wiraswasta",
		HomeDistrict: "Palu Selatan", HomeNeighborhood: "Tatura Utara",
		Address: "Jl. Anoa Lama", Role: "warga", AccountStatus: "ditolak",
		RejectionReason: "Foto KTP buram dan nama belum sesuai",
	}
	if err := store.update(func(db *database) error {
		db.Users[user.ID] = user
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{
		store: store, secret: []byte("test-secret"), uploadDir: t.TempDir(),
		templates: template.Must(template.New("dashboard").Parse(`{{define "dashboard"}}{{.Error}}|{{.Notice}}{{end}}`)),
	}
	cookieRec := httptest.NewRecorder()
	if err := app.setSession(cookieRec, user.ID); err != nil {
		t.Fatal(err)
	}
	sessionCookie := cookieRec.Result().Cookies()[0]

	// 1. Jika tidak ada KTP sama sekali dan tidak mengunggah KTP baru -> ditolak
	fields := validRegisterFields()
	fields["name"] = "Warga Perbaikan, S.Kom."
	badReq := registerRequest(t, fields, false)
	badReq.URL.Path = "/akun/verifikasi-ulang"
	badReq.AddCookie(sessionCookie)
	badRes := httptest.NewRecorder()
	app.resubmitVerification(badRes, badReq)
	if badRes.Code != http.StatusBadRequest || !strings.Contains(badRes.Body.String(), "Foto KTP wajib diunggah") {
		t.Fatalf("seharusnya gagal tanpa KTP: code=%d body=%q", badRes.Code, badRes.Body.String())
	}

	// 2. Mengirim perbaikan lengkap dengan foto KTP -> status kembali pending & alasan bersih
	goodReq := registerRequest(t, fields, true)
	goodReq.URL.Path = "/akun/verifikasi-ulang"
	goodReq.AddCookie(sessionCookie)
	goodRes := httptest.NewRecorder()
	app.resubmitVerification(goodRes, goodReq)
	if goodRes.Code != http.StatusSeeOther {
		t.Fatalf("verifikasi ulang gagal: code=%d body=%q", goodRes.Code, goodRes.Body.String())
	}
	updated := store.snapshot().Users[user.ID]
	if updated.AccountStatus != "pending" || updated.RejectionReason != "" {
		t.Fatalf("status akun setelah verifikasi ulang = %q (alasan %q), ingin pending", updated.AccountStatus, updated.RejectionReason)
	}
	if updated.Name != "Warga Perbaikan, S.Kom." || updated.KTPPath == "" {
		t.Fatalf("data perbaikan tidak tersimpan: %#v", updated)
	}
}

func TestIncomeLetterFieldsAndRendering(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		ID: "user-income", NIK: "7271010101900001", Name: "Orang Tua Kandung",
		Role: "warga", AccountStatus: "aktif",
	}
	if err := store.update(func(db *database) error {
		db.Users[user.ID] = user
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	app := &App{
		store: store, secret: []byte("test-secret"), uploadDir: t.TempDir(),
		templates: template.Must(template.New("surat_form").Parse(`{{define "surat_form"}}{{.Error}}{{end}}`)),
	}
	cookieRec := httptest.NewRecorder()
	if err := app.setSession(cookieRec, user.ID); err != nil {
		t.Fatal(err)
	}
	sessionCookie := cookieRec.Result().Cookies()[0]

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	formValues := map[string]string{
		"service_code":                     "SK-PENGHASILAN",
		"kecamatan":                        "Mantikulore",
		"kelurahan":                        "Talise",
		"applicant_nama":                   "Budi Santoso",
		"applicant_jenis_kelamin":          "Laki-laki",
		"applicant_tempat_lahir":           "Palu",
		"applicant_tanggal_lahir":          "1978-05-12",
		"applicant_agama":                  "Islam",
		"applicant_status_perkawinan":      "Kawin",
		"applicant_nik":                    "7271010101780001",
		"applicant_pekerjaan":              "Wiraswasta",
		"applicant_alamat":                 "Jl. Diponegoro No. 10",
		"specific_nama_anak":               "Ahmad Fadilah",
		"specific_jenis_kelamin_anak":      "Laki-laki",
		"specific_tempat_lahir_anak":       "Palu",
		"specific_tanggal_lahir_anak":      "2005-08-20",
		"specific_agama_anak":              "Islam",
		"specific_status_perkawinan_anak":  "Belum Kawin",
		"specific_nik_anak":                "7271010101050002",
		"specific_pekerjaan_anak":          "Pelajar/Mahasiswa",
		"specific_alamat_anak":             "Jl. Diponegoro No. 10",
		"specific_penghasilan_rata_rata":   "3500000",
		"specific_sumber_penghasilan":      "Perdagangan Sembako",
		"rt_letter_number":                 "012/RT-01/VI/2026",
	}
	for k, v := range formValues {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("rt_attachment", "pengantar.png")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDRdummy"))
	writer.Close()

	req := httptest.NewRequest("POST", "/surat", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(sessionCookie)
	rec := httptest.NewRecorder()
	app.createSurat(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("createSurat SK-PENGHASILAN gagal: code=%d body=%q", rec.Code, rec.Body.String())
	}

	snap := store.snapshot()
	var created SuratRequest
	for _, r := range snap.Requests {
		if r.ServiceCode == "SK-PENGHASILAN" {
			created = r
			break
		}
	}
	if created.ID == "" {
		t.Fatal("request SK-PENGHASILAN tidak ditemukan di database")
	}
	if created.Applicant["nama"] != "Budi Santoso" || created.Applicant["status_perkawinan"] != "Kawin" {
		t.Fatalf("applicant data salah: %#v", created.Applicant)
	}
	if created.Specific["Nama Anak"] != "Ahmad Fadilah" || created.Specific["Penghasilan Rata-Rata"] != "3500000" {
		t.Fatalf("specific data salah: %#v", created.Specific)
	}

	// Cek render document.html
	docTmpl, err := template.New("document").Funcs(template.FuncMap{
		"date": func(value time.Time) string { return value.Format("02 Jan 2006, 15:04") },
	}).ParseFiles(filepath.Join("templates", "document.html"))
	if err != nil {
		t.Fatal(err)
	}
	var docBuf bytes.Buffer
	docData := struct {
		Request     SuratRequest
		IsDraft     bool
		VerifyURL   string
		Office      LetterOffice
		OfficeEmail string
	}{
		Request: created,
		Office:  LetterOffice{District: "Mantikulore", Neighborhood: "Talise", Address: "Jl. Talise", Mayor: "Lurah Talise"},
	}
	if err := docTmpl.ExecuteTemplate(&docBuf, "document", docData); err != nil {
		t.Fatal(err)
	}
	docStr := docBuf.String()
	for _, expected := range []string{
		"I. Data Orang Tua / Wali:", "Budi Santoso", "Kawin",
		"II. Data Anak:", "Ahmad Fadilah", "Pelajar/Mahasiswa",
		"III. Data Penghasilan:", "3500000", "Perdagangan Sembako",
	} {
		if !strings.Contains(docStr, expected) {
			t.Errorf("document.html tidak memuat %q", expected)
		}
	}
}

func TestLetterRevisionAndEdit(t *testing.T) {
	store, err := openStore(filepath.Join(t.TempDir(), "data.json"))
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		ID: "user-rev", NIK: "7271010101900002", Name: "Warga Revisi",
		Role: "warga", AccountStatus: "aktif",
	}
	staff := User{
		ID: "staff-rev", NIK: "7271010101800009", Name: "Staf Pemeriksa",
		Role: "staf", AccountStatus: "aktif", OfficeDistrict: "Palu Selatan", OfficeNeighborhood: "Tatura Utara",
	}
	request := SuratRequest{
		ID: "letter-rev-1", TrackingCode: "REQ-REV-001", UserID: user.ID,
		ServiceCode: "SK-EKONOMI-LEMAH", ServiceName: "Surat Keterangan Ekonomi Lemah",
		District: "Palu Selatan", Neighborhood: "Tatura Utara",
		Applicant: map[string]string{
			"nama": "Warga Revisi", "nik": "7271010101900002", "agama": "Lainnya",
			"alamat": "Jl. Anoa", "jenis_kelamin": "Laki-laki", "tempat_lahir": "Palu", "tanggal_lahir": "1995-01-01",
		},
		Purpose: "Bantuan UKT", RTLetterNumber: "01/RT-01/2026",
		Status: "perlu_revisi", AdminNote: "kasih betul agama dan alasan nya",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := store.update(func(db *database) error {
		db.Users[user.ID] = user
		db.Users[staff.ID] = staff
		db.Requests[request.ID] = request
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	app := &App{
		store: store, secret: []byte("test-secret"), uploadDir: t.TempDir(),
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"roleName": func(role string) string { return role },
			"statusName": statusName,
			"statusClass": func(s string) string { return s },
			"fieldKey": fieldKey,
			"join": strings.Join,
			"date": func(v time.Time) string { return v.Format("02 Jan 2006, 15:04") },
			"lookupUser": func(id string) User { return store.snapshot().Users[id] },
			"isDateField": func(field string) bool { return strings.Contains(strings.ToLower(field), "tanggal") },
			"isRevisionField": isRevisionField,
			"serviceFields":   serviceFields,
			"subjectValue":     subjectValue,
			"isApplicantField": isApplicantField,
			"isRevised":        isRevisedByResident,
		}).ParseGlob(filepath.Join("templates", "*.html"))),
	}

	// 0. Staf melihat panel admin: data lengkap tampil & surat terbaru di atas
	older := request
	older.ID, older.TrackingCode, older.Status = "letter-old", "REQ-OLD-000", "selesai"
	older.CreatedAt = time.Now().Add(-48 * time.Hour)
	if err := store.update(func(db *database) error { db.Requests[older.ID] = older; return nil }); err != nil {
		t.Fatal(err)
	}
	staffRec := httptest.NewRecorder()
	if err := app.setSession(staffRec, staff.ID); err != nil {
		t.Fatal(err)
	}
	adminReq := httptest.NewRequest("GET", "/admin", nil)
	adminReq.AddCookie(staffRec.Result().Cookies()[0])
	adminRec := httptest.NewRecorder()
	app.adminPage(adminRec, adminReq)
	adminBody := adminRec.Body.String()
	if adminRec.Code != http.StatusOK {
		t.Fatalf("adminPage gagal: code=%d body=%q", adminRec.Code, adminBody)
	}
	for _, want := range []string{"Lihat data lengkap pemohon", "Bantuan UKT", "Lainnya", "Jl. Anoa", "1995-01-01", "Status perkawinan", "Belum direvisi warga"} {
		if !strings.Contains(adminBody, want) {
			t.Errorf("panel admin tidak menampilkan %q", want)
		}
	}
	if strings.Index(adminBody, "REQ-REV-001") > strings.Index(adminBody, "REQ-OLD-000") {
		t.Errorf("surat terbaru seharusnya tampil paling atas")
	}
	if err := store.update(func(db *database) error { delete(db.Requests, older.ID); return nil }); err != nil {
		t.Fatal(err)
	}

	// 1. Cek warga membuka halaman perbaikan surat GET /surat/{id}/edit
	cookieRec := httptest.NewRecorder()
	if err := app.setSession(cookieRec, user.ID); err != nil {
		t.Fatal(err)
	}
	sessionCookie := cookieRec.Result().Cookies()[0]

	getReq := httptest.NewRequest("GET", "/surat/"+request.ID+"/edit", nil)
	getReq.SetPathValue("id", request.ID)
	getReq.AddCookie(sessionCookie)
	getRec := httptest.NewRecorder()
	app.editSuratPage(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("editSuratPage gagal: code=%d body=%q", getRec.Code, getRec.Body.String())
	}
	if !strings.Contains(getRec.Body.String(), "kasih betul agama dan alasan nya") {
		t.Errorf("halaman edit tidak menampilkan catatan petugas")
	}
	if !strings.Contains(getRec.Body.String(), "field-revision-highlight") {
		t.Errorf("halaman edit tidak menandai kolom yang perlu diperbaiki dengan field-revision-highlight")
	}
	if !strings.Contains(getRec.Body.String(), "Perlu Diperbaiki") {
		t.Errorf("halaman edit tidak menampilkan badge Perlu Diperbaiki")
	}

	// 2. Warga mengirimkan perbaikan POST /surat/{id}/edit
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("kecamatan", "Palu Selatan")
	writer.WriteField("kelurahan", "Tatura Utara")
	writer.WriteField("service_code", "SK-EKONOMI-LEMAH")
	writer.WriteField("applicant_nama", "Warga Revisi")
	writer.WriteField("applicant_nik", "7271010101900002")
	writer.WriteField("applicant_agama", "Islam")
	writer.WriteField("applicant_jenis_kelamin", "Laki-laki")
	writer.WriteField("applicant_tempat_lahir", "Palu")
	writer.WriteField("applicant_tanggal_lahir", "1995-01-01")
	writer.WriteField("applicant_pekerjaan", "Mahasiswa")
	writer.WriteField("applicant_alamat", "Jl. Anoa No. 12")
	writer.WriteField("keperluan", "Keringanan Biaya Kuliah Semester 3")
	writer.WriteField("rt_letter_number", "01/RT-01/2026")
	writer.Close()

	postReq := httptest.NewRequest("POST", "/surat/"+request.ID+"/edit", &body)
	postReq.Header.Set("Content-Type", writer.FormDataContentType())
	postReq.SetPathValue("id", request.ID)
	postReq.AddCookie(sessionCookie)
	postRec := httptest.NewRecorder()
	app.updateSurat(postRec, postReq)

	if postRec.Code != http.StatusSeeOther {
		t.Fatalf("updateSurat gagal: code=%d body=%q", postRec.Code, postRec.Body.String())
	}

	updated := store.snapshot().Requests[request.ID]
	if updated.Status != "pending_verifikasi" {
		t.Fatalf("status setelah perbaikan = %q, ingin pending_verifikasi", updated.Status)
	}
	if updated.Applicant["agama"] != "Islam" || updated.Purpose != "Keringanan Biaya Kuliah Semester 3" {
		t.Fatalf("data surat yang diperbaiki salah: %#v", updated)
	}
	if updated.RevisionCount != 1 || updated.RevisedAt == nil {
		t.Fatalf("metadata revisi tidak tercatat: count=%d revised_at=%v", updated.RevisionCount, updated.RevisedAt)
	}

	// 3. Staf melihat ulang panel admin: harus jelas tertulis sudah direvisi warga
	adminReq2 := httptest.NewRequest("GET", "/admin", nil)
	adminReq2.AddCookie(staffRec.Result().Cookies()[0])
	adminRec2 := httptest.NewRecorder()
	app.adminPage(adminRec2, adminReq2)
	if !strings.Contains(adminRec2.Body.String(), "Sudah direvisi warga") {
		t.Errorf("panel admin tidak menampilkan penanda 'Sudah direvisi warga' setelah warga mengirim perbaikan")
	}
}
