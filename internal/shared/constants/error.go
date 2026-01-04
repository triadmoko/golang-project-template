package constants

import "errors"

type ErrCode int
type Lang string

const (
	LangEN Lang = "en"
	LangID Lang = "id"
)

const (
	// Global errors
	SomethingWentWrong ErrCode = iota
	InvalidInput
	ValidationFailed
	Unauthorized

	// Auth errors
	InvalidCredentials
	UserAlreadyExists
	UsernameAlreadyTaken
	FailedToHashPassword
	FailedToCreateUser
	FailedToGenerateToken

	// User errors
	UserNotFound
	FailedToUpdateUser
	FailedToGetUsers
)

var errMessages = map[ErrCode]map[Lang]string{
	// Global errors
	SomethingWentWrong: {
		LangEN: "something went wrong",
		LangID: "terjadi kesalahan",
	},
	InvalidInput: {
		LangEN: "invalid input",
		LangID: "input tidak valid",
	},
	ValidationFailed: {
		LangEN: "validation failed",
		LangID: "validasi gagal",
	},
	Unauthorized: {
		LangEN: "unauthorized",
		LangID: "tidak memiliki akses",
	},

	// Auth errors
	InvalidCredentials: {
		LangEN: "invalid email or password",
		LangID: "email atau password salah",
	},
	UserAlreadyExists: {
		LangEN: "user already exists",
		LangID: "pengguna sudah terdaftar",
	},
	UsernameAlreadyTaken: {
		LangEN: "username already taken",
		LangID: "username sudah digunakan",
	},
	FailedToHashPassword: {
		LangEN: "failed to process password",
		LangID: "gagal memproses password",
	},
	FailedToCreateUser: {
		LangEN: "failed to create user",
		LangID: "gagal membuat pengguna",
	},
	FailedToGenerateToken: {
		LangEN: "failed to generate token",
		LangID: "gagal membuat token",
	},

	// User errors
	UserNotFound: {
		LangEN: "user not found",
		LangID: "pengguna tidak ditemukan",
	},
	FailedToUpdateUser: {
		LangEN: "failed to update user",
		LangID: "gagal memperbarui pengguna",
	},
	FailedToGetUsers: {
		LangEN: "failed to get users",
		LangID: "gagal mengambil data pengguna",
	},
}

// GetError returns error message based on code and language
func GetError(code ErrCode, lang Lang) error {
	if msgs, ok := errMessages[code]; ok {
		if msg, ok := msgs[lang]; ok {
			return errors.New(msg)
		}
		// Fallback to English
		if msg, ok := msgs[LangEN]; ok {
			return errors.New(msg)
		}
	}
	return errors.New("unknown error")
}

// GetErrorMessage returns error message string based on code and language
func GetErrorMessage(code ErrCode, lang Lang) string {
	if msgs, ok := errMessages[code]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
		// Fallback to English
		if msg, ok := msgs[LangEN]; ok {
			return msg
		}
	}
	return "unknown error"
}
