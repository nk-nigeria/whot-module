package cgbdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/heroiclabs/nakama-common/runtime"
	_ "github.com/lib/pq"
	"github.com/nakama-nigeria/whot-module/entity"
)

func TestAddNewFeeGame(t *testing.T) {
	type args struct {
		ctx     context.Context
		logger  runtime.Logger
		db      *sql.DB
		feeGame entity.FeeGame
	}
	connStr := "postgresql://postgres:localdb@127.0.0.1/nakama?sslmode=disable"
	mdb, _ := sql.Open("postgres", connStr)
	defer mdb.Close()
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-add-fee-game",
			args: args{
				ctx:    context.Background(),
				logger: &entity.EmptyLogger{},
				db:     mdb,
				feeGame: entity.FeeGame{
					UserID: "x1",
					Game:   "",
					Fee:    10000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddNewFeeGame(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.feeGame)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNewFeeGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == 0 {
				t.Errorf("AddNewFeeGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddNewMultiFeGame(t *testing.T) {
	type args struct {
		ctx         context.Context
		logger      runtime.Logger
		db          *sql.DB
		listFeeGame []entity.FeeGame
	}
	connStr := "postgresql://postgres:localdb@127.0.0.1/nakama?sslmode=disable"
	mdb, _ := sql.Open("postgres", connStr)
	defer mdb.Close()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				ctx:    context.Background(),
				logger: &entity.EmptyLogger{},
				db:     mdb,
				listFeeGame: []entity.FeeGame{
					{
						UserID: "1",
						Game:   "",
						Fee:    1000,
					},
					{
						UserID: "2",
						Game:   "",
						Fee:    2000,
					},
					{
						UserID: "3",
						Game:   "",
						Fee:    3000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddNewMultiFeeGame(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.listFeeGame); (err != nil) != tt.wantErr {
				t.Errorf("AddNewMultiFeGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
