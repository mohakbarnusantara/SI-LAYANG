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
      nama: { element: document.querySelector('#standard-letter-fields [name="applicant_nama"]'), profile: "profileName" },
      nik: { element: document.querySelector('#standard-letter-fields [name="applicant_nik"]'), profile: "profileNik" },
      tempat_lahir: { element: document.querySelector('#standard-letter-fields [name="applicant_tempat_lahir"]'), profile: "profileBirthPlace" },
      tanggal_lahir: { element: document.querySelector('#standard-letter-fields [name="applicant_tanggal_lahir"]'), profile: "profileBirth" },
      jenis_kelamin: { element: document.querySelector('#standard-letter-fields [name="applicant_jenis_kelamin"]'), profile: "profileGender" },
      pekerjaan: { element: document.querySelector('#standard-letter-fields [name="applicant_pekerjaan"]'), profile: "profileOccupation" },
      alamat: { element: document.querySelector('#standard-letter-fields [name="applicant_alamat"]'), profile: "profileAddress" },
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

  const selfIncomeToggle = document.getElementById("is-self-income");
  if (selfIncomeToggle) {
    const incomeFields = {
      nama: { element: document.querySelector('#income-letter-fields [name="applicant_nama"]'), profile: "profileName" },
      nik: { element: document.querySelector('#income-letter-fields [name="applicant_nik"]'), profile: "profileNik" },
      tempat_lahir: { element: document.querySelector('#income-letter-fields [name="applicant_tempat_lahir"]'), profile: "profileBirthPlace" },
      tanggal_lahir: { element: document.querySelector('#income-letter-fields [name="applicant_tanggal_lahir"]'), profile: "profileBirth" },
      jenis_kelamin: { element: document.querySelector('#income-letter-fields [name="applicant_jenis_kelamin"]'), profile: "profileGender" },
      pekerjaan: { element: document.querySelector('#income-letter-fields [name="applicant_pekerjaan"]'), profile: "profileOccupation" },
      alamat: { element: document.querySelector('#income-letter-fields [name="applicant_alamat"]'), profile: "profileAddress" },
    };
    const fillIncomeFromProfile = (enabled) => {
      for (const field of Object.values(incomeFields)) {
        const input = field.element;
        if (!input) continue;
        if (enabled) {
          if (field.profile === "profileGender") {
            const g = input.dataset.profileGender;
            input.value = g === "P" ? "Perempuan" : (g === "L" ? "Laki-laki" : "");
          } else {
            input.value = input.dataset[field.profile] || "";
          }
        } else {
          input.value = "";
        }
        input.classList.remove("input-invalid");
        input.removeAttribute("aria-invalid");
        input.setCustomValidity("");
      }
    };
    selfIncomeToggle.addEventListener("change", () => fillIncomeFromProfile(selfIncomeToggle.checked));
  }

  const serviceSelect = document.getElementById("service-select");
  const fieldContainer = document.getElementById("specific-fields");
  const hint = document.getElementById("service-hint");
  const extraLabel = document.getElementById("extra-doc-label");
  const standardFields = document.getElementById("standard-letter-fields");
  const incomeFields = document.getElementById("income-letter-fields");

  if (serviceSelect) {
    const setContainerFieldsDisabled = (container, disabled) => {
      if (!container) return;
      container.querySelectorAll("input, select, textarea").forEach((el) => {
        el.disabled = disabled;
        if (disabled) {
          el.classList.remove("input-invalid");
          el.removeAttribute("aria-invalid");
          el.setCustomValidity("");
        }
      });
    };

    const updateServiceView = () => {
      const isIncome = serviceSelect.value === "SK-PENGHASILAN";
      if (standardFields && incomeFields) {
        if (isIncome) {
          standardFields.hidden = true;
          setContainerFieldsDisabled(standardFields, true);
          incomeFields.hidden = false;
          setContainerFieldsDisabled(incomeFields, false);
          if (fieldContainer) fieldContainer.replaceChildren();
          if (hint) hint.textContent = "Lengkapi Data Orang Tua, Data Anak, serta Data Penghasilan & Dokumen Pendukung.";
          extraLabel?.classList.remove("required-doc");
          return;
        } else {
          standardFields.hidden = false;
          setContainerFieldsDisabled(standardFields, false);
          incomeFields.hidden = true;
          setContainerFieldsDisabled(incomeFields, true);
        }
      }

      const option = serviceSelect.selectedOptions[0];
      if (fieldContainer) fieldContainer.replaceChildren();
      if (!option || !option.value) {
        if (hint) hint.textContent = "Setiap layanan dapat meminta data khusus dan berkas pendukung.";
        extraLabel?.classList.remove("required-doc");
        return;
      }
      const labels = (option.dataset.fields || "").split("|").filter(Boolean);
      if (hint) {
        hint.textContent = labels.length
          ? `Data tambahan untuk layanan ini: ${labels.join(", ")}.`
          : "Tidak ada data khusus tambahan untuk layanan ini.";
      }
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
        fieldContainer?.append(wrap);
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
    };

    serviceSelect.addEventListener("change", updateServiceView);
    if (serviceSelect.value) updateServiceView();
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

  document.addEventListener("input", (event) => {
    const field = event.target;
    if (!(field instanceof HTMLInputElement || field instanceof HTMLTextAreaElement)) return;
    const filterName = field.dataset?.filter;
    if (!filterName) return;
    const sanitize = inputFilters[filterName];
    if (!sanitize) return;
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

  const registerForm = document.getElementById("register-form");
  const registerError = document.getElementById("register-error");
  const reverifyForm = document.getElementById("reverify-form");
  const reverifyError = document.getElementById("reverify-error");
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

  const attachIdentityFormValidation = (formEl, errorEl, headingText) => {
    if (!formEl || !errorEl) return;
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

    const validateIdentityField = (field) => {
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
          if (!file) return field.required ? "Foto KTP wajib diunggah." : "";
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

    const fields = Array.from(formEl.querySelectorAll("input, select, textarea"));
    fields.forEach((field) => {
      const eventName = field.tagName === "SELECT" || field.type === "file" || field.type === "date" ? "change" : "input";
      field.addEventListener(eventName, () => {
        if (field.classList.contains("input-invalid")) {
          setFieldError(field, validateIdentityField(field));
        }
      });
    });

    formEl.addEventListener("submit", (event) => {
      const errors = [];
      let firstInvalid = null;
      fields.forEach((field) => {
        const message = validateIdentityField(field);
        setFieldError(field, message);
        if (message) {
          errors.push(message);
          if (!firstInvalid) firstInvalid = field;
        }
      });

      if (errors.length > 0) {
        event.preventDefault();
        errorEl.replaceChildren();
        const title = document.createElement("strong");
        title.textContent = headingText;
        const list = document.createElement("ul");
        list.className = "alert-list";
        errors.forEach((msg) => {
          const item = document.createElement("li");
          item.textContent = msg;
          list.append(item);
        });
        errorEl.append(title, list);
        errorEl.hidden = false;
        errorEl.scrollIntoView({ behavior: "smooth", block: "nearest" });
        firstInvalid?.focus();
      } else {
        errorEl.hidden = true;
      }
    });
  };

  attachIdentityFormValidation(registerForm, registerError, "Periksa kembali data pendaftaran Anda:");
  attachIdentityFormValidation(reverifyForm, reverifyError, "Periksa kembali data perbaikan Anda:");

  const letterForm = document.getElementById("letter-form");
  if (letterForm) {
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
      if (field.name === "applicant_nama" || field.name === "specific_nama_anak") {
        if (/\d/.test(val)) return `Kolom ${label} tidak boleh mengandung angka.`;
        if (!nameRegex.test(val)) {
          return `Kolom ${label} hanya boleh berisi huruf, spasi, serta tanda titik (.), koma (,) untuk gelar, apostrof ('), atau tanda hubung (-) di antara huruf.`;
        }
      }
      if (field.name === "applicant_nik" && field.closest("#standard-letter-fields") && !/^\d{16}$/.test(val)) {
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
  }
});
