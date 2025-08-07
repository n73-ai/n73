import { AlertCircleIcon, ChevronDown, Info, Send } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import ZustackLogo from "./zustack-logo";
import Spinner from "./spinner";
import Markdown from "./markdown";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getMessageByID, getMessagesByProjectID } from "@/api/messages";
import { useParams } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import type { Message, ErrorResponse } from "@/lib/types";
import toast from "react-hot-toast";
import { resumeProject } from "@/api/projects";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

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

export default function ChatFeed({ pStatus }: { pStatus: string }) {
  const { projectID } = useParams();
  const socketRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const [prompt, setPrompt] = useState("");
  const [buildError, setBuildError] = useState(false);
  const [buildErrorMessage, setBuildErrorMessage] = useState(true);
  const [selectedModel, setSelectedModel] = useState(models[0]);

  const handleModelSelect = (model: (typeof models)[0]) => {
    setSelectedModel(model);
  };

  // get the messages
  const { data, isLoading, isError } = useQuery<Message[]>({
    queryKey: ["messages"],
    queryFn: () => getMessagesByProjectID(projectID!),
  });

  const resumeProjectMutation = useMutation({
    mutationFn: () => resumeProject(prompt, selectedModel.apiName, projectID!),
    onSuccess: async () => {
      setPrompt("");
      await queryClient.invalidateQueries({ queryKey: ["project"] });
      addItemManually({
        role: "user",
        content: prompt,
        total_cost_usd: 0,
        duration: 0,
      });
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const handleSubmitResumeProject = (
    event: React.FormEvent<HTMLFormElement>
  ) => {
    event.preventDefault();
    if (prompt === "") {
      toast.error("The prompt is required.");
      return;
    }
    resumeProjectMutation.mutate();
  };

  const getMessageByIDMutation = useMutation({
    mutationFn: (messageID: string) => getMessageByID(messageID),
    onSuccess: async (response) => {
      if (response.role == "metadata") {
        await queryClient.invalidateQueries({ queryKey: ["project"] });
      }
      addItemManually({
        role: response.role,
        content: response.content,
        total_cost_usd: response.total_cost_usd,
        duration: response.duration,
      });
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const addItemManually = (newItem: Message) => {
    queryClient.setQueryData<Message[]>(["messages"], (oldData = []) => {
      return [...oldData, newItem];
    });
  };
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [data, pStatus]);

  useEffect(() => {
    if (data && !isLoading) {
      setTimeout(scrollToBottom, 100);
    }
  }, [isLoading, data]);

  useEffect(() => {
    let isMounted = true;
    const timeoutRef = {
      current: null as null | ReturnType<typeof setTimeout>,
    };

    const connect = () => {
      if (!isMounted) return;

      const websocketURL = import.meta.env.VITE_BACKEND_URL;
      const socket = new WebSocket(`${websocketURL}/feed/chat?id=${projectID}`);

      if (
        socketRef.current?.readyState === WebSocket.OPEN ||
        socketRef.current?.readyState === WebSocket.CONNECTING
      ) {
        socketRef.current.close();
      }

      socketRef.current = socket;

      socket.onopen = () => {
        console.log("connected");
      };

      socket.onmessage = (event) => {
        if (event.data !== "") {
          if (event.data.includes("deploy-start")) {
            queryClient.invalidateQueries({ queryKey: ["project"] });
            return;
          }
          if (event.data.includes("deploy-done")) {
            queryClient.invalidateQueries({ queryKey: ["project"] });
            return;
          }
          if (event.data.includes("build-error")) {
            setBuildError(true);
            const cleanedErrMsg = event.data.replace("build-error: ", "");
            setBuildErrorMessage(cleanedErrMsg);
            return;
          }
          getMessageByIDMutation.mutate(event.data);
        }
      };

      socket.onerror = (err) => {
        console.error("WebSocket error: ", err);
      };

      socket.onclose = () => {
        console.log("WebSocket closed. Reconnecting in 3s...");
        if (isMounted) {
          timeoutRef.current = setTimeout(connect, 3000);
        }
      };
    };

    connect();

    return () => {
      isMounted = false;
      clearTimeout(timeoutRef.current!);
      socketRef.current?.close();
    };
  }, [projectID]);

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-[10px]">
        {isLoading && <Spinner />}
        {isError && <p>An unexpected error occurred.</p>}
        {data?.map((m: Message, i: number) => (
          <div key={i}>
            {m.role === "user" ? (
              <>
                <div className="flex justify-end py-[30px]">
                  <div className="bg-secondary/60 w-[75%] px-[15px] py-[10px] rounded-md">
                    <p>{m.content}</p>
                  </div>
                </div>

                <div className="flex flex-col gap-[10px] py-[20px]">
                  <div className="flex items-center gap-5">
                    <ZustackLogo size={35} />
                    <p className="text-xl font-semibold text-secondary-foreground">
                      Zustack
                    </p>
                  </div>
                </div>
              </>
            ) : (
              <>
                <div className="py-[5px]">
                  <Markdown content={m.content} />
                </div>
                {m.role === "metadata" && (
                  <div className="flex items-center gap-[5px] text-muted-foreground">
                    <Info />
                    <p>
                      Worked for {m.duration} milliseconds and consumed USD{" "}
                      {m.total_cost_usd}
                    </p>
                  </div>
                )}
              </>
            )}
          </div>
        ))}

        {pStatus == "Building" && (
          <div className="flex items-center gap-2 text-muted-foreground py-[30px]">
            <Spinner />
            Thinking...
          </div>
        )}

        {buildError && (
          <Alert
            className="my-[20px] flex justify-between items-center"
            variant="destructive"
          >
            <div className="flex gap-[5px] items-center">
              <AlertCircleIcon />
              <AlertTitle>
                An error occurred while compiling the code.
              </AlertTitle>
            </div>
            <AlertDescription>
              <Button
                onClick={() => {
                  setPrompt(`Fix this build error: ${buildErrorMessage}`);
                  resumeProjectMutation.mutate();
                }}
                variant="outline"
              >
                Try to fix
              </Button>
            </AlertDescription>
          </Alert>
        )}

        <div ref={messagesEndRef} />
      </div>
      <div>
        <form
          onSubmit={handleSubmitResumeProject}
          className="w-full space-y-[10px] p-[10px] bg-secondary/50"
        >
          <Input
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Reply AI Zustack..."
            className="resize-none"
          />
          <div className="flex justify-end gap-[5px]">
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
                      selectedModel.apiName === model.apiName ? "bg-accent" : ""
                    }
                  >
                    {model.name}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
            <Button type="submit" variant="outline">
              <span>Send</span>
              {resumeProjectMutation.isPending ? <Spinner /> : <Send />}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
