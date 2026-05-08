const dropzone = document.getElementById("dropzone");
const fileInput = document.getElementById("fileInput");
const fileLabel = document.getElementById("fileLabel");
const uploadBtn = document.getElementById("uploadBtn");

const progressContainer = document.getElementById("progressContainer");
const progressBar = document.getElementById("progressBar");
const progressText = document.getElementById("progressText");

const spinner = document.getElementById("spinner");
const result = document.getElementById("result");

const downloadBtn = document.getElementById("downloadBtn");
const codeInput = document.getElementById("codeInput");

let selectedFile = null;

// SELECT FILE

dropzone.addEventListener("click", () => fileInput.click());

fileInput.addEventListener("change", (e) => {
  selectedFile = e.target.files[0];
  if (selectedFile) {
    fileLabel.innerText = selectedFile.name;
    fileLabel.classList.add("text-green-600");
  }
});

// DRAG & DROP

dropzone.addEventListener("dragover", (e) => {
  e.preventDefault();
  dropzone.classList.add("border-teal-500");
});

dropzone.addEventListener("dragleave", () => {
  dropzone.classList.remove("border-teal-500");
});

dropzone.addEventListener("drop", (e) => {
  e.preventDefault();
  dropzone.classList.remove("border-teal-500");

  selectedFile = e.dataTransfer.files[0];
  if (selectedFile) {
    fileLabel.innerText = selectedFile.name;
    fileLabel.classList.add("text-green-600");
  }
});

// UPLOAD

uploadBtn.addEventListener("click", () => {
  if (!selectedFile) {
    alert("Please select a file");
    return;
  }

  // limit 10MB (frontend)
  if (selectedFile.size > 10 * 1024 * 1024) {
    alert("File too large (max 10MB)");
    return;
  }

  const formData = new FormData();
  formData.append("file", selectedFile);

  const xhr = new XMLHttpRequest();
  xhr.open("POST", "/upload");

  // UI state
  progressContainer.classList.remove("hidden");
  spinner.classList.remove("hidden");
  uploadBtn.disabled = true;
  result.innerHTML = "";

  // PROGRESS

  xhr.upload.onprogress = (e) => {
    if (e.lengthComputable) {
      const percent = Math.round((e.loaded / e.total) * 100);
      progressBar.style.width = percent + "%";
      progressText.innerText = percent + "%";
    }
  };

  // DONE

  xhr.onload = () => {
    spinner.classList.add("hidden");
    uploadBtn.disabled = false;

    if (xhr.status === 200) {
      result.innerHTML = xhr.responseText;

      // extract 6-digit code
      const match = xhr.responseText.match(/\d{6}/);
      if (match) {
        codeInput.value = match[0];
      }

      // reset progress
      progressBar.style.width = "0%";
      progressText.innerText = "0%";

    } else {
      result.innerText = xhr.responseText || `Upload failed (HTTP ${xhr.status})`;
    }
  };

  xhr.onerror = () => {
    spinner.classList.add("hidden");
    uploadBtn.disabled = false;
    result.innerText = "Upload failed: cannot reach server";
  };

  xhr.send(formData);
});

// DOWNLOAD

if (downloadBtn) {
  downloadBtn.addEventListener("click", () => {
    const code = codeInput.value;

    if (!code) {
      alert("Please enter code");
      return;
    }

    window.location.href = `/download?code=${code}`;
  });
}
