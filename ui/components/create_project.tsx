"use client";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send } from "lucide-react";
import { createProject } from "@/api/projects";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import type { ErrorResponse } from "@/lib/types";
import Spinner from "@/components/spinner";

export default function CreateProject() {
  const [prompt, setPrompt] = useState("");
  const router = useRouter();

  const createProjectMut = useMutation({
    mutationFn: () => createProject(prompt),
    onSuccess: (response) => {
      router.push(`/projects/${response.project_id}`);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data || "An unexpected error occurred.");
    },
  });

  const handleCreateProject = () => {
    if (!prompt) {
      toast.error("The prompt is required to create a project.");
      return;
    }
    createProjectMut.mutate();
  };

  return (
    <div className="w-[700px] mt-[40px] space-y-[10px] p-[15px] bg-input/60 rounded-md">
      <Textarea
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        disabled={createProjectMut.isPending}
        placeholder="Ask Zustack to build..."
        className="resize-none border-none"
      />
      <div className="flex justify-end">
        <Button
          disabled={createProjectMut.isPending}
          onClick={() => handleCreateProject()}
          variant="outline"
        >
          <span>Send</span>
          {createProjectMut.isPending ? <Spinner /> : <Send />}
        </Button>
      </div>
    </div>
  );
}
