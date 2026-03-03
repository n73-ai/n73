import { query } from "@anthropic-ai/claude-agent-sdk";

async function main() {
  for await (const message of query({
    prompt: "Edit the file hello.py and print with the content of Hello World but backwords",
    options: {
      allowedTools: ["Read", "Edit", "Glob"],
      permissionMode: "acceptEdits",
      cwd: "/home/agust/work/ai/ts-claude",
      systemPrompt: "you are a good python programmer",
      // sessionId: "c55ef272-d033-4814-8860-9b1afe9ecfdd",
      //continue: true,
      resume: "c55ef272-d033-4814-8860-9b1afe9ecfdd",
      stderr: (data) => console.error("[stderr]", data),
    }
  })) {
    if (message.type === "assistant" && message.message?.content) {
      console.log("message: ", message)
      for (const block of message.message.content) {
        if (block.name == "Write") {
          const filePath = (block.input as any).file_path ?? "unknown";
          console.log("created: ", filePath)
        }

        if (block.name == "Edit") {
          const filePath = (block.input as any).file_path ?? "unknown";
          console.log("edited: ", filePath)
        }

        if (block.type == "text") {
          console.log("n83 response: ", block.text);
        }
      }
    } else if (message.type === "result") {
      console.log(`Done: ${message.subtype}`);
    }
  }
}

main();
