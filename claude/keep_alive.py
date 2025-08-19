import anyio
import asyncio
from claude_code_sdk import query

async def main():
    async for message in query(prompt="Don't answer anything."):
        print(message)

async def keep_alive():
    await asyncio.sleep(1 * 60 * 60) 
    while True:
        print(f"Ejecutando keep alive: {anyio.current_time()}")
        await main()
        await asyncio.sleep(1 * 60 * 60)  

anyio.run(keep_alive)
