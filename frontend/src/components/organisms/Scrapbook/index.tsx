import {
  Box,
  Stack,
  Typography,
} from "@mui/material";
import { useCreateScrap, useScraps } from "api/scrap/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";

import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { TextCompositor } from "components/molecules/TextCompositor";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";

interface IProps {
  userId: number;
}
export const Scrapbook: React.FC<IProps> = ({ userId }) => {
  const { data: scraps, isLoading } = useScraps(userId);
  const { mutate: sendScrap, isPending } = useCreateScrap();

  const [scrapContent, setScrapContent] = useState<string>("");

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    sendScrap({
      content: scrapContent,
      userRecipientId: userId,
    });
  };

  const items =
    scraps?.scraps.map((s) => (
      <PicAndTextHorizontal
        imageSrc={s.authorProfilePicture}
        title={`${s.authorName} em ${dayjsFromDate(s.createdAt, "D [de] MMMM [de] YYYY [às] HH[h]MM")}: `}
        text={s.content}
        url={`/profile/${s.authorId}`}
      />
    )) ?? [];

  return (
    <Stack direction={"column"} gap={1}>
      <TextCompositor
        content={scrapContent}
        header={"Deixar um recado"}
        onChangeContent={setScrapContent}
        onSubmit={handleSubmit}
        isPending={isPending}
        maxLength={2048}
      />
      <BoxWithTitle title="recados">
        {items.length > 0
          ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
          : (
            <Box sx={{ p: 1 }}>
              <Typography>Não há nenhum recado.</Typography>
            </Box>
          )}
      </BoxWithTitle>
    </Stack>
  );
};
