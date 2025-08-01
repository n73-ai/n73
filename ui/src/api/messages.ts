import { authAxios } from "@/lib/axios-instance";

export const getMessagesByProjectID = async (projectID: string) => {
  const response = await authAxios.get(`/messages/${projectID}`);
  return response.data;
};
