package postgresql

import (
	"context"
	"encoding/hex"
	"github.com/RakhimovAns/Shop/types"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type CustomerService struct {
	pool *pgxpool.Pool
}

func NewCustomerService(pool *pgxpool.Pool) *CustomerService {
	return &CustomerService{pool: pool}
}

func (s *CustomerService) Register(ctx context.Context, customer *types.Customer) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*customer.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(*customer.Password))
	if err != nil {
		log.Println(err)
		return types.ErrInvalidPassword
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

func (s *CustomerService) Login(ctx context.Context, login string, password string) (string, error) {
	var hash string
	var id int64
	err := s.pool.QueryRow(ctx, `
		select id, password from customers where phone=$1
`, login).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", types.ErrNoSuchUser
	}
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return "", types.ErrInvalidPassword
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.TokenClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		id,
	})
	TokenStr, err := token.SignedString([]byte("My Key"))
	return TokenStr, nil
}

func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	_, err := s.pool.Exec(ctx, `
		delete from customers where id=$1
`, id)
	if err == types.ErrNoSuch {
		log.Println(types.ErrNoSuchUser)
		return types.ErrNoSuchUser
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *CustomerService) GetByID(ctx context.Context, id int64) (error, *types.Customer) {
	customer := &types.Customer{}
	err := s.pool.QueryRow(ctx, `
		select  id,name,phone,password,active,created,balance from customers where id=$1
`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Password, &customer.Active, &customer.Created, &customer.Balance)
	if err == types.ErrNoSuchUser {
		log.Println(types.ErrNoSuchUser)
		return types.ErrNoSuchUser, nil
	}
	if err != nil {
		log.Println(err)
		return err, nil
	}
	return nil, customer
}

func (s *CustomerService) ChangeBalance(ctx context.Context, id int64, sum int64) error {
	_, err := s.pool.Exec(ctx, `
		update customers set balance=customers.balance-$1 where id=$2
`, sum, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *CustomerService) DepositBalance(ctx context.Context, id int64, sum int64) error {
	_, err := s.pool.Exec(ctx, `
		update customers set balance=customers.balance+$1 where id=$2
`, sum, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
