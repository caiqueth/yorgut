import { Box, Stack, TextField } from "@mui/material";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import React from "react";

interface IProps { }
export const NewCommunity: React.FC<IProps> = () => {

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("submit");
  }

  return (
    <BoxWithTitle title={"nova comunidade"}>
      <Box sx={{ p: 1 }}>
        <form onSubmit={handleSubmit}>
          <Stack direction={"column"} width={500} spacing={2}>
            <TextField label={"Nome"} />
            <TextField label={"Categoria"} />
            <Button type={"submit"}>Criar</Button>
          </Stack>
        </form>
      </Box>
    </BoxWithTitle>
  )
}