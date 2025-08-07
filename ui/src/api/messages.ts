import { authAxios } from "@/lib/axios-instance";

export const getMessageByID = async (messageID: string) => {
  const response = await authAxios.post(`/messages/solo/${messageID}`);
  return response.data;
};

export const getMessagesByProjectID = async (projectID: string) => {
  const response = await authAxios.get(`/messages/${projectID}`);
  return response.data;
};
