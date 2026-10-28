import React from "react";
import {
  Box,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import { Button } from "components/Button";
import { FormattedText } from "components/FormattedText";

// Each tip is what to type, rendered next to it by FormattedText itself, so
// the examples can't drift from what the text really does.
const tips: string[] = [
  "**negrito** ou <b>negrito</b>",
  "*itálico* ou <i>itálico</i>",
  "<u>sublinhado</u>",
  "~~riscado~~ ou <s>riscado</s>",
  "[um link](https://exemplo.com)",
  "- um item\n- outro item",
  "1. primeiro\n2. segundo",
  "> uma citação",
  "`código`",
];

const codeSx = { fontFamily: "monospace", whiteSpace: "pre-wrap", m: 0 };

interface IDialogProps {
  open: boolean;
  onClose: () => void;
}

// What formatting user text may use (see FormattedText).
export const FormattingTipsDialog: React.FC<IDialogProps> = ({ open, onClose }) => (
  <Dialog open={open} onClose={onClose} maxWidth={"sm"} fullWidth>
    <DialogTitle>dicas de formatação</DialogTitle>
    <DialogContent>
      <Typography variant={"body2"} gutterBottom>
        Dá para formatar o texto com markdown ou com HTML simples:
      </Typography>
      <Table size={"small"}>
        <TableHead>
          <TableRow>
            <TableCell>você escreve</TableCell>
            <TableCell>aparece assim</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {tips.map((tip) => (
            <TableRow key={tip}>
              <TableCell><Box component={"pre"} sx={codeSx}>{tip}</Box></TableCell>
              <TableCell><FormattedText text={tip} /></TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Typography variant={"body2"} sx={{ mt: 1 }}>
        Uma quebra de linha vira uma nova linha. Imagens, títulos e outros elementos HTML são removidos, e links só
        funcionam com https://, http:// ou mailto:.
      </Typography>
    </DialogContent>
    <DialogActions>
      <Button onClick={onClose}>fechar</Button>
    </DialogActions>
  </Dialog>
);

interface IPreviewDialogProps extends IDialogProps {
  title: string;
  text: string;
}

// The text as it will look once sent.
export const FormattingPreviewDialog: React.FC<IPreviewDialogProps> = ({ open, onClose, title, text }) => (
  <Dialog open={open} onClose={onClose} maxWidth={"sm"} fullWidth>
    <DialogTitle>{title}</DialogTitle>
    <DialogContent>
      <Box sx={{ backgroundColor: "yorgut.highlight", p: 1 }}>
        {text.trim()
          ? <Typography component={"div"}><FormattedText text={text} /></Typography>
          : <Typography variant={"body2"}>Escreva algo para visualizar.</Typography>}
      </Box>
    </DialogContent>
    <DialogActions>
      <Button onClick={onClose}>fechar</Button>
    </DialogActions>
  </Dialog>
);

interface IFormattingButtonsProps {
  // The text being written, for the preview.
  text: string;
}

// The "dicas de formatação" and "visualizar" buttons, each opening its
// dialog, for a field whose text is shown with FormattedText.
export const FormattingButtons: React.FC<IFormattingButtonsProps> = ({ text }) => {
  const [open, setOpen] = React.useState<"tips" | "preview" | null>(null);
  const close = () => setOpen(null);
  return (
    <>
      <Button onClick={() => setOpen("tips")}>dicas de formatação</Button>
      <Button onClick={() => setOpen("preview")}>visualizar</Button>
      <FormattingTipsDialog open={open === "tips"} onClose={close} />
      <FormattingPreviewDialog open={open === "preview"} onClose={close} title={"visualizar"} text={text} />
    </>
  );
};
