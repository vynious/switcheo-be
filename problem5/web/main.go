package main

import (
	"context"
	"crude/x/crude/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

var GlobalClient cosmosclient.Client

// UserRequest defines the structure for the user creation request
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type ChangePasswordRequest struct {
	Id              uint64 `json:"Id"`
	CurrentPassword string `json:"CurrentPassword"`
	NewPassword     string `json:"newPassword"`
}

func main() {
	// Initialize Cosmos client
	ctx := context.Background()
	addressPrefix := "cosmos"
	var err error
	GlobalClient, err = cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix))
	if err != nil {
		log.Fatal("Failed to initialize Cosmos client:", err)
	}

	// Start the server...
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Define other handlers here
	http.HandleFunc("/create-user", HandleCreateUser)
	http.HandleFunc("/update-user", HandleUpdateUser)
	http.HandleFunc("/update-user-password", HandleUpdateUserPassword)
	http.HandleFunc("/get-user", HandleGetUser)
	http.HandleFunc("/get-all-users", HandleGetAllUsers)
	http.HandleFunc("/get-all-users-by-address", HandlerGetUsersByAddress)
	http.HandleFunc("/delete-user", HandlerDeleteUserById)
	http.HandleFunc("/get-all-users-by-email-domain", HandlerGetUsersByEmailDomain)

	log.Println("Listening on http://localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Use globalClient instead of initializing a new client
	log.Println("Received request to create user")

	accountName := "alice"
	account, err := GlobalClient.Account(accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := account.Address("cosmos")
	if err != nil {
		log.Println("Address error")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the request body to get user data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form data: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Create a new UserRequest instance and populate it with form data
	userReq := UserRequest{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Address:  r.FormValue("address"),
	}

	// Create a message to create a user
	msg := &types.MsgCreateUser{
		Creator:  addr,
		Name:     userReq.Name,
		Email:    userReq.Email,
		Username: userReq.Username,
		Password: userReq.Password,
		Address:  userReq.Address,
	}

	ctx := context.Background()
	// Broadcast the transaction
	txResp, err := GlobalClient.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Println("Broadcasting error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the transaction response
	//json.NewEncoder(w).Encode(txResp)
	fmt.Println(txResp)
	fmt.Fprintln(w, "User created successfully")

}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update user")

	accountName := "alice"
	account, err := GlobalClient.Account(accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := account.Address("cosmos")
	if err != nil {
		log.Println("Address error")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	msg := &types.MsgUpdateUser{
		Creator:  addr,
		Id:       id,
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Address:  r.FormValue("address"),
	}

	ctx := context.Background()
	txResp, err := GlobalClient.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Println("Broadcasting error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the transaction response
	//json.NewEncoder(w).Encode(txResp)
	fmt.Println(txResp)
	fmt.Fprintln(w, "User updated successfully")
}

func HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update user password")

	accountName := "alice"
	account, err := GlobalClient.Account(accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := account.Address("cosmos")
	if err != nil {
		log.Println("Address error")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	userReq := ChangePasswordRequest{
		Id:              id,
		CurrentPassword: r.FormValue("currentPassword"),
		NewPassword:     r.FormValue("newPassword"),
	}

	msg := &types.MsgUpdateUserPassword{
		Creator:         addr,
		Id:              userReq.Id,
		CurrentPassword: userReq.CurrentPassword,
		NewPassword:     userReq.NewPassword,
	}

	ctx := context.Background()
	// Broadcast the transaction
	txResp, err := GlobalClient.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Printf("Broadcasting error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the transaction response
	//json.NewEncoder(w).Encode(txResp)
	fmt.Println(txResp)
	fmt.Fprintln(w, "Password updated successfully")
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get user detail")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	queryClient := types.NewQueryClient(GlobalClient.Context())

	// Extract user ID and fetch user details from the blockchain

	ctx := context.Background()
	user, err := queryClient.User(ctx, &types.QueryGetUserRequest{
		Id: id,
	})

	// Assuming user data is stored in 'userDetails' variable
	json.NewEncoder(w).Encode(user)
	fmt.Fprintln(w, "User queried successfully")

}

func HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get all users")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	ctx := context.Background()

	// Instantiate a query client for your `user` blockchain
	queryClient := types.NewQueryClient(GlobalClient.Context())

	// Query the blockchain using the client's `UserAll` method
	// to get all users store all users in queryResp
	queryResp, err := queryClient.UserAll(ctx, &types.QueryAllUserRequest{})
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(queryResp)
	fmt.Fprintln(w, "Users queried successfully")

}

func HandlerGetUsersByAddress(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get all users in specific area")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	address := r.FormValue("address")

	ctx := context.Background()
	// Instantiate a query client for your `user` blockchain
	queryClient := types.NewQueryClient(GlobalClient.Context())

	queryResp, err := queryClient.UserAllByAddress(ctx, &types.QueryAllUserAddressRequest{
		Address: address,
	})
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(queryResp)
	fmt.Fprintln(w, "Users queried successfully")
}

func HandlerDeleteUserById(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get delete user by id")

	accountName := "alice"
	account, err := GlobalClient.Account(accountName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addr, err := account.Address("cosmos")
	if err != nil {
		log.Println("Address error")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(r.FormValue("id"), 10, 64)
	if err != nil {
		// Handle the error if the conversion fails
		log.Printf("Error converting: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	msg := &types.MsgDeleteUser{
		Creator: addr,
		Id:      id,
	}

	ctx := context.Background()
	// Broadcast the transaction
	txResp, err := GlobalClient.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Printf("Broadcasting error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with the transaction response
	//json.NewEncoder(w).Encode(txResp)
	fmt.Println(txResp)
	fmt.Fprintln(w, "User deleted successfully")
}

func HandlerGetUsersByEmailDomain(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get all users by email domain")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	domain := r.FormValue("domain")

	ctx := context.Background()
	// Instantiate a query client for your `user` blockchain
	queryClient := types.NewQueryClient(GlobalClient.Context())

	queryResp, err := queryClient.UserAllByEmailDomain(ctx, &types.QueryAllUserEmailDomainRequest{
		Domain: domain,
	})
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(queryResp)
	fmt.Fprintln(w, "Users queried successfully")
}
