import {
  Box,
  Grid2 as Grid,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { UserProfilePicture } from "components/organisms/UserProfilePicture";
import { Menu } from "pages/Home/components/Menu";
import React, { useState } from "react";
import { useParams } from "react-router-dom";

interface IProps { }
export const Scrapbook: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const { mutate: sendScrap } = useSendScrap();
  const { data: scraps } = useScraps(userId);

  const [scrapContent, setScrapContent] = useState<string>("");

  const handleSubmit = (e: any) => {
    e.preventDefault();
    sendScrap({
      userRecipientId: userId,
      content: scrapContent,
    });
  };

  const items =
    scraps?.scraps.map((s) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={s.authorProfilePicture} height={40} />
        <Typography>{s.content}</Typography>
      </Stack>
    )) ?? [];

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <UserProfilePicture />
            <Menu />
          </Stack>
        </Grid>
        <Grid size={10}>
          <Stack direction={"column"} gap={1}>
            <form onSubmit={handleSubmit}>
              <Box sx={{ width: "100%" }}>
                <Typography>Deixar um recado</Typography>
                <TextField
                  value={scrapContent}
                  onChange={(e) => setScrapContent(e.target.value)}
                  multiline
                  rows={6}
                  fullWidth
                ></TextField>
              </Box>
              <Button type="submit">Enviar</Button>
            </form>
            <BoxWithTitle title="meus recados">
              <StripedGrid items={items} elementsPerRow={1} />
            </BoxWithTitle>
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
