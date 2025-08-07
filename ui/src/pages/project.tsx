import ChatFeed from "@/components/chat-feed";
import ProjectNavbar from "@/components/project-navbar";
import { useParams } from "react-router-dom";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { useMutation, useQuery } from "@tanstack/react-query";
import Spinner from "@/components/spinner";
import { Button } from "@/components/ui/button";
import { getProjectByID } from "@/api/projects";
import type { ErrorResponse } from "@/lib/types";
import toast from "react-hot-toast";

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
            {data?.status == "Deployed" && (
              <div className="flex flex-col h-full">
                <div className="flex-1">
                  {data?.domain == "" ? (
                    <p>
                      The project is deployed but not domain name was found,
                      contact support
                    </p>
                  ) : (
                    <iframe className="w-full h-full block" src={data.domain} />
                  )}
                </div>
              </div>
            )}

            <div className="flex justify-center items-center min-h-screen gap-[10px]">
              <Spinner />
              <h5 className="text-xl text-muted-foreground">
                Spinning up preview...
              </h5>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
