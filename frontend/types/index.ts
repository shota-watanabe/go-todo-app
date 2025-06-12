// types/index.ts
export interface Todo {
  id: number;
  content: string;
  completed: boolean;
  createdAt: string; // JSONではstringになる
  updatedAt: string;
}
