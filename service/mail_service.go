package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aman4411/protacc-backend/models"
)

type MailService struct {
	apiKey    string
	fromEmail string
	client    *http.Client
	env       string
}

type resendEmailTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type resendEmail struct {
	From    string            `json:"from"`
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html"`
	Text    string            `json:"text,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Tags    []resendEmailTag  `json:"tags,omitempty"`
	ReplyTo string            `json:"reply_to,omitempty"`
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
	fromEmail = strings.TrimSpace(fromEmail)
	if fromEmail == "" {
		fmt.Println("WARNING: FROM_EMAIL is not set, using onboarding@resend.dev")
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
	//if s.env == "dev" {
	//	utils.LogInfo("Email not sent (development mode)",
	//		"to", to,
	//		"subject", subject,
	//		"content", html)
	//	return nil
	//}

	if s.apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not set")
	}

	// Create plain text version for better deliverability
	plainText := s.htmlToPlainText(html)

	email := resendEmail{
		From:    fmt.Sprintf("ProtAcc Support <%s>", s.fromEmail),
		To:      []string{strings.TrimSpace(to)},
		Subject: subject,
		HTML:    html,
		Text:    plainText,
		ReplyTo: s.fromEmail,
		Tags: []resendEmailTag{
			{Name: "category", Value: "transactional"},
		},
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
	subject := "Verify Your ProtAcc Account - Email Confirmation Required"
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Verify Your Email - ProtAcc</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; 
					line-height: 1.6; 
					color: #374151; 
					background-color: #f8fafc;
				}
				.email-container { 
					max-width: 600px; 
					margin: 0 auto; 
					background-color: #ffffff; 
					border-radius: 16px; 
					overflow: hidden; 
					box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
				}
				.header { 
					background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); 
					padding: 40px 30px; 
					text-align: center; 
					position: relative;
					overflow: hidden;
				}
				.header::before {
					content: '';
					position: absolute;
					top: 0;
					left: 0;
					right: 0;
					bottom: 0;
					background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="white" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%%23grain)"/></svg>');
				}
				.logo { 
					font-size: 42px; 
					font-weight: 900; 
					color: white; 
					margin-bottom: 8px; 
					letter-spacing: -1px;
					position: relative;
					z-index: 1;
				}
				.tagline { 
					color: rgba(255, 255, 255, 0.9); 
					font-size: 16px; 
					font-weight: 500;
					position: relative;
					z-index: 1;
				}
				.content { 
					padding: 40px 30px; 
				}
				.greeting { 
					font-size: 24px; 
					font-weight: 700; 
					color: #1f2937; 
					margin-bottom: 16px; 
				}
				.message { 
					font-size: 16px; 
					color: #6b7280; 
					margin-bottom: 32px; 
					line-height: 1.7;
				}
				.otp-container { 
					background: linear-gradient(135deg, #f0f9ff 0%%, #e0e7ff 100%%); 
					border: 2px solid #e5e7eb; 
					border-radius: 12px; 
					padding: 32px; 
					text-align: center; 
					margin: 32px 0; 
					position: relative;
				}
				.otp-label { 
					font-size: 14px; 
					font-weight: 600; 
					color: #6366f1; 
					text-transform: uppercase; 
					letter-spacing: 1px; 
					margin-bottom: 12px; 
				}
				.otp { 
					font-size: 36px; 
					font-weight: 900; 
					color: #4f46e5; 
					letter-spacing: 8px; 
					font-family: 'Courier New', monospace; 
					background: white; 
					padding: 16px 24px; 
					border-radius: 8px; 
					display: inline-block; 
					border: 2px solid #e5e7eb;
					box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
				}
				.otp-note { 
					font-size: 14px; 
					color: #9ca3af; 
					margin-top: 16px; 
					font-style: italic; 
				}
				.security-note { 
					background-color: #fef3c7; 
					border-left: 4px solid #f59e0b; 
					padding: 16px; 
					margin: 24px 0; 
					border-radius: 0 8px 8px 0; 
				}
				.security-note p { 
					font-size: 14px; 
					color: #92400e; 
					margin: 0; 
				}
				.footer { 
					background-color: #f9fafb; 
					padding: 32px 30px; 
					text-align: center; 
					border-top: 1px solid #e5e7eb; 
				}
				.footer-text { 
					font-size: 14px; 
					color: #6b7280; 
					margin-bottom: 16px; 
				}
				.social-links { 
					margin-top: 20px; 
				}
				.social-links a { 
					display: inline-block; 
					margin: 0 8px; 
					color: #6b7280; 
					text-decoration: none; 
					font-size: 14px; 
				}
				.divider { 
					height: 1px; 
					background: linear-gradient(90deg, transparent, #e5e7eb, transparent); 
					margin: 24px 0; 
				}
				@media only screen and (max-width: 600px) {
					.email-container { margin: 10px; border-radius: 12px; }
					.header { padding: 30px 20px; }
					.content { padding: 30px 20px; }
					.footer { padding: 24px 20px; }
					.logo { font-size: 36px; }
					.otp { font-size: 32px; letter-spacing: 6px; padding: 14px 20px; }
				}
			</style>
		</head>
		<body>
			<div style="padding: 20px 0;">
				<div class="email-container">
					<!-- Header -->
					<div class="header">
						<div class="logo">ProtAcc</div>
						<div class="tagline">Professional Accounting & Compliance Services</div>
					</div>
					
					<!-- Content -->
					<div class="content">
						<div class="greeting">Hi %s,</div>
						<div class="message">
							Welcome to ProtAcc! We're excited to have you join our community of entrepreneurs and business owners.
							<br><br>
							To complete your account setup and ensure the security of your account, please verify your email address using the verification code below:
						</div>
						
						<!-- OTP Section -->
						<div class="otp-container">
							<div class="otp-label">Your Verification Code</div>
							<div class="otp">%s</div>
							<div class="otp-note">This code expires in 15 minutes</div>
						</div>
						
						<!-- Security Note -->
						<div class="security-note">
							<p><strong>🔒 Security Tip:</strong> Never share this code with anyone. ProtAcc will never ask for your verification code via phone or email.</p>
						</div>
						
						<div class="divider"></div>
						
						<div class="message">
							Once verified, you'll have full access to our comprehensive business services including company registration, tax compliance, and expert consultancy.
							<br><br>
							If you didn't create this account, please ignore this email and the account will remain unverified.
						</div>
					</div>
					
					<!-- Footer -->
					<div class="footer">
						<div class="footer-text">
							<strong>Need help?</strong> Contact our support team at <a href="mailto:info@protacc.in" style="color: #4f46e5;">info@protacc.in</a>
						</div>
						<div class="footer-text">
							Best regards,<br>
							<strong>The ProtAcc Team</strong>
						</div>
						<div class="social-links">
							<a href="https://protacc.in">Visit Website</a> • 
							<a href="https://protacc.in/services">Our Services</a> • 
							<a href="https://protacc.in/contact">Contact Us</a>
						</div>
					</div>
				</div>
			</div>
		</body>
		</html>
	`, firstName, otp)

	return s.sendEmail(toEmail, subject, html)
}

