package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func VerifyToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Load environment variables from .env file
        err := godotenv.Load()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success":         false,
                "message":         "Failed to load environment variables",
                "httpStatusCode": http.StatusInternalServerError,
                "error":           "INTERNAL_SERVER_ERROR",
            })
            c.Abort()
            return
        }

        // Retrieve the access token secret from environment variables
        accessTokenSecret := os.Getenv("JWT_SECRET")

        authHeader := c.GetHeader("Authorization")
        serviceName := os.Getenv("SERVICE_NAME")

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success":         false,
                "message":         "Authorization header is missing or invalid",
                "httpStatusCode": http.StatusUnauthorized,
                "error":           "VALIDATION_ERROR",
                "service":         serviceName,
            })
            c.Abort()
            return
        }

        token := authHeader[7:] 

        // Verify token
        claims := jwt.MapClaims{}
        jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(accessTokenSecret), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success":         false,
                "message":         "Invalid token",
                "httpStatusCode": http.StatusUnauthorized,
                "error":           "VALIDATION_ERROR",
                "service":         serviceName,
            })
            c.Abort()
            return
        }

        if !jwtToken.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success":         false,
                "message":         "Invalid token",
                "httpStatusCode": http.StatusUnauthorized,
                "error":           "VALIDATION_ERROR",
                "service":         serviceName,
            })
            c.Abort()
            return
        }

        // Set user information in the context
        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success":         false,
                "message":         "Invalid token format",
                "httpStatusCode": http.StatusUnauthorized,
                "error":           "VALIDATION_ERROR",
                "service":         serviceName,
            })
            c.Abort()
            return
        }
        c.Set("userID", userID)

        // Proceed to the next middleware or handler
        c.Next()
    }
}


// func ExtractUserID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
// 			return
// 		}

// 		token := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

// 		// Validate and extract user ID from token (e.g., JWT)
// 		userID, err := utils.ExtractUserIDFromToken(token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			return
// 		}

// 		// Set the user ID in the context for downstream handlers
// 		c.Set("userID", userID)
// 		c.Next()
// 	}
// }


// func VerifyToken() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // Load environment variables from .env file
//         err := godotenv.Load()
//         if err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{
//                 "success":         false,
//                 "message":         "Failed to load environment variables",
//                 "httpStatusCode": http.StatusInternalServerError,
//                 "error":           "INTERNAL_SERVER_ERROR",
//             })
//             c.Abort()
//             return
//         }

//         // Retrieve the access token secret from environment variables
//         accessTokenSecret := os.Getenv("JWT_SECRET")

//         authHeader := c.GetHeader("Authorization")
//         serviceName := os.Getenv("SERVICE_NAME")

//         if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
//             c.JSON(http.StatusUnauthorized, gin.H{
//                 "success":         false,
//                 "message":         "Authorization header is missing or invalid",
//                 "httpStatusCode": http.StatusUnauthorized,
//                 "error":           "VALIDATION_ERROR",
//                 "service":         serviceName,
//             })
//             c.Abort()
//             return
//         }

//         token := authHeader[7:] // Extract JWT token

//         // Verify token
//         claims := jwt.MapClaims{}
//         jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
//             return []byte(accessTokenSecret), nil
//         })

//         if err != nil {
//             c.JSON(http.StatusUnauthorized, gin.H{
//                 "success":         false,
//                 "message":         "Invalid token",
//                 "httpStatusCode": http.StatusUnauthorized,
//                 "error":           "VALIDATION_ERROR",
//                 "service":         serviceName,
//             })
//             c.Abort()
//             return
//         }

//         if !jwtToken.Valid {
//             c.JSON(http.StatusUnauthorized, gin.H{
//                 "success":         false,
//                 "message":         "Invalid token",
//                 "httpStatusCode": http.StatusUnauthorized,
//                 "error":           "VALIDATION_ERROR",
//                 "service":         serviceName,
//             })
//             c.Abort()
//             return
//         }

//         // Set user information in the context
//         c.Set("userID", claims["id"])

//         // Proceed to the next middleware or handler
//         c.Next()
//     }
// }


