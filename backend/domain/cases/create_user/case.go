package create_user

import (
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/helpers"
	"github.com/wigfri/mustore/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Response struct {
	Id string `json:"id"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	const op = "registration_user.Run"
	logger := c.Services().Logger()

	var (
		hashedPassword []byte
		userId         string
		err            error
	)

	if r.Password != "" {
		if err := helpers.ValidatePassword(r.Password); err != nil {
			return nil, err
		}
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
	}

	userId, err = createUser(c, r.Email, r.Username, string(hashedPassword), r.Role)
	if err != nil {
		return nil, err
	}

	logger.Info("successfully registered user by email", "op", op, "user_email", r.Email)
	return &Response{Id: userId}, nil
}

func createUser(c domain.Context, email, name, passwordHash, role string) (string, error) {
	id := uuid.New()
	var user *models.User

	_, err := models.IsValidRole(role)
	if err != nil {
		return "", err
	}

	user = models.NewBaseUser(id, email, name, passwordHash)

	userId, err := c.Connection().User().Insert(user)
	if err != nil {
		return "", err
	}

	scheduleRollbackUser(c, id)
	return userId, nil
}

const unverifiedUserRollbackAfter = 10 * time.Minute

func scheduleRollbackUser(c domain.Context, userId uuid.UUID) {
	time.AfterFunc(unverifiedUserRollbackAfter, func() {
		rollbackUnverifiedUserIfNeeded(c, userId)
	})
}

func rollbackUnverifiedUserIfNeeded(c domain.Context, userId uuid.UUID) {
	logger := c.Services().Logger()

	user, err := c.Connection().User().GetUser(userId)
	if err != nil {
		logger.Error("failed to fetch user in rollback", "error", err.Error())
		return
	}
	if user.IsVerified {
		return
	}
	if err := c.Connection().User().DeleteFromDb(userId); err != nil {
		logger.Error("failed to delete unverified user", "error", err.Error())
		return
	}
	logger.Info("rolled back unverified user", "user_id", userId.String())
}