func (s *MailService) SendWelcomeEmail(toEmail, firstName string) error {
	subject := "Welcome to ProtAcc - Account Successfully Verified"
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Welcome to ProtAcc</title>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body { 
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; 
					line-height: 1.6; 
					color: #374151; 
					background-color: #f8fafc;
				}
				.email-container { 
					max-width: 600px; 
					margin: 0 auto; 
					background-color: #ffffff; 
					border-radius: 16px; 
					overflow: hidden; 
					box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
				}
				.header { 
					background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); 
					padding: 40px 30px; 
					text-align: center; 
					position: relative;
					overflow: hidden;
				}
				.header::before {
					content: '';
					position: absolute;
					top: 0;
					left: 0;
					right: 0;
					bottom: 0;
					background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="60" r="0.5" fill="white" opacity="0.1"/><circle cx="90" cy="40" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%%23grain)"/></svg>');
				}
				.logo { 
					font-size: 42px; 
					font-weight: 900; 
					color: white; 
					margin-bottom: 8px; 
					letter-spacing: -1px;
					position: relative;
					z-index: 1;
				}
				.tagline { 
					color: rgba(255, 255, 255, 0.9); 
					font-size: 16px; 
					font-weight: 500;
					position: relative;
					z-index: 1;
				}
				.success-badge {
					background: rgba(255, 255, 255, 0.2);
					color: white;
					padding: 8px 16px;
					border-radius: 20px;
					font-size: 14px;
					font-weight: 600;
					margin-top: 16px;
					display: inline-block;
					position: relative;
					z-index: 1;
				}
				.content { 
					padding: 40px 30px; 
				}
				.greeting { 
					font-size: 28px; 
					font-weight: 700; 
					color: #1f2937; 
					margin-bottom: 16px; 
					text-align: center;
				}
				.message { 
					font-size: 16px; 
					color: #6b7280; 
					margin-bottom: 32px; 
					line-height: 1.7;
					text-align: center;
				}
				.cta-section {
					text-align: center;
					margin: 40px 0;
				}
				.cta-button { 
					display: inline-block; 
					padding: 16px 32px; 
					background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); 
					color: white !important; 
					text-decoration: none; 
					border-radius: 12px; 
					font-weight: 600;
					font-size: 16px;
					box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
					transition: all 0.3s ease;
				}
				.cta-button:hover {
					transform: translateY(-2px);
					box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
				}
				.features-grid {
					display: grid;
					grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
					gap: 24px;
					margin: 40px 0;
				}
				.feature-card {
					background: #f8fafc;
					padding: 24px;
					border-radius: 12px;
					text-align: center;
					border: 1px solid #e5e7eb;
				}
				.feature-icon {
					font-size: 32px;
					margin-bottom: 12px;
				}
				.feature-title {
					font-size: 18px;
					font-weight: 600;
					color: #1f2937;
					margin-bottom: 8px;
				}
				.feature-desc {
					font-size: 14px;
					color: #6b7280;
					line-height: 1.5;
				}
				.stats-section {
					background: linear-gradient(135deg, #f0f9ff 0%%, #e0e7ff 100%%);
					padding: 32px;
					border-radius: 12px;
					margin: 32px 0;
					text-align: center;
				}
				.stats-grid {
					display: grid;
					grid-template-columns: repeat(3, 1fr);
					gap: 24px;
					margin-top: 24px;
				}
				.stat-item {
					text-align: center;
				}
				.stat-number {
					font-size: 24px;
					font-weight: 900;
					color: #4f46e5;
					display: block;
				}
				.stat-label {
					font-size: 12px;
					color: #6b7280;
					text-transform: uppercase;
					letter-spacing: 1px;
					margin-top: 4px;
				}
				.next-steps {
					background: #fef3c7;
					border-left: 4px solid #f59e0b;
					padding: 24px;
					margin: 32px 0;
					border-radius: 0 12px 12px 0;
				}
				.next-steps h3 {
					color: #92400e;
					font-size: 18px;
					font-weight: 600;
					margin-bottom: 12px;
				}
				.next-steps ul {
					list-style: none;
					padding: 0;
				}
				.next-steps li {
					color: #92400e;
					font-size: 14px;
					margin-bottom: 8px;
					padding-left: 20px;
					position: relative;
				}
				.next-steps li::before {
					content: '✓';
					position: absolute;
					left: 0;
					color: #f59e0b;
					font-weight: bold;
				}
				.footer { 
					background-color: #f9fafb; 
					padding: 32px 30px; 
					text-align: center; 
					border-top: 1px solid #e5e7eb; 
				}
				.footer-text { 
					font-size: 14px; 
					color: #6b7280; 
					margin-bottom: 16px; 
				}
				.social-links { 
					margin-top: 20px; 
				}
				.social-links a { 
					display: inline-block; 
					margin: 0 8px; 
					color: #6b7280; 
					text-decoration: none; 
					font-size: 14px; 
				}
				.divider { 
					height: 1px; 
					background: linear-gradient(90deg, transparent, #e5e7eb, transparent); 
					margin: 32px 0; 
				}
				@media only screen and (max-width: 600px) {
					.email-container { margin: 10px; border-radius: 12px; }
					.header { padding: 30px 20px; }
					.content { padding: 30px 20px; }
					.footer { padding: 24px 20px; }
					.logo { font-size: 36px; }
					.greeting { font-size: 24px; }
					.features-grid { grid-template-columns: 1fr; }
					.stats-grid { grid-template-columns: 1fr; gap: 16px; }
					.cta-button { padding: 14px 28px; font-size: 15px; }
				}
			</style>
		</head>
		<body>
			<div style="padding: 20px 0;">
				<div class="email-container">
					<!-- Header -->
					<div class="header">
						<div class="logo">ProtAcc</div>
						<div class="tagline">Professional Accounting & Compliance Services</div>
						<div class="success-badge">✅ Account Verified Successfully</div>
					</div>
					
					<!-- Content -->
					<div class="content">
						<div class="greeting">Welcome, %s!</div>
						<div class="message">
							Congratulations! Your email has been verified and your ProtAcc account is now fully active. 
							You're now part of a community of successful entrepreneurs and business owners who trust us with their compliance needs.
						</div>
						
						<!-- CTA Section -->
						<div class="cta-section">
							<a href="https://protacc.in/services" class="cta-button">Explore Our Services</a>
						</div>
						
						<!-- Stats Section -->
						<div class="stats-section">
							<h3 style="color: #1f2937; font-size: 20px; font-weight: 600; margin-bottom: 8px;">Join 1,500+ Happy Clients</h3>
							<p style="color: #6b7280; font-size: 14px; margin-bottom: 0;">Trusted by businesses across India</p>
							<div class="stats-grid">
								<div class="stat-item">
									<span class="stat-number">1,500+</span>
									<span class="stat-label">Happy Clients</span>
								</div>
								<div class="stat-item">
									<span class="stat-number">50+</span>
									<span class="stat-label">Services</span>
								</div>
								<div class="stat-item">
									<span class="stat-number">99.8%%</span>
									<span class="stat-label">Success Rate</span>
								</div>
							</div>
						</div>
						
						<!-- Features Grid -->
						<div class="features-grid">
							<div class="feature-card">
								<div class="feature-icon">🏢</div>
								<div class="feature-title">Business Registration</div>
								<div class="feature-desc">Complete company registration services for all entity types</div>
							</div>
							<div class="feature-card">
								<div class="feature-icon">📊</div>
								<div class="feature-title">Tax Compliance</div>
								<div class="feature-desc">GST, Income Tax, and all compliance services</div>
							</div>
							<div class="feature-card">
								<div class="feature-icon">💼</div>
								<div class="feature-title">Expert Consultancy</div>
								<div class="feature-desc">Professional business advice and strategy</div>
							</div>
						</div>
						
						<!-- Next Steps -->
						<div class="next-steps">
							<h3>🎯 What's Next?</h3>
							<ul>
								<li>Browse our comprehensive service catalog</li>
								<li>Book a free consultation with our experts</li>
								<li>Get instant quotes for your business needs</li>
								<li>Track your orders and documents in real-time</li>
							</ul>
						</div>
						
						<div class="divider"></div>
						
						<div class="message">
							Our team of certified professionals is ready to help you navigate the complex world of business compliance. 
							Whether you need company registration, tax filing, or expert consultancy, we've got you covered.
							<br><br>
							<strong>Need immediate assistance?</strong> Our support team is just a message away!
						</div>
					</div>
					
					<!-- Footer -->
					<div class="footer">
						<div class="footer-text">
							<strong>Questions? We're here to help!</strong><br>
							📧 <a href="mailto:info@protacc.in" style="color: #4f46e5;">info@protacc.in</a> | 
							📞 <a href="tel:+919034819324" style="color: #4f46e5;">+91 9034819324</a>
						</div>
						<div class="footer-text">
							Best regards,<br>
							<strong>The ProtAcc Team</strong><br>
							<em>Your trusted compliance partner</em>
						</div>
						<div class="social-links">
							<a href="https://protacc.in">🌐 Visit Website</a> • 
							<a href="https://protacc.in/services">🛍️ Our Services</a> • 
							<a href="https://protacc.in/consultancy">💡 Get Consultation</a> • 
							<a href="https://protacc.in/contact">📞 Contact Us</a>
						</div>
					</div>
				</div>
			</div>
		</body>
		</html>
	`, firstName)

	return s.sendEmail(toEmail, subject, html)
}

func (s *MailService) SendPasswordResetEmail(toEmail, firstName, resetLink string) error {
	subject := "Reset Your ProtAcc Password"
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Reset Password - ProtAcc</title>
			<style>
				body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #374151; background-color: #f8fafc; margin: 0; padding: 20px; }
				.email-container { max-width: 600px; margin: 0 auto; background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
				.header { background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); padding: 40px 30px; text-align: center; }
				.logo { font-size: 36px; font-weight: 900; color: white; }
				.content { padding: 40px 30px; }
				.greeting { font-size: 22px; font-weight: 700; color: #1f2937; margin-bottom: 16px; }
				.message { font-size: 16px; color: #6b7280; margin-bottom: 24px; }
				.cta-button { display: inline-block; padding: 14px 28px; background: linear-gradient(135deg, #4f46e5 0%%, #7c3aed 100%%); color: white !important; text-decoration: none; border-radius: 10px; font-weight: 600; }
				.note { font-size: 14px; color: #9ca3af; margin-top: 24px; }
				.footer { background: #f9fafb; padding: 24px 30px; text-align: center; font-size: 14px; color: #6b7280; border-top: 1px solid #e5e7eb; }
			</style>
		</head>
		<body>
			<div class="email-container">
				<div class="header"><div class="logo">ProtAcc</div></div>
				<div class="content">
					<div class="greeting">Hi %s,</div>
					<div class="message">
						We received a request to reset your ProtAcc account password.
						Click the button below to choose a new password. This link expires in 1 hour.
					</div>
					<p style="text-align: center; margin: 32px 0;">
						<a href="%s" class="cta-button">Reset Password</a>
					</p>
					<p class="note">
						If you did not request a password reset, you can safely ignore this email.
						Your password will remain unchanged.
					</p>
					<p class="note" style="word-break: break-all;">Or copy this link: %s</p>
				</div>
				<div class="footer">
					<strong>The ProtAcc Team</strong><br>
					<a href="mailto:info@protacc.in" style="color: #4f46e5;">info@protacc.in</a>
				</div>
			</div>
		</body>
		</html>
	`, firstName, resetLink, resetLink)

	return s.sendEmail(toEmail, subject, html)
}

