package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"monorepo/libs/errors"
)

const LangKey = "lang"

// LanguageMiddleware extracts language from Accept-Language header
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")

		// Validate and set default
		switch lang {
		case "id", "ID":
			lang = string(errors.LangID)
		default:
			lang = string(errors.LangEN)
		}

		// Set in gin context
		c.Set(LangKey, errors.Lang(lang))

		// Set in request context
		ctx := context.WithValue(c.Request.Context(), LangKey, errors.Lang(lang))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetLangFromContext extracts language from context
func GetLangFromContext(ctx context.Context) errors.Lang {
	if lang, ok := ctx.Value(LangKey).(errors.Lang); ok {
		return lang
	}
	return errors.LangEN
}

// GetLangFromGin extracts language from gin context
func GetLangFromGin(c *gin.Context) errors.Lang {
	if lang, exists := c.Get(LangKey); exists {
		if l, ok := lang.(errors.Lang); ok {
			return l
		}
	}
	return errors.LangEN
}
