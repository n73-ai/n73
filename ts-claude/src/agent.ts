import { query } from "@anthropic-ai/claude-agent-sdk";
import { readFileSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { spawn } from "child_process";
import { processMessage } from "./scripts/processMessage";

const __dirname = dirname(fileURLToPath(import.meta.url));
const system_prompt = readFileSync(
  join(__dirname, "system_prompt.txt"),
  "utf-8"
);

function runNpmInstall(cwd: string): Promise<void> {
  return new Promise((resolve) => {
    const proc = spawn("npm", ["install", "--include=dev"], { cwd, stdio: "pipe" });
    proc.on("close", (code) => {
      if (code !== 0) console.error(`npm install exited with code ${code}`);
      resolve(); // resolve regardless so build still runs
    });
    proc.on("error", (err) => {
      console.error("npm install error:", err);
      resolve();
    });
  });
}

export async function newProject(
  prompt: string,
  model: string,
  cwd: string,
  targetUrl: string,
  jwt: string,
  projectId: string,
  storageZonePassword: string
) {
  const npmInstallPromise = runNpmInstall(cwd);

  for await (const message of query({
    prompt: prompt,
    options: {
      allowedTools: [
        "Read",
        "Write",
        "Edit",
        "MultiEdit",
        "Bash",
        "LS",
        "Glob",
        "Grep",
        "WebSearch",
        "WebFetch",
      ],
      permissionMode: "acceptEdits",
      systemPrompt: system_prompt,
      maxTurns: 50,
      cwd: cwd,
      model: model,
    },
  })) {
    await processMessage(message, targetUrl, jwt, projectId, storageZonePassword, npmInstallPromise);
  }
}

export async function resumeProject(
  prompt: string,
  model: string,
  cwd: string,
  targetUrl: string,
  jwt: string,
  session_id: string,
  projectId: string,
  storageZonePassword: string
) {
  for await (const message of query({
    prompt: prompt,
    options: {
      allowedTools: [
        "Read",
        "Write",
        "Edit",
        "MultiEdit",
        "Bash",
        "LS",
        "Glob",
        "Grep",
        "WebSearch",
        "WebFetch",
      ],
      permissionMode: "acceptEdits",
      systemPrompt: system_prompt,
      maxTurns: 50,
      cwd: cwd,
      model: model,
      resume: session_id,
    },
  })) {
    await processMessage(message, targetUrl, jwt, projectId, storageZonePassword);
  }
}
