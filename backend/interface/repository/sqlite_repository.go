package repository

import (
	"database/sql"
	"time"
	"go-todo-app/backend/domain"
)

// SQLiteViewGroupCompanyRepository は ViewGroupCompanyRepository の SQLite 実装
type SQLiteViewGroupCompanyRepository struct {
	db *sql.DB
}

// NewSQLiteViewGroupCompanyRepository は SQLiteViewGroupCompanyRepository を生成します
func NewSQLiteViewGroupCompanyRepository(db *sql.DB) *SQLiteViewGroupCompanyRepository {
	return &SQLiteViewGroupCompanyRepository{db: db}
}

// GetAll は全てのグループ会社を取得します
func (r *SQLiteViewGroupCompanyRepository) GetAll() ([]domain.ViewGroupCompany, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM view_group_companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []domain.ViewGroupCompany
	for rows.Next() {
		var c domain.ViewGroupCompany
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

// SQLiteSharedProductRepository は SharedProductRepository の SQLite 実装
type SQLiteSharedProductRepository struct {
	db *sql.DB
}

// NewSQLiteSharedProductRepository は SQLiteSharedProductRepository を生成します
func NewSQLiteSharedProductRepository(db *sql.DB) *SQLiteSharedProductRepository {
	return &SQLiteSharedProductRepository{db: db}
}

// FindByID はIDで共有商品を取得します
func (r *SQLiteSharedProductRepository) FindByID(id uint) (*domain.SharedProduct, error) {
	row := r.db.QueryRow("SELECT id, name, description, price, sku, organization_id, created_at, updated_at FROM shared_products WHERE id = ?", id)
	var p domain.SharedProduct
	var description sql.NullString
	var sku sql.NullString
	err := row.Scan(&p.ID, &p.Name, &description, &p.Price, &sku, &p.OrganizationID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Description = description.String
	p.SKU = sku.String
	return &p, nil
}

// FindByOrganizationID はオーガナイゼーションIDで共有商品を取得します
func (r *SQLiteSharedProductRepository) FindByOrganizationID(orgID uint) ([]domain.SharedProduct, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, sku, organization_id, created_at, updated_at FROM shared_products WHERE organization_id = ?", orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.SharedProduct
	for rows.Next() {
		var p domain.SharedProduct
		var description sql.NullString
		var sku sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &description, &p.Price, &sku, &p.OrganizationID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Description = description.String
		p.SKU = sku.String
		products = append(products, p)
	}
	return products, nil
}

// Create は新しい共有商品を作成します (TODO: Usecase層で呼び出されるべき)
func (r *SQLiteSharedProductRepository) Create(product *domain.SharedProduct) (*domain.SharedProduct, error) {
	stmt, err := r.db.Prepare("INSERT INTO shared_products (name, description, price, sku, organization_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(product.Name, product.Description, product.Price, product.SKU, product.OrganizationID, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	product.ID = uint(id)
	return product, nil
}

// SQLiteProjectRepository は ProjectRepository の SQLite 実装
type SQLiteProjectRepository struct {
	db *sql.DB
}

// NewSQLiteProjectRepository は SQLiteProjectRepository を生成します
func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

// Create は新しいプロジェクトを作成します
func (r *SQLiteProjectRepository) Create(project *domain.Project) (*domain.Project, error) {
	stmt, err := r.db.Prepare("INSERT INTO projects (name, organization_id, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(project.Name, project.OrganizationID, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	project.ID = uint(id)
	return project, nil
}

// FindByID はIDでプロジェクトを取得します
func (r *SQLiteProjectRepository) FindByID(id uint) (*domain.Project, error) {
	row := r.db.QueryRow("SELECT id, name, organization_id, deleted_at, created_at, updated_at FROM projects WHERE id = ?", id)
	var p domain.Project
	var deletedAt sql.NullTime
	err := row.Scan(&p.ID, &p.Name, &p.OrganizationID, &deletedAt, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}
	return &p, nil
}

// FindByOrganizationID はオーガナイゼーションIDでプロジェクトを取得します
func (r *SQLiteProjectRepository) FindByOrganizationID(orgID uint) ([]domain.Project, error) {
	rows, err := r.db.Query("SELECT id, name, organization_id, deleted_at, created_at, updated_at FROM projects WHERE organization_id = ?", orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var deletedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.OrganizationID, &deletedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Time
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// Update はプロジェクトを更新します
func (r *SQLiteProjectRepository) Update(project *domain.Project) error {
	stmt, err := r.db.Prepare("UPDATE projects SET name = ?, organization_id = ?, deleted_at = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var deletedAt sql.NullTime
	if project.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *project.DeletedAt, Valid: true}
	}

	_, err = stmt.Exec(project.Name, project.OrganizationID, deletedAt, time.Now(), project.ID)
	return err
}

// Delete はプロジェクトを論理削除します (実際にはUpdateでdeleted_atを設定)
// このメソッドはUsecase層のDeleteProjectで呼び出されることを想定
func (r *SQLiteProjectRepository) Delete(id uint) error {
	// 論理削除はUsecase層で行われるため、ここでは実際には呼ばれないか、
	// Update(project)の形で使われることを想定
	// ただし、インターフェースの要件を満たすためにダミー実装
	return nil
}


// SQLiteProjectProductRepository は ProjectProductRepository の SQLite 実装
type SQLiteProjectProductRepository struct {
	db *sql.DB
}

// NewSQLiteProjectProductRepository は SQLiteProjectProductRepository を生成します
func NewSQLiteProjectProductRepository(db *sql.DB) *SQLiteProjectProductRepository {
	return &SQLiteProjectProductRepository{db: db}
}

// Create は新しいプロジェクト商品を作成します
func (r *SQLiteProjectProductRepository) Create(product *domain.ProjectProduct) (*domain.ProjectProduct, error) {
	stmt, err := r.db.Prepare("INSERT INTO project_products (project_id, name, description, price, original_shared_product_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var originalSharedProductID sql.NullInt64
	if product.SharedProductID != nil {
		originalSharedProductID = sql.NullInt64{Int64: int64(*product.SharedProductID), Valid: true}
	}

	res, err := stmt.Exec(product.ProjectID, product.Name, product.Description, product.Price, originalSharedProductID, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	product.ID = uint(id)
	return product, nil
}

// FindByProjectID はプロジェクトIDでプロジェクト商品を取得します
func (r *SQLiteProjectProductRepository) FindByProjectID(projectID uint) ([]domain.ProjectProduct, error) {
	rows, err := r.db.Query("SELECT id, project_id, name, description, price, original_shared_product_id, created_at, updated_at FROM project_products WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.ProjectProduct
	for rows.Next() {
		var p domain.ProjectProduct
		var description sql.NullString
		var originalSharedProductID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.ProjectID, &p.Name, &description, &p.Price, &originalSharedProductID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Description = description.String
		if originalSharedProductID.Valid {
			id := uint(originalSharedProductID.Int64)
			p.SharedProductID = &id
		}
		products = append(products, p)
	}
	return products, nil
}

// FindByID はIDでプロジェクト商品を取得します
func (r *SQLiteProjectProductRepository) FindByID(id uint) (*domain.ProjectProduct, error) {
	row := r.db.QueryRow("SELECT id, project_id, name, description, price, original_shared_product_id, created_at, updated_at FROM project_products WHERE id = ?", id)
	var p domain.ProjectProduct
	var description sql.NullString
	var originalSharedProductID sql.NullInt64
	err := row.Scan(&p.ID, &p.ProjectID, &p.Name, &description, &p.Price, &originalSharedProductID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Description = description.String
	if originalSharedProductID.Valid {
		idVal := uint(originalSharedProductID.Int64)
		p.SharedProductID = &idVal
	}
	return &p, nil
}

// Update はプロジェクト商品を更新します
func (r *SQLiteProjectProductRepository) Update(product *domain.ProjectProduct) error {
	stmt, err := r.db.Prepare("UPDATE project_products SET name = ?, description = ?, price = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Price, time.Now(), product.ID)
	return err
}

// Delete はプロジェクト商品を物理削除します
func (r *SQLiteProjectProductRepository) Delete(id uint) error {
	stmt, err := r.db.Prepare("DELETE FROM project_products WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// SQLiteOrganizationRepository は OrganizationRepository の SQLite 実装
type SQLiteOrganizationRepository struct {
	db *sql.DB
}

// NewSQLiteOrganizationRepository は SQLiteOrganizationRepository を生成します
func NewSQLiteOrganizationRepository(db *sql.DB) *SQLiteOrganizationRepository {
	return &SQLiteOrganizationRepository{db: db}
}

// FindByID はIDで組織を取得します
func (r *SQLiteOrganizationRepository) FindByID(id uint) (*domain.Organization, error) {
	row := r.db.QueryRow("SELECT id, name, parent_id, created_at, updated_at FROM organizations WHERE id = ?", id)
	var o domain.Organization
	var parentID sql.NullInt64
	err := row.Scan(&o.ID, &o.Name, &parentID, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		idVal := uint(parentID.Int64)
		o.ParentID = &idVal
	}
	return &o, nil
}