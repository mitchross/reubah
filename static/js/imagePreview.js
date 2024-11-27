document.addEventListener("DOMContentLoaded", function () {
  const elements = {
    imageInput: document.getElementById("imageInput"),
    previewDiv: document.getElementById("preview"),
    uploadArea: document.getElementById("uploadArea"),
    fileStatus: document.getElementById("fileStatus"),
    fileName: document.getElementById("fileName"),
    fileSize: document.getElementById("fileSize"),
    uploadText: document.getElementById("uploadText"),
    widthInput: document.getElementById("widthInput"),
    heightInput: document.getElementById("heightInput")
  };

  function initializeImagePreview() {
    if (!elements.imageInput || !elements.uploadArea) {
      console.error("Required preview elements not found");
      return;
    }

    setupEventListeners();
  }

  function setupEventListeners() {
    elements.imageInput.addEventListener("change", (e) => {
      const file = e.target.files[0];
      if (file) handleFileSelect(file);
    });

    setupDragAndDrop();
  }

  function setupDragAndDrop() {
    elements.uploadArea.addEventListener("dragover", (e) => {
      e.preventDefault();
      elements.uploadArea.classList.add("border-indigo-500");
    });

    elements.uploadArea.addEventListener("dragleave", (e) => {
      e.preventDefault();
      elements.uploadArea.classList.remove("border-indigo-500");
    });

    elements.uploadArea.addEventListener("drop", (e) => {
      e.preventDefault();
      elements.uploadArea.classList.remove("border-indigo-500");

      if (e.dataTransfer.files.length) {
        elements.imageInput.files = e.dataTransfer.files;
        handleFileSelect(e.dataTransfer.files[0]);
      }
    });
  }

  function handleFileSelect(file) {
    if (!file.type.startsWith("image/")) {
      alert("Please select an image file");
      return;
    }

    updateFileStatus(file);
    previewImage(file);
  }

  function updateFileStatus(file) {
    if (!elements.fileStatus) return;

    elements.fileStatus.classList.remove("hidden");
    elements.fileStatus.classList.add("bg-green-50");
    elements.fileName.textContent = file.name;
    elements.fileSize.textContent = ` (${(file.size / (1024 * 1024)).toFixed(2)} MB)`;
  }

  function previewImage(file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (!elements.previewDiv) return;

      elements.previewDiv.classList.remove("hidden");
      const img = elements.previewDiv.querySelector("img");
      if (!img) return;

      img.src = e.target.result;
      img.onload = () => updateImageInfo(img);
    };
    reader.readAsDataURL(file);
  }

  function updateImageInfo(img) {
    const dimensions = `${img.naturalWidth} × ${img.naturalHeight}px`;
    const dimensionsElement = elements.previewDiv.querySelector(".image-dimensions");
    if (dimensionsElement) {
      dimensionsElement.textContent = dimensions;
    }

    elements.uploadArea.classList.add("border-green-500");
    if (elements.uploadText) {
      elements.uploadText.innerHTML = '<span class="text-green-500">File ready for processing</span>';
    }

    // Update dimension inputs with placeholders
    if (elements.widthInput) elements.widthInput.placeholder = img.naturalWidth;
    if (elements.heightInput) elements.heightInput.placeholder = img.naturalHeight;
  }

  initializeImagePreview();
});
