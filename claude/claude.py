from claude_code_sdk import query, ClaudeCodeOptions
import requests
import asyncio
from system_prompt import SYSTEM_PROMPT

async def post_json(url, json_data):
    loop = asyncio.get_running_loop()
    response = await loop.run_in_executor(None, lambda: requests.post(url, json=json_data))
    return response

async def NewProject(prompt: str, model: str, workDir: str, target_url: str):
    options = ClaudeCodeOptions(
            max_turns=10,
            model=model,
            system_prompt=SYSTEM_PROMPT,
            cwd=workDir,
            allowed_tools=[
                "Read", "Write", "Edit", "MultiEdit",
                "Bash", "LS", "Glob", "Grep", "WebSearch", "WebFetch"
                ],
            permission_mode="acceptEdits"
            )

    async for message in query(prompt=prompt, options=options):
        print(message)
        if hasattr(message, "content") and isinstance(message.content, list):
            for block in message.content:
                if hasattr(block, "text"):
                    text = block.text
                    await post_json(target_url, {"type": "text", "text": text})

        elif message.__class__.__name__ == "ResultMessage":
            await post_json(target_url, {
                "type": message.result,
                "is_error": message.is_error,
                "duration": message.duration_ms,
                "session_id": message.session_id,
                "total_cost_usd": message.total_cost_usd,
                })

async def ResumeProject(prompt: str, model: str, workDir: str, target_url: str, session_id: str):
    options = ClaudeCodeOptions(

            continue_conversation=True,
            resume=session_id,

            max_turns=10,
            model=model,
            system_prompt=SYSTEM_PROMPT,
            cwd=workDir,
            allowed_tools=[
                "Read", "Write", "Edit", "MultiEdit",
                "Bash", "LS", "Glob", "Grep", "WebSearch", "WebFetch"
                ],
            permission_mode="acceptEdits"
            )

    async for message in query(prompt=prompt, options=options):
        if hasattr(message, "content") and isinstance(message.content, list):
            for block in message.content:
                if hasattr(block, "text"):
                    text = block.text
                    print(text)
                    await post_json(target_url, {"type": "text", "text": text})

        elif message.__class__.__name__ == "ResultMessage":
            print(message)
            await post_json(target_url, {
                "type": message.result,
                "is_error": message.is_error,
                "subtype": message.subtype,
                "duration_ms": message.duration_ms,
                "session_id": message.session_id,
                "usage": {
                    "output_tokens": message.usage.get("output_tokens"),
                    "cache_creation_input_tokens": message.usage.get("cache_creation_input_tokens"),
                    "cache_read_input_tokens": message.usage.get("cache_read_input_tokens"),
                    "server_tool_use": message.usage.get("server_tool_use"),
                    "service_tier": message.usage.get("service_tier"),
                    }
                })


