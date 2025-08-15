import { getUserProjects } from "@/api/projects";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import Spinner from "./spinner";

export default function Projects() {
  const { data, isLoading, isError } = useQuery<any>({
    queryKey: ["user-projects"],
    queryFn: () => getUserProjects(),
  });

  return (
    <div className="pt-[100px]">
      <h1
        className="scroll-m-20 text-start 
      text-2xl font-bold tracking-tight text-balance pb-[20px]"
      >
        Your Projects
      </h1>
      {isLoading && <Spinner />}
      {isError && <p>An unexpected error occurred.</p>}
      {!data && !isLoading && !isError && (
        <div className="border p-[20px] px-[25px] rounded-md">
          <h1
            className="text-center
            tracking-tight text-balance text-muted-foreground"
          >
            You don't have any projects... yet!
          </h1>
        </div>
      )}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-[15px]">
        {data?.map((p: any) => (
          <Link to={`/project/${p.id}`}>
            <Card className="@container/card hover:border hover:border-muted-foreground/30 transition-all duration-200 ease-in-out">
              <CardHeader>
                <CardTitle className="pb-[10px]">{p.name}</CardTitle>
              <div className="hidden xl:block relative overflow-hidden w-full h-48 rounded-md">
                <iframe
                  src={p.domain}
                  className="absolute top-0 left-0 w-[1280px] h-[800px] scale-[0.25] origin-top-left pointer-events-none"
                />
              </div>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
