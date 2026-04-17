package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
)

type PostgresAdminRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAdminRepository(pool *pgxpool.Pool) *PostgresAdminRepository {
	return &PostgresAdminRepository{pool: pool}
}

func (r *PostgresAdminRepository) Users() repository.UserRepository             { return r }
func (r *PostgresAdminRepository) Permissions() repository.PermissionRepository { return r }
func (r *PostgresAdminRepository) Groups() repository.GroupRepository           { return r }
func (r *PostgresAdminRepository) Domains() repository.DomainRepository         { return r }

func (r *PostgresAdminRepository) ListUsers(ctx context.Context) ([]model.PortalUser, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, username, password_hash, COALESCE(email, ''), COALESCE(display_name, ''), role, active FROM users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.PortalUser, 0)
	for rows.Next() {
		var item model.PortalUser
		if err := rows.Scan(&item.ID, &item.Username, &item.PasswordHash, &item.Email, &item.DisplayName, &item.Role, &item.Active); err != nil {
			return nil, err
		}
		groupIDs, err := r.fetchUserGroupIDs(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.GroupIDs = groupIDs
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetUser(ctx context.Context, id string) (model.PortalUser, error) {
	var item model.PortalUser
	err := r.pool.QueryRow(ctx, `SELECT id, username, password_hash, COALESCE(email, ''), COALESCE(display_name, ''), role, active FROM users WHERE id = $1`, id).
		Scan(&item.ID, &item.Username, &item.PasswordHash, &item.Email, &item.DisplayName, &item.Role, &item.Active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PortalUser{}, apperrors.NewNotFound("user not found")
		}
		return model.PortalUser{}, err
	}
	groupIDs, err := r.fetchUserGroupIDs(ctx, id)
	if err != nil {
		return model.PortalUser{}, err
	}
	item.GroupIDs = groupIDs
	return item, nil
}

func (r *PostgresAdminRepository) CreateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalUser{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, email, display_name, role, active)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6, $7)
	`, user.ID, user.Username, user.PasswordHash, user.Email, user.DisplayName, user.Role, user.Active)
	if err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user")
	}
	if err := replaceUserGroups(ctx, tx, user.ID, user.GroupIDs); err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user group")
	}
	if err := tx.Commit(ctx); err != nil {
		return model.PortalUser{}, err
	}
	return r.GetUser(ctx, user.ID)
}

func (r *PostgresAdminRepository) UpdateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalUser{}, err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
		UPDATE users
		SET username = $2, password_hash = $3, email = NULLIF($4, ''), display_name = NULLIF($5, ''), role = $6, active = $7, updated_at = NOW()
		WHERE id = $1
	`, user.ID, user.Username, user.PasswordHash, user.Email, user.DisplayName, user.Role, user.Active)
	if err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user")
	}
	if commandTag.RowsAffected() == 0 {
		return model.PortalUser{}, apperrors.NewNotFound("user not found")
	}
	if err := replaceUserGroups(ctx, tx, user.ID, user.GroupIDs); err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user group")
	}
	if err := tx.Commit(ctx); err != nil {
		return model.PortalUser{}, err
	}
	return r.GetUser(ctx, user.ID)
}

func (r *PostgresAdminRepository) DeleteUser(ctx context.Context, id string) error {
	commandTag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return apperrors.NewNotFound("user not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListPermissions(ctx context.Context) ([]model.PermissionDefinition, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(description, '') FROM permissions ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.PermissionDefinition, 0)
	for rows.Next() {
		var item model.PermissionDefinition
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetPermission(ctx context.Context, id string) (model.PermissionDefinition, error) {
	var item model.PermissionDefinition
	err := r.pool.QueryRow(ctx, `SELECT id, name, COALESCE(description, '') FROM permissions WHERE id = $1`, id).
		Scan(&item.ID, &item.Name, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
		}
		return model.PermissionDefinition{}, err
	}
	return item, nil
}

func (r *PostgresAdminRepository) CreatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	_, err := r.pool.Exec(ctx, `INSERT INTO permissions (id, name, description) VALUES ($1, $2, NULLIF($3, ''))`, permission.ID, permission.Name, permission.Description)
	if err != nil {
		return model.PermissionDefinition{}, mapPostgresError(err, "permission")
	}
	return r.GetPermission(ctx, permission.ID)
}

func (r *PostgresAdminRepository) UpdatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	tag, err := r.pool.Exec(ctx, `UPDATE permissions SET name = $2, description = NULLIF($3, ''), updated_at = NOW() WHERE id = $1`, permission.ID, permission.Name, permission.Description)
	if err != nil {
		return model.PermissionDefinition{}, mapPostgresError(err, "permission")
	}
	if tag.RowsAffected() == 0 {
		return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
	}
	return r.GetPermission(ctx, permission.ID)
}

