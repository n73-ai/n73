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
// import axios from "axios";
import { AlertCircleIcon, CloudOffIcon } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";

export default function Project() {
  const { projectID } = useParams();
  const [iframeIsLoading, setIsLoading] = useState(true);
  const [iframeKey, setIframeKey] = useState(0);
  const prevStatusRef = useRef<string | undefined>(undefined);
  const fallbackTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const { data, isLoading, isError } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID!),
    refetchOnWindowFocus: false,
  });

  const reloadIframe = useCallback(() => {
    if (fallbackTimerRef.current) {
      clearTimeout(fallbackTimerRef.current);
      fallbackTimerRef.current = null;
    }
    setIsLoading(true);
    setIframeKey((k) => k + 1);
  }, []);

  // Primary: "deployed" WS message triggers reloadIframe via ChatFeed callback.
  // Fallback: if "deployed" is missed (WS disconnect, crash), reload after 15s
  // once status transitions to "idle". The delay gives the CDN time to propagate.
  useEffect(() => {
    const prev = prevStatusRef.current;
    const curr = data?.status;
    prevStatusRef.current = curr;

    const wasWorking = prev === "pending" || prev === "new_pending" || prev === "new_error";
    if (wasWorking && curr === "idle" && !data?.error_msg) {
      fallbackTimerRef.current = setTimeout(reloadIframe, 15000);
      return () => {
        if (fallbackTimerRef.current) clearTimeout(fallbackTimerRef.current);
      };
    }
  }, [data?.status, data?.error_msg, reloadIframe]);

  /*
  const { isError: isErrorIframe } = useQuery({
    queryKey: ["iframe-status", data?.fly_hostname],
    queryFn: async () => {
      const res = await axios.head(`https://${data?.domain}`);
      return res.status;
    },
    retry: false,
  });
  */

  return (
    <div className="h-[100dvh] flex flex-col">
      <ProjectNavbar />

      <div className="flex-1 min-h-0">
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel defaultSize={40}>
            {isError && <p>An unexpected error occurred.</p>}
            {isLoading && (
              <div className="flex items-center gap-2 text-muted-foreground py-[20px]">
                <Spinner /> Loading chat
              </div>
            )}
            {data && <ChatFeed p={data} onDeployed={reloadIframe} iframeKey={iframeKey} />}
          </ResizablePanel>
          <ResizableHandle />

          <ResizablePanel className="hidden lg:block">
            {(data?.status == "new_pending" || data?.status == "new_error" ||
              (data?.status == "pending" && !data?.built)) && (
              <div className="relative w-full h-full">
                <Stars />
                <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                  <h5 className="text-xl flex items-center gap-2 text-muted-foreground py-[30px]">
                    <Spinner />
                    Building your project
                  </h5>
                </div>
              </div>
            )}

            {data?.error_msg != null && data?.error_msg != "" && data?.built &&
              data?.status !== "new_pending" && data?.status !== "new_error" &&
              data?.status !== "pending" && (
              <div className="flex flex-col items-center justify-center h-full gap-3 text-center px-8">
                <AlertCircleIcon className="text-destructive h-10 w-10" />
                <div>
                  <p className="font-semibold">Deployment failed</p>
                  <p className="text-sm text-muted-foreground mt-1">
                    Click "Try to fix" in the chat to ask n83 to fix the error.
                  </p>
                </div>
              </div>
            )}

            {!data?.error_msg && data?.built && data?.status !== "new_internal_error" && (
              <>
                {iframeIsLoading && (
                  <div className="relative z-10 flex items-center justify-center h-full">
                    <div className="text-xl flex items-center gap-2 text-muted-foreground py-[30px]">
                      <Spinner />
                    </div>
                  </div>
                )}

                {data?.domain && (
                  <iframe
                    key={iframeKey}
                    onLoad={() => setIsLoading(false)}
                    style={{ display: iframeIsLoading ? "none" : "block" }}
                    className="w-full h-full block"
                    src={`https://${data.domain}?_t=${iframeKey}`}
                  />
                )}
              </>
            )}

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
