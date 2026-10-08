import { Box, Stack, TextField, Typography } from "@mui/material";
import { useCreateCommunity } from "api/community/queries";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { PICTURE_RULES_TEXT, validatePicture } from "helpers/validatePicture";
import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

interface IProps { }
export const NewCommunity: React.FC<IProps> = () => {
  const navigate = useNavigate();
  const { mutate: createCommunity, isPending, isError } = useCreateCommunity();
  const [name, setName] = useState("");
  const [category, setCategory] = useState("");
  const [picture, setPicture] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [pictureError, setPictureError] = useState<string | null>(null);
  const fileInput = useRef<HTMLInputElement>(null);

  // Release the local preview URL when it's replaced or the page closes.
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);

  const handlePickPicture = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    setPictureError(null);
    setPicture(null);
    setPreview(null);
    if (!file) return;

    const error = await validatePicture(file);
    if (error) {
      setPictureError(error);
      if (fileInput.current) fileInput.current.value = "";
      return;
    }
    setPicture(file);
    setPreview(URL.createObjectURL(file));
  };

  const handleRemovePicture = () => {
    setPicture(null);
    setPreview(null);
    if (fileInput.current) fileInput.current.value = "";
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createCommunity(
      { name, category, picture },
      { onSuccess: (response) => navigate(`/community/${response.data.id}`) },
    );
  }

  return (
    <BoxWithTitle title={"nova comunidade"}>
      <Box sx={{ p: 1 }}>
        <form onSubmit={handleSubmit}>
          <Stack direction={"column"} width={"100%"} spacing={1}>
            <Box>
              <Typography>Dê um nome à sua comunidade</Typography>
              <TextField
                size="small"
                label={"Nome"}
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                slotProps={{ htmlInput: { maxLength: 100 } }}
              />
            </Box>
            <Box>
              <Typography>Qual a categoria da comunidade?</Typography>
              <TextField
                size="small"
                label={"Categoria"}
                required
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                slotProps={{ htmlInput: { maxLength: 100 } }}
              />
            </Box>
            <Box>
              <Typography>Imagem da comunidade (opcional)</Typography>
              <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
              <Stack direction={"row"} alignItems={"center"} gap={1} sx={{ mt: 0.5 }}>
                <input ref={fileInput} type={"file"} accept={"image/png,image/jpeg"} onChange={handlePickPicture} />
                {picture && <Button onClick={handleRemovePicture}>remover</Button>}
              </Stack>
              {preview && <img src={preview} height={120} alt={"pré-visualização"} style={{ marginTop: 8 }} />}
              {pictureError && (
                <Typography variant={"body2"} color={"error"}>{pictureError}</Typography>
              )}
            </Box>
            {isError && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível criar a comunidade. Tente novamente.
              </Typography>
            )}
            <Box>
              <Button type={"submit"} loading={isPending}>Criar</Button>
            </Box>
          </Stack>
        </form>
      </Box>
    </BoxWithTitle>
  )
}
