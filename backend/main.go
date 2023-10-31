package main

import (
    "fmt"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/dgrijalva/jwt-go"
    "github.com/unidoc/unipdf/v3"
    "github.com/L0rdRL/p1/handlers"
    "github.com/L0rdRL/p1/pdf"
    "net/smtp"
    "math/rand"
    "encoding/json"

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




// User authentication middleware
func loginUser(w http.ResponseWriter, r *http.Request) {
    // Получить логин и пароль из запроса
    username := r.FormValue("username")
    password := r.FormValue("password")

    // Найдите пользователя в вашей базе данных по имени пользователя
    user := findUserByUsername(username)

    if user == nil || !comparePasswords(user.Password, password) {
        http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
        return
    }

    // Если пользователь существует и пароль совпадает, выдайте ему JWT-токен
    token := generateToken(user.Username)

    // Отправьте токен обратно пользователю
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"token": token}
    json.NewEncoder(w).Encode(response)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
    // Получить данные пользователя из запроса
    var newUser User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        http.Error(w, "Неверный запрос", http.StatusBadRequest)
        return
    }

    // Проверьте, существует ли пользователь с таким же именем
    if findUserByUsername(newUser.Username) != nil {
        http.Error(w, "Пользователь с таким именем уже существует", http.StatusConflict)
        return
    }

    // Хэшируйте пароль
    hashedPassword := hashAndSaltPassword(newUser.Password)

    // Сохраните данные пользователя в базе данных
    newUser.Password = hashedPassword
    userDB = append(userDB, newUser)

    // Ответьте успехом
    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, "Пользователь зарегистрирован успешно")
}

func generateActivationCode() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*-№;%:?*()_+=.,"
    code := make([]byte, 16)
    for i := range code {
        code[i] = charset[rand.Intn(len(charset))]
    }
    return string(code)
}

func sendActivationEmail(userEmail, activationCode string) {
    // Настройки SMTP сервера
    smtpServer := "smtp.example.com"
    smtpPort := "587"
    senderEmail := "your-email@example.com"
    senderPassword := "your-email-password"

    // Формирование письма
    subject := "Подтверждение регистрации"
    body := fmt.Sprintf("Для завершения регистрации перейдите по ссылке: http://your-website.com/activate?code=%s", activationCode)

    message := []byte("To: " + userEmail + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" +
        body + "\r\n")

    // Отправка письма
    auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
    err := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, []string{userEmail}, message)
    if err != nil {
        log.Fatal(err)
    }
}

func generateToken(userID string) (string, error) {
    // Создайте токен с данными о пользователе
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        // Здесь вы можете добавить другие данные о пользователе
    })

    // Подпишите токен с использованием секретного ключа
    tokenString, err := token.SignedString([]byte("ваш_секретный_ключ"))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

router.POST("/documents", handlers.CreateDocumentHandler)
router.PUT("/documents/:documentID", handlers.UpdateDocumentHandler)
router.DELETE("/documents/:documentID", handlers.DeleteDocument)


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
    activationCode := generateActivationCode()
    userEmail := "user@example.com"

    sendActivationEmail(userEmail, activationCode)
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