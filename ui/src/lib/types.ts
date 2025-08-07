export type ErrorResponse = {
  response: {
    data: {
      error: string;
    };
  };
};

export type Message = {
  role: string;
  content: string;
  duration: number 
  total_cost_usd: number
}
