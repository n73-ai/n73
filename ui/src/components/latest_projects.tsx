import { getLatestProjects } from "@/api/projects";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import Spinner from "./spinner";
import { GithubIcon, LinkIcon } from "lucide-react";

export default function LatestProjects() {
  const { data, isLoading, isError } = useQuery<any>({
    queryKey: ["latest-projects"],
    queryFn: () => getLatestProjects(),
  });

  return (
    <div className="py-[100px]">
      <h1
        className="scroll-m-20 text-start 
      text-2xl font-bold tracking-tight text-balance pb-[20px]"
      >
        From the Community
      </h1>
      <div className="bg-secondary/80 rounded-lg p-[10px]">
        {isLoading && <Spinner />}
        {isError && <p>An unexpected error occurred.</p>}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-[15px]">
          {data?.map((p: any) => (
            <Card className="w-full @container/card hover:border hover:border-muted-foreground/30 transition-all duration-200 ease-in-out">
              <CardHeader>
                <CardTitle className="pb-[10px]">{p.name}</CardTitle>
                <div className="hidden xl:block relative overflow-hidden w-full h-48 rounded-md">
                  <iframe
                    src={p.domain}
                    className="absolute top-0 left-0 w-[1280px] h-[800px] scale-[0.25] origin-top-left pointer-events-none"
                  />
                </div>
                <div className="flex justify-end gap-[5px] pt-[10px]">
                  <a
                    className="cursor-pointer inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-9 border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
                    href={p.gh_repo}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <GithubIcon />
                  </a>
                  <a
                    className="cursor-pointer inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive size-9 border bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50"
                    href={p.domain}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <LinkIcon />
                  </a>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}
