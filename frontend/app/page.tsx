// ★ "use client" を先頭に追加
"use client";

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import AddTodoForm from "@/components/AddTodoForm";
import TodoList from "@/components/TodoList";

// ★★★ ページ全体をクライアントコンポーネントに変更し、認証チェックを追加 ★★★

export default function Home() {
  const { user, token, isLoading } = useAuth();
  const router = useRouter();

  // 認証状態をチェックする
  useEffect(() => {
    // データロード中でなく、ユーザー情報がなければログインページへ
    if (!isLoading && !user) {
      router.push('/login');
    }
  }, [user, isLoading, router]);

  // ローディング中またはリダイレクト中は何も表示しない
  if (isLoading || !user) {
    return <div className="min-h-screen flex items-center justify-center">ローディング...</div>;
  }

  return (
    <main className="flex min-h-screen flex-col items-center p-4 sm:p-24 bg-gray-100">
      <div className="w-full max-w-2xl">
        <div className="bg-white p-6 rounded-lg shadow-lg">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-4xl font-bold text-gray-800">Todo App</h1>
            <div className='flex items-center gap-4'>
              <span className='text-gray-700'>こんにちは、{user.username}さん</span>
            </div>
          </div>
          
          <AddTodoForm token={token}/>
          {/* tokenをpropsで渡すようにする */}
          <TodoList token={token} /> 
        </div>
      </div>
    </main>
  );
}