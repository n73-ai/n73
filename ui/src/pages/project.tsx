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
import axios from "axios";
import { CloudOffIcon } from "lucide-react";

export default function Project() {
  const { projectID } = useParams();

  const { data, isLoading, isError, dataUpdatedAt } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID!),
    refetchOnWindowFocus: false,
  });

  const { isError: isErrorIframe } = useQuery({
    queryKey: ["iframe-status", data?.domain],
    queryFn: async () => {
      const res = await axios.head(data?.domain);
      return res.status;
    },
    retry: false,
  });

  console.log(data);

  return (
    <div className="h-[100dvh] flex flex-col">
      <ProjectNavbar />
      {isError && <p>An unexpected error occurred.</p>}
      {isLoading && <Spinner />}
      <div className="flex-1 min-h-0">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel defaultSize={40}>
            {data && <ChatFeed p={data} />}
          </ResizablePanel>
          <ResizableHandle />

          <ResizablePanel className="hidden lg:block">
            {(data?.status == "new_pending" || data?.status == "new_error") && (
              <Stars status={data?.status} isIframeError={isErrorIframe} />
            )}

            {/*
            {(data?.status == "pending" ||
              data?.status == "idle" ||
              data?.status == "internal_error" ||
              data?.status == "error") &&
              !isErrorIframe && (
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
                        //src={`https://${data.fly_hostname}`}
                        src="https://e1fb4217-7b30-49a5-9729-8c5f3edaf4e2.fly.dev/"
                      />
                    )}
                  </div>
                </div>
              )}
              */}

            <iframe
              key={dataUpdatedAt}
              className="w-full h-full block"
              src={`https://${data?.fly_hostname}`}
            />

            {data?.status == "new_internal_error" && (
              <div className="relative z-10 flex items-center justify-center h-full">
                <div className="flex gap-[5px] text-center text-muted-foreground">
                  <CloudOffIcon />
                  Oops, it looks like the project isn’t online.
                </div>
              </div>
            )}
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}
