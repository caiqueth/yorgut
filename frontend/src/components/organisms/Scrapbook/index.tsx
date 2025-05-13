import {
  Box,
  FormHelperText,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useCreateScrap, useScraps } from "api/scrap/queries";

import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";

interface IProps {
  userId: number;
}
export const Scrapbook: React.FC<IProps> = ({ userId }) => {

  const { data: scraps } = useScraps(userId);
  const { mutate: sendScrap } = useCreateScrap();

  const [scrapContent, setScrapContent] = useState<string>("");

  const handleSubmit = (e: any) => {
    e.preventDefault();
    sendScrap({
      content: scrapContent,
      userRecipientId: userId ?? 0,
    });
  };

  const handleChangeText = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { value } = e.target;
    if (value.length > 2048) return;
    setScrapContent(value);
  }

  const items =
    scraps?.scraps.map((s) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={`/${s.authorProfilePicture}`} height={40} />
        <Stack direction={"column"}>
          <Typography>{s.authorName} em {dayjsFromDate(s.createdAt, "D [de] MMMM [de] YYYY")}:</Typography>
          <Typography>{s.content}</Typography>
        </Stack>
      </Stack>
    )) ?? [];

  return (
    <Stack direction={"column"} gap={1}>
      <form onSubmit={handleSubmit}>
        <Box sx={{ width: "100%" }}>
          <Typography>Deixar um recado</Typography>
          <TextField
            value={scrapContent}
            onChange={handleChangeText}
            multiline
            rows={6}
            fullWidth
          ></TextField>
          <FormHelperText>{2048 - scrapContent.length} caracteres restantes.</FormHelperText>
        </Box>
        <Button type="submit">Enviar</Button>
      </form>
      <BoxWithTitle title="meus recados">
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Stack>
  );
};