// htmlToPlainText converts HTML email to plain text for better deliverability
func (s *MailService) htmlToPlainText(html string) string {
	// Simple HTML to text conversion for better spam filtering
	text := html

	// Remove style and script tags completely (Go regexp has no backreferences)
	text = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`).ReplaceAllString(text, "")

	// Convert common HTML elements to text equivalents
	text = regexp.MustCompile(`(?i)<br[^>]*>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)<p[^>]*>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)<div[^>]*>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)<h[1-6][^>]*>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</h[1-6]>`).ReplaceAllString(text, "\n")

	// Convert links to readable format
	text = regexp.MustCompile(`(?i)<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`).ReplaceAllString(text, "$2 ($1)")

	// Remove all remaining HTML tags
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")

	// Replace HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&hellip;", "...")

	// Clean up whitespace and formatting
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)

	// Add some structure for readability
	text = strings.ReplaceAll(text, "ProtAcc", "\n=== ProtAcc ===")
	text = strings.ReplaceAll(text, "Your Verification Code", "\n--- Your Verification Code ---")
	text = strings.ReplaceAll(text, "Welcome aboard", "\n=== Welcome aboard")

	return text
}

// ===== Order emails =====

// formatINR renders an amount as a rupee string, trimming trailing zeros (₹999, ₹1499.5).
func formatINR(v float64) string {
	return "₹" + strconv.FormatFloat(v, 'f', -1, 64)
}

// orderStatusLabel converts an order status enum into a human-friendly label.
func orderStatusLabel(status models.OrderStatus) string {
	switch status {
	case models.OrderStatusPendingBookingPayment:
		return "Pending booking payment"
	case models.OrderStatusBookingAmountReceived:
		return "Booking amount received"
	case models.OrderStatusProcessing:
		return "Processing"
	case models.OrderStatusDocumentsRequired:
		return "Documents required"
	case models.OrderStatusDocumentsReceived:
		return "Documents received"
	case models.OrderStatusInProgress:
		return "In progress"
	case models.OrderStatusPendingFinalPayment:
		return "Pending final payment"
	case models.OrderStatusFullPaymentReceived:
		return "Full payment received"
	case models.OrderStatusCompleted:
		return "Completed"
	case models.OrderStatusCancelled:
		return "Cancelled"
	default:
		return string(status)
	}
}

// orderDetailsHTML builds the shared order-summary block (items table + totals).
func (s *MailService) orderDetailsHTML(order *models.Order) string {
	var rows strings.Builder
	for _, it := range order.Items {
		name := "Service"
		if it.Service != nil && it.Service.Name != "" {
			name = it.Service.Name
		}
		rows.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding:12px 0;border-bottom:1px solid #eef2f7;color:#374151;">%s</td>
				<td style="padding:12px 0;border-bottom:1px solid #eef2f7;text-align:center;color:#6b7280;">%d</td>
				<td style="padding:12px 0;border-bottom:1px solid #eef2f7;text-align:right;color:#111827;font-weight:600;white-space:nowrap;">%s</td>
			</tr>`, html.EscapeString(name), it.Quantity, formatINR(it.Price)))
	}

	discountRow := ""
	if order.DiscountAmount > 0 {
		code := ""
		if order.CouponCode != nil && *order.CouponCode != "" {
			code = " (" + html.EscapeString(*order.CouponCode) + ")"
		}
		discountRow = fmt.Sprintf(`<tr><td colspan="2" style="padding:6px 0;color:#059669;">Discount%s</td><td style="padding:6px 0;text-align:right;color:#059669;font-weight:600;white-space:nowrap;">-%s</td></tr>`, code, formatINR(order.DiscountAmount))
	}

	return fmt.Sprintf(`
		<div style="background:#f8fafc;border:1px solid #eef2f7;border-radius:12px;padding:20px 22px;margin:24px 0;">
			<table width="100%%" style="border-collapse:collapse;font-size:14px;">
				<tr>
					<td style="padding:4px 0;color:#6b7280;">Order number</td>
					<td style="padding:4px 0;text-align:right;color:#111827;font-weight:700;">%s</td>
				</tr>
				<tr>
					<td style="padding:4px 0;color:#6b7280;">Order date</td>
					<td style="padding:4px 0;text-align:right;color:#111827;">%s</td>
				</tr>
				<tr>
					<td style="padding:4px 0;color:#6b7280;">Status</td>
					<td style="padding:4px 0;text-align:right;"><span style="display:inline-block;background:#eef2ff;color:#4f46e5;font-weight:600;font-size:12px;padding:3px 10px;border-radius:999px;">%s</span></td>
				</tr>
			</table>

			<table width="100%%" style="border-collapse:collapse;font-size:14px;margin-top:16px;">
				<thead>
					<tr>
						<th align="left" style="padding-bottom:8px;font-size:11px;text-transform:uppercase;letter-spacing:.5px;color:#9ca3af;">Service</th>
						<th align="center" style="padding-bottom:8px;font-size:11px;text-transform:uppercase;letter-spacing:.5px;color:#9ca3af;">Qty</th>
						<th align="right" style="padding-bottom:8px;font-size:11px;text-transform:uppercase;letter-spacing:.5px;color:#9ca3af;">Price</th>
					</tr>
				</thead>
				<tbody>%s</tbody>
			</table>

			<table width="100%%" style="border-collapse:collapse;font-size:14px;margin-top:14px;">
				%s
				<tr><td colspan="2" style="padding:6px 0;color:#374151;">Total amount</td><td style="padding:6px 0;text-align:right;color:#111827;font-weight:700;white-space:nowrap;">%s</td></tr>
				<tr><td colspan="2" style="padding:6px 0;color:#374151;">Booking amount</td><td style="padding:6px 0;text-align:right;color:#111827;font-weight:600;white-space:nowrap;">%s</td></tr>
				<tr><td colspan="2" style="padding:6px 0;color:#374151;">Remaining amount</td><td style="padding:6px 0;text-align:right;color:#111827;font-weight:600;white-space:nowrap;">%s</td></tr>
			</table>
		</div>`,
		html.EscapeString(order.OrderNumber),
		order.CreatedAt.Format("2 Jan 2006, 3:04 PM"),
		orderStatusLabel(order.Status),
		rows.String(),
		discountRow,
		formatINR(order.TotalAmount),
		formatINR(order.BookingAmount),
		formatINR(order.RemainingAmount),
	)
}

