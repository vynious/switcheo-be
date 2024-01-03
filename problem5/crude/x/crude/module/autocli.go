package crude

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "crude/api/crude/crude"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "UserAll",
					Use:       "list-user",
					Short:     "List all user",
				},
				{
					RpcMethod:      "User",
					Use:            "show-user [id]",
					Short:          "Shows a user by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateUser",
					Use:            "create-user [name] [email] [username] [password] [address]",
					Short:          "Create user",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "name"}, {ProtoField: "email"}, {ProtoField: "username"}, {ProtoField: "password"}, {ProtoField: "address"}},
				},
				{
					RpcMethod:      "UpdateUser",
					Use:            "update-user [id] [name] [email] [username] [password] [address]",
					Short:          "Update user",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "name"}, {ProtoField: "email"}, {ProtoField: "username"}, {ProtoField: "password"}, {ProtoField: "address"}},
				},
				{
					RpcMethod:      "DeleteUser",
					Use:            "delete-user [id]",
					Short:          "Delete user",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
