document.addEventListener("DOMContentLoaded", function () {
  // First verify all required elements exist
  const form = document.getElementById("uploadForm");
  if (!form) {
    console.error("Upload form not found");
    return;
  }

  const progress = document.getElementById("progress");
  const result = document.getElementById("result");
  const formatSelect = document.getElementById("formatSelect");
  const qualitySelect = document.getElementById("qualitySelect");
  const widthInput = document.getElementById("widthInput");
  const heightInput = document.getElementById("heightInput");
  const fileInput = form.querySelector('input[type="file"]');

  if (!fileInput) {
    console.error("File input not found");
    return;
  }

  const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB

  // Quick action handlers
  document.querySelectorAll("[data-action]").forEach((button) => {
    button.addEventListener("click", function (e) {
      e.preventDefault();
      const action = this.dataset.action;
      handleQuickAction(action);

      // Visual feedback
      document.querySelectorAll("[data-action]").forEach((btn) => {
        btn.classList.remove(
          "bg-indigo-50",
          "border-indigo-500",
          "text-indigo-700"
        );
      });
      this.classList.add(
        "bg-indigo-50",
        "border-indigo-500",
        "text-indigo-700"
      );
    });
  });

  // Add background removal handling
  const removeBackgroundCheckbox = document.getElementById("removeBackground");
  const bgRemovalOptions = document.getElementById("bgRemovalOptions");

  if (removeBackgroundCheckbox) {
    removeBackgroundCheckbox.addEventListener("change", function() {
      if (bgRemovalOptions) {
        bgRemovalOptions.classList.toggle("hidden", !this.checked);
      }
    });
  }

  form.addEventListener("submit", async function (e) {
    e.preventDefault();

    // Check if file is selected
    if (!fileInput.files || !fileInput.files.length) {
      showError("Please select an image file first");
      return;
    }

    const file = fileInput.files[0];
    if (file.size > MAX_FILE_SIZE) {
      showError("File size exceeds 10MB limit");
      return;
    }

    // Create FormData and explicitly append the file
    const formData = new FormData();
    formData.append("image", file); // Make sure we use "image" as the field name

    // Add other form fields
    const format = formatSelect ? formatSelect.value : "jpeg";
    formData.append("format", format);

    if (qualitySelect) {
      formData.append("quality", qualitySelect.value);
    }
    if (widthInput && widthInput.value) {
      formData.append("width", widthInput.value);
    }
    if (heightInput && heightInput.value) {
      formData.append("height", heightInput.value);
    }

    // Add background removal options
    if (removeBackgroundCheckbox && removeBackgroundCheckbox.checked) {
      formData.append("removeBackground", "true");
    }

    // Debug log
    console.log("File being sent:", file);
    console.log("Format being sent:", format);
    
    try {
      showProgress();
      const submitButton = this.querySelector('button[type="submit"]');
      if (submitButton) {
        submitButton.disabled = true;
      }

      const response = await fetch("/process", {
        method: "POST",
        body: formData,
        // Don't set Content-Type header - let the browser set it with the boundary
      });

      if (!response.ok) {
        let errorMessage = "Processing failed";
        try {
          const errorData = await response.json();
          errorMessage = errorData.error?.message || errorMessage;
        } catch (e) {
          console.error("Error parsing error response:", e);
        }
        throw new Error(errorMessage);
      }

      await handleSuccess(response, format);
    } catch (error) {
      console.error("Error:", error);
      showError(error.message);
    } finally {
      hideProgress();
      const submitButton = this.querySelector('button[type="submit"]');
      if (submitButton) {
        submitButton.disabled = false;
      }
    }
  });

  function showProgress() {
    if (progress) {
      progress.style.display = "block";
    }
    if (result) {
      result.style.display = "none";
    }
  }

  function hideProgress() {
    if (progress) {
      progress.style.display = "none";
    }
  }

  function showError(message) {
    const errorDiv = document.createElement("div");
    errorDiv.className = "bg-red-50 border-l-4 border-red-400 p-4 mb-4";
    errorDiv.innerHTML = `
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">${message}</p>
        </div>
      </div>
    `;
    form.insertBefore(errorDiv, form.firstChild);
    setTimeout(() => errorDiv.remove(), 5000);
  }

  async function handleSuccess(response, format) {
    if (!result) {
      console.error("Result element not found");
      return;
    }

    const blob = await response.blob();
    const url = URL.createObjectURL(blob);

    const resultPreview = result.querySelector(".result-preview");
    const resultInfo = result.querySelector(".result-info");

    if (resultPreview && resultInfo) {
      resultPreview.innerHTML = `<img src="${url}" alt="Processed image" class="max-w-full rounded-lg">`;
      resultInfo.innerHTML = `
        <div class="flex justify-between items-center">
          <p class="text-sm text-gray-600">Size: ${(blob.size / 1024).toFixed(2)} KB</p>
          <a href="${url}" 
             download="processed.${format}" 
             class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Download Image
          </a>
        </div>
      `;
    }

    result.style.display = "block";
  }

  function handleQuickAction(action) {
    switch (action) {
      case "optimize":
        if (formatSelect) formatSelect.value = "webp";
        if (qualitySelect) qualitySelect.value = "medium";
        if (widthInput) widthInput.value = "";
        if (heightInput) heightInput.value = "";
        break;

      case "resize":
        if (widthInput) widthInput.focus();
        break;

      case "adjust":
        if (formatSelect) formatSelect.value = "png";
        if (qualitySelect) qualitySelect.value = "lossless";
        break;

      case "convert":
        if (formatSelect) formatSelect.focus();
        break;
    }
  }
});