// orderEmailShell wraps intro copy + the order details block in the branded layout.
func (s *MailService) orderEmailShell(title, firstName, heading, intro string, order *models.Order) string {
	greetingName := strings.TrimSpace(firstName)
	if greetingName == "" {
		greetingName = "there"
	}

	ctaHTML := ""
	if frontend := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/"); frontend != "" {
		ctaHTML = fmt.Sprintf(`
			<div style="text-align:center;margin:8px 0 4px;">
				<a href="%s/orders" style="display:inline-block;background:linear-gradient(135deg,#4f46e5 0%%,#7c3aed 100%%);color:#ffffff;text-decoration:none;font-weight:600;padding:12px 28px;border-radius:12px;">View your orders</a>
			</div>`, frontend)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>%s</title></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;color:#374151;">
	<div style="max-width:600px;margin:0 auto;padding:24px 12px;">
		<div style="background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 10px 25px -5px rgba(0,0,0,0.08);">
			<div style="background:linear-gradient(135deg,#4f46e5 0%%,#7c3aed 100%%);padding:32px 30px;text-align:center;">
				<div style="font-size:32px;font-weight:900;color:#ffffff;letter-spacing:-1px;">ProtAcc</div>
				<div style="color:rgba(255,255,255,0.9);font-size:15px;font-weight:500;margin-top:4px;">%s</div>
			</div>
			<div style="padding:32px 30px;">
				<p style="font-size:16px;margin:0 0 12px;">Hi %s,</p>
				<p style="font-size:16px;line-height:1.6;margin:0 0 8px;color:#4b5563;">%s</p>
				%s
				%s
				<p style="font-size:14px;line-height:1.6;color:#6b7280;margin:20px 0 0;">If you have any questions about this order, just reply to this email and our Chartered Accountants will help you out.</p>
			</div>
			<div style="background:#f8fafc;padding:20px 30px;text-align:center;border-top:1px solid #eef2f7;">
				<p style="font-size:13px;color:#9ca3af;margin:0;">ProtAcc — your Chartered Accountants for GST, ITR &amp; company compliance.</p>
			</div>
		</div>
	</div>
</body>
</html>`, title, heading, html.EscapeString(greetingName), intro, s.orderDetailsHTML(order), ctaHTML)
}

// SendOrderPlacedEmail notifies the customer that their order was placed, with full details.
func (s *MailService) SendOrderPlacedEmail(order *models.Order) error {
	if order == nil || order.User == nil || strings.TrimSpace(order.User.Email) == "" {
		return fmt.Errorf("order placed email: missing recipient")
	}
	subject := fmt.Sprintf("Order confirmed: %s", order.OrderNumber)
	intro := "Thank you for your order! We've received it and our team will begin working on it shortly. Here are your order details:"
	html := s.orderEmailShell("Order confirmed", order.User.FirstName, "Order confirmed 🎉", intro, order)
	return s.sendEmail(order.User.Email, subject, html)
}

// SendOrderPlacedAdminEmail notifies the ProtAcc team (FROM_EMAIL) that a new
// order was booked. Sender and recipient are both FROM_EMAIL by design.
func (s *MailService) SendOrderPlacedAdminEmail(order *models.Order) error {
	if order == nil {
		return fmt.Errorf("order admin email: nil order")
	}
	if strings.TrimSpace(s.fromEmail) == "" {
		return fmt.Errorf("order admin email: FROM_EMAIL not set")
	}

	customerName, customerEmail, customerPhone := "—", "—", "—"
	if order.User != nil {
		if n := strings.TrimSpace(order.User.FirstName + " " + order.User.LastName); n != "" {
			customerName = html.EscapeString(n)
		}
		if order.User.Email != "" {
			customerEmail = html.EscapeString(order.User.Email)
		}
		if order.User.Phone != "" {
			customerPhone = html.EscapeString(order.User.Phone)
		}
	}

	customerBlock := fmt.Sprintf(`
		<div style="background:#f8fafc;border:1px solid #eef2f7;border-radius:12px;padding:20px 22px;margin:24px 0;">
			<table width="100%%" style="border-collapse:collapse;font-size:14px;">
				<tr><td style="padding:4px 0;color:#6b7280;">Customer</td><td style="padding:4px 0;text-align:right;color:#111827;font-weight:600;">%s</td></tr>
				<tr><td style="padding:4px 0;color:#6b7280;">Email</td><td style="padding:4px 0;text-align:right;color:#111827;">%s</td></tr>
				<tr><td style="padding:4px 0;color:#6b7280;">Phone</td><td style="padding:4px 0;text-align:right;color:#111827;">%s</td></tr>
			</table>
		</div>`, customerName, customerEmail, customerPhone)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>New order booked</title></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;color:#374151;">
	<div style="max-width:600px;margin:0 auto;padding:24px 12px;">
		<div style="background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 10px 25px -5px rgba(0,0,0,0.08);">
			<div style="background:linear-gradient(135deg,#4f46e5 0%%,#7c3aed 100%%);padding:32px 30px;text-align:center;">
				<div style="font-size:32px;font-weight:900;color:#ffffff;letter-spacing:-1px;">ProtAcc</div>
				<div style="color:rgba(255,255,255,0.9);font-size:15px;font-weight:500;margin-top:4px;">New order booked</div>
			</div>
			<div style="padding:32px 30px;">
				<p style="font-size:16px;line-height:1.6;margin:0 0 8px;color:#4b5563;">A new order has been booked and the booking payment is confirmed. Customer and order details are below.</p>
				%s
				%s
			</div>
			<div style="background:#f8fafc;padding:20px 30px;text-align:center;border-top:1px solid #eef2f7;">
				<p style="font-size:13px;color:#9ca3af;margin:0;">ProtAcc — internal order notification.</p>
			</div>
		</div>
	</div>
</body>
</html>`, customerBlock, s.orderDetailsHTML(order))

	subject := fmt.Sprintf("New order booked: %s", order.OrderNumber)
	return s.sendEmail(s.fromEmail, subject, htmlBody)
}

