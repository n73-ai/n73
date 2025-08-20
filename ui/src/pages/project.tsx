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
import { useState } from "react";
import axios from "axios";

export default function Project() {
  const { projectID } = useParams();
  const [iframeStatus, setIframeStatus] = useState('loading');

  const { data, isLoading, isError, dataUpdatedAt } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID!),
  });

  const { isError: isErrorIframe } = useQuery({
    queryKey: ["iframe-status", data?.domain],
    queryFn: async () => {
      const res = await axios.head(data?.domain); // HEAD = solo headers, más rápido
      return res.status;
    },
    retry: false, // no reintentar si falla
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
              {(data?.status == "Deployed" ||
                data?.status == "Building" ||
                data?.status == "Error") && !isErrorIframe && (
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

              {isErrorIframe && (
                <Stars isIframeError={isErrorIframe} />
              )}

              {(data?.status == "Building-First" ||
                data?.status == "Building-First-Error") && <Stars isIframeError={isErrorIframe} />}
            </ResizablePanel>
          </ResizablePanelGroup>
        </div>
      )}
    </div>
  );
}
