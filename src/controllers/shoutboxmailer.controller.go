package controllers

import (
	"fmt"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/models"
	"net/smtp"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ShoutboxMailer() *mongo.Collection {
	return db.GetCollection("magic899_db", "shoutboxmailer")
}

func SendAutoReplyShoutboxMailer(c fiber.Ctx) error {
	smtpEmail := "magic899shoutbox@gmail.com"
	smtpPassword := "qbyq rtay vqml tzkn"

	collection := ShoutboxMailer()

	var req models.RequestShoutbox
	if err := c.Bind().JSON(&req); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Initialize XSS sanitizer
	policy := bluemonday.UGCPolicy()

	// Sanitize user input
	req.Name = policy.Sanitize(req.Name)
	req.School_name = policy.Sanitize(req.School_name)
	req.Email = policy.Sanitize(req.Email)
	req.Position = policy.Sanitize(req.Position)
	req.Contact = policy.Sanitize(req.Contact)
	req.School_contact = policy.Sanitize(req.School_contact)
	req.Organization = policy.Sanitize(req.Organization)
	req.Title = policy.Sanitize(req.Title)
	req.Event_date = policy.Sanitize(req.Event_date)
	req.Radio_spiel = policy.Sanitize(req.Radio_spiel)

	// Load SMTP credentials from environment variables
	from := smtpEmail
	password := smtpPassword
	if from == "" || password == "" {
		return c.Status(500).SendString("Email server credentials are missing")
	}

	// SMTP Config
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authenticate SMTP
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Prevent email header injection
	if len(req.Email) > 254 || len(req.Name) > 100 {
		return c.Status(400).SendString("Invalid email format")
	}

	// Email Content for the requester
	subject := "Shoutbox Inquiry - " + req.Title
	bodyRequester := fmt.Sprintf(`
		<html>
		<body>
			<h3>Hello %s,</h3>
			<p>Thank you for your shoutbox request regarding "<b>%s</b>". We have received your details and will get back to you soon.</p>
			<p><b>Event Date:</b> %s</p>
			<p><b>Radio Spiel:</b> %s</p>
			<p>Best regards,<br>Magic 899 Team</p>
		</body>
		</html>
	`, req.Name, req.Title, req.Event_date, req.Radio_spiel)

	// Email Content for Admin
	subjectAdmin := "New Shoutbox Inquiry Received"
	bodyAdmin := fmt.Sprintf(`
		<html>
		<body>
			<h3>Hello Magic Shoutbox,</h3>
			<p>A new shoutbox inquiry has been received.</p>
			<p><b>Name:</b> %s</p>
			<p><b>Email:</b> %s</p>
			<p><b>School Name:</b> %s</p>
			<p><b>Phone Number:</b> %s</p>
			<p><b>School Phone Number:</b> %s</p>
			<p><b>Title:</b> %s</p>
			<p><b>Event Date:</b> %s</p>
			<p><b>Radio Spiel:</b> %s</p>
			<p>Best regards,<br>Magic 899 System</p>
		</body>
		</html>
	`,
		req.Name, req.Email, req.School_name, req.Contact, req.School_contact, req.Title, req.Event_date, req.Radio_spiel)

	// Function to send email
	sendEmail := func(to []string, subject, body string) error {
		message := []byte(
			"MIME-version: 1.0;\n" +
				"Content-Type: text/html; charset=\"UTF-8\";\n" +
				"From: \"Magic 899 Shoutbox\" <" + from + ">\n" +
				"To: " + to[0] + "\n" +
				"Subject: " + subject + "\n\n" + body)
		return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	}

	// Send email to requester
	err := sendEmail([]string{req.Email}, subject, bodyRequester)
	if err != nil {
		fmt.Println("Error sending email to requester:", err)
		return c.Status(500).SendString("Failed to send confirmation email")
	}

	// Send email to admin
	err = sendEmail([]string{smtpEmail}, subjectAdmin, bodyAdmin)
	if err != nil {
		fmt.Println("Error sending email to admin:", err)
		return c.Status(500).SendString("Failed to send inquiry email to admin")
	}

	// Save email request to MongoDB
	now := primitive.NewDateTimeFromTime(time.Now())
	mailer := models.RequestShoutbox{
		ID:             primitive.NewObjectID(),
		Name:           req.Name,
		School_name:    req.School_name,
		Email:          req.Email,
		Position:       req.Position,
		Contact:        req.Contact,
		School_contact: req.School_contact,
		Organization:   req.Organization,
		Title:          req.Title,
		Event_date:     req.Event_date,
		Radio_spiel:    req.Radio_spiel,
		Created_at:     now,
		Updated_at:     now,
	}

	_, err = collection.InsertOne(c.Context(), mailer)
	if err != nil {
		fmt.Println("Error saving to MongoDB:", err)
		return c.Status(500).SendString("Failed to save email request")
	}

	return c.SendString("Auto-reply email sent successfully!")
}
