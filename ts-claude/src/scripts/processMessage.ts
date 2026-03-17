import { execSync, spawn } from "child_process";
import { readFileSync, existsSync, unlinkSync, copyFileSync } from "fs";
import type { Message } from "@anthropic-ai/sdk/resources/messages";
import { postJson } from "./postJson";
import { zipDir } from "./zipDir";
import { deploy2bunny } from "./deploy2bunny";

async function waitForServer(url: string, timeoutMs = 30000): Promise<void> {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const res = await fetch(url, { signal: AbortSignal.timeout(1000) });
      if (res.ok) return;
    } catch {}
    await new Promise((r) => setTimeout(r, 500));
  }
  throw new Error(`Server at ${url} did not start within ${timeoutMs}ms`);
}

async function takeScreenshot(): Promise<void> {
  const staticServer = spawn(
    "npx",
    ["vite", "preview", "--host", "0.0.0.0", "--port", "5173"],
    { cwd: "/app/ui-only", stdio: "pipe" }
  );
  try {
    await waitForServer("http://0.0.0.0:5173");
    execSync("/app/src/scripts/screenshot http://0.0.0.0:5173", { stdio: "pipe" });
  } catch (e) {
    console.error("screenshot failed:", e);
  } finally {
    staticServer.kill();
  }
}

interface ResultMessage {
  type: "result";
  duration_ms: number;
  session_id: string;
  total_cost_usd: number;
}

export async function processMessage(
  message: Message | ResultMessage,
  targetUrl: string,
  jwt: string,
  projectId: string,
  storageZonePassword: string,
  npmInstallPromise?: Promise<void>
): Promise<void> {
  if (message.type === "assistant" && message.message?.content) {
    for (const block of message.message.content) {
      if (block.name == "Write") {
        const filePath = (block.input as any).file_path ?? "unknown";
        await postJson(
          targetUrl,
          { type: "text", text: `Created **${filePath}**` },
          jwt
        );
      }

      if (block.name == "Edit") {
        const filePath = (block.input as any).file_path ?? "unknown";
        await postJson(
          targetUrl,
          { type: "text", text: `Edited **${filePath}**` },
          jwt
        );
      }

      if (block.type == "text") {
        await postJson(targetUrl, { type: "text", text: block.text }, jwt);
      }
    }
  } else if (message.type === "result") {
    const result = message as ResultMessage;

    if (npmInstallPromise) {
      await npmInstallPromise;
    }

    let isBuildError = false;
    let buildErrorMsg: string | null = null;

    try {
      execSync("npm install", { cwd: "/app/ui-only", stdio: "pipe" });
      execSync("npm run build", { cwd: "/app/ui-only", stdio: "pipe" });
    } catch (e: any) {
      isBuildError = true;
      const stderr = e.stderr?.toString().trim() || "";
      const stdout = e.stdout?.toString().trim() || "";
      buildErrorMsg = [stderr, stdout].filter(Boolean).join("\n") || String(e);
    }

    if (!isBuildError) {
      await takeScreenshot();
    }

    let zipData: string | null = null;
    if (!isBuildError) {
      if (existsSync("/app/screenshot.png")) {
        copyFileSync("/app/screenshot.png", "/app/ui-only/dist/screenshot.png");
      }
      await zipDir("/app/ui-only", "/app/project.zip");
      zipData = readFileSync("/app/project.zip").toString("base64");
    }

    await postJson(
      targetUrl,
      {
        type: "result",
        build_error: isBuildError,
        build_error_msg: isBuildError ? buildErrorMsg : null,
        file: zipData,
        duration: result.duration_ms,
        session_id: result.session_id,
        total_cost_usd: result.total_cost_usd,
      },
      jwt
    );

    if (!isBuildError && storageZonePassword) {
      try {
        await deploy2bunny(storageZonePassword, projectId, "/app/ui-only/dist");
      } catch (e) {
        console.error("deploy2bunny failed:", e);
      }
      await postJson(targetUrl, { type: "deployed" }, jwt);
    }

    if (existsSync("/app/project.zip")) unlinkSync("/app/project.zip");
  }
}
