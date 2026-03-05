package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_SMTP_HOST")
	portStr := os.Getenv("EMAIL_SMTP_PORT")

	if from == "" || password == "" || host == "" || portStr == "" {
		return fmt.Errorf("variáveis de email não configuradas: FROM=%s, PASS=%s, HOST=%s, PORT=%s", from, password, host, portStr)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("porta smtp invalida: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, from, password)
	d.TLSConfig = &tls.Config{ServerName: host}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Erro ao enviar email para %s: %v", to, err)
		return err
	}

	log.Printf("Email enviado com sucesso para %s", to)
	return nil
}