func (r *PostgresAdminRepository) DeletePermission(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM permissions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("permission not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListGroups(ctx context.Context) ([]model.PortalGroup, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(description, '') FROM groups ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.PortalGroup, 0)
	for rows.Next() {
		var item model.PortalGroup
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		permissionIDs, err := r.fetchGroupPermissionIDs(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.PermissionIDs = permissionIDs
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetGroup(ctx context.Context, id string) (model.PortalGroup, error) {
	var item model.PortalGroup
	err := r.pool.QueryRow(ctx, `SELECT id, name, COALESCE(description, '') FROM groups WHERE id = $1`, id).
		Scan(&item.ID, &item.Name, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PortalGroup{}, apperrors.NewNotFound("group not found")
		}
		return model.PortalGroup{}, err
	}
	permissionIDs, err := r.fetchGroupPermissionIDs(ctx, id)
	if err != nil {
		return model.PortalGroup{}, err
	}
	item.PermissionIDs = permissionIDs
	return item, nil
}

func (r *PostgresAdminRepository) CreateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalGroup{}, err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `INSERT INTO groups (id, name, description) VALUES ($1, $2, NULLIF($3, ''))`, group.ID, group.Name, group.Description); err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group")
	}
	if err := replaceGroupPermissions(ctx, tx, group.ID, group.PermissionIDs); err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group permission")
	}
	if err := tx.Commit(ctx); err != nil {
		return model.PortalGroup{}, err
	}
	return r.GetGroup(ctx, group.ID)
}

func (r *PostgresAdminRepository) UpdateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalGroup{}, err
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, `UPDATE groups SET name = $2, description = NULLIF($3, ''), updated_at = NOW() WHERE id = $1`, group.ID, group.Name, group.Description)
	if err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group")
	}
	if tag.RowsAffected() == 0 {
		return model.PortalGroup{}, apperrors.NewNotFound("group not found")
	}
	if err := replaceGroupPermissions(ctx, tx, group.ID, group.PermissionIDs); err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group permission")
	}
	if err := tx.Commit(ctx); err != nil {
		return model.PortalGroup{}, err
	}
	return r.GetGroup(ctx, group.ID)
}

func (r *PostgresAdminRepository) DeleteGroup(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM groups WHERE id = $1`, id)
	if err != nil {
		return mapPostgresError(err, "group")
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("group not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListDomains(ctx context.Context) ([]model.RegisteredDomain, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(description, ''), owner_group_id, status, metadata FROM domains ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.RegisteredDomain, 0)
	for rows.Next() {
		item, err := scanDomain(rows)
		if err != nil {
			return nil, err
		}
		reviews, err := r.fetchDomainReviews(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.Review = reviews
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) ListApprovedDomains(ctx context.Context) ([]model.RegisteredDomain, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(description, ''), owner_group_id, status, metadata FROM domains WHERE status = 'APPROVED' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.RegisteredDomain, 0)
	for rows.Next() {
		item, err := scanDomain(rows)
		if err != nil {
			return nil, err
		}
		reviews, err := r.fetchDomainReviews(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		item.Review = reviews
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetDomain(ctx context.Context, id string) (model.RegisteredDomain, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, name, COALESCE(description, ''), owner_group_id, status, metadata FROM domains WHERE id = $1`, id)
	item, err := scanDomain(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
		}
		return model.RegisteredDomain{}, err
	}
	reviews, err := r.fetchDomainReviews(ctx, id)
	if err != nil {
		return model.RegisteredDomain{}, err
	}
	item.Review = reviews
	return item, nil
}

