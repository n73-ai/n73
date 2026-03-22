import { noAuthAxios } from "@/lib/axios-instance";

export interface FlyioIncident {
  id: string;
  name: string;
  status: string;
  impact: string;
  resolved: boolean;
  created_at: string;
  updated_at: string;
}

export const getFlyioStatus = async (): Promise<{ incidents: FlyioIncident[] }> => {
  const response = await noAuthAxios.get("/flyio/status");
  return response.data;
};
