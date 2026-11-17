package service

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// sessionCleanupTimeout bounds one cleanup run.
const sessionCleanupTimeout = time.Minute

// SessionService holds the rules of session management (server ADR 0014): a user sees and ends
// their own sessions; the superadmin sees and ends any user's. Revoking methods return the token
// hashes of the sessions they ended, so the handler can close their live sockets.
type SessionService struct {
	sessionRepository *repository.SessionRepository
	authRepository    *repository.AuthRepository
	sessionConfig     config.SessionConfig
}

// NewSessionService creates a SessionService on dbPool with the retention and cleanup interval of
// sessionConfig.
func NewSessionService(dbPool *pgxpool.Pool, sessionConfig config.SessionConfig) *SessionService {
	return &SessionService{
		sessionRepository: repository.NewSessionRepository(dbPool),
		authRepository:    repository.NewAuthRepository(dbPool),
		sessionConfig:     sessionConfig,
	}
}

// GetMine returns actor's active sessions, most recently used first, with currentSessionID marked.
// "Log in as" sessions of the account are the superadmin's and are left out, except the current one.
func (sessionService *SessionService) GetMine(ctx context.Context, actor model.AuthUser, currentSessionID int64) ([]dto.SessionResponse, error) {
	sessions, err := sessionService.sessionRepository.FindAllActive(ctx, actor.ID)
	if err != nil {
		return nil, err
	}
	sessions = slices.DeleteFunc(sessions, func(session model.Session) bool {
		return session.IsLoginAs() && session.ID != currentSessionID
	})
	return dto.NewSessionResponses(sessions, currentSessionID), nil
}

// DeleteMine revokes actor's session id and returns its token hash. Any session that is not one
// of actor's active password-login sessions answers SessionNotFound, so ids of other users'
// sessions cannot be probed.
func (sessionService *SessionService) DeleteMine(ctx context.Context, actor model.AuthUser, id int64) (string, error) {
	hash, err := sessionService.sessionRepository.RevokeOwn(ctx, id, actor.ID, model.RevokeReasonUserRevoked)
	return hash, toSessionResponseError(err)
}

// DeleteMyOthers revokes actor's active password-login sessions except currentSessionID and
// returns their token hashes.
func (sessionService *SessionService) DeleteMyOthers(ctx context.Context, actor model.AuthUser, currentSessionID int64) ([]string, error) {
	return sessionService.sessionRepository.RevokeOthers(ctx, actor.ID, currentSessionID, model.RevokeReasonUserRevoked)
}

// GetByUser returns every active session of userID, "log in as" ones included, for the superadmin.
func (sessionService *SessionService) GetByUser(ctx context.Context, actor model.AuthUser, userID int64) ([]dto.SessionResponse, error) {
	if err := sessionService.checkSessionUser(ctx, actor, userID); err != nil {
		return nil, err
	}
	sessions, err := sessionService.sessionRepository.FindAllActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.NewSessionResponses(sessions, 0), nil
}

// DeleteByID revokes any user's active session id for the superadmin and returns its token hash.
func (sessionService *SessionService) DeleteByID(ctx context.Context, actor model.AuthUser, id int64) (string, error) {
	if actor.Role != model.RoleSuperAdmin {
		return "", response.Forbidden
	}
	hash, err := sessionService.sessionRepository.RevokeByID(ctx, id, actor.ID, model.RevokeReasonAdminRevoked)
	if err == nil {
		slog.Info("session revoked by superadmin", "sessionId", id, "revokedBy", actor.ID)
	}
	return hash, toSessionResponseError(err)
}

// DeleteByUser revokes every active session of userID for the superadmin and returns their token hashes.
func (sessionService *SessionService) DeleteByUser(ctx context.Context, actor model.AuthUser, userID int64) ([]string, error) {
	if err := sessionService.checkSessionUser(ctx, actor, userID); err != nil {
		return nil, err
	}
	hashes, err := sessionService.sessionRepository.RevokeAllOfUser(ctx, userID, actor.ID, model.RevokeReasonAdminRevoked)
	if err != nil {
		return nil, err
	}
	slog.Info("user sessions revoked by superadmin", "userId", userID, "revokedBy", actor.ID, "count", len(hashes))
	return hashes, nil
}

// checkSessionUser lets only the superadmin manage another user's sessions, and only of a user that
// exists. The route is superAdminOnly too; this keeps the rule when the service is reused.
func (sessionService *SessionService) checkSessionUser(ctx context.Context, actor model.AuthUser, userID int64) error {
	if actor.Role != model.RoleSuperAdmin {
		return response.Forbidden
	}
	_, err := sessionService.authRepository.FindUserByID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return response.SessionUserNotFound
	}
	return err
}

// RunCleanup marks overdue sessions expired and purges ended ones older than SESSION_RETENTION_DAYS,
// once at start and then every SESSION_CLEANUP_MINUTES, until ctx ends.
func (sessionService *SessionService) RunCleanup(ctx context.Context) {
	ticker := time.NewTicker(sessionService.sessionConfig.CleanupInterval)
	defer ticker.Stop()
	for {
		sessionService.CleanUp(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// CleanUp runs one cleanup: overdue active sessions become expired, and logged out, revoked and
// expired sessions that ended more than the retention period ago are deleted. Errors are logged.
func (sessionService *SessionService) CleanUp(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, sessionCleanupTimeout)
	defer cancel()
	expired, err := sessionService.sessionRepository.ExpireOverdue(ctx)
	if err != nil {
		slog.Error("session cleanup: expire", "error", err)
		return
	}
	purged, err := sessionService.sessionRepository.PurgeEnded(ctx, sessionService.sessionConfig.Retention)
	if err != nil {
		slog.Error("session cleanup: purge", "error", err)
		return
	}
	if expired > 0 || purged > 0 {
		slog.Info("session cleanup", "expired", expired, "purged", purged)
	}
}

// toSessionResponseError maps repository errors of session management to responses.
func toSessionResponseError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return response.SessionNotFound
	}
	return err
}
