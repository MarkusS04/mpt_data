// Package middleware provides middlewares for http mux server
package middleware

import (
	"context"
	"mpt_data/database"
	"net/http"

	"gorm.io/gorm"
)

type key int

const (
	txKey key = iota
	rollbackKey
)

// TransactionMiddleware provides database transaction handling for API
func TransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := database.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		ctx := context.WithValue(r.Context(), txKey, tx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		if shouldRollback(ctx) {
			tx.Rollback()
		} else {
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	})
}

// GetTx loads a gorm transaction
func GetTx(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txKey).(*gorm.DB)
	return tx
}

func shouldRollback(ctx context.Context) bool {
	rollback, _ := ctx.Value(rollbackKey).(bool)
	return rollback
}

// SetRollback sets wheter a transaction will always be rolled back, for testing only
func SetRollback(ctx context.Context, rollback bool) context.Context {
	return context.WithValue(ctx, rollbackKey, rollback)
}
