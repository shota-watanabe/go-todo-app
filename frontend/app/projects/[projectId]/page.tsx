// frontend/app/projects/[projectId]/page.tsx
"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { getProjectProducts, promoteProjectProductToShared, ProjectProduct } from '@/lib/sharedDataApi';

interface ProjectDetailProps {
  params: {
    projectId: string; // URLパスから取得されるprojectId
  };
}

export default function ProjectDetailPage({ params }: ProjectDetailProps) {
  const { user, token, isLoading } = useAuth();
  const router = useRouter();
  const projectId = Number(params.projectId);
  const [projectProducts, setProjectProducts] = useState<ProjectProduct[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loadingData, setLoadingData] = useState(true);
  const [promotingProductId, setPromotingProductId] = useState<number | null>(null);
  const [promoteSuccessMessage, setPromoteSuccessMessage] = useState<string | null>(null);


  const fetchProjectProducts = async () => {
    setLoadingData(true);
    setError(null);
    setPromoteSuccessMessage(null);
    if (!token) {
      setError("認証トークンがありません。");
      setLoadingData(false);
      return;
    }
    const { data, error: fetchError } = await getProjectProducts(projectId, token);
    if (fetchError) {
      setError(`プロジェクト商品の取得に失敗しました: ${fetchError}`);
    } else {
      setProjectProducts(data || []);
    }
    setLoadingData(false);
  };

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login');
      return;
    }
    
    if (user && token && projectId) {
      fetchProjectProducts();
    }
  }, [user, token, isLoading, router, projectId]);

  const handlePromoteProduct = async (projectProductId: number) => {
    if (!token) {
      setError("認証トークンがありません。ログインし直してください。");
      return;
    }

    setPromotingProductId(projectProductId);
    setError(null);
    setPromoteSuccessMessage(null);

    // TODO: 格上げは親オーガナイゼーションの管理者のみが許可されるべき。
    // フロントエンド側でも権限チェックのUIを設けるべきだが、今回はバックエンドに任せる。
    const { data, error: promoteError } = await promoteProjectProductToShared(projectProductId, token);
    if (promoteError) {
      setError(`商品の格上げに失敗しました: ${promoteError}`);
    } else {
      setPromoteSuccessMessage(`商品「${data?.name}」を共有データに格上げしました！`);
      fetchProjectProducts(); // データ更新のため再取得
    }
    setPromotingProductId(null);
  };

  if (isLoading || !user) {
    return <div className="min-h-screen flex items-center justify-center ml-64">認証情報を確認中...</div>;
  }

  if (!projectId) {
    return <div className="min-h-screen flex items-center justify-center ml-64 text-red-500">プロジェクトIDが見つかりません。</div>;
  }

  return (
    <div className="p-4 sm:p-8 bg-gray-100 min-h-screen">
      <div className="w-full max-w-4xl mx-auto bg-white p-6 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-6">プロジェクトID: {projectId} の商品一覧</h1>
        
        {loadingData ? (
          <p>プロジェクト商品読み込み中...</p>
        ) : error ? (
          <p className="text-red-500 mb-4">{error}</p>
        ) : (
          <>
            {promoteSuccessMessage && (
              <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
                <span className="block sm:inline">{promoteSuccessMessage}</span>
              </div>
            )}

            {projectProducts.length === 0 ? (
              <p className="text-gray-600">このプロジェクトには商品がありません。</p>
            ) : (
              <ul className="space-y-4">
                {projectProducts.map((product) => (
                  <li key={product.id} className="bg-gray-50 p-4 rounded-md shadow-sm flex justify-between items-center">
                    <div>
                      <h2 className="text-xl font-semibold text-gray-700">{product.name}</h2>
                      <p className="text-gray-600">{product.description}</p>
                      <p className="text-green-700 font-bold">¥{product.price.toLocaleString()}</p>
                      <p className="text-sm text-gray-500">
                        コピー元ID: {product.original_shared_product_id ? product.original_shared_product_id : 'なし'}
                      </p>
                    </div>
                    <button
                      onClick={() => handlePromoteProduct(product.id)}
                      disabled={promotingProductId === product.id}
                      className="px-4 py-2 bg-purple-500 text-white rounded-md hover:bg-purple-600 disabled:bg-purple-300 transition-colors duration-200"
                    >
                      {promotingProductId === product.id ? '格上げ中...' : '共有データに格上げ'}
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </>
        )}
      </div>
    </div>
  );
}