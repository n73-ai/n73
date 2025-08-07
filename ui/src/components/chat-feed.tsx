import { ChevronDown, ChevronUp, Info, Send } from "lucide-react";
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
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronLeft, User } from "lucide-react";

export default function ChatFeed({ pStatus }: { pStatus: string }) {
  const { projectID } = useParams();
  const socketRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const [prompt, setPrompt] = useState("");
  // hard code
  const [email, setEmail] = useState("hej@agustfricke.com");

  // get the messages
  const { data, isLoading, isError } = useQuery<Message[]>({
    queryKey: ["messages"],
    queryFn: () => getMessagesByProjectID(projectID),
  });

  const resumeProjectMutation = useMutation({
    mutationFn: () => resumeProject(prompt, projectID),
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

  // add the new message to the data array
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

  // Auto-scroll inicial cuando se cargan los datos
  useEffect(() => {
    if (data && !isLoading) {
      setTimeout(scrollToBottom, 100); // Pequeño delay para asegurar que el DOM esté actualizado
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
      const socket = new WebSocket(`${websocketURL}/feed/chat?email=${email}`);

      if (
        socketRef.current?.readyState === WebSocket.OPEN ||
        socketRef.current?.readyState === WebSocket.CONNECTING
      ) {
        socketRef.current.close();
      }

      socketRef.current = socket;

      socket.onopen = () => {
        //console.log("connected");
      };

      socket.onmessage = (event) => {
        if (event.data !== "") {
          console.log("event.data", event.data);
          console.log("event", event);
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
  }, [email]);

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
                  <div className="bg-secondary/60 px-[15px] py-[10px] rounded-md">
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
              <span>Send</span>
              {resumeProjectMutation.isPending ? <Spinner /> : <Send />}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
