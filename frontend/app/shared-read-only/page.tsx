"use client";

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getViewGroupCompanies, ViewGroupCompany } from '@/lib/sharedDataApi';
import { useRouter } from 'next/navigation';

export default function SharedReadOnlyPage() {
  const { user, token, isLoading } = useAuth();
  const router = useRouter();
  const [companies, setCompanies] = useState<ViewGroupCompany[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loadingData, setLoadingData] = useState(true);

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login');
      return;
    }
    
    if (user && token) {
      const fetchCompanies = async () => {
        setLoadingData(true);
        const { data, error } = await getViewGroupCompanies(token);
        if (error) {
          setError(`グループ会社一覧の取得に失敗しました: ${error}`);
        } else {
          setCompanies(data || []);
        }
        setLoadingData(false);
      };
      fetchCompanies();
    }
  }, [user, token, isLoading, router]);

  if (isLoading || !user) {
    return <div className="min-h-screen flex items-center justify-center ml-64">認証情報を確認中...</div>;
  }

  return (
    <div className="p-4 sm:p-8 bg-gray-100 min-h-screen">
      <div className="w-full max-w-4xl mx-auto bg-white p-6 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-6">グループ会社一覧 (参照のみ)</h1>
        {loadingData ? (
          <p>データ読み込み中...</p>
        ) : error ? (
          <p className="text-red-500">{error}</p>
        ) : companies.length === 0 ? (
          <p className="text-gray-600">登録されているグループ会社がありません。</p>
        ) : (
          <ul className="space-y-3">
            {companies.map((company) => (
              <li key={company.id} className="bg-gray-50 p-4 rounded-md shadow-sm">
                <h2 className="text-xl font-semibold text-gray-700">{company.name}</h2>
                <p className="text-sm text-gray-500">作成日時: {new Date(company.created_at).toLocaleString()}</p>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}