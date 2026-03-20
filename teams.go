package flow

import (
	"context"
	"encoding/json"
	"fmt"
)

type TeamService struct{ http *httpClient }

func (s *TeamService) GetTeam(ctx context.Context) (map[string]interface{}, error) {
	data, err := s.http.get(ctx, "/merchants/me/team", nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("flow: unmarshal team: %w", err)
	}
	return result, nil
}

func (s *TeamService) GetMyRole(ctx context.Context) (map[string]interface{}, error) {
	data, err := s.http.get(ctx, "/merchants/me/team/my-role", nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("flow: unmarshal role: %w", err)
	}
	return result, nil
}

func (s *TeamService) Invite(ctx context.Context, email, role string) (map[string]interface{}, error) {
	data, err := s.http.post(ctx, "/merchants/me/team/invite", map[string]string{"email": email, "role": role})
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("flow: unmarshal invite: %w", err)
	}
	return result, nil
}

func (s *TeamService) RevokeInvite(ctx context.Context, inviteID string) (map[string]interface{}, error) {
	data, err := s.http.post(ctx, "/merchants/me/team/invite/"+inviteID+"/revoke", nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("flow: unmarshal revoke: %w", err)
	}
	return result, nil
}

func (s *TeamService) UpdateRole(ctx context.Context, memberID, role string) (map[string]interface{}, error) {
	data, err := s.http.do(ctx, "PATCH", "/merchants/me/team/"+memberID+"/role", map[string]string{"role": role}, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("flow: unmarshal role update: %w", err)
	}
	return result, nil
}

func (s *TeamService) Remove(ctx context.Context, memberID string) error {
	_, err := s.http.do(ctx, "DELETE", "/merchants/me/team/"+memberID, nil, nil)
	return err
}
