package service

import (
	"fmt"
	"log"
	"time"

	"github.com/bjorndonald/golang-backend-template/constants"
	"github.com/bjorndonald/golang-backend-template/internal/helpers"
	"github.com/bjorndonald/golang-backend-template/internal/models"
	"github.com/bjorndonald/golang-backend-template/internal/otp"
	"github.com/bjorndonald/golang-backend-template/resend"
)

type EmailServicer interface {
	SendNewUserEmail(name, email, url string)
	SendForgotPasswordEmail(name, email string)
	SendOTPEmail(name, email string)

	SendNewDeviceEmail(name, email, url string, agent models.UserAgent)
	SendNewLocationEmail(name, email, url string, location models.GeoLocation)
}

type EmailService struct {
	Client *resend.Client
}

var (
	constant = constants.New()
)

func NewEmailService() EmailServicer {
	client := resend.NewClient(constant.ResendApiKey)
	return &EmailService{Client: client}
}

func (s *EmailService) Send(name, email, subject, content string) error {
	to := []string{
		email,
	}

	_, err := s.Client.Send(to, constant.SendFromEmail, constant.SendFromName, subject, content)

	return err
}

func (s *EmailService) SendOTPEmail(name, email string) {
	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	type NewUser struct {
		Email string
		Name  string
		OTP   string
	}

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	messageBody, err := helpers.ParseTemplateFile("otp.html", NewUser{Email: email, Name: name, OTP: otpToken})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "OTP Code", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

func (s *EmailService) SendForgotPasswordEmail(name, email string) {
	c := constants.New()
	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	type NewUser struct {
		Email string
		Name  string
		OTP   string
		Url   string
	}

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	verificationUrl := fmt.Sprintf("%s/auth/reset-password/%s/%s", c.ClientUrl)
	messageBody, err := helpers.ParseTemplateFile("reset_password.html", NewUser{Email: email, Name: name, Url: verificationUrl, OTP: otpToken})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
	err = s.Send(name, email, "Reset password", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// Sends email to new user
func (s *EmailService) SendNewUserEmail(name, email, url string) {
	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	type NewUser struct {
		Email string
		Name  string
		Url   string
	}

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	verificationUrl := fmt.Sprintf("%s/api/v1/auth/verify/%s/%s", url, email, otpToken)

	messageBody, err := helpers.ParseTemplateFile("verify_account.html", NewUser{Email: email, Name: name, Url: verificationUrl})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "Verify your account", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// Sends email to user notifying them of login from new device
func (s *EmailService) SendNewDeviceEmail(name, email, url string, agent models.UserAgent) {
	constant := constants.New()
	// otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)

	// if err != nil {
	// 	log.Printf("Error sending email: %v", err.Error())
	// }

	type NewUser struct {
		Email       string
		Name        string
		Url         string
		Platform    string
		OS          string
		BrowserName string
		Model       string
		Mobile      bool
		Time        string
	}

	// verificationUrl := fmt.Sprintf("%s/api/v1/auth/verify/device/%s/%s", url, email, otpToken)
	forgotPasswordUrl := fmt.Sprintf("%s/auth/forgot-password", constant.ClientUrl)

	messageBody, err := helpers.ParseTemplateFile(
		"verify_device.html",
		NewUser{
			Email:       email,
			Name:        name,
			Url:         forgotPasswordUrl,
			Platform:    agent.Platform,
			OS:          agent.OS,
			BrowserName: agent.BrowserName,
			Model:       agent.Model,
			Mobile:      agent.Mobile,
			Time:        time.Now().String(),
		})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "New Device Detected", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// Sends email to user notifying them of login from new location
func (s *EmailService) SendNewLocationEmail(name, email, url string, location models.GeoLocation) {
	// otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	// if err != nil {
	// 	log.Printf("Error sending email: %v", err.Error())
	// }

	type NewUser struct {
		Email   string
		Name    string
		Url     string
		City    string
		Country string
		Region  string
		Time    string
	}

	// verificationUrl := fmt.Sprintf("%s/api/v1/auth/verify/location/%s/%s", url, email, otpToken)
	forgotPasswordUrl := fmt.Sprintf("%s/auth/forgot-password", constant.ClientUrl)

	messageBody, err := helpers.ParseTemplateFile("verify_location.html",
		NewUser{
			Email:   email,
			Name:    name,
			Url:     forgotPasswordUrl,
			City:    location.City,
			Country: location.Country,
			Region:  location.Region,
			Time:    time.Now().String(),
		})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "New Location Detected", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}
