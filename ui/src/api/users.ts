import { noAuthAxios } from "@/lib/axios-instance";

export const authLink = async (email: string) => {
  const response = await noAuthAxios.post(`/users/auth/link`, {
    email,
  });
  return response.data;
};

export const authVerify = async (tokenString: string) => {
  const response = await noAuthAxios.post(`/users/auth/verify/${tokenString}`);
  return response.data;
};
