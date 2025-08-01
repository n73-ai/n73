import { Send } from "lucide-react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import ZustackLogo from "./zustack-logo";
import Spinner from "./spinner";

export default function ChatFeed() {
  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-[10px]">
        <div className="flex justify-end">
          <div className="bg-secondary/60 px-[15px] py-[10px] rounded-md">
            <p>Create a hello world</p>
          </div>
        </div>

        <div className="flex items-center gap-[5px]">
          <ZustackLogo size={35} />
          <p className="text-xl font-semibold text-secondary-foreground">
            Zustack
          </p>
        </div>

        <div className="flex items-center gap-2 text-muted-foreground pt-[20px]">
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
