package account

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*postgresRepository, error) {

	log.Println("url for account", url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Println("error creating postgres account repository")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("error pinging postgres account repository")
		return nil, err
	}

	return &postgresRepository{db: db}, nil

}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {

	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id , name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	account := &Account{}
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE  id = $1", id)

	if err := row.Scan(&account.ID, &account.Name); err != nil {
		return nil, err
	}

	return account, nil

}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {

	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM ORDER BY id  DESC OFFSET $1 LIMIT $2", skip, take)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		account := &Account{}

		if err := rows.Scan(&account.ID, &account.Name); err == nil {
			accounts = append(accounts, *account)
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	return accounts, nil

}
