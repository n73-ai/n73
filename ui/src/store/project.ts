import { create } from "zustand";
import { persist } from "zustand/middleware";

type ProjectState = {
  isPending: boolean;
  errorMsg: string;
  fixingErr: boolean;
};

type State = {
  projects: Record<string, ProjectState>;
};

type Actions = {
  setProjectState: (
    projectID: string, 
    isPending: boolean, 
    errorMsg: string, 
    fixingErr: boolean
  ) => void;
  getProjectState: (projectID: string) => ProjectState;
  clearProjectState: (projectID: string) => void;
  clearAllProjects: () => void;
};

const defaultProjectState: ProjectState = {
  isPending: false,
  fixingErr: false,
  errorMsg: "",
};

export const useProjectStore = create(
  persist<State & Actions>(
    (set, get) => ({
      projects: {},
      
      setProjectState: (projectID, isPending, errorMsg, fixingErr) =>
        set((state) => ({
          projects: {
            ...state.projects,
            [projectID]: {
              isPending,
              errorMsg,
              fixingErr,
            },
          },
        })),

      getProjectState: (projectID) => {
        const state = get();
        return state.projects[projectID] || defaultProjectState;
      },

      clearProjectState: (projectID) =>
        set((state) => {
          const { [projectID]: _, ...restProjects } = state.projects;
          return { projects: restProjects };
        }),

      clearAllProjects: () => set({ projects: {} }),
    }),
    {
      name: "project-states",
    }
  )
);

export const useProjectStateById = (projectID: string) => {
  const setProjectState = useProjectStore((state) => state.setProjectState);
  //const getProjectState = useProjectStore((state) => state.getProjectState);
  const clearProjectState = useProjectStore((state) => state.clearProjectState);
  
  const projectState = useProjectStore((state) => 
    state.projects[projectID] || defaultProjectState
  );

  return {
    ...projectState,
    setProjectState: (isPending: boolean, errorMsg: string, fixingErr: boolean) =>
      setProjectState(projectID, isPending, errorMsg, fixingErr),
    clearProjectState: () => clearProjectState(projectID),
  };
};
