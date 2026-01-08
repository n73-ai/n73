import * as React from "react";
import { Settings, SquareArrowOutUpRightIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Separator } from "./ui/separator";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  deleteProject,
  editProject,
  getProjectByID,
  updateProjectOwner,
} from "@/api/projects";
import type { ErrorResponse } from "@/lib/types";
import toast from "react-hot-toast";
import Spinner from "./spinner";
import { useAuthStore } from "@/store/auth";

export default function SettingsDialog() {
  const [isOpen, setIsOpen] = React.useState(false);
  const [name, setName] = React.useState("");
  const [owner, setOwner] = React.useState("");
  const { email } = useAuthStore();

  const { projectID } = useParams();

  const editProjectMutation = useMutation({
    mutationFn: () => editProject(projectID!, name),
    onSuccess: () => {
      //setIsOpen(false);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const navigate = useNavigate();

  const deleteProjectMutation = useMutation({
    mutationFn: () => deleteProject(projectID!),
    onSuccess: () => {
      setIsOpen(false);
      navigate("/");
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const updateProjectOwnerMutation = useMutation({
    mutationFn: () => updateProjectOwner(projectID!, owner),
    onSuccess: () => {
      setIsOpen(false);
      navigate("/");
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const { data, isLoading, isError } = useQuery({
    queryKey: ["project", projectID],
    queryFn: () => getProjectByID(projectID!),
  });

  React.useEffect(() => {
    if (data) {
      setName(data.name);
      setOwner(email);
    }
  }, [data]);

  return (
    <Dialog onOpenChange={(open: boolean) => setIsOpen(open)} open={isOpen}>
      <DialogTrigger asChild>
        <Button
          onClick={() => setIsOpen(!isOpen)}
          variant="outline"
          className="flex gap-[5px]"
        >
          <Settings />
        </Button>
      </DialogTrigger>
      <DialogContent
        onOpenAutoFocus={(e) => e.preventDefault()}
        className="overflow-hidden p-0 md:max-h-[800px] md:max-w-[800px] lg:max-w-[1000px]"
      >
        <DialogTitle className="sr-only">Settings</DialogTitle>
        <DialogDescription className="sr-only">
          Customize your settings here.
        </DialogDescription>
        <main className="flex h-auto flex-1 flex-col overflow-hidden">
          <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
            <div className="flex items-center gap-2 px-4">
              <h1 className="text-xl font-semibold">Project Settings</h1>
            </div>
          </header>
          <div className="flex flex-1 flex-col gap-4 overflow-y-auto p-4 pt-0">
            <Separator className="my-[2px]" />

            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm font-medium">
                Project name
              </Label>
              {isError ? (
                <div className="flex justify-center">
                  <p>An unexpected error occurred.</p>
                </div>
              ) : isLoading ? (
                <div className="flex justify-center">
                  <Spinner />
                </div>
              ) : (
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  id="name"
                  type="text"
                  placeholder="Project name"
                />
              )}
              <div className="flex justify-end">
                <Button
                  variant="secondary"
                  disabled={editProjectMutation.isPending}
                  onClick={() => {
                    editProjectMutation.mutate();
                  }}
                >
                  Save
                  {editProjectMutation.isPending && <Spinner />}
                </Button>
              </div>
            </div>

            <Separator className="my-[2px]" />

            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm font-medium">
                Ownership
              </Label>
              <Input
                value={owner}
                onChange={(e) => setOwner(e.target.value)}
                id="name"
                type="text"
                placeholder="name@domain.com"
              />
              <div className="flex justify-end">
                <Button
                  variant="secondary"
                  disabled={updateProjectOwnerMutation.isPending}
                  onClick={() => {
                    updateProjectOwnerMutation.mutate();
                  }}
                >
                  Save
                  {updateProjectOwnerMutation.isPending && <Spinner />}
                </Button>
              </div>
            </div>

            <Separator className="my-[2px]" />

            <div className="space-y-2">
              <div className="flex justify-between">
                <Label htmlFor="name" className="text-sm font-medium">
                  Domain
                </Label>
                <Button
                  onClick={() => {
                    window.open(
                      "https://developers.cloudflare.com/pages/configuration/custom-domains/",
                      "_blank"
                    );
                  }}
                  variant="secondary"
                >
                  Configure custom domain
                  <SquareArrowOutUpRightIcon />
                </Button>
              </div>
              <Button
                variant="link"
                onClick={() => {
                  window.open(data?.domain, "_blank");
                }}
              >
                {data?.domain && (
                  <>
                  <span>{data.domain}</span>
                  <SquareArrowOutUpRightIcon />
                  </>
                )}
              </Button>
            </div>

            <Separator className="my-[2px]" />

            <div className="flex justify-between items-end">
              <div className="space-y-2">
                <Label htmlFor="name" className="text-sm font-medium">
                  Delete project
                </Label>
                <p className="text-sm text-muted-foreground">
                  Permanently delete this project.
                </p>
              </div>
              <Button
                variant="destructive"
                disabled={deleteProjectMutation.isPending}
                onClick={() => {
                  deleteProjectMutation.mutate();
                }}
              >
                Delete project
                {deleteProjectMutation.isPending && <Spinner />}
              </Button>
            </div>
          </div>
        </main>
      </DialogContent>
    </Dialog>
  );
}
