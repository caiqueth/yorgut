import { Box, FormHelperText, TextField, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import React from "react";

interface IProps {
  content: string;
  header?: string;
  onChangeContent: (value: string) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  maxLength: number;
}
export const TextCompositor: React.FC<IProps> = ({
  content,
  header,
  onChangeContent,
  onSubmit,
  maxLength = 256,
}) => {

  const handleSubmit = (e: any) => {
    e.preventDefault();
    onSubmit(e);
  };

  const handleChangeContent = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { value } = e.target;
    if (value.length > maxLength) return;
    onChangeContent(value);
  }

  return (

    <form onSubmit={handleSubmit}>
      <Box sx={{ width: "100%" }}>
        <Typography>{header}</Typography>
        <TextField
          value={content}
          onChange={handleChangeContent}
          multiline
          rows={6}
          fullWidth
        ></TextField>
        <FormHelperText>{maxLength - content.length} caracteres restantes.</FormHelperText>
      </Box>
      <Button type="submit">Enviar</Button>
    </form>

  )
}