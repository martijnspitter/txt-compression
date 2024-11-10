package reader

import "os"

type Writer struct {
	headerText     string
	compressedText string
	outputFile     string
}

const (
	headerStart = "[H-START]"
	headerEnd   = "[H-END]"
)

func NewWriter(outputFile, headerText, compressedText string) *Writer {
	return &Writer{
		headerText,
		compressedText,
		outputFile,
	}
}

func (w *Writer) createHeader() string {
	return headerStart + w.headerText + headerEnd
}

func (w *Writer) WriteFile() error {
	return os.WriteFile(w.outputFile, []byte(w.createHeader()+w.compressedText), 0644)
}
