import { readdirSync, statSync, readFileSync } from "fs";
import { join, relative } from "path";

const REGION = "se";

export async function deploy2bunny(
  bunnyToken: string,
  storageZoneName: string,
  distDir: string
): Promise<void> {
  const baseURL = `https://${REGION}.storage.bunnycdn.com/${storageZoneName}/`;
  await uploadDirectory(bunnyToken, baseURL, distDir, distDir);
}

async function uploadDirectory(
  token: string,
  baseURL: string,
  distDir: string,
  currentDir: string
): Promise<void> {
  const entries = readdirSync(currentDir);
  for (const entry of entries) {
    const fullPath = join(currentDir, entry);
    if (statSync(fullPath).isDirectory()) {
      await uploadDirectory(token, baseURL, distDir, fullPath);
    } else {
      const relPath = relative(distDir, fullPath).replace(/\\/g, "/");
      const url = baseURL + relPath;
      const fileData = readFileSync(fullPath);
      const resp = await fetch(url, {
        method: "PUT",
        headers: {
          AccessKey: token,
          "Content-Type": "application/octet-stream",
          Accept: "application/json",
        },
        body: fileData,
      });
      if (resp.status !== 201) {
        const body = await resp.text();
        throw new Error(
          `Failed to upload ${relPath}: status ${resp.status}, ${body}`
        );
      }
    }
  }
}
