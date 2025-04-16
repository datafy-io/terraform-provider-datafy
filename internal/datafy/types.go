package datafy

import (
	"time"
)

type CreateAccountRequest struct {
	AccountName string
}

type CreateAccountResponse struct {
	Account Account
}

type UpdateAccountRequest struct {
	AccountId   string
	AccountName string
}

type UpdateAccountResponse struct {
	Account Account
}

type GetAccountRequest struct {
	AccountId string
}

type GetAccountResponse struct {
	Account Account
}

type DeleteAccountRequest struct {
	AccountId string
}

type DeleteAccountResponse struct {
}

type Account struct {
	AccountId       string `json:"accountId"`
	AccountName     string `json:"accountName"`
	ParentAccountId string `json:"parentAccountId"`
}

type CreateAccountRoleArnRequest struct {
	AccountId string
	Arn       string
}

type CreateAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type GetAccountRoleArnRequest struct {
	AccountId string
}

type GetAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type UpdateAccountRoleArnRequest struct {
	AccountId string
	Arn       string
}

type UpdateAccountRoleArnResponse struct {
	AccountRoleArn AccountRoleArn
}

type DeleteAccountRoleArnRequest struct {
	AccountId string
}

type DeleteAccountRoleArnResponse struct {
}

type AccountRoleArn struct {
	RoleArn string `json:"roleArn"`
}

type CreateAccountTokenRequest struct {
	AccountId   string
	Description string
	Ttl         time.Duration
	RoleIds     []string
}

type CreateAccountTokenResponse struct {
	AccountToken AccountToken
}

type GetAccountTokenRequest struct {
	AccountId string
	TokenId   string
}

type GetAccountTokenResponse struct {
	AccountToken AccountToken
}

type DeleteAccountTokenRequest struct {
	AccountId string
	TokenId   string
}

type DeleteAccountTokenResponse struct {
}

type AccountToken struct {
	AccountId   string    `json:"accountId"`
	TokenId     string    `json:"tokenId"`
	Description string    `json:"description"`
	Secret      string    `json:"secret,omitempty"`
	Expires     time.Time `json:"expires"`
	CreatedAt   time.Time `json:"createdAt"`
	RoleIds     []string  `json:"roleIds"`
}
