package handlers

import accountpb "github.com/Edbeer/proto/api/account/v1"

type AccountService interface {}

type accountHandler struct {
	accountpb.UnimplementedAccountServiceServer
	service AccountService
}

func newAccountHandler(service AccountService) *accountHandler {
	return &accountHandler{
		service: service,
	}
}