func (r *PostgresAdminRepository) CreateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	metadata, _ := json.Marshal(domain.Metadata)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO domains (id, name, description, owner_group_id, status, metadata)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6::jsonb)
	`, domain.ID, domain.Name, domain.Description, domain.OwnerGroupID, domain.Status, string(metadata))
	if err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain")
	}
	return r.GetDomain(ctx, domain.ID)
}

func (r *PostgresAdminRepository) UpdateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	metadata, _ := json.Marshal(domain.Metadata)
	tag, err := r.pool.Exec(ctx, `
		UPDATE domains
		SET name = $2, description = NULLIF($3, ''), owner_group_id = $4, status = $5, metadata = $6::jsonb, updated_at = NOW()
		WHERE id = $1
	`, domain.ID, domain.Name, domain.Description, domain.OwnerGroupID, domain.Status, string(metadata))
	if err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain")
	}
	if tag.RowsAffected() == 0 {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	return r.GetDomain(ctx, domain.ID)
}

func (r *PostgresAdminRepository) DeleteDomain(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM domains WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("domain not found")
	}
	return nil
}

func (r *PostgresAdminRepository) AddDomainReview(ctx context.Context, id string, comment model.ReviewComment, nextStatus string) (model.RegisteredDomain, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.RegisteredDomain{}, err
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, `UPDATE domains SET status = $2, updated_at = NOW() WHERE id = $1`, id, nextStatus)
	if err != nil {
		return model.RegisteredDomain{}, err
	}
	if tag.RowsAffected() == 0 {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO domain_reviews (id, domain_id, author, author_id, decision, comment, created_at)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7)
	`, comment.ID, id, comment.Author, comment.AuthorID, comment.Decision, comment.Comment, comment.CreatedAt); err != nil {
		return model.RegisteredDomain{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return model.RegisteredDomain{}, err
	}
	return r.GetDomain(ctx, id)
}

func (r *PostgresAdminRepository) fetchUserGroupIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT group_id FROM user_groups WHERE user_id = $1 ORDER BY group_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) fetchGroupPermissionIDs(ctx context.Context, groupID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT permission_id FROM group_permissions WHERE group_id = $1 ORDER BY permission_id`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) fetchDomainReviews(ctx context.Context, domainID string) ([]model.ReviewComment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, author, COALESCE(author_id, ''), decision, comment, created_at
		FROM domain_reviews
		WHERE domain_id = $1
		ORDER BY created_at
	`, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.ReviewComment, 0)
	for rows.Next() {
		var item model.ReviewComment
		if err := rows.Scan(&item.ID, &item.Author, &item.AuthorID, &item.Decision, &item.Comment, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func replaceUserGroups(ctx context.Context, tx pgx.Tx, userID string, groupIDs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM user_groups WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return insertRelationIDs(ctx, tx, `INSERT INTO user_groups (user_id, group_id) VALUES ($1, $2)`, userID, groupIDs)
}

func replaceGroupPermissions(ctx context.Context, tx pgx.Tx, groupID string, permissionIDs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM group_permissions WHERE group_id = $1`, groupID); err != nil {
		return err
	}
	return insertRelationIDs(ctx, tx, `INSERT INTO group_permissions (group_id, permission_id) VALUES ($1, $2)`, groupID, permissionIDs)
}

func insertRelationIDs(ctx context.Context, tx pgx.Tx, query, primaryID string, relatedIDs []string) error {
	ordered := append([]string(nil), relatedIDs...)
	sort.Strings(ordered)
	for _, relatedID := range ordered {
		if _, err := tx.Exec(ctx, query, primaryID, relatedID); err != nil {
			return err
		}
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanDomain(s scanner) (model.RegisteredDomain, error) {
	var item model.RegisteredDomain
	var metadataRaw []byte
	if err := s.Scan(&item.ID, &item.Name, &item.Description, &item.OwnerGroupID, &item.Status, &metadataRaw); err != nil {
		return model.RegisteredDomain{}, err
	}
	item.Metadata = map[string]any{}
	if len(metadataRaw) > 0 {
		if err := json.Unmarshal(metadataRaw, &item.Metadata); err != nil {
			return model.RegisteredDomain{}, fmt.Errorf("decode domain metadata: %w", err)
		}
	}
	return item, nil
}

func mapPostgresError(err error, entity string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperrors.NewConflict(entity + " already exists")
		case "23503":
			return apperrors.NewValidation(entity + " references a missing related record")
		}
	}
	return err
}
