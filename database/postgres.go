package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"kevinPicon/go/rest-ws/models"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	log.Println("Inserting user")
	query := `INSERT INTO users (id, name, last_name, cc, age, birth_date, password, email, address, suburb, voting_place, civil_status, phone, ecan)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	fmt.Println("user: ", user)
	_, err := repo.db.ExecContext(ctx, query, user.Id, user.Name, user.LastName, user.Cc, user.Age, user.BirthDate, user.Password, user.Email, user.Address, user.Suburb, user.VotingPlace, user.CivilStatus, user.Phone, user.Ecan)
	if err != nil {
		fmt.Println("error inserting user: 2", err)
		return fmt.Errorf("error inserting user: %w", err)
	}

	if len(user.Children) > 0 {
		childQuery := `INSERT INTO children (id, user_id, name, last_name, age) 
						VALUES ($1, $2, $3, $4, $5)`

		for _, child := range user.Children {
			_, err = repo.db.ExecContext(ctx, childQuery, child.Id, user.Id, child.Name, child.LastName, child.Age)
			if err != nil {
				fmt.Println("error inserting children: 2", err)
				return fmt.Errorf("error inserting children: %w", err)
			}
		}

	}

	return nil
}
func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.UserPayload, error) {
	var user models.UserPayload
	err := repo.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = repo.db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &user, nil
}
func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	rows, err := repo.db.QueryContext(ctx, "SELECT id, password, name, email FROM users WHERE email = $1",
		email)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Password, &user.Name, &user.Email); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
func (repo *PostgresRepository) GetUsers(ctx context.Context) ([]*models.UserPayload, error) {
	fmt.Println("Init get users")

	rows, err := repo.db.QueryContext(ctx, `
    SELECT id, name, email, age, cc, birth_date, address, suburb, voting_place, civil_status, phone, ecan FROM users`)
	if err != nil {
		fmt.Println("error querying users: ", err)
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	users := []*models.UserPayload{}
	for rows.Next() {
		var user models.UserPayload
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age, &user.Cc, &user.BirthDate, &user.Address, &user.Suburb, &user.VotingPlace, &user.CivilStatus, &user.Phone, &user.Ecan); err != nil {
			fmt.Println("error scanning users: ", err)
			return nil, fmt.Errorf("error scanning users: %w", err)
		}
		children, err := repo.GetChildrenByUserID(ctx, user.Id)
		fmt.Println("children: ", children)
		if err != nil {
			fmt.Println("error getting children: ", err)
			return nil, fmt.Errorf("error getting children: %w", err)
		}
		user.Children = children
		services, err := repo.GetServicesByUserID(ctx, user.Id)
		fmt.Println("services: ", services)
		if err != nil {
			fmt.Println("error getting services: ", err)
			return nil, fmt.Errorf("error getting services: %w", err)
		}
		user.Services = services
		fmt.Println("user: ", user)
		users = append(users, &user)
	}
	return users, nil
}
func (repo *PostgresRepository) GetServicesByUserID(ctx context.Context, userId string) ([]*models.ServicePayload, error) {
	fmt.Println("Init get services")
	rows, err := repo.db.QueryContext(ctx, `
	SELECT service FROM services WHERE user_id = $1`, userId)
	if err != nil {
		fmt.Println("error querying children: ", err)
		return nil, fmt.Errorf("error querying children: %w", err)
	}
	defer rows.Close()
	services := []*models.ServicePayload{}
	for rows.Next() {
		var service models.ServicePayload
		if err := rows.Scan(&service.ServiceName); err != nil {
			fmt.Println("error scanning children: ", err)
			return nil, fmt.Errorf("error scanning children: %w", err)
		}
		services = append(services, &service)

	}
	return services, nil
}
func (repo *PostgresRepository) GetChildrenByUserID(ctx context.Context, userId string) ([]*models.ChildPayload, error) {
	fmt.Println("Init get children")
	rows, err := repo.db.QueryContext(ctx, `
	SELECT name, last_name, age FROM children WHERE user_id = $1`, userId)
	if err != nil {
		fmt.Println("error querying children: ", err)
		return nil, fmt.Errorf("error querying children: %w", err)
	}
	defer rows.Close()
	children := []*models.ChildPayload{}
	for rows.Next() {
		var child models.ChildPayload
		if err := rows.Scan(&child.Name, &child.LastName, &child.Age); err != nil {
			fmt.Println("error scanning children: ", err)
			return nil, fmt.Errorf("error scanning children: %w", err)
		}
		fmt.Println("child: ", child)
		children = append(children, &child)
	}
	fmt.Println("children2: ", children)
	return children, nil
}
func (repo *PostgresRepository) CreateUserService(ctx context.Context, service models.Service) error {
	fmt.Println("Init update user service")
	_, err := repo.db.ExecContext(ctx, ` INSERT INTO services (id, user_id, service) VALUES ($1, $2, $3)`, service.Id, service.UserId, service.ServiceName)
	if err != nil {
		fmt.Println("error updating service: ", err)
		return fmt.Errorf("error updating service: %w", err)
	}
	return nil
}
func (repo *PostgresRepository) UpdateEcan(ctx context.Context, id string) error {
	fmt.Println("Init update ecan")
	_, err := repo.db.ExecContext(ctx, ` UPDATE users SET ecan = true WHERE id = $1`, id)
	if err != nil {
		fmt.Println("error updating ecan: ", err)
		return fmt.Errorf("error updating ecan: %w", err)
	}
	return nil
}
func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
