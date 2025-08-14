import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useState } from "react";
import { Button } from "./ui/button";
import { EditIcon } from "lucide-react";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useParams } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import type { ErrorResponse } from "@/lib/types";
import { editProject } from "@/api/projects";
import toast from "react-hot-toast";
import Spinner from "./spinner";

export default function EditProject() {
  const [isOpen, setIsOpen] = useState(false);
  const [name, setName] = useState("");

  const { projectID } = useParams();

  const editProjectMutation = useMutation({
    mutationFn: () => editProject(projectID!, ""),
    onSuccess: () => {
      setIsOpen(false);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  return (
    <AlertDialog
      onOpenChange={(open: boolean) => setIsOpen(open)}
      open={isOpen}
    >
      <Button
        onClick={() => setIsOpen(!isOpen)}
        variant="outline"
        className="flex gap-[5px]"
      >
        <EditIcon />
      </Button>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Project</AlertDialogTitle>
          <AlertDialogDescription>
            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm font-medium">
                Project name
              </Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                id="name"
                type="text"
                placeholder="Project name"
              />
            </div>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button
            variant="secondary"
            disabled={editProjectMutation.isPending}
            onClick={() => {
              editProjectMutation.mutate();
            }}
          >
            Save changes
            {editProjectMutation.isPending && <Spinner />}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
