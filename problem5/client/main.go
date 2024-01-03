package main

import (
	"context"
	"fmt"
	"log"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	// Importing the types package of your blog blockchain
	"crude/x/crude/types"
)

func main() {
	ctx := context.Background()
	addressPrefix := "cosmos"

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix))
	if err != nil {
		log.Fatal(err)
	}

	// Account `alice` was initialized during `ignite chain serve`
	accountName := "alice"

	// Get account from the keyring
	account, err := client.Account(accountName)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	// Define a message to create a user
	msg := &types.MsgCreateUser{
		Creator:  addr,
		Name:     "sh13aw3123n12312",
		Email:    "sh12a311123wn1@sh12312awn1.com",
		Username: "sha1231wn12311",
		Password: "shawn",
		Address:  "sha13w312312123n1",
	}

	// Broadcast a transaction from account `alice` with the message
	// to create a user store response in txResp
	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgCreateUser:\n\n")
	fmt.Println(txResp)

	// Instantiate a query client for your `user` blockchain
	queryClient := types.NewQueryClient(client.Context())

	// Query the blockchain using the client's `UserAll` method
	// to get all users store all users in queryResp
	queryResp, err := queryClient.UserAll(ctx, &types.QueryAllUserRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Print response from querying all the posts
	fmt.Print("\n\nAll users:\n\n")
	fmt.Println(queryResp)

	x, err := queryClient.User(ctx, &types.QueryGetUserRequest{
		Id: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nSingle user:\n\n")
	fmt.Print(x)

	msgUpdatePassword := &types.MsgUpdateUserPassword{
		Creator:         addr,
		Id:              1,
		CurrentPassword: "shawn",
		NewPassword:     "shawnthiah",
	}

	txUpdatePasswordRes, err := client.BroadcastTx(ctx, account, msgUpdatePassword)
	if err != nil {
		log.Fatal(err)
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgUpdateUserPassword:\n\n")
	fmt.Println(txUpdatePasswordRes)

	y, err := queryClient.User(ctx, &types.QueryGetUserRequest{
		Id: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nSingle user:\n\n")
	fmt.Print(y)

}
