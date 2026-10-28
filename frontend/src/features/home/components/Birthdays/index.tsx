import React from "react";

import { Box, Typography } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { dayjsFromDate } from "helpers/formatting/date";
import { StripedGrid } from "components/StripedGrid";
import { useMeBirthdays } from "api/me/queries";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";

interface IBirthdaysProps { };
export const Birthdays: React.FC<IBirthdaysProps> = () => {
  const { data, isLoading, isFetching } = useMeBirthdays();

  const items =
    data?.birthdays.map((f) => (
      <PicAndTextHorizontal
        key={f.id}
        title={f.firstName}
        text={
          f.daysUntil === 0
            ? "hoje! 🎂"
            : f.daysUntil === 1
              ? "amanhã"
              : dayjsFromDate(f.nextBirthday, "D [de] MMMM")
        }
        imageSrc={f.profilePicture}
        url={`/profile/${f.id}`}
        imageSize={56}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`próximos aniversários`} showProgress={isFetching && !isLoading}>
      {items.length === 0 && !isLoading
        ? (
          <Box sx={{ p: 1 }}>
            <Typography>Não há nenhum aniversário próximo.</Typography>
          </Box>
        ) : (
          <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton imageSize={56} lineWidths={[80, 96]} />} />
        )}
    </BoxWithTitle>
  );
};
