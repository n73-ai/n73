import { createProject } from "@/api/projects";
import Spinner from "@/components/spinner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import ZustackLogo from "@/components/zustack-logo";
import type { ErrorResponse } from "@/lib/types";
import { useMutation } from "@tanstack/react-query";
import { ChevronDown, Send } from "lucide-react";
import { useState } from "react";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";
import Projects from "@/components/projects";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function Landing() {
  const [prompt, setPrompt] = useState("");
  const [name, setName] = useState("");
  const navigate = useNavigate();

  const createProjectMut = useMutation({
    mutationFn: () => createProject(prompt, name),
    onSuccess: (response) => {
      navigate(`/project/${response.project_id}`);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const handleCreateProject = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!prompt) {
      toast.error("The prompt is required to create a project.");
      return;
    }
    createProjectMut.mutate();
  };

  return (
    <section className="container mx-auto px-[10px] 2xl:px-[200px]">
      <div className="flex justify-center items-center gap-[20px] pt-[200px]">
        <ZustackLogo size={300} />
        <div className="flex flex-col gap-[20px]">
          <h1 className="scroll-m-20 text-center text-4xl font-extrabold tracking-tight text-balance">
            Build. Preview. Ship with AI.
          </h1>
          <p className="leading-7">
            Create apps and websites by chatting with AI
          </p>
          <form onSubmit={handleCreateProject}>
            <div className="flex flex-col gap-[10px]">
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Name your project"
              />
              <Input
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="Ask Zustack to build something cool"
              />
            </div>
            <div className="flex justify-end gap-[10px] pt-[10px]">
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Button variant="outline" className="flex gap-[5px]">
                    Claude Opus 4.1
                    <ChevronDown />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem>Claude Opus 4.1</DropdownMenuItem>
                  <DropdownMenuItem>Claude Opus 4</DropdownMenuItem>
                  <DropdownMenuItem>Claude Sonnet 4</DropdownMenuItem>
                  <DropdownMenuItem>Claude Sonnet 3.7</DropdownMenuItem>
                  <DropdownMenuItem>Claude Haiku 3.5</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              <Button type="submit" variant="outline">
                <span>Create</span>
                {createProjectMut.isPending ? <Spinner /> : <Send />}
              </Button>
            </div>
          </form>
        </div>
      </div>
      <Projects />
    </section>
  );
}
