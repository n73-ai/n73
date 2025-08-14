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
import { TrashIcon } from "lucide-react";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useMutation } from "@tanstack/react-query";
import { deleteProject } from "@/api/projects";
import { useNavigate, useParams } from "react-router-dom";
import type { ErrorResponse } from "@/lib/types";
import toast from "react-hot-toast";
import Spinner from "./spinner";

export default function DeleteProject() {
  const [isOpen, setIsOpen] = useState(false);
  const {projectID} = useParams()
  const navigate = useNavigate()

  const deleteProjectMutation = useMutation({
    mutationFn: () => deleteProject(projectID!),
    onSuccess: () => {
      setIsOpen(false)
      navigate("/")
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
      <TrashIcon/>
      </Button>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Project</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to delete this project?
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel
          disabled={deleteProjectMutation.isPending}
          >Cancel</AlertDialogCancel>
          <Button
          variant="destructive"
          disabled={deleteProjectMutation.isPending}
          onClick={() => {
            deleteProjectMutation.mutate()
          }}
          >Delete project
          {deleteProjectMutation.isPending && <Spinner />}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
