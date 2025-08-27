package repositories

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/MinhT933/file-converter/internal/domain"
)

type UserRepository struct {  // struct, không phải interface
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
    query := `
        INSERT INTO users (user_id, email, name, avatar_url, provider, 
                          email_verified, created_at, is_premium, updated_at, provider_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
    
    _, err := r.db.ExecContext(ctx, query,
        user.UserID, user.Email, user.Name, user.AvatarURL, user.Provider,
        user.EmailVerified, user.CreatedAt, user.IsPremium, user.UpdatedAt, user.ProviderID)
    
    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT user_id, email, name, avatar_url, provider, email_verified, 
               created_at, is_premium, updated_at, provider_id
        FROM users WHERE email = $1`
    
    var user domain.User
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.UserID, &user.Email, &user.Name, &user.AvatarURL, &user.Provider,
        &user.EmailVerified, &user.CreatedAt, &user.IsPremium, &user.UpdatedAt, &user.ProviderID)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    
    return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
    query := `
        SELECT user_id, email, name, avatar_url, provider, email_verified, 
               created_at, is_premium, updated_at, provider_id
        FROM users WHERE user_id = $1`
    
    var user domain.User
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &user.UserID, &user.Email, &user.Name, &user.AvatarURL, &user.Provider,
        &user.EmailVerified, &user.CreatedAt, &user.IsPremium, &user.UpdatedAt, &user.ProviderID)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    
    return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
    query := `
        UPDATE users 
        SET email = $2, name = $3, avatar_url = $4, provider = $5, 
            email_verified = $6, is_premium = $7, updated_at = $8, provider_id = $9
        WHERE user_id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        user.UserID, user.Email, user.Name, user.AvatarURL, user.Provider,
        user.EmailVerified, user.IsPremium, user.UpdatedAt, user.ProviderID)
    
    return err
}

func (r *UserRepository) Delete(ctx context.Context, userID string) error {
    query := `DELETE FROM users WHERE user_id = $1`
    _, err := r.db.ExecContext(ctx, query, userID)
    return err
}

func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
    // Nếu chưa cần dùng, có thể trả về empty slice
    return []*domain.User{}, nil
}
