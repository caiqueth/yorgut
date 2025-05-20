import { Box, Stack, Typography } from "@mui/material";
import { useCommunityHotTopics } from "api/community/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Link } from "react-router-dom";

interface IProps {
  communityId: number;
}
export const HotTopics: React.FC<IProps> = ({ communityId }) => {
  const { data: topics } = useCommunityHotTopics(communityId);

  const items =
    topics?.topics?.map((t) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={`/${t.authorProfilePicture}`} height={40} />
        <Stack direction={"column"} sx={{ overflow: "hidden" }}>
          <Link to={`forum/topic/${t.id}`}><Typography>{t.title}</Typography></Link>
          <Typography variant={"body2"}>{t.replies} respostas</Typography>
          <Stack
            direction={"row"}
            alignItems={"center"}
            spacing={1}
            sx={{
              maxWidth: "100%",
              whiteSpace: "nowrap",
            }}>
            <Typography variant={"body2"}>{dayjsFromDate(t.lastReplyDate, "DD [de] MM [de] YYYY")}:</Typography>
            <Box sx={{ textOverflow: "ellipsis" }}>
              <Typography variant={"body2"}>{t.lastReplyPreview}</Typography>
            </Box>
          </Stack>
        </Stack>
      </Stack>
    )) ?? [];

  return (
    <BoxWithTitle title={"fórum"}>
      <StripedGrid items={items} elementsPerRow={1} />
    </BoxWithTitle>
  )
}