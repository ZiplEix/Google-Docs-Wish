package auth

// var (
// 	googleOauthConfig *oauth2.Config
// )

// func InitOauth() {
// 	googleOauthConfig = &oauth2.Config{
// 		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
// 		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
// 		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
// 		Scopes:       []string{"email", "profile"},
// 		Endpoint:     google.Endpoint,
// 	}
// }

// func googleLogin(c *fiber.Ctx) error {
// 	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	return c.Redirect(url)
// }

// func googleCallback(c *fiber.Ctx) error {
// 	code := c.Query("code")

// 	token, err := googleOauthConfig.Exchange(c.Context(), code)
// 	if err != nil {
// 		return err
// 	}

// 	idToken, ok := token.Extra("id_token").(string)
// 	if !ok {
// 		return fmt.Errorf("no id_token")
// 	}

// 	tokenInfo, err := database.FirebaseAuthClient.VerifyIDToken(c.Context(), idToken)
// 	if err != nil {
// 		return fmt.Errorf("error verifying ID token: %v", err)
// 	}

// 	// Créer un nouvel utilisateur Firebase s'il n'existe pas encore
// 	userRecord, err := database.FirebaseAuthClient.GetUser(c.Context(), tokenInfo.UID)
// 	if err != nil {
// 		if auth.IsUserNotFound(err) {
// 			userRecord, err = database.FirebaseAuthClient.CreateUser(c.Context(), &auth.UserToCreate{})
// 			if err != nil {
// 				return fmt.Errorf("error creating user: %v", err)
// 			}
// 		}
// 	} else {
// 		fmt.Printf("Successfully fetched user data: %v\n", userRecord)
// 	}

// 	fmt.Printf("User ID : %s, Email : %s\n", tokenInfo.UID, tokenInfo.Claims["email"])

// 	return c.SendStatus(fiber.StatusOK)
// }
