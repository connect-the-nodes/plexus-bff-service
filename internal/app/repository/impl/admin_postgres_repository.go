package impl

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"plexus-bff-service-go/internal/app/apperrors"
	"plexus-bff-service-go/internal/app/config"
	"plexus-bff-service-go/internal/app/model"
	"plexus-bff-service-go/internal/app/repository"
)

const (
	praxisStatusDraft           = "DRAFT"
	praxisStatusPendingApproval = "PENDING_APPROVAL"
	praxisStatusApproved        = "APPROVED"
	praxisStatusRejected        = "REJECTED"
	praxisPermissionCategory    = "Custom"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type PostgresAdminRepository struct {
	pool              *pgxpool.Pool
	schema            string
	customerAccountID int64
	defaultActorID    int64
	defaultActorName  string
}

func NewPostgresAdminRepository(cfg config.DatabaseConfig, pool *pgxpool.Pool) (*PostgresAdminRepository, error) {
	schema := strings.TrimSpace(cfg.Schema)
	if schema == "" {
		schema = "public"
	}
	if !identifierPattern.MatchString(schema) {
		return nil, fmt.Errorf("invalid database schema name %q", schema)
	}

	return &PostgresAdminRepository{
		pool:              pool,
		schema:            schema,
		customerAccountID: cfg.CustomerAccountID,
		defaultActorID:    cfg.DefaultActorID,
		defaultActorName:  cfg.DefaultActorName,
	}, nil
}

func (r *PostgresAdminRepository) Users() repository.UserRepository             { return r }
func (r *PostgresAdminRepository) Permissions() repository.PermissionRepository { return r }
func (r *PostgresAdminRepository) Groups() repository.GroupRepository           { return r }
func (r *PostgresAdminRepository) Domains() repository.DomainRepository         { return r }

func (r *PostgresAdminRepository) ListUsers(ctx context.Context) ([]model.PortalUser, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT admin_user_id, username, password_value, COALESCE(email_address, ''), COALESCE(display_name, ''), portal_role::text, is_active
		FROM %s
		ORDER BY username
	`, r.table("admin_user")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.PortalUser, 0)
	for rows.Next() {
		var (
			item   model.PortalUser
			dbID   int64
			role   string
			groups []string
		)
		if err := rows.Scan(&dbID, &item.Username, &item.PasswordHash, &item.Email, &item.DisplayName, &role, &item.Active); err != nil {
			return nil, err
		}
		groups, err = r.fetchUserGroupIDs(ctx, dbID)
		if err != nil {
			return nil, err
		}
		item.ID = formatID(dbID)
		item.Role = role
		item.GroupIDs = groups
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetUser(ctx context.Context, id string) (model.PortalUser, error) {
	dbID, err := parseID(id, "user")
	if err != nil {
		return model.PortalUser{}, err
	}

	var (
		item model.PortalUser
		role string
	)
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT username, password_value, COALESCE(email_address, ''), COALESCE(display_name, ''), portal_role::text, is_active
		FROM %s
		WHERE admin_user_id = $1
	`, r.table("admin_user")), dbID).
		Scan(&item.Username, &item.PasswordHash, &item.Email, &item.DisplayName, &role, &item.Active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PortalUser{}, apperrors.NewNotFound("user not found")
		}
		return model.PortalUser{}, err
	}

	groupIDs, err := r.fetchUserGroupIDs(ctx, dbID)
	if err != nil {
		return model.PortalUser{}, err
	}

	item.ID = id
	item.Role = role
	item.GroupIDs = groupIDs
	return item, nil
}

