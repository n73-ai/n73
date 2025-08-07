import ChatFeed from "@/components/chat-feed";
import ProjectNavbar from "@/components/project-navbar";
import { useParams } from "react-router-dom";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { useQuery } from "@tanstack/react-query";
import Spinner from "@/components/spinner";
import { Button } from "@/components/ui/button";
import { getProjectByID } from "@/api/projects";

export default function Project() {
  const { projectID } = useParams();

  const { data, isLoading, isError } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID),
  });

  return (
    <div className="h-screen flex flex-col">
      <ProjectNavbar />
      <div className="flex-1 min-h-0">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel defaultSize={40}>
            <ChatFeed pStatus={data?.status} />
          </ResizablePanel>
          <ResizableHandle />
          <ResizablePanel>
            {false && (
              <div className="flex justify-center items-center min-h-screen">
                <h5 className="text-xl text-muted-foreground">Building</h5>
                <Spinner />
              </div>
            )}

            <div className="flex justify-center items-center min-h-screen">
              <Button variant="outline">
                Deploy Project 
              </Button>
            </div>

            <div className="flex flex-col h-full">
              <div className="flex-1">
                <iframe className="w-full h-full block" src="" />
              </div>
            </div>

          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
