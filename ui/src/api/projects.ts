import { authAxios } from "@/lib/axios-instance";

export const createProject = async (prompt: string) => {
  const response = await authAxios.post(`/projects`, {
    prompt,
  });
  return response.data;
};
