package receipt

import (
	"os"
	"path/filepath"
	"time"
)

var ReceiptDirectory string = filepath.Join("inventory-service", "uploads")

type Receipt struct {
	ReceiptName string    `json:"name"`
	UploadDate  time.Time `json:"uploadDate"`
}

func GetReceipts() ([]Receipt, error) {
	files, err := os.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}
	receipts := make([]Receipt, 0)
	for _, f := range files {
		fileInfo, err := f.Info()
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: fileInfo.ModTime()})
	}
	return receipts, nil
}
