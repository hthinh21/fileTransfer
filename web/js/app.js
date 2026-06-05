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
      try {
        const data = JSON.parse(xhr.responseText);

        result.innerHTML = `Your code: <span class="font-bold text-2xl tracking-widest text-teal-600">${data.code}</span>`;
        document.getElementById("countdownContainer").classList.remove("hidden");
        startCountdown(600);

      } catch (e) {
        result.innerHTML = xhr.responseText;
      }

      // reset & hide progress
      progressBar.style.width = "0%";
      progressText.innerText = "0%";
      progressContainer.classList.add("hidden");

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

function startCountdown(totalSeconds) {
  const countdown = document.getElementById("countdown");
  const end = Date.now() + totalSeconds * 1000;

  const interval = setInterval(() => {
    const remaining = Math.max(0, Math.round((end - Date.now()) / 1000));

    const m = Math.floor(remaining / 60).toString().padStart(2, "0");
    const s = (remaining % 60).toString().padStart(2, "0");
    countdown.innerText = `${m}:${s}`;

    if (remaining <= 0) {
      clearInterval(interval);
      countdown.innerText = "Expired";
      countdown.classList.remove("text-red-500");
      countdown.classList.add("text-gray-400");
      result.innerText = "Code has expired.";
    }
  }, 1000);
}

