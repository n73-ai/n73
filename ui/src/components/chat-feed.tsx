import {
  AlertCircleIcon,
  ChevronDown,
  GithubIcon,
  LinkIcon,
  Send,
} from "lucide-react";
import { Button } from "./ui/button";
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
import { useModelStore } from "@/store/models";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Textarea } from "./ui/textarea";

const models = [
  { name: "Claude Sonnet 4", apiName: "claude-sonnet-4-20250514" },
  { name: "Claude Sonnet 3.7", apiName: "claude-3-7-sonnet-20250219" },
  { name: "Claude Haiku 3.5", apiName: "claude-3-5-haiku-20241022" },
  { name: "Claude Haiku 3", apiName: "claude-3-haiku-20240307" },
];

export default function ChatFeed({ p }: { p: any }) {
  const { projectID } = useParams();
  const socketRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const [prompt, setPrompt] = useState("");

  const { model, setModel } = useModelStore();

  const handleModelSelect = (modelObj: (typeof models)[0]) => {
    setModel(modelObj.apiName);
  };

  const selectedModel = models.find((m) => m.apiName === model) || models[0];

  const { data, isLoading, isError } = useQuery<Message[]>({
    queryKey: ["messages", projectID],
    queryFn: () => getMessagesByProjectID(projectID!),
  });

  const resumeProjectMutation = useMutation({
    mutationFn: () => resumeProject(prompt, selectedModel.apiName, projectID!),
    onSuccess: async () => {
      queryClient.invalidateQueries({ queryKey: ["project"] });
      setPrompt("");
      addItemManually({
        role: "user",
        content: prompt,
        total_cost_usd: 0,
        duration: 0,
        model: "",
      });
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
    resumeProjectMutation.mutate();
  };

  const handleSubmitResumeProject = (
    event: React.FormEvent<HTMLFormElement>
  ) => {
    event.preventDefault();
    submitLogic();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      submitLogic();
    }
  };

  const getMessageByIDMutation = useMutation({
    mutationFn: (messageID: string) => getMessageByID(messageID),
    onSuccess: async (response) => {
      addItemManually({
        role: response.role,
        content: response.content,
        total_cost_usd: response.total_cost_usd,
        duration: response.duration,
        model: response.model,
      });
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const addItemManually = (newItem: Message) => {
    queryClient.setQueryData<Message[]>(
      ["messages", projectID],
      (oldData = []) => {
        return [...oldData, newItem];
      }
    );
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [data, p?.status]);

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
          if (event.data.includes("error") || event.data.includes("idle")) {
            queryClient.invalidateQueries({ queryKey: ["project"] });
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
        {isError && <p>An unexpected error occurred.</p>}
        {data?.map((m: Message, i: number) => (
          <div key={i}>
            {m.role === "user" ? (
              <>
                <div className="flex justify-end pb-[30px]">
                  <div className="bg-secondary/60 w-[75%] px-[15px] py-[10px] rounded-md">
                    <p>{m.content}</p>
                  </div>
                </div>

                <div className="flex flex-col gap-[10px] py-[20px]">
                  <div className="flex items-center gap-5">
                    <ZustackLogo size={35} />
                    <p className="text-xl font-semibold text-secondary-foreground">
                      n73
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
                  <div className="flex justify-start text-muted-foreground pb-[20px]">
                    <Accordion type="single" collapsible>
                      <AccordionItem value="item-1">
                        <AccordionTrigger>Response metadata</AccordionTrigger>
                        <AccordionContent>
                          Processing time: {(m.duration / 1000).toFixed(2)}{" "}
                          seconds
                        </AccordionContent>
                        <AccordionContent>
                          Cost: ${m.total_cost_usd.toFixed(4)} USD
                        </AccordionContent>
                        <AccordionContent>Model: {m.model}</AccordionContent>
                      </AccordionItem>
                    </Accordion>
                  </div>
                )}
              </>
            )}
          </div>
        ))}

        {(p.status == "new_pending" || p.status == "pending" || isLoading) && (
          <div className="flex items-center gap-2 text-muted-foreground pb-[20px]">
            <Spinner />
          </div>
        )}

        {(p.status == "internal_error" || p.status == "new_internal_error") && (
          <Alert
            className="flex justify-between items-center"
            variant="destructive"
          >
            <AlertDescription>
              The project deployment encountered an issue. The latest production
              push failed due to an internal error, and our team has already
              been notified.
            </AlertDescription>
          </Alert>
        )}

        {p.error_msg != "" && (
          <Alert
            className="flex justify-between items-center"
            variant="destructive"
          >
            <div className="flex gap-[5px] items-center">
              <AlertCircleIcon />
              <AlertTitle>Code compilation failed</AlertTitle>
            </div>
            <AlertDescription>
              <Button
                onClick={() => {
                  setPrompt(`Fix this build error: ${p.error_msg}`);
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
          <Textarea
            onKeyDown={handleKeyDown}
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder="Reply to n73..."
            className="resize-none"
          />

          <div className="flex justify-end gap-[5px]">
            <a
              className="cursor-pointer inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-9 border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
              href={p?.domain}
              target="_blank"
              rel="noopener noreferrer"
            >
              <LinkIcon />
            </a>

            {p?.gh_repo != "" && (
              <a
                className="cursor-pointer inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[33px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-9 border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
                href={p.gh_repo}
                target="_blank"
                rel="noopener noreferrer"
              >
                <GithubIcon />
              </a>
            )}

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