func (r *PostgresAdminRepository) CreateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalUser{}, err
	}
	defer tx.Rollback(ctx)

	var dbID int64
	err = tx.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			customer_account_id,
			external_id,
			username,
			password_value,
			display_name,
			email_address,
			portal_role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7::%s, $8)
		RETURNING admin_user_id
	`, r.table("admin_user"), r.enumType("portal_user_role")),
		r.customerAccountID,
		externalID("user", user.Username),
		user.Username,
		user.PasswordHash,
		coalesceString(user.DisplayName, user.Username),
		user.Email,
		normalizePortalRole(user.Role),
		user.Active,
	).Scan(&dbID)
	if err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user")
	}

	if err := replaceUserGroups(ctx, tx, r.table("admin_user_group"), dbID, user.GroupIDs); err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user group assignment")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PortalUser{}, err
	}
	return r.GetUser(ctx, formatID(dbID))
}

func (r *PostgresAdminRepository) UpdateUser(ctx context.Context, user model.PortalUser) (model.PortalUser, error) {
	dbID, err := parseID(user.ID, "user")
	if err != nil {
		return model.PortalUser{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalUser{}, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET username = $2,
			password_value = $3,
			display_name = $4,
			email_address = NULLIF($5, ''),
			portal_role = $6::%s,
			is_active = $7,
			updated_at = timezone('utc', now())
		WHERE admin_user_id = $1
	`, r.table("admin_user"), r.enumType("portal_user_role")),
		dbID,
		user.Username,
		user.PasswordHash,
		coalesceString(user.DisplayName, user.Username),
		user.Email,
		normalizePortalRole(user.Role),
		user.Active,
	)
	if err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user")
	}
	if tag.RowsAffected() == 0 {
		return model.PortalUser{}, apperrors.NewNotFound("user not found")
	}

	if err := replaceUserGroups(ctx, tx, r.table("admin_user_group"), dbID, user.GroupIDs); err != nil {
		return model.PortalUser{}, mapPostgresError(err, "user group assignment")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PortalUser{}, err
	}
	return r.GetUser(ctx, user.ID)
}

func (r *PostgresAdminRepository) DeleteUser(ctx context.Context, id string) error {
	dbID, err := parseID(id, "user")
	if err != nil {
		return err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE admin_user_id = $1`, r.table("admin_user")), dbID)
	if err != nil {
		return mapPostgresError(err, "user")
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("user not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListPermissions(ctx context.Context) ([]model.PermissionDefinition, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT admin_permission_id, permission_name, COALESCE(permission_description, '')
		FROM %s
		ORDER BY permission_name
	`, r.table("admin_permission")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.PermissionDefinition, 0)
	for rows.Next() {
		var item model.PermissionDefinition
		var dbID int64
		if err := rows.Scan(&dbID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		item.ID = formatID(dbID)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetPermission(ctx context.Context, id string) (model.PermissionDefinition, error) {
	dbID, err := parseID(id, "permission")
	if err != nil {
		return model.PermissionDefinition{}, err
	}

	var item model.PermissionDefinition
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT permission_name, COALESCE(permission_description, '')
		FROM %s
		WHERE admin_permission_id = $1
	`, r.table("admin_permission")), dbID).
		Scan(&item.Name, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
		}
		return model.PermissionDefinition{}, err
	}
	item.ID = id
	return item, nil
}

func (r *PostgresAdminRepository) CreatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	var dbID int64
	err := r.pool.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (permission_code, permission_name, permission_description, permission_category)
		VALUES ($1, $2, NULLIF($3, ''), $4)
		RETURNING admin_permission_id
	`, r.table("admin_permission")),
		permissionCode(permission.Name),
		permission.Name,
		permission.Description,
		praxisPermissionCategory,
	).Scan(&dbID)
	if err != nil {
		return model.PermissionDefinition{}, mapPostgresError(err, "permission")
	}
	return r.GetPermission(ctx, formatID(dbID))
}

func (r *PostgresAdminRepository) UpdatePermission(ctx context.Context, permission model.PermissionDefinition) (model.PermissionDefinition, error) {
	dbID, err := parseID(permission.ID, "permission")
	if err != nil {
		return model.PermissionDefinition{}, err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET permission_name = $2,
			permission_description = NULLIF($3, ''),
			updated_at = timezone('utc', now())
		WHERE admin_permission_id = $1
	`, r.table("admin_permission")), dbID, permission.Name, permission.Description)
	if err != nil {
		return model.PermissionDefinition{}, mapPostgresError(err, "permission")
	}
	if tag.RowsAffected() == 0 {
		return model.PermissionDefinition{}, apperrors.NewNotFound("permission not found")
	}
	return r.GetPermission(ctx, permission.ID)
}

func (r *PostgresAdminRepository) DeletePermission(ctx context.Context, id string) error {
	dbID, err := parseID(id, "permission")
	if err != nil {
		return err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE admin_permission_id = $1`, r.table("admin_permission")), dbID)
	if err != nil {
		return mapPostgresError(err, "permission")
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("permission not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListGroups(ctx context.Context) ([]model.PortalGroup, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT admin_group_id, group_name, COALESCE(group_description, '')
		FROM %s
		ORDER BY group_name
	`, r.table("admin_group")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.PortalGroup, 0)
	for rows.Next() {
		var item model.PortalGroup
		var dbID int64
		if err := rows.Scan(&dbID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		permissionIDs, err := r.fetchGroupPermissionIDs(ctx, dbID)
		if err != nil {
			return nil, err
		}
		item.ID = formatID(dbID)
		item.PermissionIDs = permissionIDs
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetGroup(ctx context.Context, id string) (model.PortalGroup, error) {
	dbID, err := parseID(id, "group")
	if err != nil {
		return model.PortalGroup{}, err
	}

	var item model.PortalGroup
	err = r.pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT group_name, COALESCE(group_description, '')
		FROM %s
		WHERE admin_group_id = $1
	`, r.table("admin_group")), dbID).
		Scan(&item.Name, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PortalGroup{}, apperrors.NewNotFound("group not found")
		}
		return model.PortalGroup{}, err
	}

	permissionIDs, err := r.fetchGroupPermissionIDs(ctx, dbID)
	if err != nil {
		return model.PortalGroup{}, err
	}

	item.ID = id
	item.PermissionIDs = permissionIDs
	return item, nil
}

