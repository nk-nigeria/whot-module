package cgbdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/heroiclabs/nakama-common/runtime"
	pb "github.com/nakama-nigeria/cgp-common/proto/whot"
	"github.com/nakama-nigeria/whot-module/entity"
)

func TestAddOrUpdateChipJackpot(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger runtime.Logger
		db     *sql.DB
		game   string
		chips  int64
	}
	connStr := "postgresql://postgres:localdb@127.0.0.1/nakama?sslmode=disable"
	mdb, _ := sql.Open("postgres", connStr)
	defer mdb.Close()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			// TODO: Add test cases.
			name: "AddOrUpdateChipJackpot",
			args: args{
				ctx:    context.Background(),
				logger: &entity.EmptyLogger{},
				db:     mdb,
				game:   entity.ModuleName,
				chips:  1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddOrUpdateChipJackpot(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.game, tt.args.chips); (err != nil) != tt.wantErr {
				t.Errorf("AddOrUpdateChipJackpot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncChipJackpot(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger runtime.Logger
		db     *sql.DB
		game   string
		chips  int64
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
			name: "IncChipJackpot",
			args: args{
				ctx:    context.Background(),
				logger: &entity.EmptyLogger{},
				db:     mdb,
				game:   entity.ModuleName,
				chips:  -2000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if err := IncChipJackpot(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.game, tt.args.chips); (err != nil) != tt.wantErr {
			// 	t.Errorf("IncChipJackpot() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}

func TestGetJackpot(t *testing.T) {
	type args struct {
		ctx    context.Context
		logger runtime.Logger
		db     *sql.DB
		game   string
	}
	connStr := "postgresql://postgres:localdb@127.0.0.1/nakama?sslmode=disable"
	mdb, _ := sql.Open("postgres", connStr)
	defer mdb.Close()
	tests := []struct {
		name    string
		args    args
		want    *pb.Jackpot
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestGetJackpot",
			args: args{
				ctx:    context.TODO(),
				logger: &entity.EmptyLogger{},
				db:     mdb,
				game:   entity.ModuleName,
			},
			want: &pb.Jackpot{
				Chips: 1000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJackpot(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.game)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJackpot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Chips != tt.want.Chips {
				t.Errorf("GetJackpot() = %v, want %v", got, tt.want)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetJackpot() = %v, want %v", got, tt.want)
			// }
		})
	}
}
