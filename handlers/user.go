package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"kevinPicon/go/rest-ws/models"
	"kevinPicon/go/rest-ws/repository"
	"kevinPicon/go/rest-ws/server"
	"net/http"
	"strings"
	"time"
)

const (
	HASH_COST = 8
)

type ComplaintsRequest struct {
	Complaint string `json:"complaint"`
}
type SignUpRequest struct {
	Name        string         `json:"name"`
	LastName    string         `json:"last_name"`
	Cc          string         `json:"cc"`
	Age         string         `json:"age"`
	BirthDate   string         `json:"birth_date"`
	Password    string         `json:"password"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	Suburb      string         `json:"suburb"`
	VotingPlace string         `json:"voting_place"`
	CivilStatus string         `json:"civil_status"`
	Phone       string         `json:"phone"`
	Ecan        bool           `json:"ecan"`
	Children    []ChildRequest `json:"children"`
}
type ChildRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Age      string `json:"age"`
}
type ServiceRequest struct {
	ServiceName string `json:"service_name"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
type SignUpResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
type errorResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Init SignUpHandler()")
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//validate user dont exist
		userV, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			fmt.Println("Error getting user by email", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userV.Email != "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse{
				Message: "User already exist",
				Status:  false,
			})
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			fmt.Println("Error hashing password")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			fmt.Println("Error generating id")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var children []models.Child
		for _, childReq := range request.Children {
			child := models.Child{
				Id:       ksuid.New().String(),
				Name:     childReq.Name,
				LastName: childReq.LastName,
				Age:      childReq.Age,
			}

			children = append(children, child)
		}

		var user = models.User{
			Id:          id.String(),
			Email:       request.Email,
			Name:        request.Name,
			LastName:    request.LastName,
			Cc:          request.Cc,
			Age:         request.Age,
			BirthDate:   request.BirthDate,
			Password:    string(hashedPass),
			Address:     request.Address,
			Suburb:      request.Suburb,
			VotingPlace: request.VotingPlace,
			CivilStatus: request.CivilStatus,
			Phone:       request.Phone,
			Ecan:        request.Ecan,
			Children:    children,
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			fmt.Println("Error inserting user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		})
	}

}
func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error de datos", http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, "Error GUE", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Usuario no existe", http.StatusUnauthorized)
			return
		}
		fmt.Println(user)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error de credenciales", http.StatusUnauthorized)
			return
		}
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			tokenString,
		},
		)
	}
}
func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			user, err := repository.GetUserById(r.Context(), claims.UserId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
func UsersHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		users, err := repository.GetUsers(r.Context())
		if err != nil {
			fmt.Println("Error getting users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
func ServiceHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var request = ServiceRequest{}
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http.Error(w, "Error de datos", http.StatusBadRequest)
				return
			}
			var service = models.Service{
				Id:          ksuid.New().String(),
				UserId:      claims.UserId,
				ServiceName: request.ServiceName,
			}
			err = repository.CreateUserService(r.Context(), service)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("OK")
	}
}
func UpdateEcan(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			err = repository.UpdateEcan(r.Context(), claims.UserId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("OK")
	}
}
func InsertComplaints(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var request = ComplaintsRequest{}
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http.Error(w, "Error de datos", http.StatusBadRequest)
				return
			}
			var complaint = models.Complaint{
				Id:        ksuid.New().String(),
				UserId:    claims.UserId,
				Complaint: request.Complaint,
			}
			err = repository.InsertComplaints(r.Context(), complaint)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("OK")
	}
}
