package private

import (
	"fmt"
	"go-gaurd/api/security"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
)

func (p *ProfileController) UploadFile(c *fiber.Ctx) (error, string) {
	file, err := c.FormFile("picture")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "picture is required",
		}), ""
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		}), ""
	}
	defer src.Close()

	objectName := fmt.Sprintf("profile/%d-%s", time.Now().UnixNano(), file.Filename)
	fmt.Println(p.MinioDB.Client)
	_, err = p.MinioDB.Client.PutObject(
		c.Context(),
		"profile",
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		}), ""
	}

	url := fmt.Sprintf("%s/profile/%s", p.MinioDB.Client.EndpointURL().String(), objectName)
	return nil, url
}

func (ac *ProfileController) ValidateBody(c *fiber.Ctx, req interface{}) error {
	log.Println("Validating request body")

	if err := c.BodyParser(req); err != nil {
		log.Printf("Body parsing failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"success": false,
		})
	}

	if err := ac.validate.Struct(req); err != nil {
		log.Printf("Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	log.Println("Request body validation successful")
	return nil
}

func (ac *ProfileController) GetClaimFromToken(c *fiber.Ctx) (string, string, string, *jwt.NumericDate, error) {
	authHeader := c.Get("Authorization")
	token, err := security.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("Failed to extract token: %v", err)
		return "", "", "", nil, err
	}
	userId, role, jti, expiresAt, err := security.ValidateAccessToken(token)
	if err != nil {
		log.Printf("Failed to validate access token: %v", err)
		return "", "", "", nil, err
	}
	return userId, role, jti, expiresAt, nil
}
