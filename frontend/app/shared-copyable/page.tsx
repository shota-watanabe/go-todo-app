"use client";

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getSharedProducts, getProjects, copySharedProductToProject, SharedProduct, Project } from '@/lib/sharedDataApi';
import { useRouter } from 'next/navigation';

export default function SharedCopyablePage() {
  const { user, token, isLoading } = useAuth();
  const router = useRouter();
  const [sharedProducts, setSharedProducts] = useState<SharedProduct[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loadingData, setLoadingData] = useState(true);
  const [copyingProductId, setCopyingProductId] = useState<number | null>(null);
  const [copySuccessMessage, setCopySuccessMessage] = useState<string | null>(null);

  // 仮の組織ID (実際はユーザー情報から取得)
  const organizationId = user?.organization_id || 1; // TODO: userオブジェクトにorganization_idを追加すること

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login');
      return;
    }
    
    if (user && token) {
      const fetchData = async () => {
        setLoadingData(true);
        setError(null);
        setCopySuccessMessage(null);

        // 共有商品を取得
        const { data: productsData, error: productsError } = await getSharedProducts(organizationId, token);
        if (productsError) {
          setError(`共有商品の取得に失敗しました: ${productsError}`);
        } else {
          setSharedProducts(productsData || []);
        }

        // プロジェクト一覧を取得
        const { data: projectsData, error: projectsError } = await getProjects(organizationId, token);
        if (projectsError) {
          setError((prev) => `${prev ? prev + "\n" : ""}プロジェクト一覧の取得に失敗しました: ${projectsError}`);
        } else {
          setProjects(projectsData || []);
          if (projectsData && projectsData.length > 0) {
            setSelectedProjectId(projectsData[0].id); // 最初のプロジェクトをデフォルト選択
          }
        }
        setLoadingData(false);
      };
      fetchData();
    }
  }, [user, token, isLoading, router, organizationId]);

  const handleCopyProduct = async (productId: number) => {
    if (!selectedProjectId) {
      setError("コピー先のプロジェクトを選択してください。");
      return;
    }
    if (!token) {
      setError("認証トークンがありません。ログインし直してください。");
      return;
    }

    setCopyingProductId(productId);
    setError(null);
    setCopySuccessMessage(null);

    const { data, error } = await copySharedProductToProject(selectedProjectId, productId, token);
    if (error) {
      setError(`商品のコピーに失敗しました: ${error}`);
    } else {
      setCopySuccessMessage(`商品「${data?.name}」をプロジェクトにコピーしました！`);
    }
    setCopyingProductId(null);
  };

  if (isLoading || !user) {
    return <div className="min-h-screen flex items-center justify-center ml-64">認証情報を確認中...</div>;
  }

  return (
    <div className="p-4 sm:p-8 bg-gray-100 min-h-screen">
      <div className="w-full max-w-4xl mx-auto bg-white p-6 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-6">共有商品 (コピー可)</h1>
        
        {loadingData ? (
          <p>データ読み込み中...</p>
        ) : error ? (
          <p className="text-red-500 mb-4">{error}</p>
        ) : (
          <>
            <div className="mb-6">
              <label htmlFor="project-select" className="block text-gray-700 text-sm font-bold mb-2">
                コピー先のプロジェクトを選択:
              </label>
              <select
                id="project-select"
                className="shadow border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                value={selectedProjectId || ''}
                onChange={(e) => setSelectedProjectId(Number(e.target.value))}
              >
                {projects.length === 0 ? (
                  <option value="">プロジェクトがありません</option>
                ) : (
                  projects.map((project) => (
                    <option key={project.id} value={project.id}>
                      {project.name}
                    </option>
                  ))
                )}
              </select>
            </div>

            {copySuccessMessage && (
              <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
                <span className="block sm:inline">{copySuccessMessage}</span>
              </div>
            )}

            {sharedProducts.length === 0 ? (
              <p className="text-gray-600">共有商品がありません。</p>
            ) : (
              <ul className="space-y-4">
                {sharedProducts.map((product) => (
                  <li key={product.id} className="bg-gray-50 p-4 rounded-md shadow-sm flex justify-between items-center">
                    <div>
                      <h2 className="text-xl font-semibold text-gray-700">{product.name}</h2>
                      <p className="text-gray-600">{product.description}</p>
                      <p className="text-green-700 font-bold">¥{product.price.toLocaleString()}</p>
                      <p className="text-sm text-gray-500">SKU: {product.sku}</p>
                    </div>
                    <button
                      onClick={() => handleCopyProduct(product.id)}
                      disabled={!selectedProjectId || copyingProductId === product.id}
                      className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:bg-blue-300 transition-colors duration-200"
                    >
                      {copyingProductId === product.id ? 'コピー中...' : 'プロジェクトにコピー'}
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