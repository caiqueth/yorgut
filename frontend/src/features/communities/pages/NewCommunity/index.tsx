import { Box, Checkbox, FormControlLabel, Stack, TextField, Typography } from "@mui/material";
import { useCreateCommunity } from "api/community/queries";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { CommunityCategorySelect } from "components/CommunityCategorySelect";
import { FormattingButtons } from "components/FormattingDialogs";
import { CommunityCategory } from "api/generated/model";
import { PICTURE_RULES_TEXT } from "helpers/validatePicture";
import { usePicturePicker } from "hooks/usePicturePicker";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export const NewCommunity: React.FC = () => {
  const navigate = useNavigate();
  const { mutate: createCommunity, isPending, error } = useCreateCommunity();
  const [name, setName] = useState("");
  const [category, setCategory] = useState<CommunityCategory | "">("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const { inputRef, picture, preview, error: pictureError, pick, clear } = usePicturePicker();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createCommunity(
      { name, category, description, private: isPrivate, picture: picture ?? undefined },
      { onSuccess: (response) => navigate(`/community/${response.data!.id}`) },
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
                label={"Nome"}
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                slotProps={{ htmlInput: { maxLength: 100 } }}
              />
            </Box>
            <Box>
              <Typography>Qual a categoria da comunidade?</Typography>
              <CommunityCategorySelect label={"Categoria"} value={category} onChange={setCategory} />
            </Box>
            <Box>
              <Typography>Descreva a comunidade (opcional)</Typography>
              <TextField
                label={"Descrição"}
                multiline
                minRows={4}
                fullWidth
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                slotProps={{ htmlInput: { maxLength: 2048 } }}
              />
              <FormattingButtons text={description} />
            </Box>
            <Box>
              <Typography>Imagem da comunidade (opcional)</Typography>
              <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
              <Stack direction={"row"} alignItems={"center"} gap={1} sx={{ mt: 0.5 }}>
                <input ref={inputRef} type={"file"} style={{ maxWidth: "100%" }} accept={"image/png,image/jpeg"} onChange={pick} />
                {picture && <Button onClick={clear}>remover</Button>}
              </Stack>
              {preview && <img src={preview} height={120} alt={"pré-visualização"} style={{ marginTop: 8 }} />}
              {pictureError && (
                <Typography variant={"body2"} color={"error"}>{pictureError}</Typography>
              )}
            </Box>
            <Box>
              <FormControlLabel
                control={
                  <Checkbox size={"small"} checked={isPrivate} onChange={(e) => setIsPrivate(e.target.checked)} />
                }
                label={"Comunidade privada: só membros veem o fórum, e você aprova quem entra"}
              />
            </Box>
            {error && (
              <Typography variant={"body2"} color={"error"}>
                {error.message}
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
