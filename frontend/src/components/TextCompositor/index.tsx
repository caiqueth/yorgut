import { Box, FormHelperText, Stack, TextField, Typography } from "@mui/material";
import { Button } from "components/Button";
import { FormattingButtons } from "components/FormattingDialogs";
import React from "react";

interface IProps {
  content: string;
  header?: string;
  onChangeContent: (value: string) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  isPending: boolean;
  maxLength: number;
  // Extra condition to send, besides a non-blank text (e.g. a topic's title).
  canSubmit?: boolean;
  // The text is rendered with FormattedText: the send button goes to the
  // right of the box, with buttons for formatting tips and a preview under it.
  formatted?: boolean;
}
export const TextCompositor: React.FC<IProps> = ({
  content,
  header,
  onChangeContent,
  onSubmit,
  isPending,
  maxLength = 256,
  canSubmit = true,
  formatted = false,
}) => {

  // Text made only of spaces or line breaks counts as blank.
  const isBlank = content.trim() === "";
  const isDisabled = isBlank || !canSubmit;

  // The text is kept until the caller clears it (on success), so a failed
  // request doesn't lose what was written.
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isDisabled || isPending) return;
    onSubmit(e);
  };

  const handleChangeContent = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { value } = e.target;
    if (value.length > maxLength) return;
    onChangeContent(value);
  }

  const textBox = (
    <Box sx={{ width: "100%" }}>
      <TextField
        value={content}
        onChange={handleChangeContent}
        multiline
        rows={6}
        fullWidth
      ></TextField>
      <FormHelperText>{maxLength - content.length} caracteres restantes.</FormHelperText>
    </Box>
  );
  const sendButton = <Button type="submit" loading={isPending} disabled={isDisabled}>enviar</Button>;

  if (!formatted) {
    return (
      <form onSubmit={handleSubmit}>
        <Typography>{header}</Typography>
        {textBox}
        {sendButton}
      </form>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      <Typography>{header}</Typography>
      {/* On phones the buttons go under the box, in a row that wraps. */}
      <Stack direction={{ xs: "column", sm: "row" }} alignItems={{ xs: "stretch", sm: "flex-start" }}>
        {textBox}
        <Stack direction={{ xs: "row", sm: "column" }} flexWrap={"wrap"} sx={{ flexShrink: 0 }}>
          {sendButton}
          <FormattingButtons text={content} />
        </Stack>
      </Stack>
    </form>
  );
}
