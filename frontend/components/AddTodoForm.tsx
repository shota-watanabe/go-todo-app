// components/AddTodoForm.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation"; // App Router用のルーター

const API_URL = "http://localhost:8080/api";

export default function AddTodoForm() {
  const [content, setContent] = useState("");
  const router = useRouter(); // ページをリフレッシュするために使用

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;

    try {
      const res = await fetch(`${API_URL}/todos`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });

      if (res.ok) {
        setContent("");
        router.refresh(); // サーバーコンポーネントを再実行させ、リストを更新する
      }
    } catch (error) {
      console.error("Failed to create todo:", error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <input
        type="text"
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="New Todo..."
        className="flex-grow p-2 border border-gray-500 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
      />
      <button
        type="submit"
        className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        Add
      </button>
    </form>
  );
}