import { authAxios } from "@/lib/axios-instance";

export const deployProject = async (projectID: string) => {
  const response = await authAxios.post(`/projects/deploy/${projectID}`);
  return response.data;
};

export const getUserProjects = async () => {
  const response = await authAxios.get(`/projects/user`);
  return response.data;
};

export const resumeProject = async (prompt: string, projectID: string) => {
  const response = await authAxios.post(`/projects/resume/${projectID}`, {
    prompt,
  });
  return response.data;
};

export const getProjectByID = async (projectID: string) => {
  const response = await authAxios.get(`/projects/${projectID}`);
  return response.data;
};

export const createProject = async (prompt: string, name: string) => {
  const response = await authAxios.post(`/projects`, {
    name,
    prompt,
  });
  return response.data;
};
