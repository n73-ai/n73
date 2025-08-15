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
import { getProjectByID } from "@/api/projects";
import Stars from "@/components/stars";

export default function Project() {
  const { projectID } = useParams();

  const { data, isLoading, isError, dataUpdatedAt } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID!),
  });

  return (
    <div className="h-screen flex flex-col">
      <ProjectNavbar />
      {isError && <p>An unexpected error occurred.</p>}
      {isLoading ? (
        <Spinner />
      ) : (
        <div className="flex-1 min-h-0">
          <ResizablePanelGroup direction="horizontal">
            <ResizablePanel defaultSize={40}>
              <ChatFeed
                pStatus={data?.status}
                domain={data?.domain}
                slug={data?.slug}
              />

            </ResizablePanel>
            <ResizableHandle />
            <ResizablePanel>
              {(data?.status == "Deployed" || data?.status == "Building") && (
                <div className="flex flex-col h-full">
                  <div className="flex-1">
                    {data?.domain == "" ? (
                      <p className="text-center mt-[400px]">
                        The project is deployed but not domain name was found.
                      </p>
                    ) : (
                      <iframe
                        key={dataUpdatedAt}
                        className="w-full h-full block"
                        src={data.domain}
                      />
                    )}
                  </div>
                </div>
              )}

              {data?.status == "Building-First" && <Stars />}
            </ResizablePanel>
          </ResizablePanelGroup>
        </div>
      )}
    </div>
  );
}
