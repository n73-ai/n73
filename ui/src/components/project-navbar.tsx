import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronLeft, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link, useParams } from "react-router-dom";
import { useAuthStore } from "@/store/auth";
import SettingsDialog from "./settings-dialog";
import { useMutation } from "@tanstack/react-query";
import { publishProject } from "@/api/projects";
import toast from "react-hot-toast";
import type { ErrorResponse } from "@/lib/types";
import Spinner from "./spinner";

export default function ProjectNavbar() {
  const { logout, email } = useAuthStore();
  const { projectID } = useParams()

  const publishProjectMutation = useMutation({
    mutationFn: () => publishProject(projectID!),
    onSuccess: () => {
      toast.success("Your project has been publish!")
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  return (
    <nav className="px-[10px] h-[60px] flex items-center bg-secondary/50">
      <div className="flex justify-between w-full">
        <div className="flex gap-[10px]">
          <Link to="/">
            <Button variant="outline" className="flex gap-[5px]">
              <ChevronLeft />
            </Button>
          </Link>
          <SettingsDialog />
          <Button 
            variant="outline"
            onClick={() => publishProjectMutation.mutate()}
            >
            Publish
            {publishProjectMutation.isPending && (
              <Spinner />
            )}
          </Button>
        </div>
        <div className="flex gap-[10px]">
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="outline" className="flex gap-[5px]">
                <User />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>{email}</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => logout()}>
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </nav>
  );
}
