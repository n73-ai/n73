import { create } from "zustand";
import type { FlyioIncident } from "@/api/flyio";

type State = {
  incidents: FlyioIncident[];
};

type Actions = {
  setIncidents: (incidents: FlyioIncident[]) => void;
};

export const useFlyioStore = create<State & Actions>((set) => ({
  incidents: [],
  setIncidents: (incidents) => set({ incidents }),
}));
