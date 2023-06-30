package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"goes/webServiceWithGoPluralSight/cors"
)

const receiptPath = "receipts"

func SetupRoute(apiBasePath string) {
	receiptsHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))
}

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		err := r.ParseMultipartForm(2 << 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return // cors.Middleware handles the logic for us
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	stat, err := os.Stat(file.Name())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)
}
