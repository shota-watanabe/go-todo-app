// components/TodoList.tsx
"use client";

import { Todo } from "@/types";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext"; // useAuthをインポート

interface TodoListProps {
  initialTodos: Todo[];
}

export default function TodoList({ token }: TodoListProps) {
  // 親から渡された初期データをstateとして管理
  const [todos, setTodos] = useState<Todo[]>([]);
  const router = useRouter();
  const API_URL = "http://localhost:8080/api";
  const { logout } = useAuth(); // 401エラー時にログアウトさせる

  // コンポーネントマウント時にTodoを取得
  useEffect(() => {
    const fetchTodos = async () => {
      if (!token) return;

      try {
        const res = await fetch(`${API_URL}/todos`, {
          headers: {
            // ★ ヘッダーにトークンを追加
            Authorization: `Bearer ${token}`,
          },
          method: "GET",
        });
        if (res.ok) {
          const data = await res.json();
          setTodos(data);
        } else if (res.status === 401) {
          // トークンが無効ならログアウト
          logout();
        }
      } catch (error) {
        console.error("Failed to fetch todos", error);
      }
    };

    fetchTodos();
  }, [token, logout]); // tokenが変更されたら再取得

  const handleDelete = async (id: number) => {
    // ユーザーに確認ダイアログを表示
    if (!window.confirm("このTodoを本当に削除しますか？")) {
      return;
    }

    try {
      // バックエンドの削除APIを呼び出す
      const res = await fetch(`${API_URL}/todos/${id}`, {
        method: "DELETE",
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (res.ok) {
        // 成功したら、ページをリフレッシュしてリストを更新
        alert("Todoを削除しました。");
        router.refresh();
      } else {
        // 失敗したらエラーメッセージを表示
        alert("Todoの削除に失敗しました。");
      }
    } catch (error) {
      console.error("削除処理中にエラー:", error);
      alert("通信エラーが発生しました。");
    }
  };

  return (
    <ul className="mt-6 space-y-2">
      {todos?.map((todo) => (
        <li
          key={todo.id}
          className="flex items-center justify-between p-3 bg-gray-50 rounded-md"
        >
          <span
            className={`text-gray-800 ${
              todo.completed ? "line-through text-gray-400" : ""
            }`}
          >
            {todo.content}
          </span>
          <button
            onClick={() => handleDelete(todo.id)} // クリック時にhandleDeleteを呼び出す
            className="text-red-500 cursor-pointer"
          >
            削除
          </button>
        </li>
      ))}
    </ul>
  );
}
