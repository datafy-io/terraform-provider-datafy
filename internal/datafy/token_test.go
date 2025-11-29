package datafy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountToken(t *testing.T) {
	expires := time.Now().Add(30 * time.Minute).UTC().Round(time.Second)
	created := time.Now().UTC().Round(time.Second)
	expected := AccountToken{
		AccountId:   "acc-123",
		TokenId:     "tok-abc",
		Description: "read-only",
		Secret:      "secret-xyz",
		Expires:     expires,
		CreatedAt:   created,
		RoleIds:     []string{"role-1", "role-2"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/tokens" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.CreateAccountToken(context.Background(), &CreateAccountTokenRequest{
		AccountId:   expected.AccountId,
		Description: expected.Description,
		Ttl:         30 * time.Minute,
		RoleIds:     expected.RoleIds,
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.AccountToken)
}

func TestGetAccountToken(t *testing.T) {
	expires := time.Now().Add(30 * time.Minute).UTC().Round(time.Second)
	created := time.Now().UTC().Round(time.Second)
	expected := AccountToken{
		AccountId:   "acc-123",
		TokenId:     "tok-abc",
		Description: "read-only",
		Expires:     expires,
		CreatedAt:   created,
		RoleIds:     []string{"role-1", "role-2"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/tokens/tok-abc" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.GetAccountToken(context.Background(), &GetAccountTokenRequest{AccountId: expected.AccountId, TokenId: expected.TokenId})

	assert.NoError(t, err)
	assert.Equal(t, expected, out.AccountToken)
}

func TestDeleteAccountToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/api/v1/accounts/acc-123/tokens/tok-abc" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := NewClient("dummy", ts.URL)
	out, err := c.DeleteAccountToken(context.Background(), &DeleteAccountTokenRequest{AccountId: "acc-123", TokenId: "tok-abc"})

	assert.NoError(t, err)
	assert.NotNil(t, out)
}
