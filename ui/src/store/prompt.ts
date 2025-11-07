import { create } from "zustand";
import { persist } from "zustand/middleware";

type State = {
  prompt: string;
};

type Actions = {
  setPrompt: (prompt: string) => void;
};

export const usePromptStore = create(
  persist<State & Actions>(
    (set) => ({
      prompt: "", 
      setPrompt: (prompt: string) =>
        set(() => ({
          prompt,
        })),
    }),
    {
      name: "prompt", 
    }
  )
);
