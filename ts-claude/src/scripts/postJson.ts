export async function postJson(url: string, jsonData: unknown, jwt: string): Promise<Response> {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${jwt}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(jsonData),
    signal: AbortSignal.timeout(30000),
  });

  console.log("response status: ", response.status)
  if (response.status === 401) throw new Error("Invalid or expired JWT token");
  if (!response.ok) throw new Error(`Request failed with status ${response.status}: ${await response.text()}`);

  return response;
}