func (r *PostgresAdminRepository) CreateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalGroup{}, err
	}
	defer tx.Rollback(ctx)

	slug := slugify(group.Name, "-")
	var dbID int64
	err = tx.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			customer_account_id,
			external_id,
			group_name,
			group_slug,
			group_description
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		RETURNING admin_group_id
	`, r.table("admin_group")),
		r.customerAccountID,
		externalID("group", group.Name),
		group.Name,
		slug,
		group.Description,
	).Scan(&dbID)
	if err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group")
	}

	if err := replaceGroupPermissions(ctx, tx, r.table("admin_group_permission"), dbID, group.PermissionIDs); err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group permission assignment")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PortalGroup{}, err
	}
	return r.GetGroup(ctx, formatID(dbID))
}

func (r *PostgresAdminRepository) UpdateGroup(ctx context.Context, group model.PortalGroup) (model.PortalGroup, error) {
	dbID, err := parseID(group.ID, "group")
	if err != nil {
		return model.PortalGroup{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.PortalGroup{}, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET group_name = $2,
			group_slug = $3,
			group_description = NULLIF($4, ''),
			updated_at = timezone('utc', now())
		WHERE admin_group_id = $1
	`, r.table("admin_group")),
		dbID,
		group.Name,
		slugify(group.Name, "-"),
		group.Description,
	)
	if err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group")
	}
	if tag.RowsAffected() == 0 {
		return model.PortalGroup{}, apperrors.NewNotFound("group not found")
	}

	if err := replaceGroupPermissions(ctx, tx, r.table("admin_group_permission"), dbID, group.PermissionIDs); err != nil {
		return model.PortalGroup{}, mapPostgresError(err, "group permission assignment")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PortalGroup{}, err
	}
	return r.GetGroup(ctx, group.ID)
}

