from claude_code_sdk import query, ClaudeCodeOptions
import requests
import asyncio
import os
from system_prompt import SYSTEM_PROMPT
import zipfile
import os
from pathlib import Path

def zip_directory(folder_path: str | Path, output_zip: str | Path):
    folder_path = Path(folder_path).resolve()
    output_zip = Path(output_zip)

    if not folder_path.is_dir():
        raise ValueError(f"Directory not found: {folder_path}")

    with zipfile.ZipFile(output_zip, 'w', zipfile.ZIP_DEFLATED) as zipf:
        for root, dirs, files in os.walk(folder_path):
            # --- EXCLUSION ---
            if "node_modules" in dirs:
                dirs.remove("node_modules")   # evita entrar a node_modules
            # -----------------

            for file in files:
                file_path = Path(root) / file
                arcname = file_path.relative_to(folder_path.parent)
                zipf.write(file_path, arcname)

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


async def process_message(message, target_url, jwt):
    """Process a single message and handle file operations and text content."""
    print(message)
    
    # Detect file operations first
    if hasattr(message, "content") and isinstance(message.content, list):
        for block in message.content:
            # Check for tool use blocks (file operations)
            if hasattr(block, "name") and hasattr(block, "input"):
                if block.name == "Write":
                    # File being created
                    file_path = block.input.get("file_path", "unknown")
                    clean_path = file_path.replace("../../ui/", "").replace("../ui/", "")
                    await post_json(target_url, {
                       "type": "text", 
                        "text": f"Created **{clean_path}**"
                    }, jwt)
                elif block.name == "Edit":
                    # File being edited
                    file_path = block.input.get("file_path", "unknown")
                    clean_path = file_path.replace("../../ui/", "").replace("../ui/", "")
                    await post_json(target_url, {
                        "type": "text", 
                        "text": f"Edited **{clean_path}**"
                    }, jwt)
            
            # Process text content
            if hasattr(block, "text"):
                text = block.text
                await post_json(target_url, {
                    "type": "text", 
                    "text": text
                }, jwt)
    
    # Handle result messages
    elif message.__class__.__name__ == "ResultMessage":
        zip_directory("/app/project", "project.zip")
        
        # Read zip file and encode as base64
        import base64
        with open("/app/project.zip", 'rb') as f:
            zip_data = base64.b64encode(f.read()).decode('utf-8')
        
        await post_json(target_url, {
            "file": zip_data,
            "type": "result",
            "is_error": message.is_error,
            "duration": message.duration_ms,
            "session_id": message.session_id,
            "total_cost_usd": message.total_cost_usd,
        }, jwt)
        
        # Clean up zip file
        if os.path.exists("/app/project.zip"):
            os.remove("/app/project.zip")

async def NewProject(prompt: str, model: str, workDir: str, target_url: str, jwt: str):
    options = ClaudeCodeOptions(
        max_turns=50,
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
        await process_message(message, target_url, jwt)


async def ResumeProject(prompt: str, model: str, workDir: str, target_url: str, session_id: str, jwt: str):
    options = ClaudeCodeOptions(
        continue_conversation=True,
        resume=session_id,
        max_turns=50,
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
        await process_message(message, target_url, jwt)
