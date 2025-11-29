package datafy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountRoleArn(t *testing.T) {
	expected := AccountRoleArn{RoleArn: "arn:aws:iam::123456789012:role/test"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/role-arn" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.CreateAccountRoleArn(context.Background(), &CreateAccountRoleArnRequest{
		AccountId: "acc-123",
		Arn:       expected.RoleArn,
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.AccountRoleArn)
}

func TestGetAccountRoleArn(t *testing.T) {
	expected := AccountRoleArn{RoleArn: "arn:aws:iam::123456789012:role/test"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/role-arn" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.GetAccountRoleArn(context.Background(), &GetAccountRoleArnRequest{AccountId: "acc-123"})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.AccountRoleArn)
}

func TestUpdateAccountRoleArn(t *testing.T) {
	expected := AccountRoleArn{RoleArn: "arn:aws:iam::123456789012:role/updated"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Update uses Create under the hood
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/role-arn" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.UpdateAccountRoleArn(context.Background(), &UpdateAccountRoleArnRequest{AccountId: "acc-123", Arn: expected.RoleArn})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.AccountRoleArn)
}

func TestDeleteAccountRoleArn(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/role-arn" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.DeleteAccountRoleArn(context.Background(), &DeleteAccountRoleArnRequest{AccountId: "acc-123"})

	assert.NoError(t, err)
	assert.NotNil(t, out)
}
