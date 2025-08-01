import { Info, Send } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import ZustackLogo from "./zustack-logo";
import Spinner from "./spinner";
import Markdown from "./markdown";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { getMessagesByProjectID } from "@/api/messages";

export default function ChatFeed() {

  const {projectID} = useParams()

  const { data, isLoading, isError } = useQuery<any[]>({
    queryKey: ["messages", projectID],
    queryFn: () => getMessagesByProjectID(projectID),
  });

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-[10px]">
        {isLoading && (
          <p>Loading..</p>
        )}
        {isError && (
          <p>error..</p>
        )}

        {data?.map((m:any, i) => (
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
                      <div className="flex items-center gap-[5px] text-muted-foreground">
                        <Info />
                        <p>Worked for {m.metadata?.duration} miliseconds</p>
                        <p>and consumed USD {m.metadata?.total_cost_usd}</p>
                      </div>
                    </div>
                  )}
                <div className="">
                  <Markdown content={m.content} />
                </div>
              </>
            )}
          </div>
        ))}

        <div className="flex items-center gap-2 text-muted-foreground pt-[10px]">
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
