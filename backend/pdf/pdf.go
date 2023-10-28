package pdf


import (
	"encoding/json"
    "fmt"
    "net/http"
    "strings"
)

ype Document struct {
    Link    string `json:"link"`
    Status  string `json:"status"`
    Type    string `json:"type"`
    // You can add more fields as needed
}

// Create a new PDF
func createPDF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // Your PDF creation code here
    pdf := unipdf.NewPdf()
	pdf.AppendPage()

    pdfBytes, err := pdf.Bytes()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename=output.pdf")
    w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
    _, err = w.Write(pdfBytes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return

	}
}

func updatePDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    pdfID := ps.ByName("pdfID")
	//Find an existing PDF by pdfID. Implement PDF search logic in your system.
	

	pdf := unipdf.NewPdf()
    pdf.AppendPage()
	// Implement the code to search for PDF in your system.
	if unipdf == nil {
        // PDF не найден, возвращаем ошибку
        http.Error(w, "PDF не найден", http.StatusNotFound)
        return
    }
	// Update your PDF data. At this point you need to

	err := savePDF(uniPDF)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}