import { Box, Stack, TextField, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import React from "react";
import { useParams } from "react-router-dom";

interface IProps {}
export const UserTestimonialsList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const { data: testimonials } = {
    data: {
      testimonials: [
        {
          author: "Test Author 1",
          text: "Alo alo 1",
        },
        {
          author: "Test Author 2",
          text: "Alo alo 2",
        },
      ],
    },
  };

  const items = testimonials.testimonials.map((t) => (
    <Stack
      alignItems={"center"}
      flexDirection={"row"}
      sx={{ textAlign: "left" }}
      gap={1}
    >
      <Typography>{t.author}</Typography>
      <Typography>{t.text}</Typography>
    </Stack>
  ));

  return (
    <Box>
      <Box>
        <Stack direction={"column"}>
          Escrever depoimento:
          <TextField multiline rows={8} fullWidth />
          <Button disabled>Enviar</Button>
        </Stack>
      </Box>
      <BoxWithTitle title={"depoimentos recebidos"}>
        <StripedGrid items={items} elementsPerRow={1} />
      </BoxWithTitle>
    </Box>
  );
};
