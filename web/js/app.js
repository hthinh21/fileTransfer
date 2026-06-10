const dropzone = document.getElementById("dropzone");
const fileInput = document.getElementById("fileInput");
const fileLabel = document.getElementById("fileLabel");
const uploadBtn = document.getElementById("uploadBtn");

const progressContainer = document.getElementById("progressContainer");
const progressBar = document.getElementById("progressBar");
const progressText = document.getElementById("progressText");

const spinner = document.getElementById("spinner");
const result = document.getElementById("result");

const uploadSection = document.getElementById("uploadSection");
const resultSection = document.getElementById("resultSection");

const downloadBtn = document.getElementById("downloadBtn");
const codeInput = document.getElementById("codeInput");

let selectedFile = null;
let currentCode = null;
let countdownInterval = null;
let notifySocket = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;

function connectNotifySocket(code) {
  if (notifySocket) {
    notifySocket.onclose = null; // tránh vòng lặp đóng-mở chồng nhau khi tự thay thế kết nối cũ
    notifySocket.close();
  }

  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  notifySocket = new WebSocket(`${protocol}://${window.location.host}/ws?code=${code}`);

  notifySocket.onopen = () => {
    reconnectAttempts = 0;
  };

  notifySocket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.event === "downloaded" && code === currentCode) {
      if (countdownInterval) clearInterval(countdownInterval);
      document.getElementById("countdownContainer").classList.add("hidden");
      result.innerHTML = `File đã được tải về.`;
      currentCode = null;
    }
  };

  notifySocket.onclose = () => {
    // chỉ thử kết nối lại nếu code này vẫn đang còn hiệu lực và chưa vượt quá số lần cho phép
    if (code !== currentCode || reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) return;

    reconnectAttempts++;
    const delay = Math.min(1000 * 2 ** reconnectAttempts, 10000);
    setTimeout(() => connectNotifySocket(code), delay);
  };
}

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

const uploadError = document.getElementById("uploadError");

uploadBtn.addEventListener("click", () => {
  uploadError.classList.add("hidden");
  uploadError.innerText = "";


  if (!selectedFile) {
    uploadError.innerText = "Vui lòng chọn file để upload";
    uploadError.classList.remove("hidden");
    return;
  }

  // limit 10MB (frontend)
  if (selectedFile.size > 10 * 1024 * 1024) {
    uploadError.innerText = "File quá lớn (tối đa 10MB)";
    uploadError.classList.remove("hidden");
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
    uploadError.classList.add("hidden");
    spinner.classList.add("hidden");
    uploadBtn.disabled = false;

    if (xhr.status === 200) {
      try {
        const data = JSON.parse(xhr.responseText);

        currentCode = data.code;
        connectNotifySocket(data.code);

        result.innerHTML = `Code: <span class="font-bold text-2xl tracking-widest text-teal-600">${data.code}</span>`;
        document.getElementById("countdownContainer").classList.remove("hidden");
        countdownInterval = startCountdown(600);

        // ẩn upload, hiện result
        uploadSection.classList.add("hidden");
        resultSection.classList.remove("hidden");

      } catch (e) {
        result.innerHTML = xhr.responseText;
      }

      // reset progress
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
  const downloadError = document.getElementById("downloadError");

  downloadBtn.addEventListener("click", async () => {
    const code = codeInput.value.trim();

    if (!code) {
      alert("Please enter code");
      return;
    }

    downloadError.classList.add("hidden");
    downloadError.innerText = "";

    try {
      const res = await fetch(`/download?code=${code}`);

      if (!res.ok) {
        const message = await res.text();
        downloadError.innerText = message || "Mã không hợp lệ hoặc file đã hết hạn.";
        downloadError.classList.remove("hidden");
        return;
      }

      const disposition = res.headers.get("Content-Disposition");
      let filename = "download";
      const match = disposition && disposition.match(/filename="?([^"]+)"?/);
      if (match) filename = match[1];

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      downloadError.innerText = "Cannot reach server";
      downloadError.classList.remove("hidden");
    }
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

  return interval;
}

