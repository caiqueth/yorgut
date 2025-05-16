import { Box, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";

interface IProps {
  communityId: number;
}
export const HotTopics: React.FC<IProps> = ({ communityId }) => {
  const topics = {
    topics: [
      { id: 1, ownerProfilePicture: "", title: "Título do Tópico 1", replies: 5, lastReplyPreview: "oadjaspdapsjdp a", lastReplyDate: new Date() },
      { id: 2, ownerProfilePicture: "", title: "Título do Tópico 2", replies: 4, lastReplyPreview: "oadjasp das dsaikd9[ 1 0edjasdjasijd as] dapsasdkpas odkjsap[djsaidpjasi[ dasjdp a", lastReplyDate: new Date() },
      { id: 3, ownerProfilePicture: "", title: "Título do Tópico 3", replies: 31, lastReplyPreview: "oadjasp das da asd asdasdas dasdkqw daspsjdp a", lastReplyDate: new Date() },
      { id: 4, ownerProfilePicture: "", title: "Título do Tópico 4", replies: 993129, lastReplyPreview: "oadaslkdas[ dp das dapsjdp a", lastReplyDate: new Date() },
      { id: 5, ownerProfilePicture: "", title: "Título do Tópico 5", replies: 21, lastReplyPreview: "oadjasp das dapsjdp a", lastReplyDate: new Date() },
    ]
  }

  const items =
    topics?.topics.map((t) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        <img src={`/${t.ownerProfilePicture}`} height={40} />
        <Stack direction={"column"} sx={{ overflow: "hidden" }}>
          <Typography>{t.title}</Typography>
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