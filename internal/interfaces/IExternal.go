package interfaces

import (
	"context"

	"ewallet-ums/external"
)

type IWallet interface {
	CreateWallet(ctx context.Context, userID int) (*external.Wallet, error)
}
