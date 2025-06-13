// frontend/lib/sharedDataApi.ts

// バックエンドのベースURL
const API_BASE_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

// APIレスポンスの共通型定義
interface ApiResponse<T> {
  data?: T;
  error?: string;
}

// データモデルの型定義 (バックエンドのGoのDomain層と一致させる)
export interface ViewGroupCompany {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface SharedProduct {
  id: number;
  name: string;
  description: string;
  price: number;
  sku: string;
  organization_id: number;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: number;
  name: string;
  organization_id: number;
  deleted_at: string | null; // 論理削除のためnull許容
  created_at: string;
  updated_at: string;
}

export interface ProjectProduct {
  id: number;
  project_id: number;
  name: string;
  description: string;
  price: number;
  original_shared_product_id: number | null; // コピー元がない場合もあるためnull許容
  created_at: string;
  updated_at: string;
}

// Helper function for API calls
async function fetchApi<T>(
  url: string,
  method: string = "GET",
  token?: string,
  body?: any
): Promise<ApiResponse<T>> {
  try {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const errorData = await response.json();
      return { error: errorData.error || response.statusText };
    }

    const data: T = await response.json();
    return { data };
  } catch (err: any) {
    console.error(`API call error to ${url}:`, err);
    return { error: err.message || "An unknown error occurred" };
  }
}

// --- API Functions ---

// 参照のみのグループ会社一覧を取得
export const getViewGroupCompanies = async (
  token?: string
): Promise<ApiResponse<ViewGroupCompany[]>> => {
  return fetchApi<ViewGroupCompany[]>("/shared/group-companies", "GET", token);
};

// 共有商品一覧を取得
// TODO: バックエンドのorganization_idの取得方法に合わせて調整
export const getSharedProducts = async (
  organizationId: number,
  token?: string
): Promise<ApiResponse<SharedProduct[]>> => {
  return fetchApi<SharedProduct[]>(
    `/shared/products?organization_id=${organizationId}`,
    "GET",
    token
  );
};

// プロジェクトを作成
export const createProject = async (
  name: string,
  token?: string
): Promise<ApiResponse<Project>> => {
  return fetchApi<Project>("/projects", "POST", token, { name });
};

// プロジェクト一覧を取得 (現時点ではバックエンドにエンドポイントがないため、組織IDでフィルタリングする想定)
// TODO: バックエンドでプロジェクト一覧API (GetProjectsByOrganizationID) が公開されたら修正
export const getProjects = async (
  organizationId: number,
  token?: string
): Promise<ApiResponse<Project[]>> => {
  // 現状バックエンドに全件取得APIがないため、一旦ダミーデータを返すか、
  // もしくはこの関数内で /projects/:projectID に繰り返しリクエストを出すなどの対応が必要
  // 今回は一旦、バックエンドのGetProjectsByOrganizationIDに対応する形で実装を想定
  return fetchApi<Project[]>(
    `/projects?organization_id=${organizationId}`, // バックエンドにこのようなエンドポイントがあることを想定
    "GET",
    token
  );
};


// 共有商品をプロジェクトにコピー
export const copySharedProductToProject = async (
  projectId: number,
  sharedProductId: number,
  token?: string
): Promise<ApiResponse<ProjectProduct>> => {
  return fetchApi<ProjectProduct>(
    `/projects/${projectId}/products/copy`,
    "POST",
    token,
    { shared_product_id: sharedProductId }
  );
};

// 特定プロジェクトの商品一覧を取得
export const getProjectProducts = async (
  projectId: number,
  token?: string
): Promise<ApiResponse<ProjectProduct[]>> => {
  return fetchApi<ProjectProduct[]>(
    `/projects/${projectId}/products`,
    "GET",
    token
  );
};

// プロジェクトの商品を共有データに格上げ
export const promoteProjectProductToShared = async (
  projectProductId: number,
  token?: string
): Promise<ApiResponse<SharedProduct>> => {
  return fetchApi<SharedProduct>(
    `/projects/products/${projectProductId}/promote`,
    "POST",
    token
  );
};