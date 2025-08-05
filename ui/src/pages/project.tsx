import ChatFeed from "@/components/chat-feed";
import ProjectNavbar from "@/components/project-navbar";
import { useParams } from "react-router-dom";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

export default function Project() {
  const { projectID } = useParams();

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
            <div className="flex flex-col h-full">
              <div className="flex-1">
                <iframe
                  className="w-full h-full block"
                  src=""
                />
              </div>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
