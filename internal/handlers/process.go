package handlers

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/dendianugerah/reubah/internal/processor"
	"github.com/dendianugerah/reubah/internal/processor/resize"
	"github.com/dendianugerah/reubah/internal/validator"
	"github.com/dendianugerah/reubah/pkg/errors"
)

func ProcessImage(w http.ResponseWriter, r *http.Request) {
	// Add debug logging
	log.Printf("Starting image processing request")

	// Log the content type of the request
	log.Printf("Request Content-Type: %s", r.Header.Get("Content-Type"))

	// Parse multipart form with increased size limit
	err := r.ParseMultipartForm(32 << 20)
	if err != nil { // 32MB limit
		log.Printf("Error parsing multipart form: %v", err)
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Unable to parse form", err))
		return
	}

	// Log form values
	log.Printf("Form values: %+v", r.Form)
	log.Printf("File headers: %+v", r.MultipartForm.File)

	// Get the uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "No file uploaded", err))
		return
	}
	defer file.Close()

	// Log file details
	log.Printf("Successfully received file: %s, size: %d, content-type: %s",
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"))

	// Validate MIME type
	if err := validator.ValidateMIMEType(file); err != nil {
		log.Printf("MIME type validation failed: %v", err)
		errors.SendError(w, errors.New(errors.ErrInvalidMIME, "Invalid file type", err))
		return
	}

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Invalid image file", err))
		return
	}

	// Get processing options from form
	format := r.FormValue("format")
	if format == "" {
		format = "jpeg" // Default format
	}

	// Parse dimensions
	width := 0
	height := 0
	if widthStr := r.FormValue("width"); widthStr != "" {
		width, err = strconv.Atoi(widthStr)
		if err != nil {
			errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Invalid width value", err))
			return
		}
	}

	if heightStr := r.FormValue("height"); heightStr != "" {
		height, err = strconv.Atoi(heightStr)
		if err != nil {
			errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Invalid height value", err))
			return
		}
	}

	// Get resize mode
	resizeMode := r.FormValue("resizeMode")
	if resizeMode == "" {
		resizeMode = "fit"
	}

	// Parse resize mode
	parsedResizeMode, err := resize.ParseResizeMode(resizeMode)
	if err != nil {
		errors.SendError(w, errors.New(errors.ErrInvalidFormat, "Invalid resize mode", err))
		return
	}

	// Parse quality
	qualityValue := 85 // default quality
	if quality := r.FormValue("quality"); quality != "" {
		switch quality {
		case "low":
			qualityValue = 60
		case "medium":
			qualityValue = 75
		case "high":
			qualityValue = 90
		case "lossless":
			qualityValue = 100
		}
	}

	// Create processing options
	opts := processor.ProcessOptions{
		Width:            width,
		Height:           height,
		ResizeMode:       parsedResizeMode,
		OutputFormat:     format,
		Quality:          qualityValue,
		Optimize:         r.FormValue("optimize") == "true",
		RemoveBackground: r.FormValue("removeBackground") == "true",
	}

	// Process the image
	proc := processor.NewImageProcessor()
	processedImage, err := proc.ProcessImageData(img, opts)
	if err != nil {
		log.Printf("Processing error: %v", err)
		errors.SendError(w, errors.New(errors.ErrProcessingFailed, "Failed to process image", err))
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", format))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=processed.%s", format))

	// Write the processed image directly to the response
	if err := processedImage.Write(w); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}
