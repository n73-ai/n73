from claude_code_sdk import query, ClaudeCodeOptions
import requests
import asyncio
from system_prompt import SYSTEM_PROMPT

async def post_json(url, json_data, jwt):
    headers = {
        'Authorization': f'Bearer {jwt}',
        'Content-Type': 'application/json'
    }
    
    loop = asyncio.get_running_loop()
    
    try:
        response = await loop.run_in_executor(
            None, 
            lambda: requests.post(url, json=json_data, headers=headers, timeout=30)
        )
        
        if response.status_code == 401:
            raise ValueError("Invalid or expired JWT token")
        elif not response.ok:
            raise requests.RequestException(f"Request failed with status {response.status_code}: {response.text}")
        
        return response
    
    except requests.RequestException as e:
        raise Exception(f"HTTP request failed: {e}")


async def NewProject(prompt: str, model: str, workDir: str, target_url: str, jwt: str):
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
                    await post_json(target_url, {"type": "text", "text": text}, jwt)

        elif message.__class__.__name__ == "ResultMessage":
            await post_json(target_url, {
                "type": "result",
                "is_error": message.is_error,
                "duration": message.duration_ms,
                "session_id": message.session_id,
                "total_cost_usd": message.total_cost_usd,
                }, jwt)

async def ResumeProject(prompt: str, model: str, workDir: str, target_url: str, session_id: str, jwt: str):
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
        print(message)
        if hasattr(message, "content") and isinstance(message.content, list):
            for block in message.content:
                if hasattr(block, "text"):
                    text = block.text
                    await post_json(target_url, {"type": "text", "text": text}, jwt)

        elif message.__class__.__name__ == "ResultMessage":
            await post_json(target_url, {
                "type": "result",
                "is_error": message.is_error,
                "duration": message.duration_ms,
                "session_id": message.session_id,
                "total_cost_usd": message.total_cost_usd,
                }, jwt)


