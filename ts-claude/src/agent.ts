import { query } from "@anthropic-ai/claude-agent-sdk";
import { readFileSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";
import { processMessage } from "./scripts/processMessage";

const __dirname = dirname(fileURLToPath(import.meta.url));
const system_prompt = readFileSync(
  join(__dirname, "system_prompt.txt"),
  "utf-8"
);

export async function newProject(
  prompt: string,
  model: string,
  cwd: string,
  targetUrl: string,
  jwt: string
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
    },
  })) {
    await processMessage(message, targetUrl, jwt);
  }
}

export async function resumeProject(
  prompt: string,
  model: string,
  cwd: string,
  targetUrl: string,
  jwt: string,
  session_id: string
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
    await processMessage(message, targetUrl, jwt);
  }
}
