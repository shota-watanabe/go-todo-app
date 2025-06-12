"use client";

import AuthForm from "@/components/AuthForm";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Link from "next/link";

export default function RegisterPage() {
  const { user } = useAuth();
  const router = useRouter();

  // もし既にログインしていたら、Todoページへリダイレクト
  useEffect(() => {
    if (user) {
      router.push("/");
    }
  }, [user, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            新規登録
          </h2>
        </div>
        <AuthForm
          apiEndpoint="/register"
          buttonText="登録する"
          onSuccess={() => {}} // 登録成功時の処理はAuthForm内で行う
        />
        <p className="mt-2 text-center text-sm text-gray-600">
          すでにアカウントをお持ちですか？{" "}
          <Link
            href="/login"
            className="font-medium text-indigo-600 hover:text-indigo-500"
          >
            ログインはこちら
          </Link>
        </p>
      </div>
    </div>
  );
}