func (r *PostgresAdminRepository) DeleteGroup(ctx context.Context, id string) error {
	dbID, err := parseID(id, "group")
	if err != nil {
		return err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE admin_group_id = $1`, r.table("admin_group")), dbID)
	if err != nil {
		return mapPostgresError(err, "group")
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("group not found")
	}
	return nil
}

func (r *PostgresAdminRepository) ListDomains(ctx context.Context) ([]model.RegisteredDomain, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT domain_id, domain_name, COALESCE(domain_description, ''), owner_group_id, lifecycle_status::text
		FROM %s
		ORDER BY domain_name
	`, r.table("domain")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.RegisteredDomain, 0)
	for rows.Next() {
		item, dbID, err := r.scanDomain(rows)
		if err != nil {
			return nil, err
		}
		reviews, err := r.fetchDomainReviews(ctx, dbID)
		if err != nil {
			return nil, err
		}
		item.Review = reviews
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) ListApprovedDomains(ctx context.Context) ([]model.RegisteredDomain, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT domain_id, domain_name, COALESCE(domain_description, ''), owner_group_id, lifecycle_status::text
		FROM %s
		WHERE lifecycle_status = $1::%s
		ORDER BY domain_name
	`, r.table("domain"), r.enumType("lifecycle_status")), praxisStatusApproved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.RegisteredDomain, 0)
	for rows.Next() {
		item, dbID, err := r.scanDomain(rows)
		if err != nil {
			return nil, err
		}
		reviews, err := r.fetchDomainReviews(ctx, dbID)
		if err != nil {
			return nil, err
		}
		item.Review = reviews
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) GetDomain(ctx context.Context, id string) (model.RegisteredDomain, error) {
	dbID, err := parseID(id, "domain")
	if err != nil {
		return model.RegisteredDomain{}, err
	}

	row := r.pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT domain_id, domain_name, COALESCE(domain_description, ''), owner_group_id, lifecycle_status::text
		FROM %s
		WHERE domain_id = $1
	`, r.table("domain")), dbID)
	item, _, err := r.scanDomain(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
		}
		return model.RegisteredDomain{}, err
	}

	reviews, err := r.fetchDomainReviews(ctx, dbID)
	if err != nil {
		return model.RegisteredDomain{}, err
	}
	item.Review = reviews
	return item, nil
}

func (r *PostgresAdminRepository) CreateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	var dbID int64
	err := r.pool.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			customer_account_id,
			external_id,
			domain_name,
			domain_code,
			domain_description,
			owner_group_id,
			owner_role_name,
			created_by_user_id,
			created_by_display_name,
			lifecycle_status,
			submitted_at
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9, $10::%s, $11)
		RETURNING domain_id
	`, r.table("domain"), r.enumType("lifecycle_status")),
		r.customerAccountID,
		externalID("domain", domain.Name),
		domain.Name,
		domainCode(domain.Name),
		domain.Description,
		mustParseID(domain.OwnerGroupID),
		ownerRoleName(domain.Name),
		r.nullableActorID(),
		r.actorName(),
		normalizeDomainStatus(domain.Status),
		submittedAt(normalizeDomainStatus(domain.Status)),
	).Scan(&dbID)
	if err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain")
	}
	return r.GetDomain(ctx, formatID(dbID))
}

func (r *PostgresAdminRepository) UpdateDomain(ctx context.Context, domain model.RegisteredDomain) (model.RegisteredDomain, error) {
	dbID, err := parseID(domain.ID, "domain")
	if err != nil {
		return model.RegisteredDomain{}, err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET domain_name = $2,
			domain_code = $3,
			domain_description = NULLIF($4, ''),
			owner_group_id = $5,
			owner_role_name = $6,
			lifecycle_status = $7::%s,
			submitted_at = $8,
			updated_at = timezone('utc', now())
		WHERE domain_id = $1
	`, r.table("domain"), r.enumType("lifecycle_status")),
		dbID,
		domain.Name,
		domainCode(domain.Name),
		domain.Description,
		mustParseID(domain.OwnerGroupID),
		ownerRoleName(domain.Name),
		normalizeDomainStatus(domain.Status),
		submittedAt(normalizeDomainStatus(domain.Status)),
	)
	if err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain")
	}
	if tag.RowsAffected() == 0 {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}
	return r.GetDomain(ctx, domain.ID)
}

