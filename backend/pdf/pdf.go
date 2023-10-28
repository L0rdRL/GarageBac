package pdf


import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/julienschmidt/httprouter"
    "github.com/pkg/errors"
    "github.com/unidoc/unipdf/v3"
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
    w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes))
    _, err = w.Write(pdfBytes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func updatePDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    pdfID := ps.ByName("pdfID")
    // Find an existing PDF by pdfID. Implement PDF search logic in your system.
    foundPDF, err := findPDFByID(pdfID)
    if err != nil {
        // Handle the error when searching for a PDF
        http.Error(w, "An error occurred while searching for a PDF", http.StatusInternalServerError)
        return
    }
    if foundPDF == nil {
        // PDF not found, return an error
        http.Error(w, "PDF not found", http.StatusNotFound)
        return
    }
    // Update your PDF data. At this point, you need to

    err = savePDF(foundPDF)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "PDF updated successfully")
}

func findPDFByID(pdfID string) ([]byte, error) {
    // Получаем коллекцию PDF-документов
    collection := client.Database("pdfsdb").Collection("pdfs")

    // Поиск PDF по его идентификатору
    filter := bson.M{"_id": pdfID}
    var pdf Document

    if err := collection.FindOne(context.TODO(), filter).Decode(&pdf); err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("PDF не найден")
        }
        return nil, err
    }

    return pdf.Content, nil
}
