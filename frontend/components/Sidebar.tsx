"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext"; // ログアウト機能のため

export default function Sidebar() {
  const pathname = usePathname();
  const { logout } = useAuth(); // ログアウト関数を取得

  const navItems = [
    { name: "TODO", href: "/" },
    { name: "グループ会社一覧 (参照のみ)", href: "/shared-read-only" },
    { name: "共有商品 (コピー可)", href: "/shared-copyable" },
    { name: "プロジェクト管理", href: "/projects" },
    // 必要に応じて他のメニューを追加
  ];

  return (
    <div className="flex flex-col w-64 bg-gray-800 text-white p-4 h-screen fixed left-0 top-0 overflow-y-auto">
      <div className="text-2xl font-bold mb-8 text-center border-b border-gray-700 pb-4">
        Todo Admin
      </div>
      <nav className="flex-grow">
        <ul className="space-y-2">
          {navItems.map((item) => (
            <li key={item.name}>
              <Link href={item.href}>
                <p
                  className={`block px-4 py-2 rounded-md transition-colors duration-200 ${
                    pathname === item.href
                      ? "bg-blue-600 text-white"
                      : "hover:bg-gray-700 text-gray-300"
                  }`}
                >
                  {item.name}
                </p>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      <div className="mt-auto pt-4 border-t border-gray-700">
        <button
          onClick={logout}
          className="w-full px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors duration-200"
        >
          ログアウト
        </button>
      </div>
    </div>
  );
}
