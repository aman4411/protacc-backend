package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
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
							<strong>Need help?</strong> Contact our support team at <a href="mailto:support@protacc.in" style="color: #4f46e5;">support@protacc.in</a>
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
							📧 <a href="mailto:support@protacc.in" style="color: #4f46e5;">support@protacc.in</a> | 
							📞 <a href="tel:+919876543210" style="color: #4f46e5;">+91 98765 43210</a>
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
					<a href="mailto:support@protacc.in" style="color: #4f46e5;">support@protacc.in</a>
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