func (r *PostgresAdminRepository) DeleteDomain(ctx context.Context, id string) error {
	dbID, err := parseID(id, "domain")
	if err != nil {
		return err
	}

	tag, err := r.pool.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE domain_id = $1`, r.table("domain")), dbID)
	if err != nil {
		return mapPostgresError(err, "domain")
	}
	if tag.RowsAffected() == 0 {
		return apperrors.NewNotFound("domain not found")
	}
	return nil
}

func (r *PostgresAdminRepository) AddDomainReview(ctx context.Context, id string, comment model.ReviewComment, nextStatus string) (model.RegisteredDomain, error) {
	dbID, err := parseID(id, "domain")
	if err != nil {
		return model.RegisteredDomain{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.RegisteredDomain{}, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, fmt.Sprintf(`
		UPDATE %s
		SET lifecycle_status = $2::%s,
			updated_at = timezone('utc', now())
		WHERE domain_id = $1
	`, r.table("domain"), r.enumType("lifecycle_status")), dbID, normalizeDomainStatus(nextStatus))
	if err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain")
	}
	if tag.RowsAffected() == 0 {
		return model.RegisteredDomain{}, apperrors.NewNotFound("domain not found")
	}

	if _, err := tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (
			domain_id,
			external_id,
			author_user_id,
			author_display_name,
			review_decision,
			comment_text,
			reviewed_at
		)
		VALUES ($1, $2, $3, $4, $5::%s, $6, $7)
	`, r.table("domain_review_comment"), r.enumType("review_decision")),
		dbID,
		externalID("comment", comment.ID),
		nullableParsedID(comment.AuthorID),
		comment.Author,
		normalizeDomainStatus(comment.Decision),
		comment.Comment,
		comment.CreatedAt.UTC(),
	); err != nil {
		return model.RegisteredDomain{}, mapPostgresError(err, "domain review")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.RegisteredDomain{}, err
	}
	return r.GetDomain(ctx, id)
}

func (r *PostgresAdminRepository) fetchUserGroupIDs(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT admin_group_id
		FROM %s
		WHERE admin_user_id = $1
		ORDER BY admin_group_id
	`, r.table("admin_user_group")), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, formatID(id))
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) fetchGroupPermissionIDs(ctx context.Context, groupID int64) ([]string, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT admin_permission_id
		FROM %s
		WHERE admin_group_id = $1
		ORDER BY admin_permission_id
	`, r.table("admin_group_permission")), groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, formatID(id))
	}
	return items, rows.Err()
}

