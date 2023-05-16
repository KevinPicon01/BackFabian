package models

type User struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	LastName    string  `json:"last_name"`
	Cc          string  `json:"cc"`
	Age         string  `json:"age"`
	BirthDate   string  `json:"birth_date"`
	Password    string  `json:"password"`
	Email       string  `json:"email"`
	Address     string  `json:"address"`
	Suburb      string  `json:"suburb"`
	VotingPlace string  `json:"voting_place"`
	CivilStatus string  `json:"civil_status"`
	Phone       string  `json:"phone"`
	Ecan        bool    `json:"ecan"`
	Children    []Child `json:"children"`
}

type Child struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Age      string `json:"age"`
	CreateAt string `json:"create_at"`
}
type Service struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	ServiceName string `json:"service_name"`
}
type UserPayload struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Age         int               `json:"age"`
	Cc          string            `json:"cc"`
	BirthDate   string            `json:"birth_date"`
	Address     string            `json:"address"`
	Suburb      string            `json:"suburb"`
	VotingPlace string            `json:"voting_place"`
	CivilStatus string            `json:"civil_status"`
	Phone       string            `json:"phone"`
	Ecan        string            `json:"ecan"`
	Children    []*ChildPayload   `json:"children,omitempty"`
	Services    []*ServicePayload `json:"services,omitempty"`
}

type ChildPayload struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Age      int    `json:"age"`
}

type ServicePayload struct {
	ServiceName string `json:"service_name"`
}
