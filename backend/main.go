package main

import (
    "fmt"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/dgrijalva/jwt-go"
    "github.com/unidoc/unipdf/v3"
)

type User struct {
    Username   string `json:"username"`
    Password   string `json:"password"`
    FirstName  string `json:"first_name"`
    Surname    string `json:"surname"`
    Patronymic string `json:"patronymic"`
    Position   string `json:"position"`
    Role       string `json:"role"`
}

type ServerConfig struct {
    Address      string `json:"address"`
    Port         int    `json:"port"`
    MongoDBURI   string `json:"mongodb_uri"`

}

type Document struct {
    Link    string `json:"link"`
    Status  string `json:"status"`
    Type    string `json:"type"`
    // You can add more fields as needed
}

var documentDB []Document


var secretKey = []byte("your-secret-key")


var userDB []User = []User{
    {
        Username:   "admin",
        Password:   "admin_password",
        FirstName:  "Admin",
        Surname:    "User",
        Patronymic: "None",
        Position:   "Administrator",
        Role:       "admin",
    },
    // Add other user entries
    

}

// User registration
func registerUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var newUser User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Hash and salt the password (you should use a proper library for this)
    hashedPassword := hashAndSaltPassword(newUser.Password)

    // Store the user information in the database
    newUser.Password = hashedPassword
    userDB = append(userDB, newUser)

    // Respond with success message
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, "User registered successfully")
}

// User authentication middleware
func authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Verify the username and password (you should use a proper library for this)
        username, password := extractUsernamePassword(tokenString)
        user := findUserByUsername(username)

        if user == nil || !comparePasswords(user.Password, password) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if user.Role != "admin" {
            http.Error(w, "Access Denied", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

router.POST("/documents", authenticate(createDocument))
router.PUT("/documents/:documentID", authenticate(updateDocument))
router.DELETE("/documents/:documentID", authenticate(deleteDocument))


func extractUsernamePassword(tokenString string) (string, string) {
    // Parse the username and password from the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    // Example: "Basic username:password"
    parts := strings.Split(token.Header["Authorization"], " ")
    if len(parts)!= 2 {
        return "", ""
    }

    // Split the token to extract the username and password
    usernameAndPassword := strings.Split(parts[1], ":")
    if len(usernameAndPassword)!= 2 {
        return "", ""
    }

    // You should implement proper parsing logic here
    return usernameAndPassword[0], usernameAndPassword[1]

    // For demonstration purposes, we assume a simple "username:password" format.
    parts := strings.Split(tokenString, " ")
    if len(parts) == 2 {
        token := parts[1]
        credentials := strings.Split(token, ":")
        if len(credentials) == 2 {
            username := credentials[0]
            password := credentials[1]
            return username, password
        }
    }
    return "", ""
}

func comparePasswords(hashedPassword, plainPassword string) bool {
    // Compare hashed and plain passwords (you should use a proper library for this)
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil

    // Implement proper password comparison logic here
    return hashedPassword == plainPassword
}

router.POST("/register", registerUser)
router.GET("/protected", authenticate(http.HandlerFunc(protected)))


func findUserByUsername(username string) *User {
    // Find the user in the database by username
    for _, user := range userDB {
        if user.Username == username {
            return &user
        }
    }
    return nil
}
// User authorization middleware
func authorize(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return secretKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}




// Show PDFs with sorting and filters
func showPDFs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // Your PDF listing code here
    pdf, err := unipdf.New()
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}

// Delete a PDF
func deletePDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    pdfID := ps.ByName("pdfID")
    // Your PDF deletion code here
    pdf, err := unipdf.New()
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}

func loadConfig() ServerConfig {
    configFile, err := os.Open("config.json")
    if err != nil {
        return ServerConfig{
            Address:    "localhost",
            Port:       8080,
            MongoDBURI: "mongodb://localhost:27017",
            // Другие значения по умолчанию
        }
    }
    defer configFile.Close()

    decoder := json.NewDecoder(configFile)
    var config ServerConfig
    err = decoder.Decode(&config)
    if err != nil {
        log.Printf("Ошибка разбора файла конфигурации: %v. Использую значения по умолчанию.", err)
        return ServerConfig{
            Address:    "localhost",
            Port:       8080,
            MongoDBURI: "mongodb://localhost:27017",
            // Другие значения по умолчанию
        }
    }

    return config
}

func saveConfig(config ServerConfig) error {
    configFile, err := os.Create("config.json")
    if err != nil {
        return err
    }
    defer configFile.Close()

    encoder := json.NewEncoder(configFile)
    if err := encoder.Encode(config); err != nil {
        return err
    }

    return nil
}


// Где-то в вашем коде, когда вы хотите сохранить конфигурацию
config := loadConfig() // Загрузка текущей конфигурации
config.Address = "new_address"
config.Port = 8081
config.MongoDBURI = "new_mongodb_uri"

if err := saveConfig(config); err != nil {
    log.Printf("Ошибка сохранения конфигурации: %v", err)
}




func main() {
    router := httprouter.New()

    // Создаем канал для ожидания сигналов завершения работы сервера
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    // Routes for PDF operations
    router.POST("/pdfs", createPDF)
    router.PUT("/pdfs/:pdfID", updatePDF)
    router.GET("/pdfs", showPDFs)
    router.DELETE("/pdfs/:pdfID", deletePDF)

    // Start the server
    // Устанавливаем соединение с MongoDB
	// Устанавливаем соединение с MongoDB
    clientOptions := options.Client().ApplyURI("mongodb+srv://arcklitrl:101899Aeza@cl1.uapna8x.mongodb.net")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Ошибка при подключении к MongoDB: %v", err)
        return
    }
    defer client.Disconnect(context.TODO())

    // Создаем базу данных и коллекцию
    db := client.Database("pdfsdb")
    collection := db.Collection("pdfs")

    fmt.Println("База данных и коллекция успешно созданы.")

    // Запускаем сервер в горутине
    go func() {
        fmt.Println("Server started on :8080")
        http.ListenAndServe(":8080", authorize(router))
    }()

    // Ожидание сигнала завершения работы сервера
    <-c

    // Закрытие базы данных и другие необходимые действия для корректного завершения
    client.Disconnect(context.TODO())
    fmt.Println("Сервер завершает работу.")
    os.Exit(0)
}

// Create a new document
func createDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    var newDoc Document
    err := json.NewDecoder(r.Body).Decode(&newDoc)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Store the document information in the database
    documentDB = append(documentDB, newDoc)

    // Respond with success message
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, "Document created successfully")
}

// Update an existing document
func updateDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    documentID := ps.ByName("documentID")
    // Find the document by ID and update its information
    for i, doc := range documentDB {
        if doc.Link == documentID {
            documentDB[i].Link = newDoc.Link
            documentDB[i].Status = newDoc.Status
            documentDB[i].Type = newDoc.Type
            break
        }
    }

    // Implement your update logic here
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Document updated successfully")

}

// List documents with sorting and filters
func listDocuments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // Implement your document listing code with sorting and filters here
    document, err := unipdf.New()
    if err!= nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}

// Delete a document
func deleteDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    documentID := ps.ByName("documentID")
    // Delete the specified document
    for i, doc := range documentDB {
        if doc.Link == documentID {
            documentDB = append(documentDB[:i], documentDB[i+1:]...)
            break
        }
    }    

    // Implement your document deletion logic here
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Document deleted successfully")
}

router.POST("/documents", createDocument)
router.PUT("/documents/:documentID", updateDocument)
router.GET("/documents", listDocuments)
router.DELETE("/documents/:documentID", deleteDocument)
