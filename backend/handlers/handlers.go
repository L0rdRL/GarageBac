package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
)

type Document struct {
	Link   string `json:"link"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

var documentsDB = []Document

func init() {
	documentsDB = append(documentsDB, Document{Link: "document1", Status: "Draft", Type: "PDF"})
	documentsDB = append(documentsDB, Document{Link: "document2", Status: "Published", Type: "Word"})
}

func CreateDocumentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newDoc Document
	err := json.NewDecoder(r.Body).Decode(&newDoc)
	if err != nil {
		http.Error(w, "Невреный запрос", http.StatusBadRequest)
		return
	}
	documentsDB = append(documentsDB, newDoc)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Документ успешно создан")
}

func UpdateDocumentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	documentID := ps.ByName("documentID")
	var updateDoc Document
	found := false
	for i, doc := range documentsDB {
		if doc.Link == documentID {
			documentsDB[i].Link = updateDoc.Link
			documentsDB[i].Status = updateDoc.Status
			documentsDB[i].Type = updateDoc.Type
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Документ не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Документ успешно обновлен")
}

func ListDocumentsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var documents []Document
	for _, doc := range documentsDB {
		documents = append(documents, doc)
	}
	sortingField := r.URL.Query().Get("sort_by")
	filterType := r.URL.Query().Get("type")

	// sort
	if sortingField == "asc" {
		sort.Slice(documents, func(i, j int) bool {
			return documents[i].Link < documents[j].Link
		})
	} else if sortingField == "desc" {
		sort.Slice(documents, func(i, j int) bool {
			return documents[i].Link > documents[j].Link
		})
	}
	// filter
	if filterType != "" {
		var filteredDocuments []Document
		for _, doc := range documents {
			if doc.Type == filterType {
				filteredDocuments = append(filteredDocuments, doc)
			}
		}
		documents = filteredDocuments
	}
	json.NewEncoder(w).Encode(documents)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	documentID := ps.ByName("documentID")

	for i, doc := range documentsDB {
		if doc.ID == documentID {
			documentsDB = append(documentsDB[:i], documentsDB[i+1:]...)
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Document deleted successfully"))
            return
		}
	}
	http.Error(w, "Document not found", http.StatusNotFound)
}

