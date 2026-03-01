import { Box, Stack, Typography } from "@mui/material";
import { useCommunityHotTopics } from "api/community/queries";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Link, useNavigate } from "react-router-dom";

interface IProps {
  communityId: number;
}
export const HotTopics: React.FC<IProps> = ({ communityId }) => {
  const { data: topics, isLoading } = useCommunityHotTopics(communityId);

  const navigate = useNavigate();

  const handleClickSeeAll = () => {
    navigate(`forum`)
  }

  const handleCreateTopic = () => {
    navigate(`forum/new-topic`)
  }

  const items =
    topics?.topics?.map((t) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={`${t.authorProfilePicture}`} height={40} />
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
            <Typography variant={"body2"}>{dayjsFromDate(t.lastReplyDate, "DD [de] MMM [de] YYYY")}:</Typography>
            <Box sx={{ textOverflow: "ellipsis" }}>
              <Typography variant={"body2"}>{t.lastReplyPreview}</Typography>
            </Box>
          </Stack>
        </Stack>
      </Stack>
    )) ?? [];

  return (
    <BoxWithTitle title={"fórum"}>
      <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
      <Stack direction={"row"} justifyContent={"end"}>
        <Button onClick={handleClickSeeAll}>ver todos</Button>
        <Button onClick={handleCreateTopic}>+ novo tópico</Button>
      </Stack>
    </BoxWithTitle>
  )
}