// SendOrderCompletedEmail notifies the customer that their order is complete, with full details.
func (s *MailService) SendOrderCompletedEmail(order *models.Order) error {
	if order == nil || order.User == nil || strings.TrimSpace(order.User.Email) == "" {
		return fmt.Errorf("order completed email: missing recipient")
	}
	subject := fmt.Sprintf("Your order is complete: %s", order.OrderNumber)
	intro := "Great news — your order is now complete! Thank you for trusting ProtAcc. Here's a summary of the completed order:"
	html := s.orderEmailShell("Order complete", order.User.FirstName, "Your order is complete ✅", intro, order)
	return s.sendEmail(order.User.Email, subject, html)
}

// ===== Lead / Contact admin notifications (sent to FROM_EMAIL) =====

// adminRow renders one label/value row; blank values are shown as "—".
func adminRow(label, value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		v = "—"
	}
	return fmt.Sprintf(`<tr><td style="padding:6px 0;color:#6b7280;vertical-align:top;white-space:nowrap;padding-right:16px;">%s</td><td style="padding:6px 0;color:#111827;">%s</td></tr>`,
		html.EscapeString(label), html.EscapeString(v))
}

// sendAdminNotification wraps a rows table in the branded layout and sends it to FROM_EMAIL.
func (s *MailService) sendAdminNotification(subject, heading, intro, rowsHTML string) error {
	if strings.TrimSpace(s.fromEmail) == "" {
		return fmt.Errorf("admin notification: FROM_EMAIL not set")
	}
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>%s</title></head>
<body style="margin:0;padding:0;background-color:#f8fafc;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;color:#374151;">
	<div style="max-width:600px;margin:0 auto;padding:24px 12px;">
		<div style="background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 10px 25px -5px rgba(0,0,0,0.08);">
			<div style="background:linear-gradient(135deg,#4f46e5 0%%,#7c3aed 100%%);padding:32px 30px;text-align:center;">
				<div style="font-size:32px;font-weight:900;color:#ffffff;letter-spacing:-1px;">ProtAcc</div>
				<div style="color:rgba(255,255,255,0.9);font-size:15px;font-weight:500;margin-top:4px;">%s</div>
			</div>
			<div style="padding:32px 30px;">
				<p style="font-size:16px;line-height:1.6;margin:0 0 8px;color:#4b5563;">%s</p>
				<div style="background:#f8fafc;border:1px solid #eef2f7;border-radius:12px;padding:20px 22px;margin:20px 0;">
					<table width="100%%" style="border-collapse:collapse;font-size:14px;">%s</table>
				</div>
			</div>
			<div style="background:#f8fafc;padding:20px 30px;text-align:center;border-top:1px solid #eef2f7;">
				<p style="font-size:13px;color:#9ca3af;margin:0;">ProtAcc — internal notification.</p>
			</div>
		</div>
	</div>
</body>
</html>`, subject, heading, intro, rowsHTML)
	return s.sendEmail(s.fromEmail, subject, htmlBody)
}

// SendLeadNotificationEmail notifies the team of a new consultancy enquiry.
func (s *MailService) SendLeadNotificationEmail(lead *models.BusinessLead) error {
	if lead == nil {
		return fmt.Errorf("lead notification: nil lead")
	}
	name := strings.TrimSpace(lead.FirstName + " " + lead.LastName)
	company, businessType, budget, message := "", "", "", ""
	if lead.CompanyName != nil {
		company = *lead.CompanyName
	}
	if lead.BusinessType != nil {
		businessType = *lead.BusinessType
	}
	if lead.BudgetRange != nil {
		budget = *lead.BudgetRange
	}
	if lead.Message != nil {
		message = *lead.Message
	}
	var rows strings.Builder
	rows.WriteString(adminRow("Name", name))
	rows.WriteString(adminRow("Email", lead.Email))
	rows.WriteString(adminRow("Phone", lead.Phone))
	rows.WriteString(adminRow("Company", company))
	rows.WriteString(adminRow("Business type", businessType))
	rows.WriteString(adminRow("Services interested", strings.Join(lead.ServicesInterested, ", ")))
	rows.WriteString(adminRow("Budget range", budget))
	rows.WriteString(adminRow("Preferred contact", string(lead.PreferredContactMethod)))
	rows.WriteString(adminRow("Message", message))

	subject := fmt.Sprintf("New consultancy enquiry: %s", name)
	return s.sendAdminNotification(subject, "New consultancy enquiry", "A new consultancy request has been submitted:", rows.String())
}

// SendContactNotificationEmail notifies the team of a new contact-form message.
func (s *MailService) SendContactNotificationEmail(msg *models.ContactMessage) error {
	if msg == nil {
		return fmt.Errorf("contact notification: nil message")
	}
	company, serviceInterest := "", ""
	if msg.Company != nil {
		company = *msg.Company
	}
	if msg.ServiceInterest != nil {
		serviceInterest = *msg.ServiceInterest
	}
	var rows strings.Builder
	rows.WriteString(adminRow("Name", msg.Name))
	rows.WriteString(adminRow("Email", msg.Email))
	rows.WriteString(adminRow("Phone", msg.Phone))
	rows.WriteString(adminRow("Company", company))
	rows.WriteString(adminRow("Subject", msg.Subject))
	rows.WriteString(adminRow("Service interest", serviceInterest))
	rows.WriteString(adminRow("Preferred contact", msg.PreferredContact))
	rows.WriteString(adminRow("Message", msg.Message))

	subject := fmt.Sprintf("New contact message: %s", msg.Subject)
	return s.sendAdminNotification(subject, "New contact message", "A new message has been submitted through the contact form:", rows.String())
}
