"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import { jwtDecode } from "jwt-decode"; // トークンをデコードするライブラリ

// JWTからデコードされるユーザー情報の型
interface User {
  username: string;
  user_id: number;
}

// AuthContextで提供する値の型
interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (token: string) => void;
  logout: () => void;
  isLoading: boolean;
}

// Contextオブジェクトを作成
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// アプリケーション全体をラップするプロバイダーコンポーネント
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true); // 初期ロード状態
  const router = useRouter();

  // アプリ起動時にlocalStorageからトークンを読み込む
  useEffect(() => {
    try {
      const storedToken = localStorage.getItem("jwt_token");
      if (storedToken) {
        const decodedUser = jwtDecode<User>(storedToken);
        setUser(decodedUser);
        setToken(storedToken);
      }
    } catch (error) {
      console.error("Invalid token:", error);
      localStorage.removeItem("jwt_token"); // 不正なトークンは削除
    } finally {
      setIsLoading(false);
    }
  }, []);

  const login = (newToken: string) => {
    try {
      const decodedUser = jwtDecode<User>(newToken);
      localStorage.setItem("jwt_token", newToken);
      setUser(decodedUser);
      setToken(newToken);
      router.push("/"); // ログイン後にTodoページへリダイレクト
    } catch (error) {
      console.error("Failed to decode token on login:", error);
    }
  };

  const logout = () => {
    localStorage.removeItem("jwt_token");
    setUser(null);
    setToken(null);
    router.push("/login"); // ログアウト後にログインページへリダイレクト
  };

  const value = { user, token, login, logout, isLoading };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// Contextを簡単に使うためのカスタムフック
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
