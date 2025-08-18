import { create } from "zustand";
import { persist } from "zustand/middleware";

type State = {
  exp: number;
  access: string;
  email: string;
  isAuth: boolean;
  hydrated: boolean;
};

type Actions = {
  setAuthState: (
    access: string,
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
      exp: 0,
      email: "",
      isAuth: false,
      hydrated: false,
      setAuthState: (access, exp, email, isAuth) =>
        set(() => ({
          access,
          exp,
          email,
          isAuth,
        })),
      logout: () =>
        set(() => ({ access: "", email: "", isAuth: false })),
    }),
    {
      name: "auth",
      onRehydrateStorage: () => (state) => {
        state?.setAuthState(state.access, state.exp, state.email, state.isAuth);
      },
    }
  )
);
