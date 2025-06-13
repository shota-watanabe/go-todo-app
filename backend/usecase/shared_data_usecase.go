package usecase

import (
	"errors"
	"time"
	"go-todo-app/backend/domain" // your_module_nameをgo.modのmodule名に置き換え
)

// SharedDataUseCase は共有データとプロジェクトデータの操作を定義するインターフェース
type SharedDataUseCase interface {
	// GroupCompany関連
	GetViewGroupCompanies() ([]domain.ViewGroupCompany, error)

	// SharedProduct関連
	GetSharedProducts(orgID uint) ([]domain.SharedProduct, error)
	// TODO: 親オーガナイゼーション管理者によるSharedProductのCRUD (Create/Update/Delete)
	CreateSharedProduct(product *domain.SharedProduct) (*domain.SharedProduct, error) // SharedProductの作成を追加

	// Project関連
	CreateProject(name string, orgID uint) (*domain.Project, error)
	GetProjectByID(projectID uint) (*domain.Project, error)
	GetProjectsByOrganizationID(orgID uint) ([]domain.Project, error)
	DeleteProject(projectID uint) error // 論理削除

	// ProjectProduct関連
	CopySharedProductToProject(sharedProductID, projectID uint) (*domain.ProjectProduct, error)
	GetProjectProducts(projectID uint) ([]domain.ProjectProduct, error)
	// TODO: ProjectProductのCRUD (プロジェクト内で完結)
	PromoteProjectProductToShared(projectProductID uint) (*domain.SharedProduct, error) // プロジェクトのデータを共有データに格上げ
}

// sharedDataInteractor は SharedDataUseCase の実装
type sharedDataInteractor struct {
	// ここで定義されたインターフェース型を参照するように変更
	viewGroupCompanyRepo ViewGroupCompanyRepository
	sharedProductRepo    SharedProductRepository
	projectRepo          ProjectRepository
	projectProductRepo   ProjectProductRepository
	organizationRepo     OrganizationRepository
}

// NewSharedDataUseCase は SharedDataUseCase を生成します
func NewSharedDataUseCase(
	vgRepo ViewGroupCompanyRepository, // ここもusecaseパッケージのインターフェースに
	spRepo SharedProductRepository,
	projRepo ProjectRepository,
	ppRepo ProjectProductRepository,
	orgRepo OrganizationRepository,
) SharedDataUseCase {
	return &sharedDataInteractor{
		viewGroupCompanyRepo: vgRepo,
		sharedProductRepo:    spRepo,
		projectRepo:          projRepo,
		projectProductRepo:   ppRepo,
		organizationRepo:     orgRepo,
	}
}

// GetViewGroupCompanies は全てのグループ会社を取得します
func (interactor *sharedDataInteractor) GetViewGroupCompanies() ([]domain.ViewGroupCompany, error) {
	// TODO: ここでブラックリストのチェックを行う（OrganizationRepositoryと連携）
	return interactor.viewGroupCompanyRepo.GetAll()
}

// GetSharedProducts は指定されたオーガナイゼーションの共有商品を取得します
func (interactor *sharedDataInteractor) GetSharedProducts(orgID uint) ([]domain.SharedProduct, error) {
	// TODO: ここでブラックリストのチェックを行う（OrganizationRepositoryと連携）
	return interactor.sharedProductRepo.FindByOrganizationID(orgID)
}

// CreateSharedProduct は新しい共有商品を作成します
func (interactor *sharedDataInteractor) CreateSharedProduct(product *domain.SharedProduct) (*domain.SharedProduct, error) {
	return interactor.sharedProductRepo.Create(product)
}


// CreateProject は新しいプロジェクトを作成します
func (interactor *sharedDataInteractor) CreateProject(name string, orgID uint) (*domain.Project, error) {
	project := &domain.Project{
		Name:          name,
		OrganizationID: orgID,
	}
	return interactor.projectRepo.Create(project)
}

// GetProjectByID は指定されたIDのプロジェクトを取得します
func (interactor *sharedDataInteractor) GetProjectByID(projectID uint) (*domain.Project, error) {
	return interactor.projectRepo.FindByID(projectID)
}

// GetProjectsByOrganizationID は指定されたオーガナイゼーションのプロジェクトを取得します
func (interactor *sharedDataInteractor) GetProjectsByOrganizationID(orgID uint) ([]domain.Project, error) {
	return interactor.projectRepo.FindByOrganizationID(orgID)
}

// DeleteProject はプロジェクトを論理削除します
func (interactor *sharedDataInteractor) DeleteProject(projectID uint) error {
	project, err := interactor.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	now := time.Now()
	project.DeletedAt = &now // 論理削除フラグを設定
	return interactor.projectRepo.Update(project)
}


// CopySharedProductToProject は共有商品をプロジェクトにコピーします
func (interactor *sharedDataInteractor) CopySharedProductToProject(sharedProductID, projectID uint) (*domain.ProjectProduct, error) {
	// 共有商品が存在するか確認
	sharedProduct, err := interactor.sharedProductRepo.FindByID(sharedProductID)
	if err != nil {
		return nil, err
	}
	if sharedProduct == nil {
		return nil, errors.New("shared product not found")
	}

	// プロジェクトが存在するか確認
	project, err := interactor.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil || project.DeletedAt != nil { // 論理削除されているプロジェクトは対象外
		return nil, errors.New("project not found or already deleted")
	}

	// プロジェクト固有の商品としてコピー
	projectProduct := &domain.ProjectProduct{
		ProjectID:            projectID,
		Name:                 sharedProduct.Name,
		Description:          sharedProduct.Description,
		Price:                sharedProduct.Price,
		SharedProductID: &sharedProduct.ID,
	}

	return interactor.projectProductRepo.Create(projectProduct)
}

// GetProjectProducts は指定されたプロジェクトの全ての商品を取得します
func (interactor *sharedDataInteractor) GetProjectProducts(projectID uint) ([]domain.ProjectProduct, error) {
	// TODO: プロジェクトがブラックリストに入っていないかチェック
	return interactor.projectProductRepo.FindByProjectID(projectID)
}

// PromoteProjectProductToShared はプロジェクトの商品を共有商品に格上げします
func (interactor *sharedDataInteractor) PromoteProjectProductToShared(projectProductID uint) (*domain.SharedProduct, error) {
	// TODO: この操作は親オーガナイゼーションのAdminアカウントのみが許可されるようにする
	// プロジェクトの商品を取得
	projectProduct, err := interactor.projectProductRepo.FindByID(projectProductID)
	if err != nil {
		return nil, err
	}
	if projectProduct == nil {
		return nil, errors.New("project product not found")
	}

	// 所属プロジェクトのオーガナイゼーションIDを取得
	project, err := interactor.projectRepo.FindByID(projectProduct.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found for project product")
	}

	// 新しい共有商品として登録
	sharedProduct := &domain.SharedProduct{
		Name:           projectProduct.Name,
		Description:    projectProduct.Description,
		Price:          projectProduct.Price,
		OrganizationID: project.OrganizationID, // プロジェクトが所属するオーガナイゼーションのID
		// SKUはユニーク制約があるため、ここでは自動生成または入力させる必要がある
		// 簡単のため、ここでは仮に空文字列にするが、本番ではロジックが必要
		SKU: "AUTO_GENERATED_SKU_" + time.Now().Format("20060102150405"),
	}

	// SharedProductRepository の Create メソッドを呼び出す
	createdSharedProduct, err := interactor.sharedProductRepo.Create(sharedProduct)
	if err != nil {
		return nil, errors.New("failed to promote product to shared: " + err.Error())
	}

	return createdSharedProduct, nil
}