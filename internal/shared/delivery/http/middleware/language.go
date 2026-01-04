package middleware

import (
	"app/internal/shared/constants"
	"context"

	"github.com/gin-gonic/gin"
)

const LangKey = "lang"

// LanguageMiddleware extracts language from Accept-Language header
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")

		// Validate and set default
		switch lang {
		case "id", "ID":
			lang = string(constants.LangID)
		default:
			lang = string(constants.LangEN)
		}

		// Set in gin context
		c.Set(LangKey, constants.Lang(lang))

		// Set in request context
		ctx := context.WithValue(c.Request.Context(), LangKey, constants.Lang(lang))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetLangFromContext extracts language from context
func GetLangFromContext(ctx context.Context) constants.Lang {
	if lang, ok := ctx.Value(LangKey).(constants.Lang); ok {
		return lang
	}
	return constants.LangEN
}

// GetLangFromGin extracts language from gin context
func GetLangFromGin(c *gin.Context) constants.Lang {
	if lang, exists := c.Get(LangKey); exists {
		if l, ok := lang.(constants.Lang); ok {
			return l
		}
	}
	return constants.LangEN
}
