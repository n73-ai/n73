import { authAxios } from "@/lib/axios-instance";

export const getLatestProjects = async () => {
  const response = await authAxios.get(`/projects/latest`);
  return response.data;
};

export const getUserProjects = async () => {
  const response = await authAxios.get(`/projects/user`);
  return response.data;
};

export const resumeProject = async (
  prompt: string,
  model: string,
  projectID: string
) => {
  const response = await authAxios.post(`/projects/resume/${projectID}`, {
    prompt,
    model,
  });
  return response.data;
};

export const getProjectByID = async (projectID: string) => {
  const response = await authAxios.get(`/projects/${projectID}`);
  return response.data;
};

export const createProject = async (
  prompt: string,
  name: string,
  model: string
) => {
  const response = await authAxios.post(`/projects`, {
    name,
    prompt,
    model,
  });
  return response.data;
};
