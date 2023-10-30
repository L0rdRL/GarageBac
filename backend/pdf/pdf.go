package pdf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/unidoc/unipdf/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Document struct {
	Link   string `json:"link"`
	Status string `json:"status"`
	Type   string `json:"type"`
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

func showPDF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pdfID := r.URL.Query().Get("pdfID")
	if pdfID == "" {
		http.Error(w, "ID PDF не предоставлен", http.StatusBadRequest)
		return
	}

	// PDF, sort and filter it based on your criteria.
	pdf, err := findPDFByID(pdfID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pdf == nil {
		http.Error(w, "PDF не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=output.pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdf)))
	_, err = w.Write(pdf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Your list of PDFs, sorted and filtered, should be presented as a JSON response.
	// Example for demonstration:
	pdfList := []Document{
		{Link: "pdf1.pdf", Status: "published", Type: "manual"},
		{Link: "pdf2.pdf", Status: "draft", Type: "report"},
		// Add other items to the list
	}
	jsonResponse, err := json.Marshal(pdfList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func deletePDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pdfID := ps.ByName("pdfID")

	// Найдите и удалите PDF по идентификатору pdfID
	for i, doc := range documentDB {
		if doc.Link == pdfID {
			// Удалите PDF из вашей базы данных (или хранилища)
			documentDB = append(documentDB[:i], documentDB[i+1:]...)
			break
		}
	}

	// Отправьте ответ об успешном удалении
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "PDF удален успешно")
}
