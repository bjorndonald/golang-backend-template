package service

import (
	"fmt"
	"log"
	"time"

	"github.com/bjorndonald/lasgcce/constants"
	"github.com/bjorndonald/lasgcce/internal/helpers"
	"github.com/bjorndonald/lasgcce/internal/otp"
	"github.com/bjorndonald/lasgcce/resend"
)

type EmailServicer interface {
	SendVerificationEmail(name, email, url string)
	SendForgotPasswordEmail(name, email string)
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

func (s *EmailService) SendForgotPasswordEmail(name, email string) {
	//TODO: send otp email to user
}

// Sends account verification email
func (s *EmailService) SendVerificationEmail(name, email, url string) {

	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	type OTP struct {
		Otp  string
		Name string
		Url  string
	}

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	verificationUrl := fmt.Sprintf("%s/api/v1/auth/verify/%s/%s", url, email, otpToken)

	messageBody, err := helpers.ParseTemplateFile("verify_account.html", OTP{Otp: otpToken, Name: name, Url: verificationUrl})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "Verify your account", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}
