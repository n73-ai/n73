import { execSync } from "child_process";
import { readFileSync, existsSync, unlinkSync } from "fs";
import type { Message } from "@anthropic-ai/sdk/resources/messages";
import { postJson } from "./postJson";
import { zipDir } from "./zipDir";
import { deploy2bunny, uploadScreenshot } from "./deploy2bunny";

interface ResultMessage {
  type: "result";
  duration_ms: number;
  session_id: string;
  total_cost_usd: number;
}

export async function processMessage(
  message: Message | ResultMessage,
  targetUrl: string,
  jwt: string
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

    execSync("/app/src/scripts/screenshot http://0.0.0.0:5173", {
      stdio: "pipe",
    });

    let isBuildError = false;
    let buildErrorMsg: string | null = null;

    try {
      execSync("npm run build", { cwd: "/app/ui-only", stdio: "pipe" });
    } catch (e) {
      isBuildError = true;
      buildErrorMsg = String(e);
    }

    let zipData: string | null = null;
    if (!isBuildError) {
      await zipDir("/app/ui-only", "/app/project.zip");
      zipData = readFileSync("/app/project.zip").toString("base64");
    }

    const imageBase64 = readFileSync("/app/screenshot.png").toString("base64");

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
        image: imageBase64,
      },
      jwt
    );

    // if else depends of bunny status!!
    const isNewProject = true;
    if (isNewProject) {
      // need to know if the creation of storage and pull zone has been created successfully
      // after than do the upload
      // upload()
    } else {
      // delete()
      // upload()
      // purge()
    }

    // here do the deploy of dist dir to bunny net 
    const distDir = "~/ui-only/dist"
    // this comes from http request
    const bunnyStorageZonePassword = ""
    deploy2bunny(bunnyStorageZonePassword, distDir)

    // here do upload of screenshot to bunny net
    const imagePath = "~/screenshot.png"
    uploadScreenshot(bunnyStorageZonePassword, imagePath)

    if (existsSync("/app/project.zip")) unlinkSync("/app/project.zip");
  }
}
