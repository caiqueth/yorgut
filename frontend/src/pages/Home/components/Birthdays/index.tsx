import React from "react";
import { Link } from "react-router-dom";

import { Box, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { dayjsFromDate } from "helpers/formatting/date";
import { StripedGrid } from "components/molecules/StripedGrid";
import { useMeFriends } from "api/me/queries";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";

interface ItemProps {
  id: number;
  firstName: string;
  profilePicture: string;
  birthDate: Date;
}

interface IBirthdaysProps { };
export const Birthdays: React.FC<IBirthdaysProps> = () => {
  const { data: friends, isLoading } = useMeFriends();

  const items =
    friends?.friends.map((f) => (
      <PicAndTextHorizontal
        title={f.firstName}
        text={dayjsFromDate(f.birthDate, "D [de] MMMM")}
        imageSrc={f.profilePicture}
        url={`/profile/${f.id}`}
        imageSize={56}
      />
    )) ?? [];

  return (
    <BoxWithTitle title={`próximos aniversários`} isLoading={isLoading}>
      {items.length > 0
        ? <StripedGrid items={items} elementsPerRow={1} />
        : (
          <Box sx={{ p: 1 }}>
            <Typography>Não há nenhum aniversário próximo.</Typography>
          </Box>
        )}
    </BoxWithTitle>
  );
};
