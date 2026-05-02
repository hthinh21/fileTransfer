const dropzone = document.getElementById("dropzone");
const fileInput = document.getElementById("fileInput");
const uploadBtn = document.getElementById("uploadBtn");

const progressContainer = document.getElementById("progressContainer");
const progressBar = document.getElementById("progressBar");
const progressText = document.getElementById("progressText");

const spinner = document.getElementById("spinner");
const result = document.getElementById("result");

let selectedFile = null;

// click chọn file
dropzone.addEventListener("click", () => fileInput.click());

fileInput.addEventListener("change", (e) => {
  selectedFile = e.target.files[0];
  dropzone.innerHTML = `<p class="text-green-600">${selectedFile.name}</p>`;
});

// drag & drop
dropzone.addEventListener("dragover", (e) => {
  e.preventDefault();
  dropzone.classList.add("border-teal-500");
});

dropzone.addEventListener("dragleave", () => {
  dropzone.classList.remove("border-teal-500");
});

dropzone.addEventListener("drop", (e) => {
  e.preventDefault();
  selectedFile = e.dataTransfer.files[0];
  dropzone.innerHTML = `<p class="text-green-600">${selectedFile.name}</p>`;
});

// upload
uploadBtn.addEventListener("click", () => {
  if (!selectedFile) {
    alert("Please select a file");
    return;
  }

  const formData = new FormData();
  formData.append("file", selectedFile);

  const xhr = new XMLHttpRequest();
  xhr.open("POST", "/upload");

  progressContainer.classList.remove("hidden");
  spinner.classList.remove("hidden");
  result.innerHTML = "";

  xhr.upload.onprogress = (e) => {
    if (e.lengthComputable) {
      const percent = Math.round((e.loaded / e.total) * 100);
      progressBar.style.width = percent + "%";
      progressText.innerText = percent + "%";
    }
  };

  xhr.onload = () => {
    spinner.classList.add("hidden");

    if (xhr.status === 200) {
      result.innerHTML = xhr.responseText;
    } else {
      result.innerHTML = "Upload failed";
    }
  };

  xhr.send(formData);
});