document.addEventListener("DOMContentLoaded", function () {
  const imageInput = document.getElementById("imageInput");
  const previewDiv = document.getElementById("preview");
  const uploadArea = document.getElementById("uploadArea");
  const fileStatus = document.getElementById("fileStatus");
  const fileName = document.getElementById("fileName");
  const fileSize = document.getElementById("fileSize");
  const uploadText = document.getElementById("uploadText");

  imageInput.addEventListener("change", function (e) {
    const file = e.target.files[0];
    if (file) {
      handleFileSelect(file);
    }
  });

  function handleFileSelect(file) {
    if (!file.type.startsWith("image/")) {
      alert("Please select an image file");
      return;
    }

    // Show and update file status
    fileStatus.classList.remove("hidden");
    fileStatus.classList.add("bg-green-50");
    fileName.textContent = file.name;
    fileSize.textContent = ` (${(file.size / (1024 * 1024)).toFixed(2)} MB)`;

    const reader = new FileReader();
    reader.onload = function (e) {
      // Show preview
      previewDiv.classList.remove("hidden");
      const img = previewDiv.querySelector("img");
      img.src = e.target.result;

      // Show image dimensions when loaded
      img.onload = function () {
        const dimensions = `${img.naturalWidth} × ${img.naturalHeight}px`;
        previewDiv.querySelector(".image-dimensions").textContent = dimensions;

        // Indicate that the file is ready for processing
        uploadArea.classList.add("border-green-500"); // Change border color
        uploadText.innerHTML = `<span class="text-green-500">File ready for processing</span>`;

        const widthInput = document.getElementById("widthInput");
        const heightInput = document.getElementById("heightInput");
        if (widthInput) widthInput.placeholder = img.naturalWidth;
        if (heightInput) heightInput.placeholder = img.naturalHeight;
      };
    };
    reader.readAsDataURL(file);
  }

  uploadArea.addEventListener("dragover", function (e) {
    e.preventDefault();
    this.classList.add("border-indigo-500");
  });

  uploadArea.addEventListener("dragleave", function (e) {
    e.preventDefault();
    this.classList.remove("border-indigo-500");
  });

  uploadArea.addEventListener("drop", function (e) {
    e.preventDefault();
    this.classList.remove("border-indigo-500");

    if (e.dataTransfer.files.length) {
      imageInput.files = e.dataTransfer.files;
      handleFileSelect(e.dataTransfer.files[0]);
    }
  });
});
