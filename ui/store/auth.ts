import { create } from "zustand";
import { persist } from "zustand/middleware";

type State = {
  exp: number;
  access: string;
  userId: string;
  email: string;
  isAuth: boolean;
};

type Actions = {
  setAuthState: (
    access: string,
    userId: string,
    exp: number,
    email: string,
    isAuth: boolean
  ) => void;
  logout: () => void;
};

export const useAuthStore = create(
  persist<State & Actions>(
    (set) => ({
      access: "",
      userId: "",
      exp: 0,
      email: "",
      isAuth: false,
      setAuthState: (
        access: string,
        userId: string,
        exp: number,
        email: string,
        isAuth: boolean
      ) =>
        set(() => ({
          access,
          userId,
          exp,
          email,
          isAuth,
        })),
      logout: () =>
        set(() => ({ access: "", userId: "", email: "", isAuth: false })),
    }),
    {
      name: "auth",
    }
  )
);
