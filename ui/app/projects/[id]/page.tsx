import ChatFeed from "@/components/chat_feed";
import Preview from "@/components/preview";
import ProjectNavbar from "@/components/navigation/project-navbar";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

export default function Page() {
  return (
    <div className="h-screen flex flex-col">
      <ProjectNavbar />
      <div className="flex-1 min-h-0">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel defaultSize={30}>
            <ChatFeed />
          </ResizablePanel>
          <ResizableHandle />
          <ResizablePanel>
            <Preview />
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
