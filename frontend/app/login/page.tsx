"use client";

import AuthForm from "@/components/AuthForm";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import Link from "next/link";

export default function LoginPage() {
  const { login, user } = useAuth();
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
            ログイン
          </h2>
        </div>
        <AuthForm
          apiEndpoint="/login"
          buttonText="ログイン"
          onSuccess={login}
        />
        <p className="mt-2 text-center text-sm text-gray-600">
          アカウントをお持ちでないですか？{" "}
          <Link
            href="/register"
            className="font-medium text-indigo-600 hover:text-indigo-500"
          >
            新規登録はこちら
          </Link>
        </p>
      </div>
    </div>
  );
}