func (r *PostgresAdminRepository) fetchDomainReviews(ctx context.Context, domainID int64) ([]model.ReviewComment, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT domain_review_comment_id, author_display_name, COALESCE(author_user_id::text, ''), review_decision::text, comment_text, reviewed_at
		FROM %s
		WHERE domain_id = $1
		ORDER BY reviewed_at, domain_review_comment_id
	`, r.table("domain_review_comment")), domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.ReviewComment, 0)
	for rows.Next() {
		var item model.ReviewComment
		var dbID int64
		if err := rows.Scan(&dbID, &item.Author, &item.AuthorID, &item.Decision, &item.Comment, &item.CreatedAt); err != nil {
			return nil, err
		}
		item.ID = formatID(dbID)
		items = append(items, item)
	}
	return items, rows.Err()
}

func replaceUserGroups(ctx context.Context, tx pgx.Tx, table string, userID int64, groupIDs []string) error {
	if _, err := tx.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE admin_user_id = $1`, table), userID); err != nil {
		return err
	}
	ordered, err := sortAndParseIDs(groupIDs, "group")
	if err != nil {
		return err
	}
	for _, groupID := range ordered {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (admin_user_id, admin_group_id)
			VALUES ($1, $2)
		`, table), userID, groupID); err != nil {
			return err
		}
	}
	return nil
}

func replaceGroupPermissions(ctx context.Context, tx pgx.Tx, table string, groupID int64, permissionIDs []string) error {
	if _, err := tx.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE admin_group_id = $1`, table), groupID); err != nil {
		return err
	}
	ordered, err := sortAndParseIDs(permissionIDs, "permission")
	if err != nil {
		return err
	}
	for _, permissionID := range ordered {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (admin_group_id, admin_permission_id)
			VALUES ($1, $2)
		`, table), groupID, permissionID); err != nil {
			return err
		}
	}
	return nil
}

type domainScanner interface {
	Scan(dest ...any) error
}

func (r *PostgresAdminRepository) scanDomain(s domainScanner) (model.RegisteredDomain, int64, error) {
	var (
		item       model.RegisteredDomain
		dbID       int64
		ownerGroup int64
		status     string
	)
	if err := s.Scan(&dbID, &item.Name, &item.Description, &ownerGroup, &status); err != nil {
		return model.RegisteredDomain{}, 0, err
	}
	item.ID = formatID(dbID)
	item.OwnerGroupID = formatID(ownerGroup)
	item.Status = normalizeDomainStatus(status)
	item.Metadata = map[string]any{}
	item.Review = []model.ReviewComment{}
	return item, dbID, nil
}

func (r *PostgresAdminRepository) table(name string) string {
	return quoteIdent(r.schema) + "." + quoteIdent(name)
}

func (r *PostgresAdminRepository) enumType(name string) string {
	return quoteIdent(r.schema) + "." + quoteIdent(name)
}

func (r *PostgresAdminRepository) nullableActorID() any {
	if r.defaultActorID == 0 {
		return nil
	}
	return r.defaultActorID
}

func (r *PostgresAdminRepository) actorName() string {
	if strings.TrimSpace(r.defaultActorName) == "" {
		return "System"
	}
	return r.defaultActorName
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func parseID(value, entity string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed <= 0 {
		return 0, apperrors.NewValidation(entity + " id must be a positive integer")
	}
	return parsed, nil
}

func mustParseID(value string) int64 {
	parsed, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return parsed
}

func nullableParsedID(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return nil
	}
	return parsed
}

func sortAndParseIDs(values []string, entity string) ([]int64, error) {
	parsed := make([]int64, 0, len(values))
	for _, value := range values {
		id, err := parseID(value, entity)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, id)
	}
	sort.Slice(parsed, func(i, j int) bool { return parsed[i] < parsed[j] })
	return parsed, nil
}

func formatID(value int64) string {
	return strconv.FormatInt(value, 10)
}

func normalizePortalRole(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "ADMIN":
		return "ADMIN"
	default:
		return "USER"
	}
}

func normalizeDomainStatus(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case praxisStatusApproved:
		return praxisStatusApproved
	case praxisStatusRejected:
		return praxisStatusRejected
	case praxisStatusPendingApproval:
		return praxisStatusPendingApproval
	default:
		return praxisStatusDraft
	}
}

func permissionCode(name string) string {
	return "perm_" + slugify(name, "_")
}

func domainCode(name string) string {
	return slugify(name, "_")
}

func ownerRoleName(domainName string) string {
	base := strings.TrimSpace(domainName)
	if base == "" {
		return "Domain Approver"
	}
	return base + " Approver"
}

func externalID(prefix, seed string) string {
	slug := slugify(seed, "-")
	if slug == "" {
		slug = slugify(strconv.FormatInt(time.Now().UnixNano(), 10), "-")
	}
	return prefix + "-" + slug
}

func slugify(value, separator string) string {
	var builder strings.Builder
	lastWasSep := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		isLetter := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if isLetter || isDigit {
			builder.WriteRune(r)
			lastWasSep = false
			continue
		}
		if !lastWasSep && builder.Len() > 0 {
			builder.WriteString(separator)
			lastWasSep = true
		}
	}
	result := strings.Trim(builder.String(), separator)
	if result == "" {
		return "generated"
	}
	return result
}

func coalesceString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}

func submittedAt(status string) any {
	if status != praxisStatusPendingApproval {
		return nil
	}
	return time.Now().UTC()
}

func mapPostgresError(err error, entity string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperrors.NewConflict(entity + " already exists")
		case "23503":
			return apperrors.NewValidation(entity + " references a missing or protected record")
		case "22P02":
			return apperrors.NewValidation(entity + " contains an invalid value")
		}
	}
	return err
}
