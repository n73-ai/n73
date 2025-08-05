import { Info, Send } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import ZustackLogo from "./zustack-logo";
import Spinner from "./spinner";
import Markdown from "./markdown";

const data = [
  {
    role: "user",
    content: "Build a Hello World in html",
    model: "",
    duration: 0,
    is_error: false,
    total_cost_usd: 0,
  },
  {
    role: "assistant",
    content: "I'll create a beautiful Hello World app with a modern, clean design. Let me start by setting up the design system and then building the interface.",
    model: "",
    duration: 0,
    is_error: false,
    total_cost_usd: 0,
  },
  {
    role: "assistant",
    content: "Created a beautiful Hello World app with a modern dark theme featuring:",
    model: "",
    duration: 0,
    is_error: false,
    total_cost_usd: 0,
  },
  {
    role: "assistant",
    content: "- **Dark gradient background** with subtle depth",
    model: "",
    duration: 0,
    is_error: false,
    total_cost_usd: 0,
  },
  {
    role: "assistant",
    content: "- **Animated floating text** with gradient coloring",
    model: "",
    duration: 0,
    is_error: false,
    total_cost_usd: 0,
  },
  {
    role: "metadata",
    content: "",
    model: "Sonnet 3.7",
    duration: 140,
    is_error: false,
    total_cost_usd: 0.34,
  },
];

export default function ChatFeed() {
  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-[10px]">
        {data?.map((m: any, i) => (
          <div key={i}>
            {m.role === "user" ? (
              <div className="flex justify-end">
                <div className="bg-secondary/60 px-[15px] py-[10px] rounded-md">
                  <p>{m.content}</p>
                </div>
              </div>
            ) : (
              <>
                {m.role === "assistant" &&
                  (i === 0 || data[i - 1].role !== "assistant") && (
                    <div className="flex flex-col gap-[10px] py-[20px]">
                      <div className="flex items-center gap-5">
                        <ZustackLogo size={35} />
                        <p className="text-xl font-semibold text-secondary-foreground">
                          Zustack
                        </p>
                      </div>
                    </div>
                  )}
                <div className="py-[5px]">
                  <Markdown content={m.content} />
                </div>
                {m.role === "metadata" && (
                  <div className="flex items-center gap-[5px] text-muted-foreground">
                    <Info />
                    <p>Worked for {m.duration} miliseconds</p>
                    <p>and consumed USD {m.total_cost_usd}</p>
                  </div>
                )}
              </>
            )}
          </div>
        ))}

        <div className="flex items-center gap-2 text-muted-foreground py-[30px]">
          <Spinner />
          Thinking...
        </div>
      </div>
      <div>
        <div className="w-full space-y-[10px] p-[10px] bg-secondary/50">
          <Input
            placeholder="Ask Zustack AI to build..."
            className="resize-none"
          />
          <div className="flex justify-end">
            <Button variant="outline">
              <span>Send</span>
              <Send />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
