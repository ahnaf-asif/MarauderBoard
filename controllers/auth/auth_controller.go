package auth

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func RegisterAuthRoutes(auth fiber.Router) {
	auth.Get("/login", func(ctx *fiber.Ctx) error {
		return ctx.Render("auth/login", fiber.Map{})
	})

	auth.Post("/login", func(ctx *fiber.Ctx) error {
		email := ctx.FormValue("email")
		password := ctx.FormValue("password")

		if email == "" || password == "" {
			return ctx.Render("partials/error-message", "all fields are required")
		}

		user, err := models.GetUserByEmail(database.DB, email)
		if err != nil {
			log.Println("Error getting user by email:", err)
			return ctx.Render("partials/error-message", "Invalid email or password")
		}

		if *user.Provider == "google" {
			return ctx.Render("partials/error-message", "This account is registered with Google. Please login with Google.")
		}

		if models.CheckPasswordHash(password, *user.Password) {
			log.Println("Error checking password hash:", err)
			return ctx.Render("partials/error-message", "Invalid email or password")
		}

		user_data := UserSessionData{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Provider:  *user.Provider,
			Avatar:    *user.Avatar,
		}

		user_data_json, err := json.Marshal(user_data)
		if err != nil {
			log.Println("Error marshalling user data:", err)
			return ctx.Render("partials/error-message", "Invalid email or password")
		}

		if err := goth_fiber.StoreInSession(
			"user_data", string(user_data_json), ctx,
		); err != nil {
			log.Println("Error storing user session:", err)
			return ctx.Render("partials/error-message", "Invalid email or password")
		} else {
			log.Println("User data stored in session:", string(user_data_json))
		}

		redirect_uri := ctx.Cookies("redirect_uri")
		if redirect_uri == "" {
			redirect_uri = "/"
		}

		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Login successful! Redirecting...",
			"Redirect": redirect_uri,
		})
	})

	auth.Get("/register", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/", fiber.StatusFound)
	})

	auth.Post("/register", func(ctx *fiber.Ctx) error {
		email := ctx.FormValue("email")
		password := ctx.FormValue("password")
		first_name := ctx.FormValue("first_name")
		last_name := ctx.FormValue("last_name")
		confirm_password := ctx.FormValue("confirm_password")

		if email == "" ||
			password == "" ||
			first_name == "" ||
			last_name == "" ||
			confirm_password == "" {
			return ctx.Render("partials/error-message", "All fields are required")
		}

		if password != confirm_password {
			return ctx.Render("partials/error-message", "Passwords do not match")
		}

		hashed_password, err := models.GeneratePasswordHash(password)
		if err != nil {
			log.Println("Error hashing password:", err)
			return ctx.Render(
				"partials/error-message", "Error creating account. Please try again.",
			)
		}

		avatar_url := fmt.Sprintf(
			"https://ui-avatars.com/api/?name=%s+%s&background=random",
			first_name, last_name,
		)

		provider := "manual"

		user := models.User{
			Email:     email,
			Password:  &hashed_password,
			FirstName: first_name,
			LastName:  last_name,
			Provider:  &provider,
			Avatar:    &avatar_url,
		}

		db_user, err := models.AddNewUser(database.DB, user)
		if err != nil {
			log.Println("Error creating user:", err)
			return ctx.Render("partials/error-message", "Error creating account. Please try again")
		}

		user_data := UserSessionData{
			ID:        db_user.ID,
			Email:     db_user.Email,
			FirstName: db_user.FirstName,
			LastName:  db_user.LastName,
			Provider:  *db_user.Provider,
			Avatar:    *db_user.Avatar,
		}

		user_data_json, err := json.Marshal(user_data)
		if err != nil {
			log.Println("Error marshalling user data:", err)
			return ctx.Render(
				"partials/error-message", "Error creating account. Please try again",
			)
		}

		if err := goth_fiber.StoreInSession("user_data", string(user_data_json), ctx); err != nil {
			log.Println("Error storing user session:", err)
			return ctx.Render(
				"partials/error-message", "Error creating account. Please try again",
			)
		} else {
			log.Println("User data stored in session(manual):", string(user_data_json))
		}

		redirect_uri := ctx.Cookies("redirect_uri", "/")
		if redirect_uri == "" {
			redirect_uri = "/"
		}

		return ctx.Render("partials/success-message-with-redirect", fiber.Map{
			"Message":  "Account created successfully! Redirecting...",
			"Redirect": redirect_uri,
		})
	})

	auth.Get(":provider", func(ctx *fiber.Ctx) error {
		return goth_fiber.BeginAuthHandler(ctx)
	})

	auth.Get(":provider/callback", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Println("Error completing user auth:", err)
		}

		found_user, err := models.GetUserByEmail(database.DB, user.Email)
		if err != nil {
			log.Println("User not found, creating new user:", err)
			new_user := models.User{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Provider:  &user.Provider,
				Email:     user.Email,
				Avatar:    &user.AvatarURL,
			}

			db_user, err := models.AddNewUser(database.DB, new_user)
			if err != nil {
				log.Fatal("AddNewUser Error happened here:", err)
			}
			found_user = db_user
		}

		user_data := UserSessionData{
			ID:        found_user.ID,
			Email:     found_user.Email,
			FirstName: found_user.FirstName,
			LastName:  found_user.LastName,
			Provider:  *found_user.Provider,
			Avatar:    user.AvatarURL,
		}

		user_data_json, err := json.Marshal(user_data)
		if err != nil {
			log.Fatal(err)
		}

		if err := goth_fiber.StoreInSession("user_data", string(user_data_json), ctx); err != nil {
			log.Fatal(err)
		} else {
			log.Println("User data stored in session:", string(user_data_json))
		}

		redirect_uri := ctx.Cookies("redirect_uri", "/")
		return ctx.Redirect(redirect_uri, fiber.StatusFound)
	})
}
