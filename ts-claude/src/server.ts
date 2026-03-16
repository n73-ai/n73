import { Hono } from "hono";
import { serve } from "@hono/node-server";
import { logger } from "hono/logger";
import { cors } from "hono/cors";
import { prettyJSON } from "hono/pretty-json";
import { newProject, resumeProject } from "./agent";

// --- App ---
const app = new Hono();

// Middleware
app.use("*", logger());
app.use("*", cors());
app.use("*", prettyJSON());

interface Payload {
  prompt: string;
  model: string;
  work_dir: string;
  webhook_url: string;
  session_id: string;
  jwt: string;
  project_id: string;
  storage_zone_password: string;
}

// GET all users
app.post("/claude/new", async (c) => {
  const body = await c.req.json<Payload>();

  const { prompt, model, work_dir, webhook_url, jwt, project_id, storage_zone_password } = body;

  // validate required fields
  if (!prompt || !model || !jwt) {
    return c.json({ error: "Missing required fields" }, 400);
  }

  newProject(prompt, model, work_dir, webhook_url, jwt, project_id, storage_zone_password)

  return c.json({ status: "ok" });
});

app.post("/claude/resume", async (c) => {
  const body = await c.req.json<Payload>();

  const { prompt, model, work_dir, webhook_url, session_id, jwt, project_id, storage_zone_password } = body;

  // validate required fields
  if (!prompt || !model || !session_id || !jwt) {
    return c.json({ error: "Missing required fields" }, 400);
  }

  resumeProject(prompt, model, work_dir, webhook_url, jwt, session_id, project_id, storage_zone_password)

  return c.json({ status: "ok", session_id });
});

app.get("/health", (c) => {
  return c.json({ status: "ok", timestamp: new Date().toISOString() });
});

app.notFound((c) => {
  return c.json({ error: `Route ${c.req.path} not found` }, 404);
});

const PORT = 5000;

serve({ fetch: app.fetch, port: PORT, hostname: "0.0.0.0" }, (info) => {
  console.log(`Server running at http://0.0.0.0:${info.port}`);
});
