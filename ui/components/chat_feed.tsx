"use client";
import { useMutation } from "@tanstack/react-query";
import Message from "./message";
import MessageInput from "./message_input";
import { resumeProject } from "@/api/projects";
import { useRouter } from "next/router";
import toast from "react-hot-toast";
import type { ErrorResponse } from "@/lib/types";
import { useState } from "react";

export default function ChatFeed() {

  const router = useRouter();
  const { id } = router.query;

  const [prompt, setPrompt] = useState("")

  const resumeProjectMut = useMutation({
    mutationFn: () => {
      if (!id) throw new Error("Missing the project ID.");
      if (!prompt) toast.error("The prompt is required to create a new project.")
      return resumeProject(prompt, id as string)
    },
    onSuccess: (response) => {
      router.push(`/projects/${response.project_id}`);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data || "An unexpected error occurred.");
    },
  });

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-[10px]">
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
        <Message />
      </div>
      <div>
        <MessageInput />
      </div>
    </div>
  );
}
