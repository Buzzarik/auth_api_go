package postgres

import (
	"auth/internal/config"
	"auth/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
};

func New(cnf config.ConfigPostgres) (*Postgres, error){
	const op = "StoragePostgres.New";

	//создаем конфиг строку
	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cnf.Host,
			cnf.Port,
			cnf.UserName,
			cnf.Password,
			cnf.DbName,
			cnf.Sslmode);
	
	db, err := sql.Open(cnf.Driver, conn_str);

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err);
	}

	db.SetMaxOpenConns(int(cnf.MaxOpenConns));
    db.SetMaxIdleConns(int(cnf.MaxIdleConns));
    db.SetConnMaxIdleTime(cnf.MaxIdleTime);

	//NOTE: проверка на подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second);
    defer cancel();

	err = db.PingContext(ctx);
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err);;
    }

	return &Postgres{db: db}, nil;
}

//TODO:
func (s *Postgres) GetByPhoneNumber(phoneNumber string) (*models.User, error){
	const op = "Postgres.GetByPhoneNumber";

	stmt, err := s.db.Prepare(`
		SELECT id, created_at, name, phone_number, hash_password
		FROM users
		WHERE phone_number = $1
	`);

	if err != nil {
		return nil, fmt.Errorf("(Prepare)%s: %w", op, err);
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3);
	defer cancel();

	var user models.User;
	err = stmt.QueryRowContext(ctx, phoneNumber).Scan(&user.ID, 
													&user.CreatedAt, 
													&user.Name,
													&user.PhoneNumber,
													&user.HashPassword);
	// нет такого ресурса
	if err == sql.ErrNoRows {
		return nil, nil;
	}
	if err != nil {
		return nil, fmt.Errorf("(QueryRowContext)%s: %w", op, err);
	}

	return &user, nil;
}

//TODO:
func (s *Postgres) SetUser(user *models.User) error {
	//различать ошибки CONSTAINT от других
	const op = "Postgres.SetUser";

	stmt, err := s.db.Prepare(`        
		INSERT INTO users (name, phone_number, hash_password)
		VALUES ($1, $2, $3)
	`);
	
	if err != nil {
		return fmt.Errorf("(Prepare)%s: %w", op, err);
	}

	//ставим таймаут на запрос
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel();
	
	err = stmt.QueryRowContext(ctx, user.Name, user.PhoneNumber, user.HashPassword).Err();

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return fmt.Errorf("CONSTRAINT");
	}

	if (err != nil){
		return fmt.Errorf("(QueryRowContext)%s: %w", op, err);
	}

	return nil;
}

func (s *Postgres) runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin();
	if err != nil {
		return err;
	}

	err = fn(tx);
	if err == nil {
		return tx.Commit();
	}

	rollbackErr := tx.Rollback();
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr);
	}

	return err;
}
//TODO
func (s *Postgres) SetToken(token *models.Token) error {
	return s.runInTx(s.db, func(tx *sql.Tx) error {
		const op = "Postgres.SetToken";

		stmt1, err2 := s.db.Prepare(`
			INSERT INTO tokens (id_user, hash, expiry, id_api)
			VALUES ($1, $2, $3, $4)
		`);

		stmt2, err := s.db.Prepare(`
			SELECT hash, expiry, id_api FROM tokens
			WHERE id_user = $1 AND id_api = $2
		`);

		if err != nil || err2 != nil {
			return fmt.Errorf("(Prepare1)%s: %w", op, errors.Join(err, err2));
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
		defer cancel();

		err = stmt1.QueryRowContext(ctx, 
				token.IdUser, token.Hash,
				token.Expiry, token.IdAPI).Err();

		if err == nil {
			return nil;
		}

		//если это не CONSTRAINT
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code != "23505" {
			return err;
		}

		//обновляем token, так как уже создан этим приложением
		err = stmt2.QueryRowContext(ctx, token.IdUser, token.IdAPI).Scan(
			&token.Hash,
			&token.Expiry,
			&token.IdAPI,
		);

		if (err != nil){
			return fmt.Errorf("(QueryRowContext)%s: %w", op, err);
		}

		return nil;
	});
}