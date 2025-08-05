import {
  Card,
  CardAction,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Link } from "react-router-dom";

const data = [
  {
    id: "69",
    status: "building",
    name: "hello-world",
  },
  {
    id: "420",
    status: "deployed",
    name: "lovably-recreated-spark",
  },
];

export default function Projects() {
  return (
    <div className="pt-[200px]">
      <h1 className="scroll-m-20 text-start 
      text-3xl font-bold tracking-tight text-balance pb-[20px]">
        Your Projects
      </h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-[15px]">
        {data?.map((p) => (
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
