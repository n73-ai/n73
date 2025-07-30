import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send } from "lucide-react";

export default function MessageInput() {
  return (
    <div className="w-full space-y-[10px] p-[10px] bg-secondary/50">
      <Textarea
        placeholder="Ask Zustack AI..."
        className="resize-none"
      />
      <div className="flex justify-end">
        <Button variant="outline">
          <span>Send</span>
          <Send />
        </Button>
      </div>
    </div>
  );
}
