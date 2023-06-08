package customer

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/RakhimovAns/Shop/cmd/help"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

var ErrNoSuchUser = errors.New("no such user")

// var ErrInternal = errors.New("internal error")
var ErrInvalidPassword = errors.New("invalid password")

type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password *string   `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
	Balance  int64     `json:"balance"`
}
type Token struct {
	Token string `json:"token"`
}

func (s *Service) Register(ctx context.Context, customer *Customer) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*customer.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(*customer.Password))
	if err != nil {
		log.Println(err)
		return ErrInvalidPassword
	}
	_, err = s.pool.Exec(ctx, `
insert into customers(name,phone,password,balance) values ($1,$2,$3,$4) on conflict (phone) do update set name=excluded.name
`, customer.Name, customer.Phone, hex.EncodeToString(hash), customer.Balance)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) Login(ctx context.Context, login string, password string) (string, error) {
	var hash string
	var id int64
	err := s.pool.QueryRow(ctx, `
select id, password from customers where phone=$1
`, login).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &help.TokenClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		id,
	})
	TokenStr, err := token.SignedString([]byte("My Key"))
	return TokenStr, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	_, err := s.pool.Exec(ctx, `
delete from customers where id=$1
`, id)
	if err == ErrNoSuchUser {
		log.Println(ErrNoSuchUser)
		return ErrNoSuchUser
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (error, *Customer) {
	customer := &Customer{}
	err := s.pool.QueryRow(ctx, `
select  id,name,phone,password,active,created,balance from customers where id=$1
`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Password, &customer.Active, &customer.Created, &customer.Balance)
	if err == ErrNoSuchUser {
		log.Println(ErrNoSuchUser)
		return ErrNoSuchUser, nil
	}
	if err != nil {
		log.Println(err)
		return err, nil
	}
	return nil, customer
}

func (s *Service) ChangeBalance(ctx context.Context, id int64, sum int64) error {
	_, err := s.pool.Exec(ctx, `
update customers set balance=customers.balance-$1 where id=$2
`, sum, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (s *Service) DepositBalance(ctx context.Context, id int64, sum int64) error {
	_, err := s.pool.Exec(ctx, `
update customers set balance=customers.balance+$1 where id=$2
`, sum, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
