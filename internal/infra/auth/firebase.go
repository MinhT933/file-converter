package auth

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
)

type FirebaseProvider struct {
	client *auth.Client
}

func NewFirebaseProvider(client *auth.Client) *FirebaseProvider {
	//& có nghĩa là chúng ta đang nhận một con trỏ đến firebase provider
	// Điều này giúp tiết kiệm bộ nhớ và cho phép chúng ta thay đổi giá trị
	return &FirebaseProvider{
		client: client,
	}
}

// (p *FirebaseProvider) là một phương thức của FirebaseProvider
// Nó nhận một context và một token, sau đó xác thực token và trả về thông tin
func (p *FirebaseProvider) VerifyToken(ctx context.Context, token string) (*User, error) {
	//ctx có nghĩa là context, là một đối tượng dùng để quản lý thời gian sống của một request
	// Nó giúp chúng ta hủy bỏ request nếu cần thiết

	firbaseToken, err := p.client.VerifyIDToken(ctx, token)
	if err != nil {
		log.Printf("❌ Lỗi xác thực token: %v", err)
		return nil, err
	}

	return &User{
		UID:         firbaseToken.UID,
		Email:       getStringClaim(firbaseToken.Claims, "email"),
		DisplayName: getStringClaim(firbaseToken.Claims, "name"),
	}, nil
}

func (p *FirebaseProvider) Login(ctx context.Context, email, password string) (*User, string, error) {
	return nil, "", NewAuthError(ErrCodeInternalError, "Email/password login not supported by Firebase Admin SDK")
}

func (p *FirebaseProvider) GoogleLogin(ctx context.Context, accessToken, idToken string) (*User, string, error) {
	if idToken == "" {
		return nil, "", NewAuthError(ErrCodeInternalError, "Google ID token is required for Google login")
	}

	user, err := p.VerifyToken(ctx, idToken)
	if err != nil {
		log.Printf("💀 Lỗi xác thực Google ID token: %v", err)
		return nil, "", err
	}

	return user, idToken, nil
}

func (p *FirebaseProvider) FacebookLogin(ctx context.Context, accessToken string) (*User, string, error) {
	if accessToken == "" {
		return nil, "", NewAuthError(ErrCodeInternalError, "Facebook access token is required for Facebook login")
	}

	// Firebase Admin SDK does not support Facebook login directly
	return nil, "", NewAuthError(ErrCodeInternalError, "Facebook login not supported by Firebase Admin SDK")
}

func (p *FirebaseProvider) GetUser(ctx context.Context, uid string) (*User, error) {
	user, err := p.client.GetUser(ctx, uid)
	if err != nil {
		log.Printf("❌ Lỗi lấy thông tin user: %v", err)
		return nil, err
	}

	return &User{
		UID:           user.UID,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		PhotoURL:      user.PhotoURL,
		EmailVerified: user.EmailVerified,
		PhoneNumber:   user.PhoneNumber,
	}, nil
}

func (p *FirebaseProvider) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// Chuyển đổi CreateUserRequest thành UserToCreate
	// Firebase admin SDK yêu cầu một struct cụ thể để tạo user
	// Chúng ta sẽ sử dụng UserToCreate để tạo user mới
	// Nếu bạn muốn sử dụng UserToCreate, bạn cần định nghĩa nó trong package auth
	//& ở trước auth.UserToCreate{} là để tạo một con trỏ đến struct của UserToCreate
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.DisplayName)

	if req.PhotoURL != "" {
		params = params.PhotoURL(req.PhotoURL)
	}
	if req.PhoneNumber != "" {
		params = params.PhoneNumber(req.PhoneNumber)
	}

	user, err := p.client.CreateUser(ctx, params)
	if err != nil {
		log.Printf("❌ Lỗi tạo user: %v", err)
		return nil, err
	}

	return &User{
		UID:         user.UID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		PhotoURL:    user.PhotoURL,
	}, nil
}

func (p *FirebaseProvider) UpdateUser(ctx context.Context, uid string, req *UpdateUserRequest) (*User, error) {
	params := (&auth.UserToUpdate{})

	if req.Email != nil {
		params = params.Email(*req.Email)
	}
	if req.DisplayName != nil {
		params = params.DisplayName(*req.DisplayName)
	}
	if req.PhotoURL != nil {
		params = params.PhotoURL(*req.PhotoURL)
	}
	if req.PhoneNumber != nil {
		params = params.PhoneNumber(*req.PhoneNumber)
	}
	if req.Disabled != nil {
		params = params.Disabled(*req.Disabled)
	}
	user, err := p.client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("❌ Lỗi cập nhật user: %v", err)
		return nil, err
	}
	return &User{
		UID:         user.UID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		PhotoURL:    user.PhotoURL,
	}, nil
}

func (p *FirebaseProvider) SocialLogin(ctx context.Context, provider, accessToken, idToken string) (*User, string, error) {
	switch provider {
	case "google":
		return p.GoogleLogin(ctx, accessToken, idToken)
	case "facebook":
		return p.FacebookLogin(ctx, accessToken)
	default:
		return nil, "", NewAuthError(ErrCodeInternalError, "Unsupported social login provider")
	}
}

func (p *FirebaseProvider) SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) error {
	err := p.client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		return NewAuthError(ErrCodeInternalError, "Failed to set custom claims")
	}
	return nil
}
