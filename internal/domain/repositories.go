//Định nghĩa những công việc cần làm

package domain

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context) ([]*User, error)
}

type ConversionRepository interface {
	Create(ctx context.Context, conversion *Conversion) (string, error)
	FindByUserID(ctx context.Context, userID string) ([]Conversion, error)
	FindByID(ctx context.Context, conversionID string) (*Conversion, error)
	UpdateConversionStatus(ctx context.Context, id string, status string) error
}
