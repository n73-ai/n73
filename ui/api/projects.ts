import { authAxios } from "@/lib/axios_instance";

export const resumeProject = async (prompt: string, id: string) => {
  const response = await authAxios.post(`/project/${id}`, {
    prompt,
  });
  return response.data;
};

export const createProject = async (prompt: string) => {
  const response = await authAxios.post(`/project`, {
    prompt,
  });
  return response.data;
};
