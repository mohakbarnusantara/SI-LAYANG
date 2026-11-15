document.addEventListener("DOMContentLoaded", () => {
  const toggle = document.querySelector(".menu-toggle");
  const nav = document.querySelector(".main-nav");
  if (toggle && nav) {
    toggle.addEventListener("click", () => {
      const open = nav.classList.toggle("nav-open");
      toggle.setAttribute("aria-expanded", String(open));
    });
  }

  const selfToggle = document.getElementById("is-self");
  if (selfToggle) {
    const fields = {
      nama: { element: document.querySelector('[name="applicant_nama"]'), profile: "profileName" },
      nik: { element: document.querySelector('[name="applicant_nik"]'), profile: "profileNik" },
      tempat_lahir: { element: document.querySelector('[name="applicant_tempat_lahir"]'), profile: "profileBirthPlace" },
      tanggal_lahir: { element: document.querySelector('[name="applicant_tanggal_lahir"]'), profile: "profileBirth" },
      jenis_kelamin: { element: document.querySelector('[name="applicant_jenis_kelamin"]'), profile: "profileGender" },
      pekerjaan: { element: document.querySelector('[name="applicant_pekerjaan"]'), profile: "profileOccupation" },
      alamat: { element: document.querySelector('[name="applicant_alamat"]'), profile: "profileAddress" },
    };
    const fillFromProfile = (enabled) => {
      for (const field of Object.values(fields)) {
        const input = field.element;
        if (!input) continue;
        input.value = enabled ? input.dataset[field.profile] || "" : "";
        input.classList.remove("input-invalid");
        input.removeAttribute("aria-invalid");
        input.setCustomValidity("");
      }
    };
    selfToggle.addEventListener("change", () => fillFromProfile(selfToggle.checked));
  }

  const serviceSelect = document.getElementById("service-select");
  const fieldContainer = document.getElementById("specific-fields");
  const hint = document.getElementById("service-hint");
  const extraLabel = document.getElementById("extra-doc-label");
  if (serviceSelect && fieldContainer) {
    serviceSelect.addEventListener("change", () => {
      const option = serviceSelect.selectedOptions[0];
      fieldContainer.replaceChildren();
      if (!option || !option.value) {
        hint.textContent = "Setiap layanan dapat meminta data khusus dan berkas pendukung.";
        extraLabel?.classList.remove("required-doc");
        return;
      }
      const labels = (option.dataset.fields || "").split("|").filter(Boolean);
      hint.textContent = labels.length
        ? `Data tambahan untuk layanan ini: ${labels.join(", ")}.`
        : "Tidak ada data khusus tambahan untuk layanan ini.";
      labels.forEach((label) => {
        const wrap = document.createElement("label");
        wrap.textContent = label;
        const input = document.createElement(label.toLowerCase().includes("tanggal") ? "input" : "textarea");
        const key = label.toLowerCase().replace(/[^a-z0-9]+/g, "_").replace(/^_|_$/g, "");
        input.name = `specific_${key}`;
        input.placeholder = `Masukkan ${label.toLowerCase()}`;
        input.required = true;
        input.dataset.fieldLabel = label;
        if (input.tagName === "INPUT") input.type = "date";
        else input.rows = 2;
        wrap.append(input);
        fieldContainer.append(wrap);
      });
      if (option.dataset.extra) {
        const label = extraLabel?.querySelector("input");
        if (label) label.setAttribute("aria-label", option.dataset.extra);
        const caption = extraLabel?.querySelector(".label-help");
        if (caption) caption.textContent = `(${option.dataset.extra} · JPG, PNG, atau PDF · maks. 5 MB)`;
        extraLabel?.classList.add("required-doc");
      } else {
        extraLabel?.classList.remove("required-doc");
        const caption = extraLabel?.querySelector(".label-help");
        if (caption) caption.textContent = "(jika dipersyaratkan · maks. 5 MB)";
      }
    });
  }

  const districts = {
    Mantikulore: ["Talise", "Talise Valangguni", "Tondo", "Layana Indah", "Tanamodindi"],
    "Palu Barat": ["Balaroa", "Kamonji", "Lere", "Siranindi"],
    "Palu Selatan": ["Birobuli Selatan", "Birobuli Utara", "Petobo", "Tatura Selatan"],
    "Palu Timur": ["Besusu Barat", "Besusu Tengah", "Besusu Timur", "Lasoani"],
    "Palu Utara": ["Kayumalue Ngapa", "Mamboro", "Taipa"],
    Tatanga: ["Duyu", "Nunu", "Tatura Utara"],
    Tawaeli: ["Baiya", "Lasoani", "Panau"],
    Ulujadi: ["Buluri", "Donggala Kodi", "Silae"],
  };
  const districtSelect = document.getElementById("district-select");
  const neighborhoodSelect = document.getElementById("neighborhood-select");
  if (districtSelect && neighborhoodSelect) {
    const populateNeighborhoods = (preserveSelection) => {
      const previousNeighborhood = preserveSelection ? neighborhoodSelect.value : "";
      neighborhoodSelect.value = "";
      neighborhoodSelect.replaceChildren();
      const selectedDistrict = districtSelect.selectedOptions[0];
      const options = selectedDistrict?.dataset.neighborhoods
        ? selectedDistrict.dataset.neighborhoods.split("|").filter(Boolean)
        : districts[districtSelect.value] || [];
      const placeholder = document.createElement("option");
      placeholder.value = "";
      placeholder.textContent = options.length ? "Pilih kelurahan" : "Pilih kecamatan dahulu";
      neighborhoodSelect.append(placeholder);
      options.forEach((name) => {
        const option = document.createElement("option");
        option.value = name;
        option.textContent = name;
        neighborhoodSelect.append(option);
      });
      if (options.includes(previousNeighborhood)) {
        neighborhoodSelect.value = previousNeighborhood;
      }
      neighborhoodSelect.disabled = options.length === 0;
    };

    districtSelect.addEventListener("change", () => {
      populateNeighborhoods(false);
    });
    if (districtSelect.value) populateNeighborhoods(true);
    window.addEventListener("pageshow", () => {
      if (districtSelect.value) populateNeighborhoods(true);
    });
  }

  const inputFilters = {
    "person-name": (value) =>
      value
        .replace(/[\u2018\u2019`]/g, "'")
        .replace(/[^A-Za-z .,'\-]/g, "")
        .replace(/^[ .,'\-]+/, "")
        .replace(/([.'\-])[.'\-]+/g, "$1")
        .replace(/,{2,}/g, ",")
        .replace(/ +,/g, ",")
        .replace(/ {2,}/g, " "),
    digits: (value) => value.replace(/[^0-9]/g, ""),
    phone: (value) => {
      const hasPlus = value.startsWith("+");
      const digits = value.replace(/[^0-9]/g, "").slice(0, 15);
      return (hasPlus ? "+" : "") + digits;
    },
    email: (value) => value.replace(/\s+/g, ""),
    "place-name": (value) =>
      value
        .replace(/[\u2018\u2019`]/g, "'")
        .replace(/[^A-Za-z .,'\-]/g, "")
        .replace(/^[ .,'\-]+/, "")
        .replace(/([.,'\-])[.,'\-]+/g, "$1")
        .replace(/ {2,}/g, " "),
    occupation: (value) =>
      value
        .replace(/[^A-Za-z ./,\-]/g, "")
        .replace(/^[ ./,\-]+/, "")
        .replace(/([./,\-])[./,\-]+/g, "$1")
        .replace(/ {2,}/g, " "),
    "alpha-text": (value) =>
      value
        .replace(/[\u2018\u2019`]/g, "'")
        .replace(/[^A-Za-z '\-]/g, "")
        .replace(/^[ '\-]+/, "")
        .replace(/(['\-])['\-]+/g, "$1")
        .replace(/ {2,}/g, " "),
    address: (value) =>
      value
        .replace(/[\u2018\u2019`]/g, "'")
        .replace(/[^A-Za-z0-9 .,\-/'()\n]/g, "")
        .replace(/^[ .,\-/'()\n]+/, "")
        .replace(/ {2,}/g, " "),
    "letter-number": (value) => value.replace(/[^A-Za-z0-9/.\- ]/g, "").replace(/ {2,}/g, " "),
    "tracking-code": (value) => value.toUpperCase().replace(/[^A-Z0-9\-]/g, ""),
  };

  document.querySelectorAll("[data-filter]").forEach((field) => {
    const sanitize = inputFilters[field.dataset.filter];
    if (!sanitize) return;
    field.addEventListener("input", () => {
      const raw = field.value;
      const start = field.selectionStart;
      const cleaned = sanitize(raw);
      if (cleaned !== raw) {
        field.value = cleaned;
        if (typeof start === "number" && field.type !== "email") {
          const nextPos = sanitize(raw.slice(0, start)).length;
          field.setSelectionRange(nextPos, nextPos);
        }
      }
      field.setCustomValidity?.("");
    });
  });

  const registerForm = document.getElementById("register-form");
  const registerError = document.getElementById("register-error");
  const registerName = document.getElementById("register-name");
  if (registerName) {
    registerName.addEventListener("invalid", () => {
      const value = registerName.value.trim();
      let message = "Nama lengkap wajib diisi.";
      if (value) {
        message = /\d/.test(value)
          ? "Nama lengkap tidak boleh mengandung angka."
          : "Nama hanya boleh berisi huruf, spasi, serta tanda titik (.), koma (,) untuk gelar, apostrof ('), atau tanda hubung (-) di antara huruf.";
      }
      registerName.setCustomValidity(message);
    });
  }

  const registerKtp = document.getElementById("register-ktp");
  if (registerKtp) {
    registerKtp.addEventListener("change", () => {
      const file = registerKtp.files?.[0];
      let message = "";
      if (file && !["image/jpeg", "image/png"].includes(file.type)) message = "Foto KTP harus berformat JPG atau PNG.";
      else if (file && file.size > 5 * 1024 * 1024) message = "Ukuran foto KTP maksimal 5 MB.";
      registerKtp.setCustomValidity(message);
      if (message) registerKtp.reportValidity();
    });
    registerKtp.addEventListener("invalid", () => {
      if (!registerKtp.files?.length) registerKtp.setCustomValidity("Foto KTP wajib diunggah.");
    });
  }

  if (registerForm && registerError) {
    const nameRegex = /^[A-Za-z]+(?:(?:\.?, ?|\. ?|[ '\-])[A-Za-z]+)*\.?$/;
    const setFieldError = (field, message) => {
      if (!field) return;
      const label = field.closest("label");
      const existing = label?.querySelector(".field-error");
      if (message) {
        field.classList.add("input-invalid");
        field.setAttribute("aria-invalid", "true");
        if (label) {
          const note = existing || document.createElement("small");
          note.className = "field-error";
          note.textContent = message;
          if (!existing) label.append(note);
        }
      } else {
        field.classList.remove("input-invalid");
        field.removeAttribute("aria-invalid");
        existing?.remove();
      }
    };

    const validateRegisterField = (field) => {
      const val = (field.value || "").trim();
      switch (field.name) {
        case "name":
          if (!val) return "Nama lengkap wajib diisi.";
          if (/\d/.test(val)) return "Nama lengkap tidak boleh mengandung angka.";
          if (!nameRegex.test(val)) {
            return "Nama lengkap hanya boleh berisi huruf, spasi, serta tanda titik (.), koma (,) untuk gelar, apostrof ('), atau tanda hubung (-) di antara huruf.";
          }
          return "";
        case "nik":
          if (!val) return "NIK wajib diisi.";
          if (!/^\d{16}$/.test(val)) return "NIK harus terdiri dari 16 digit angka.";
          return "";
        case "phone":
          if (!val) return "Nomor HP / WhatsApp wajib diisi.";
          if (!/^\+?[0-9]{9,15}$/.test(val)) return "Nomor HP / WhatsApp harus berupa 9–15 digit angka (boleh diawali +).";
          return "";
        case "email":
          if (!val) return "Email aktif wajib diisi.";
          if (field.validity.typeMismatch || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val)) return "Format email tidak valid.";
          return "";
        case "birth_place":
          if (val && !/[A-Za-z]/.test(val)) return "Tempat lahir harus berisi huruf.";
          if (val && /[^A-Za-z .,'\-]/.test(val)) return "Tempat lahir hanya boleh berisi huruf dan tanda baca umum.";
          return "";
        case "birth_date":
          return val ? "" : "Tanggal lahir wajib diisi.";
        case "gender":
          return val ? "" : "Jenis kelamin wajib dipilih.";
        case "occupation":
          if (val && !/[A-Za-z]/.test(val)) return "Pekerjaan harus berisi huruf.";
          if (val && /[^A-Za-z ./,\-]/.test(val)) return "Pekerjaan hanya boleh berisi huruf, spasi, atau tanda /, -, ., ,.";
          return "";
        case "home_district":
          return val ? "" : "Kecamatan domisili wajib dipilih.";
        case "home_neighborhood":
          return val ? "" : "Kelurahan domisili wajib dipilih.";
        case "address":
          if (!val) return "Alamat domisili wajib diisi.";
          if (/[^A-Za-z0-9 .,\-/'()\n\r]/.test(val)) return "Alamat mengandung karakter yang tidak diizinkan.";
          return "";
        case "ktp": {
          const file = field.files?.[0];
          if (!file) return "Foto KTP wajib diunggah.";
          if (!["image/jpeg", "image/png"].includes(file.type)) return "Foto KTP harus berformat JPG atau PNG.";
          if (file.size > 5 * 1024 * 1024) return "Ukuran foto KTP maksimal 5 MB.";
          return "";
        }
        case "password":
          if (!field.value) return "Kata sandi wajib diisi.";
          if (field.value.length < 8) return "Kata sandi minimal 8 karakter.";
          return "";
        default:
          return "";
      }
    };

    const fields = Array.from(registerForm.querySelectorAll("input, select, textarea"));
    fields.forEach((field) => {
      const eventName = field.tagName === "SELECT" || field.type === "file" || field.type === "date" ? "change" : "input";
      field.addEventListener(eventName, () => {
        if (field.classList.contains("input-invalid")) {
          setFieldError(field, validateRegisterField(field));
        }
      });
    });

    registerForm.addEventListener("submit", (event) => {
      const errors = [];
      let firstInvalid = null;
      fields.forEach((field) => {
        const message = validateRegisterField(field);
        setFieldError(field, message);
        if (message) {
          errors.push(message);
          if (!firstInvalid) firstInvalid = field;
        }
      });

      if (errors.length > 0) {
        event.preventDefault();
        registerError.replaceChildren();
        const title = document.createElement("strong");
        title.textContent = "Periksa kembali data pendaftaran Anda:";
        const list = document.createElement("ul");
        list.className = "alert-list";
        errors.forEach((msg) => {
          const item = document.createElement("li");
          item.textContent = msg;
          list.append(item);
        });
        registerError.append(title, list);
        registerError.hidden = false;
        registerError.scrollIntoView({ behavior: "smooth", block: "nearest" });
        firstInvalid?.focus();
      } else {
        registerError.hidden = true;
      }
    });
  }

  const letterForm = document.getElementById("letter-form");
  const previewButton = document.getElementById("preview-draft");
  const previewPanel = document.getElementById("draft-preview");
  const previewError = document.getElementById("draft-error");
  const paper = document.getElementById("draft-paper");
  const backToForm = document.getElementById("back-to-form");
  const downloadDraft = document.getElementById("download-draft");
  if (letterForm && previewButton && previewPanel && previewError && paper) {
    const addRow = (table, label, value) => {
      const row = document.createElement("tr");
      const heading = document.createElement("th");
      const cell = document.createElement("td");
      heading.textContent = label;
      cell.textContent = value || "—";
      row.append(heading, cell);
      table.append(row);
    };

    const renderDraft = (draft) => {
      paper.replaceChildren();
      const heading = document.createElement("header");
      heading.className = "draft-document-header";
      const logo = document.createElement("img");
      logo.className = "draft-emblem";
      logo.src = "/static/lambang-kota-palu.png";
      logo.alt = "Lambang Kota Palu";
      const letterhead = document.createElement("div");
      letterhead.className = "draft-letterhead-text";
      const government = document.createElement("strong");
      government.className = "draft-government-name";
      government.textContent = "PEMERINTAH KOTA PALU";
      const district = document.createElement("strong");
      district.className = "draft-district-name";
      district.textContent = `KECAMATAN ${draft.kecamatan}`;
      const neighborhood = document.createElement("strong");
      neighborhood.className = "draft-office-name";
      neighborhood.textContent = `KELURAHAN ${draft.kelurahan}`;
      const address = document.createElement("span");
      address.className = "draft-address";
      address.textContent = `Alamat: ${draft.alamat_kantor}${draft.email_kantor ? ` | Email: ${draft.email_kantor}` : ""}`;
      letterhead.append(government, district, neighborhood, address);
      heading.append(logo, letterhead);

      const rule = document.createElement("div");
      rule.className = "draft-divider";
      rule.append(document.createElement("span"));
      const title = document.createElement("h1");
      title.textContent = draft.jenis_surat;
      const number = document.createElement("p");
      number.className = "draft-number";
      number.textContent = "Nomor: ____________________";
      const intro = document.createElement("p");
      intro.className = "draft-intro";
      intro.textContent = "Yang bertanda tangan di bawah ini menerangkan bahwa:";
      const table = document.createElement("table");
      table.className = "draft-data-table";
      [
        ["Nama lengkap", draft.nama_lengkap],
        ["NIK", draft.nik],
        ["Tempat / tanggal lahir", draft.tempat_tgl_lahir],
        ["Jenis kelamin", draft.jenis_kelamin],
        ["Pekerjaan", draft.pekerjaan],
        ["Alamat lengkap", draft.alamat_lengkap],
        ["Tujuan / keperluan", draft.keperluan],
      ].forEach(([label, value]) => addRow(table, label, value));
      (draft.data_tambahan || []).forEach((field) => addRow(table, field.label, field.value));

      const closing = document.createElement("p");
      closing.className = "draft-closing";
      closing.textContent = "Demikian surat keterangan ini dibuat sebagai draf untuk diperiksa oleh petugas kelurahan.";
      const signature = document.createElement("section");
      signature.className = "draft-signature";
      const date = document.createElement("p");
      date.textContent = `Palu, ${draft.tanggal_draf}`;
      const official = document.createElement("p");
      official.textContent = "Lurah";
      const officialName = document.createElement("strong");
      officialName.textContent = draft.nama_lurah;
      const nip = document.createElement("p");
      nip.textContent = `NIP: ${draft.nip_lurah}`;
      signature.append(date, official, document.createElement("br"), officialName, nip);
      paper.append(heading, rule, title, number, intro, table, closing, signature);
    };

    const nameRegex = /^[A-Za-z]+(?:(?:\.?, ?|\. ?|[ '\-])[A-Za-z]+)*\.?$/;
    const getFieldLabel = (field) => {
      if (field.dataset.fieldLabel) return field.dataset.fieldLabel;
      const label = field.closest("label");
      if (!label) return "ini";
      const clone = label.cloneNode(true);
      clone.querySelectorAll("input, select, textarea, .label-help, .field-error").forEach((el) => el.remove());
      return clone.textContent.trim() || "ini";
    };

    const getLetterFieldError = (field) => {
      if (field.disabled || !field.required) return "";
      const val = (field.value || "").trim();
      const label = getFieldLabel(field);
      if (!val) {
        return field.tagName === "SELECT"
          ? `Kolom ${label} masih kosong, silakan pilih terlebih dahulu.`
          : `Kolom ${label} masih kosong, silakan isi terlebih dahulu.`;
      }
      if (field.name === "applicant_nama") {
        if (/\d/.test(val)) return "Nama lengkap tidak boleh mengandung angka.";
        if (!nameRegex.test(val)) {
          return "Nama lengkap hanya boleh berisi huruf, spasi, serta tanda titik (.), koma (,) untuk gelar, apostrof ('), atau tanda hubung (-) di antara huruf.";
        }
      }
      if (field.name === "applicant_nik" && !/^\d{16}$/.test(val)) {
        return "NIK harus terdiri dari 16 digit angka.";
      }
      return "";
    };

    const clearLetterFieldState = (field) => {
      field.setCustomValidity("");
      if (field.classList.contains("input-invalid") && !getLetterFieldError(field)) {
        field.classList.remove("input-invalid");
        field.removeAttribute("aria-invalid");
      }
    };

    letterForm.addEventListener("input", (event) => {
      if (event.target instanceof HTMLElement) clearLetterFieldState(event.target);
    });
    letterForm.addEventListener("change", (event) => {
      if (event.target instanceof HTMLElement) clearLetterFieldState(event.target);
    });

    const validateLetterForm = () => {
      const fields = Array.from(letterForm.querySelectorAll("input, select, textarea"));
      let firstInvalid = null;
      fields.forEach((field) => {
        const message = getLetterFieldError(field);
        field.setCustomValidity(message);
        if (message) {
          field.classList.add("input-invalid");
          field.setAttribute("aria-invalid", "true");
          if (!firstInvalid) firstInvalid = field;
        } else {
          field.classList.remove("input-invalid");
          field.removeAttribute("aria-invalid");
        }
      });
      if (firstInvalid) {
        firstInvalid.scrollIntoView({ behavior: "smooth", block: "center" });
        firstInvalid.focus({ preventScroll: true });
        firstInvalid.reportValidity();
        return false;
      }
      return true;
    };

    letterForm.addEventListener("submit", (event) => {
      if (!validateLetterForm()) event.preventDefault();
    });

    previewButton.addEventListener("click", async () => {
      previewError.hidden = true;
      if (!validateLetterForm()) return;

      const data = new FormData(letterForm);
      const additional = {
        kecamatan: data.get("kecamatan"),
        kelurahan: data.get("kelurahan"),
        keperluan: data.get("keperluan"),
        is_self: selfToggle?.checked ? "true" : "false",
      };
      if (selfToggle?.checked) {
        let occupation = String(data.get("applicant_pekerjaan") || "").trim();
        if (!occupation) {
          occupation = (window.prompt("Pekerjaan belum ada di profil. Masukkan pekerjaan untuk melanjutkan:") || "").trim();
          if (!occupation) return;
          const occupationInput = letterForm.querySelector('[name="applicant_pekerjaan"]');
          if (occupationInput) occupationInput.value = occupation;
        }
        additional.pekerjaan = occupation;
      } else {
        [
          "nama",
          "nik",
          "tempat_lahir",
          "tanggal_lahir",
          "jenis_kelamin",
          "pekerjaan",
          "alamat",
        ].forEach((key) => {
          additional[`pemohon_${key}`] = String(data.get(`applicant_${key}`) || "").trim();
        });
        if (!additional.pemohon_pekerjaan) {
          additional.pemohon_pekerjaan = (window.prompt("Pekerjaan pemohon belum diisi. Masukkan pekerjaan untuk melanjutkan:") || "").trim();
          if (!additional.pemohon_pekerjaan) return;
          const occupationInput = letterForm.querySelector('[name="applicant_pekerjaan"]');
          if (occupationInput) occupationInput.value = additional.pemohon_pekerjaan;
        }
      }
      letterForm.querySelectorAll('[name^="specific_"]').forEach((input) => {
        additional[input.name] = input.value.trim();
      });

      previewButton.disabled = true;
      previewButton.textContent = "Menyiapkan pratinjau…";
      try {
        const response = await fetch("/api/v1/letters/draft-preview", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            user_id: letterForm.dataset.userId,
            jenis_surat_id: serviceSelect?.value,
            data_tambahan: additional,
          }),
        });
        const result = await response.json();
        if (!response.ok || result.status !== "success" || !result.draft_data) {
          throw new Error(result.message || "Pratinjau draf tidak dapat dibuat.");
        }
        renderDraft(result.draft_data);
        letterForm.hidden = true;
        document.querySelector(".page-heading")?.setAttribute("hidden", "");
        previewPanel.hidden = false;
        document.body.classList.add("draft-mode");
      } catch (error) {
        previewError.textContent = error instanceof Error ? error.message : "Pratinjau draf tidak dapat dibuat.";
        previewError.hidden = false;
      } finally {
        previewButton.disabled = false;
        previewButton.textContent = "Pratinjau draf";
      }
    });

    backToForm?.addEventListener("click", () => {
      previewPanel.hidden = true;
      letterForm.hidden = false;
      document.querySelector(".page-heading")?.removeAttribute("hidden");
      document.body.classList.remove("draft-mode");
    });
    downloadDraft?.addEventListener("click", () => window.print());
  }
});
