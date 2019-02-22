package xlsxReader

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ntfox0001/svrLib/log"
)

type XlsxReader struct {
	xlsxFile *excelize.File
}

func NewXlsxReader(filename string) (*XlsxReader, error) {
	if xlsx, err := excelize.OpenFile(filename); err != nil {
		log.Error("No found file", "err", err.Error())
		return nil, err
	} else {
		reader := XlsxReader{
			xlsxFile: xlsx,
		}

		return &reader, nil
	}
}
