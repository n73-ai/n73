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
import { deployProject, getProjectByID } from "@/api/projects";
import type { ErrorResponse } from "@/lib/types";
import toast from "react-hot-toast";

export default function Project() {
  const { projectID } = useParams();

  const { data, isLoading, isError } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID),
  });

  const deployProjectMutation = useMutation({
    mutationFn: () => deployProject(projectID),
    onSuccess: () => {},
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
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
            {data.status == "Deployed" ? (
              <div className="flex flex-col h-full">
                <div className="flex-1">
                  <iframe
                    className="w-full h-full block"
                    src="https://1c366a91.super-cool-app.pages.dev/"
                  />
                </div>
              </div>
            ) : (
              <>
                {deployProjectMutation.isPending && (
                  <div className="flex justify-center items-center min-h-screen">
                    <h5 className="text-xl text-muted-foreground">Building</h5>
                    <Spinner />
                  </div>
                )}

                {!deployProjectMutation.isPending && (
                  <div className="flex justify-center items-center min-h-screen">
                    <Button
                      onClick={() => {
                        deployProjectMutation.mutate();
                      }}
                      variant="outline"
                    >
                      Deploy Project
                    </Button>
                  </div>
                )}
              </>
            )}
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
