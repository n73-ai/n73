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
import { useAuthStore } from "@/store/auth";

const models = [
  { name: "Claude Opus 4.1", apiName: "claude-opus-4-1-20250805" },
  { name: "Claude Opus 4", apiName: "claude-opus-4-20250514" },
  { name: "Claude Sonnet 4", apiName: "claude-sonnet-4-20250514" },
  { name: "Claude Sonnet 3.7", apiName: "claude-3-7-sonnet-20250219" },
  { name: "Claude Sonnet 3.5 v2", apiName: "claude-3-5-sonnet-20241022" },
  { name: "Claude Sonnet 3.5", apiName: "claude-3-5-sonnet-20240620" },
  { name: "Claude Haiku 3.5", apiName: "claude-3-5-haiku-20241022" },
  { name: "Claude Haiku 3", apiName: "claude-3-haiku-20240307" },
];

export default function Landing() {
  const [prompt, setPrompt] = useState("");
  const [name, setName] = useState("");
  const [selectedModel, setSelectedModel] = useState(models[0]);
  const navigate = useNavigate();

  const { isAuth } = useAuthStore();

  const createProjectMut = useMutation({
    mutationFn: () => createProject(prompt, name, selectedModel.apiName),
    onSuccess: (response) => {
      navigate(`/project/${response.project_id}`);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const handleCreateProject = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!isAuth) {
      navigate("/signup");
      return;
    }
    if (!prompt) {
      toast.error("The prompt is required to create a project.");
      return;
    }
    if (!name) {
      toast.error("The project name is required to create a project.");
      return;
    }
    if (!selectedModel) {
      toast.error("Please select a model.");
      return;
    }
    createProjectMut.mutate();
  };

  const handleModelSelect = (model: (typeof models)[0]) => {
    setSelectedModel(model);
  };

  return (
    <section className="container mx-auto px-[10px] 2xl:px-[200px]">
      <div className="flex justify-center items-center gap-[20px] pt-[150px]">
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
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="flex gap-[5px]">
                    {selectedModel.name}
                    <ChevronDown />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  {models.map((model) => (
                    <DropdownMenuItem
                      key={model.apiName}
                      onClick={() => handleModelSelect(model)}
                      className={
                        selectedModel.apiName === model.apiName
                          ? "bg-accent"
                          : ""
                      }
                    >
                      {model.name}
                    </DropdownMenuItem>
                  ))}
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
      {isAuth && <Projects />}
    </section>
  );
}
