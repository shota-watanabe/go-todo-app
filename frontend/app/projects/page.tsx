"use client";

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getProjects, createProject, Project } from '@/lib/sharedDataApi';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function ProjectsPage() {
  const { user, token, isLoading } = useAuth();
  const router = useRouter();
  const [projects, setProjects] = useState<Project[]>([]);
  const [newProjectName, setNewProjectName] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loadingData, setLoadingData] = useState(true);
  const [creatingProject, setCreatingProject] = useState(false);

  // 仮の組織ID (実際はユーザー情報から取得)
  const organizationId = user?.organization_id || 1; // TODO: userオブジェクトにorganization_idを追加すること

  const fetchProjects = async () => {
    setLoadingData(true);
    setError(null);
    const { data, error: fetchError } = await getProjects(organizationId, token);
    if (fetchError) {
      setError(`プロジェクト一覧の取得に失敗しました: ${fetchError}`);
    } else {
      // deleted_at が設定されていない（論理削除されていない）プロジェクトのみ表示
      setProjects(data?.filter(p => !p.deleted_at) || []);
    }
    setLoadingData(false);
  };

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login');
      return;
    }
    
    if (user && token) {
      fetchProjects();
    }
  }, [user, token, isLoading, router, organizationId]);

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newProjectName.trim() || !token) {
      setError("プロジェクト名を入力してください。");
      return;
    }

    setCreatingProject(true);
    setError(null);

    const { data, error: createError } = await createProject(newProjectName, token);
    if (createError) {
      setError(`プロジェクト作成に失敗しました: ${createError}`);
    } else {
      setNewProjectName('');
      fetchProjects(); // プロジェクト一覧を再取得
    }
    setCreatingProject(false);
  };

  if (isLoading || !user) {
    return <div className="min-h-screen flex items-center justify-center ml-64">認証情報を確認中...</div>;
  }

  return (
    <div className="p-4 sm:p-8 bg-gray-100 min-h-screen">
      <div className="w-full max-w-4xl mx-auto bg-white p-6 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-800 mb-6">プロジェクト管理</h1>

        <form onSubmit={handleCreateProject} className="mb-8 p-4 border rounded-md bg-gray-50">
          <h2 className="text-xl font-semibold mb-3 text-gray-700">新しいプロジェクトを作成</h2>
          <div className="flex flex-col sm:flex-row gap-4">
            <input
              type="text"
              placeholder="プロジェクト名"
              value={newProjectName}
              onChange={(e) => setNewProjectName(e.target.value)}
              className="flex-grow shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              disabled={creatingProject}
            />
            <button
              type="submit"
              className="px-6 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 disabled:bg-green-300 transition-colors duration-200"
              disabled={creatingProject}
            >
              {creatingProject ? '作成中...' : 'プロジェクトを作成'}
            </button>
          </div>
          {error && <p className="text-red-500 text-sm mt-2">{error}</p>}
        </form>

        <h2 className="text-2xl font-bold text-gray-800 mb-4">プロジェクト一覧</h2>
        {loadingData ? (
          <p>プロジェクト読み込み中...</p>
        ) : projects.length === 0 ? (
          <p className="text-gray-600">プロジェクトがありません。</p>
        ) : (
          <ul className="space-y-4">
            {projects.map((project) => (
              <li key={project.id} className="bg-gray-50 p-4 rounded-md shadow-sm flex justify-between items-center">
                <div>
                  <h3 className="text-xl font-semibold text-gray-700">{project.name}</h3>
                  <p className="text-sm text-gray-500">作成日時: {new Date(project.created_at).toLocaleString()}</p>
                </div>
                <Link href={`/projects/${project.id}`} legacyBehavior>
                  <a className="px-4 py-2 bg-indigo-500 text-white rounded-md hover:bg-indigo-600 transition-colors duration-200">
                    詳細を見る
                  </a>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}