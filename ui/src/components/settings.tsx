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
import { Settings as SettingsIcon } from "lucide-react";

export default function Settings() {
  const [isOpen, setIsOpen] = useState(false);

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
        <SettingsIcon />
      </Button>
      <AlertDialogContent className="w-[1000px]">
        <AlertDialogHeader>
          <AlertDialogTitle>Project Settings</AlertDialogTitle>
          <AlertDialogDescription>
            Manage your project details, visibility, and preferences.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <Button>Save changes</Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
