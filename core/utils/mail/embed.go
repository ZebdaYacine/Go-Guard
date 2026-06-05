package mail

import (
	"bytes"
	"html/template"
)

type OTPData struct {
	OTP    string
	Expiry int
}

func RenderTemplate(path string, data OTPData) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}
