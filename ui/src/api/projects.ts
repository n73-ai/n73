import { authAxios, noAuthAxios } from "@/lib/axios-instance";

export const publishProject = async (projectID: string) => {
  const response = await authAxios.post(`/projects/publish/${projectID}`);
  return response.data;
};

export const editProject = async (projectID: string, name: string) => {
  const response = await authAxios.put(`/projects/${projectID}`, {
    name,
  });
  return response.data;
};

export const updateProjectOwner = async (projectID: string, newEmailOwner: string) => {
  const response = await authAxios.post(`/projects/transfer/${projectID}/${newEmailOwner}`);
  return response.data;
};

export const deleteProject = async (projectID: string) => {
  const response = await authAxios.delete(`/projects/${projectID}`);
  return response.data;
};

export const getLatestProjects = async () => {
  const response = await noAuthAxios.get(`/projects/latest`);
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
