import axios from "axios";
import type { AxiosRequestHeaders } from 'axios';
import { useAuthStore } from "@/store/auth";

const baseURL = import.meta.env.VITE_BACKEND_URL;

export const noAuthAxios = axios.create({
  baseURL,
});

export const authAxios = axios.create({
  baseURL,
  withCredentials: false,
});

authAxios.interceptors.request.use(
  async (config) => {
    try {
      const { access, exp, logout } = useAuthStore.getState();
      if (access) {
        config.headers = {
          Authorization: `Bearer ${access}`,
        } as AxiosRequestHeaders;

        const expirationDate = new Date(exp * 1000);

        const now = new Date();

        if (expirationDate.getTime() - now.getTime() < 0) {
          logout()
          throw new axios.Cancel("Token expired");
        }
      }
      return config;
    } catch (error) {
      return Promise.reject(error);
    }
  },
  (error) => {
    return Promise.reject(error);
  }
);
