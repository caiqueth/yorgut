import React from "react";
import { Dialog, DialogActions, DialogContent, DialogTitle, Typography } from "@mui/material";
import { Button } from "components/Button";

interface IConfirmDialogProps {
  open: boolean;
  title: string;
  text: React.ReactNode;
  // The confirm button's label ("desfazer", "sair").
  confirmLabel: string;
  onConfirm: () => void;
  onClose: () => void;
}

// Asks before an action that's easy to hit by mistake. Confirming closes it.
export const ConfirmDialog: React.FC<IConfirmDialogProps> = ({
  open, title, text, confirmLabel, onConfirm, onClose,
}) => (
  <Dialog open={open} onClose={onClose} maxWidth={"xs"} fullWidth>
    <DialogTitle>{title}</DialogTitle>
    <DialogContent>
      <Typography variant={"body2"}>{text}</Typography>
    </DialogContent>
    <DialogActions>
      <Button onClick={onClose}>cancelar</Button>
      <Button onClick={() => { onClose(); onConfirm(); }}>{confirmLabel}</Button>
    </DialogActions>
  </Dialog>
);
