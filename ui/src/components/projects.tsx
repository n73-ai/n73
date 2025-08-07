import { getUserProjects } from "@/api/projects";
import {
  Card,
  CardAction,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import Spinner from "./spinner";

export default function Projects() {
  const { data, isLoading, isError } = useQuery<any>({
    queryKey: ["user-projects"],
    queryFn: () => getUserProjects(),
  });

  return (
    <div className="py-[200px]">
      <h1
        className="scroll-m-20 text-start 
      text-3xl font-bold tracking-tight text-balance pb-[20px]"
      >
        Your Projects
      </h1>
      {isLoading && <Spinner />}
      {isError && <p>An unexpected error occurred.</p>}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-[15px]">
        {data?.map((p: any) => (
          <Link to={`/project/${p.id}`}>
            <Card className="@container/card hover:border hover:border-zinc-700 transition-all duration-200 ease-in-out">
              <CardHeader>
                <CardTitle className="font-semibold">{p.name}</CardTitle>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}
