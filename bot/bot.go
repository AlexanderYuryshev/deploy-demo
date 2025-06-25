package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID    string
	Name  sql.NullString
	Email sql.NullString
}

type Post struct {
	ID      int
	Title   string
	Created time.Time
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DATABASE := os.Getenv("POSTGRES_DATABASE")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	connStr := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DATABASE,
	)
	log.Printf("Connecting to database: %s@%s:5432/%s", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DATABASE)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	log.Println("Successfully connected to database")

	mainMenuButton := tgbotapi.NewKeyboardButton("Список пользователей")
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(mainMenuButton))
	keyboard.ResizeKeyboard = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, db, update.Message, keyboard)
		} else if update.CallbackQuery != nil {
			handleCallback(bot, db, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, db *sql.DB, msg *tgbotapi.Message, keyboard tgbotapi.ReplyKeyboardMarkup) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)

	if msg.Text == "/start" {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Добро пожаловать! Используйте кнопки для навигации")
		reply.ReplyMarkup = keyboard
		bot.Send(reply)
		return
	}

	switch msg.Text {
	case "Список пользователей":
		users, err := getUsers(db)
		if err != nil {
			log.Printf("Error getting users: %v", err)
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Ошибка получения пользователей")
			bot.Send(reply)
			return
		}

		if len(users) == 0 {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Пользователи не найдены")
			bot.Send(reply)
			return
		}

		inlineKeyboard := createUsersKeyboard(users)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Выберите пользователя:")
		reply.ReplyMarkup = inlineKeyboard
		bot.Send(reply)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, db *sql.DB, callback *tgbotapi.CallbackQuery) {
	callbackData := callback.Data
	chatID := callback.Message.Chat.ID

	if strings.HasPrefix(callbackData, "user_") {
		userID := strings.TrimPrefix(callbackData, "user_")
		posts, err := getUserPosts(db, userID)
		if err != nil {
			log.Printf("Error getting posts: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения постов"))
			return
		}

		response := formatPosts(posts, userID)
		bot.Send(tgbotapi.NewMessage(chatID, response))
	}

	callbackCfg := tgbotapi.NewCallback(callback.ID, "")
	bot.Send(callbackCfg)
}

func getUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, name, email FROM "users"`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUserPosts(db *sql.DB, userID string) ([]Post, error) {
	query := `SELECT id, name, created_at FROM posts WHERE created_by_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Created)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func createUsersKeyboard(users []User) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, user := range users {
		displayName := "Без имени"
		if user.Name.Valid && user.Name.String != "" {
			displayName = user.Name.String
		} else if user.Email.Valid && user.Email.String != "" {
			displayName = user.Email.String
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(displayName, "user_"+user.ID)
		row := tgbotapi.NewInlineKeyboardRow(btn)
		rows = append(rows, row)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func formatPosts(posts []Post, userID string) string {
	if len(posts) == 0 {
		return fmt.Sprintf("У пользователя %s нет постов", userID)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Посты пользователя %s:\n\n", userID))
	for _, post := range posts {
		builder.WriteString(fmt.Sprintf("📝 *%s*\n", post.Title))
		builder.WriteString(fmt.Sprintf("🆔 ID: %d\n", post.ID))
		builder.WriteString(fmt.Sprintf("📅 Дата: %s\n\n", post.Created.Format("02.01.2006 15:04")))
	}
	return builder.String()
}
