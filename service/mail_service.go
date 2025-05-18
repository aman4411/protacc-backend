package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aman4411/protacc-backend/utils"
)

type MailService struct {
	apiKey    string
	fromEmail string
	client    *http.Client
	env       string
}

type resendEmail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

type resendError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Name       string `json:"name"`
}

func NewMailService() *MailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		env = "dev" // Default to development if not set
	}

	if apiKey == "" && env != "dev" {
		fmt.Println("WARNING: RESEND_API_KEY is not set")
	}
	if fromEmail == "" {
		fmt.Println("WARNING: FROM_EMAIL is not set")
		// Set a default sender domain from Resend
		fromEmail = "onboarding@resend.dev"
	}

	return &MailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{},
		env:       env,
	}
}

func (s *MailService) sendEmail(to, subject, html string) error {
	// In development, just log the email instead of sending it
	if s.env == "dev" {
		utils.LogInfo("Email not sent (development mode)",
			"to", to,
			"subject", subject,
			"content", html)
		return nil
	}

	if s.apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not set")
	}

	email := resendEmail{
		From:    fmt.Sprintf("ProtAcc Team <%s>", s.fromEmail),
		To:      to,
		Subject: subject,
		HTML:    html,
	}

	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("error marshaling email: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Read error response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading error response: %v", err)
		}

		var resendErr resendError
		if err := json.Unmarshal(body, &resendErr); err != nil {
			return fmt.Errorf("error sending email, status code: %d, body: %s", resp.StatusCode, string(body))
		}

		return fmt.Errorf("resend API error: %s (%s)", resendErr.Message, resendErr.Name)
	}

	return nil
}

func (s *MailService) SendVerificationEmail(toEmail, firstName, otp string) error {
	subject := "Verify your email address"
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.otp { font-size: 32px; font-weight: bold; color: #2563eb; letter-spacing: 2px; }
				.footer { margin-top: 30px; font-size: 14px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Welcome to ProtAcc!</h2>
				<p>Hi %s,</p>
				<p>Thank you for signing up. Please use the following OTP to verify your email address:</p>
				<p class="otp">%s</p>
				<p>This OTP will expire in 15 minutes.</p>
				<p>If you didn't request this verification, please ignore this email.</p>
				<div class="footer">
					<p>Best regards,<br>The ProtAcc Team</p>
				</div>
			</div>
		</body>
		</html>
	`, firstName, otp)

	return s.sendEmail(toEmail, subject, html)
}

func (s *MailService) SendWelcomeEmail(toEmail, firstName string) error {
	subject := "Welcome to ProtAcc!"
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.cta-button { 
					display: inline-block; 
					padding: 12px 24px; 
					background-color: #2563eb; 
					color: white; 
					text-decoration: none; 
					border-radius: 6px; 
					margin: 20px 0;
				}
				.footer { margin-top: 30px; font-size: 14px; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Welcome to ProtAcc!</h2>
				<p>Hi %s,</p>
				<p>Thank you for verifying your email address. Your account is now active and you can start using our services.</p>
				<a href="https://protacc.netlify.app/services" class="cta-button">Explore Our Services</a>
				<p>If you have any questions or need assistance, feel free to contact our support team.</p>
				<div class="footer">
					<p>Best regards,<br>The ProtAcc Team</p>
				</div>
			</div>
		</body>
		</html>
	`, firstName)

	return s.sendEmail(toEmail, subject, html)
}
