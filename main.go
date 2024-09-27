package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type OperationBody struct {
	Operation  string                 `json:"operation"`
	Parameters map[string]interface{} `json:"parameters"`
}

type ImagePayload struct {
	fileName string
	payload  []byte
	err      error
}

func main() {

	log.Println("Starting image processor service...")

	http.HandleFunc("/upload", uploadHandler)

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)

	log.Println("Server started at port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Could not listen and serve", err)
		return
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if err := validateRequest(r); err != nil {
		writeHttpError(w, err, http.StatusMethodNotAllowed)
		return
	}

	if err := parseMultipartFormFromRequest(r); err != nil {
		writeHttpError(w, err, http.StatusInternalServerError)
		return
	}

	opBody, err := parseOperationBodyFromRequest(r)
	if err != nil {
		writeHttpError(w, err, http.StatusBadRequest)
		return
	}

	opProcessor, err := getProcessorForOperation(opBody.Operation)
	if err != nil {
		writeHttpError(w, err, http.StatusBadRequest)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]

	if err := validateFileBatchLength(len(fileHeaders)); err != nil {
		writeHttpError(w, err, http.StatusBadRequest)
		return
	}

	var imgsToBeProcessed []ImagePayload

	for _, fileHeader := range fileHeaders {

		imgBytes, err := readImageBytesFromFileHeader(fileHeader)
		if err != nil {
			writeHttpError(w, err, http.StatusBadRequest)
			return
		}

		imgPayload := ImagePayload{
			fileName: fileHeader.Filename,
			payload:  imgBytes,
		}

		imgsToBeProcessed = append(imgsToBeProcessed, imgPayload)

	}

	var wg sync.WaitGroup

	processedImgs := make(chan ImagePayload, len(imgsToBeProcessed))

	for _, originalImg := range imgsToBeProcessed {

		wg.Add(1)

		go processImg(opProcessor, originalImg, opBody, &wg, processedImgs)

	}

	wg.Wait()

	close(processedImgs)

	responseBuffer, multiPartWriter, err := buildMultiPartWriter()

	if err != nil {
		writeHttpError(w, err, http.StatusInternalServerError)
		return
	}

	for processedImg := range processedImgs {

		if err := buildResponse(multiPartWriter, processedImg); err != nil {
			writeHttpError(w, err, http.StatusInternalServerError)
			return
		}

	}

	if err := closeMultiPartWriterAndWriteResponse(w, multiPartWriter, responseBuffer); err != nil {
		writeHttpError(w, err, http.StatusInternalServerError)
		return
	}

	return

}

func writeHttpError(w http.ResponseWriter, err error, httpStatus int) {

	errorResponse := ErrorResponse{Error: err.Error()}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(httpStatus)

	_ = json.NewEncoder(w).Encode(errorResponse)

	http.Error(w, fmt.Sprintf("Error: %v", err), httpStatus)
}

func validateRequest(r *http.Request) error {

	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed")
	}

	return nil
}

func parseMultipartFormFromRequest(r *http.Request) error {

	if err := r.ParseMultipartForm(MaxUploadSize * MaxBatchSize); err != nil {
		return fmt.Errorf("could not parse multipart form")
	}

	return nil
}

func parseOperationBodyFromRequest(r *http.Request) (OperationBody, error) {

	operationJSON := r.FormValue("operation")

	var operationBody OperationBody

	if err := json.Unmarshal([]byte(operationJSON), &operationBody); err != nil {
		return OperationBody{}, err
	}

	return operationBody, nil
}

func getProcessorForOperation(operation string) (func(file []byte, parameters map[string]interface{}) ([]byte, error), error) {

	op := SupportedOperations[operation]

	if op == nil {
		return op, fmt.Errorf("invalid operation")
	}

	return op, nil

}

func validateFileBatchLength(fhLen int) error {

	if fhLen < 1 {
		return fmt.Errorf("at least one image must be provided")
	}

	if fhLen > 3 {
		return fmt.Errorf("cannot process %d files", fhLen)
	}

	return nil

}

func readImageBytesFromFileHeader(fileHeader *multipart.FileHeader) ([]byte, error) {

	file, err := fileHeader.Open()

	if err != nil {
		return nil, err
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	if err := validateFileSize(fileHeader.Size); err != nil {
		return nil, err
	}

	if err := validateImageFormat(fileHeader); err != nil {
		return nil, err
	}

	imgBytes, err := readImgFromFile(file)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil

}

func validateFileSize(fSize int64) error {

	if fSize > MaxUploadSize {
		return fmt.Errorf("cannot process a file larger than %d bytes", MaxUploadSize)
	}

	return nil

}

func validateImageFormat(fileHeader *multipart.FileHeader) error {

	format := filepath.Ext(fileHeader.Filename)

	if !SupportedImageFormats[format] {
		return fmt.Errorf("invalid file format: %s", format)
	}

	return nil

}

func readImgFromFile(file multipart.File) ([]byte, error) {

	imgBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return imgBytes, nil

}

func processImg(opProcessor func(file []byte, parameters map[string]interface{}) ([]byte, error), originalImg ImagePayload, opBody OperationBody, wg *sync.WaitGroup, processedImgsChan chan ImagePayload) {

	defer wg.Done() // Notify the WaitGroup that this goroutine is done

	imgBytes, err := opProcessor(originalImg.payload, opBody.Parameters)

	imgPayload := ImagePayload{
		fileName: originalImg.fileName,
		payload:  imgBytes,
		err:      err,
	}

	processedImgsChan <- imgPayload

}

func buildResponse(writer *multipart.Writer, payload ImagePayload) error {

	filename := payload.fileName

	formFileWriter, err := writer.CreateFormFile(filename, filename)

	if err != nil {
		return err
	}

	processedImgBytes := payload.payload

	if _, err := formFileWriter.Write(processedImgBytes); err != nil {
		return err
	}

	return nil

}

func buildMultiPartWriter() (*bytes.Buffer, *multipart.Writer, error) {

	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	_, err := writer.CreateFormField("images")

	if err != nil {
		return &bytes.Buffer{}, nil, err
	}

	return &buf, writer, nil
}

func closeMultiPartWriterAndWriteResponse(w http.ResponseWriter, multiPartWriter *multipart.Writer, buf *bytes.Buffer) error {

	err := multiPartWriter.Close()

	if err != nil {
		return err
	}

	err = writeResponse(w, multiPartWriter.FormDataContentType(), buf)

	if err != nil {
		return err
	}

	return nil

}

func writeResponse(w http.ResponseWriter, formDataContentType string, buf *bytes.Buffer) error {

	w.Header().Set("Content-Type", formDataContentType)

	_, err := w.Write(buf.Bytes())

	if err != nil {
		return err
	}

	return nil

}
