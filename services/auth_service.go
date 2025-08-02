package services


import (
	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
)

type AuthService interface {
	Register(req *models.UserRegistrationRequest) (*models.UserResponse, error)
	Login(req *models.UserLoginRequest) (string, *models.UserResponse, error)
	GetUserProfile(userID uint) (*models.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

//func NewAuthService(userRepo repository.UserRepository) AuthService {
//	return &authService{
//		userRepo: userRepo,
//	}
//}

//func (s *authService) Register(req *models.UserRegistrationRequest) (*models.UserResponse, error) {
//	// Check if user already exists
//	_, err := s.userRepo.GetByEmail(req.Email)
//	if err == nil {
//		return nil, errors.New("user with this email already exists")
//	}
//
//	// Hash password
//	hashedPassword, err := utils.HashPassword(req.Password)
//	if err != nil {
//		return nil, err
//	}
//
//	// Create user
//	user := &models.User{
//		Email:     req.Email,
//		Password:  hashedPassword,
//		FirstName: req.FirstName,
//		LastName:  req.LastName,
//	}
//
//	err = s.userRepo.Create(user)
//	if err != nil {
//		return nil, err
//	}
//
//	return &models.UserResponse{
//		ID:        user.ID,
//		Email:     user.Email,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		CreatedAt: user.CreatedAt,
//	}, nil
//}

//func (s *authService) Login(req *models.UserLoginRequest) (string, *models.UserResponse, error) {
//	user, err := s.userRepo.GetByEmail(req.Email)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return "", nil, errors.New("invalid credentials")
//		}
//		return "", nil, err
//	}
//
//	if !utils.CheckPassword(req.Password, user.Password) {
//		return "", nil, errors.New("invalid credentials")
//	}
//
//	token, err := utils.GenerateToken(user.ID, user.Email)
//	if err != nil {
//		return "", nil, err
//	}
//
//	userResponse := &models.UserResponse{
//		ID:        user.ID,
//		Email:     user.Email,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		CreatedAt: user.CreatedAt,
//	}
//
//	return token, userResponse, nil
//}

func (s *authService) GetUserProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}, nil
}
