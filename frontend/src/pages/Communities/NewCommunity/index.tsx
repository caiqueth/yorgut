import { Box, Stack, TextField } from "@mui/material";
import { useCreateCommunity } from "api/community/queries";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import React, { useState } from "react";

interface IProps { }
export const NewCommunity: React.FC<IProps> = () => {
  const { mutate: createCommunity } = useCreateCommunity();
  const [name, setName] = useState("");
  const [category, setCategory] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createCommunity({
      name,
      category,
    })
  }

  const handleChangeName = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setName(e.target.value);
  }

  const handleChangeCategory = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setCategory(e.target.value);
  }

  return (
    <BoxWithTitle title={"nova comunidade"}>
      <Box sx={{ p: 1 }}>
        <form onSubmit={handleSubmit}>
          <Stack direction={"column"} width={500} spacing={2}>
            <TextField label={"Nome"} value={name} onChange={handleChangeName} />
            <TextField label={"Categoria"} value={category} onChange={handleChangeCategory} />
            <Button type={"submit"}>Criar</Button>
          </Stack>
        </form>
      </Box>
    </BoxWithTitle>
  )
}