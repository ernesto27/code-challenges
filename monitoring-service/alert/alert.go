package main

import (
	"fmt"
	"monitoring/db"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

func SendEmail(to string, body string) error {
	apiKey := ""

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Monitor system <onboarding@resend.dev>",
		To:      []string{"ernestoponce27@gmail.com"},
		Html:    body,
		Subject: "Monitor system alert",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	fmt.Println(sent.Id)
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	myDB, dbErr := db.NewMysql(host, user, password, port, database)
	if dbErr != nil {
		panic(dbErr)
	}

	notifications, err := myDB.GetNotifications()
	if err != nil {
		panic(err)
	}

	for _, notification := range notifications {
		err := SendEmail("ernestoponce27@gmail.com", notification.Message)
		if err != nil {
			fmt.Println(err)
		} else {
			err := myDB.UpdateNotificationSent(notification.ID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
