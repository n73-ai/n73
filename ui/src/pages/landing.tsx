import { createProject } from "@/api/projects";
import Spinner from "@/components/spinner";
import { Button } from "@/components/ui/button";
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
import { useModelStore } from "@/store/models";
import { Textarea } from "@/components/ui/textarea";
import LatestProjects from "@/components/latest_projects";
import { usePromptStore } from "@/store/prompt";

const models = [
  { name: "Claude Sonnet 4.5", apiName: "claude-sonnet-4-5-20250929" },
  { name: "Claude Sonnet 4", apiName: "claude-sonnet-4-20250514" },
  { name: "Claude Haiku 4.5", apiName: "claude-haiku-4-5-20251001" },
];

export default function Landing() {
  //const [prompt, setPrompt] = useState("");
  const [name, _] = useState("");
  const navigate = useNavigate();
  const { model, setModel } = useModelStore();

  const handleModelSelect = (modelObj: (typeof models)[0]) => {
    setModel(modelObj.apiName);
  };

  const selectedModel = models.find((m) => m.apiName === model) || models[0];

  const { isAuth } = useAuthStore();
  const { prompt, setPrompt } = usePromptStore()

  const createProjectMut = useMutation({
    mutationFn: () => createProject(prompt, name, selectedModel.apiName),
    onSuccess: (response) => {
      navigate(`/project/${response.project_id}`);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const submitLogic = () => {
    if (prompt === "") {
      toast.error("The prompt is required.");
      return;
    }
    if (!isAuth) {
      navigate("/login")
      return
    }
    createProjectMut.mutate();
  };

  const handleCreateProject = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    submitLogic();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      submitLogic();
    }
  };

  return (
    <section className="container mx-auto px-[10px] 2xl:px-[200px] pt-[150px]">
      <div className="flex flex-col gap-[20px]">
        <div className="flex justify-center items-center gap-[10px] ">
          <div className="hidden lg:block">
            <ZustackLogo size={50} />
          </div>
          <h5 className="text-4xl font-extrabold">
            Build. Preview. Ship 
          </h5>
        </div>

        <div className="flex justify-center">
          <form onSubmit={handleCreateProject} className="w-[700px]">
            <div className="flex flex-col gap-[10px]">
              <Textarea
                onKeyDown={handleKeyDown}
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="Ask n73 to build something cool"
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
                <span>Send</span>
                {createProjectMut.isPending ? <Spinner /> : <Send />}
              </Button>
            </div>
          </form>
        </div>
      </div>
      {isAuth && <Projects />}
      <LatestProjects />
    </section>
  );
}
