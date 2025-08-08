import { create } from "zustand";
import { persist } from "zustand/middleware";

type State = {
  model: string;
};

type Actions = {
  setModel: (model: string) => void;
};

export const useModelStore = create(
  persist<State & Actions>(
    (set) => ({
      model: "claude-sonnet-4-20250514", 
      setModel: (model: string) =>
        set(() => ({
          model,
        })),
    }),
    {
      name: "model", 
    }
  )
);
