package controllers

import (
	"encoding/json"
	"magic-server-2026/src/db"
	"magic-server-2026/src/gen"
	"magic-server-2026/src/middlewares"
	"magic-server-2026/src/models"
	"magic-server-2026/src/utils"
	"magic-server-2026/src/validators"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// --------------------
// LOGIN
// --------------------
func Login(c fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	body := c.Body() // returns []byte
	if err := json.Unmarshal(body, &input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	collection := db.Client.Database("magic899").Collection("users")
	var user models.User
	err := collection.FindOne(c.Context(), bson.M{"email": input.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !validators.CheckPasswordHash(input.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	// Generate RSA keys if missing
	if user.RSAPrivate == "" || user.RSAPublic == "" {
		priv, pub, err := gen.GenerateRSAKeys()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to generate RSA keys")
		}
		user.RSAPrivate = priv
		user.RSAPublic = pub
		_, _ = collection.UpdateByID(c.Context(), user.ID, bson.M{
			"$set": bson.M{"rsa_private": priv, "rsa_public": pub},
		})
	}

	// Generate AES-GCM session ID
	sessionID, err := gen.GenerateSecureSessionID()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate session ID")
	}

	// Generate tokens
	accessToken, _ := gen.GenerateAccessToken(user.ID.Hex(), sessionID, user.Role)
	refreshToken, _ := gen.GenerateRefreshToken(user.ID.Hex(), sessionID)
	csrfToken, _ := gen.GenerateBindedCSRF(sessionID)
	rspToken, _ := gen.GenerateBindedRSP(sessionID)

	// Encrypt access token with user RSA public key
	encryptedAccess, _ := utils.EncryptWithPublicKey([]byte(accessToken), user.RSAPublic)

	// Update session in DB
	update := bson.M{
		"$set": bson.M{
			"session_id":    sessionID,
			"last_login_at": time.Now(),
		},
	}
	_, err = collection.UpdateByID(c.Context(), user.ID, update)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update session")
	}

	// Set cookies
	c.Cookie(&fiber.Cookie{Name: "session_id", Value: sessionID, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "access_token", Value: string(encryptedAccess), HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: refreshToken, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "bind_csrf", Value: csrfToken, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/csrf", MaxAge: 5})
	c.Cookie(&fiber.Cookie{Name: "bind_rsp", Value: rspToken, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/rsp", MaxAge: 5})

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user": fiber.Map{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// --------------------
// LOGOUT
// --------------------
func Logout(c fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		// Revoke bind tokens
		middlewares.TokenRevoker.RevokeCSRF(c.Cookies("bind_csrf"))
		middlewares.TokenRevoker.RevokeRSP(c.Cookies("bind_rsp"))
	}

	// Remove session_id from DB
	collection := db.Client.Database("magic899").Collection("users")
	collection.UpdateOne(c.Context(), bson.M{"session_id": sessionID}, bson.M{"$unset": bson.M{"session_id": ""}})

	// Clear cookies
	c.ClearCookie("session_id")
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	c.ClearCookie("bind_csrf")
	c.ClearCookie("bind_rsp")

	return c.JSON(fiber.Map{"message": "logout successful"})
}

// --------------------
// REFRESH ACCESS TOKEN
// --------------------
func RefreshAccessToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	sessionID := c.Cookies("session_id")
	if refreshToken == "" || sessionID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing tokens")
	}

	claims, err := middlewares.ValidateRefreshToken(refreshToken)
	if err != nil || claims.SessionID != sessionID {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	collection := db.Client.Database("magic899").Collection("users")
	var user models.User
	err = collection.FindOne(c.Context(), bson.M{"session_id": sessionID}).Decode(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid session")
	}

	// Generate new tokens
	newAccessToken, _ := gen.GenerateAccessToken(user.ID.Hex(), sessionID, user.Role)
	newCSRF, _ := gen.GenerateBindedCSRF(sessionID)
	newRSP, _ := gen.GenerateBindedRSP(sessionID)

	// Encrypt access token with RSA
	encryptedAccess, _ := utils.EncryptWithPublicKey([]byte(newAccessToken), user.RSAPublic)

	// Revoke old bind tokens
	middlewares.TokenRevoker.RevokeCSRF(c.Cookies("bind_csrf"))
	middlewares.TokenRevoker.RevokeRSP(c.Cookies("bind_rsp"))

	// Set new cookies
	c.Cookie(&fiber.Cookie{Name: "access_token", Value: string(encryptedAccess), HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/"})
	c.Cookie(&fiber.Cookie{Name: "bind_csrf", Value: newCSRF, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/csrf", MaxAge: 5})
	c.Cookie(&fiber.Cookie{Name: "bind_rsp", Value: newRSP, HTTPOnly: true, Secure: true, SameSite: "Lax", Path: "/rsp", MaxAge: 5})

	return c.JSON(fiber.Map{"access_token": string(encryptedAccess)})
}

// --------------------
// GET USER PROFILE (protected)
// --------------------
func GetProfile(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	collection := db.Client.Database("magic899").Collection("users")
	var user models.User
	err := collection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(fiber.Map{
		"id":         user.ID.Hex(),
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"avatar":     user.Avatar,
		"isVerified": user.IsVerified,
	})
}
