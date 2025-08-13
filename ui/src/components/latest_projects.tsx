import { getLatestProjects } from "@/api/projects";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import Spinner from "./spinner";

export default function LatestProjects() {
  const { data, isLoading, isError } = useQuery<any>({
    queryKey: ["latest-projects"],
    queryFn: () => getLatestProjects(),
  });

  return (
    <div className="pt-[100px]">
      <h1
        className="scroll-m-20 text-start 
      text-2xl font-bold tracking-tight text-balance pb-[20px]"
      >
        From the Community
      </h1>
      {isLoading && <Spinner />}
      {isError && <p>An unexpected error occurred.</p>}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-[15px]">
        {data?.map((p: any) => (
          <a target="_blank" rel="noopener noreferrer" href={p.domain}>
            <Card className="@container/card hover:border hover:border-zinc-700 transition-all duration-200 ease-in-out">
              <CardHeader>
                <CardTitle className="font-semibold">{p.name}</CardTitle>
              </CardHeader>
            </Card>
          </a>
        ))}
      </div>
    </div>
  );
}
